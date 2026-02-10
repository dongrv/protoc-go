// Package protoc provides a Go API for compiling Protocol Buffer files on Windows
// where wildcard patterns are not supported by the protoc command.
//
// This package solves the problem of compiling multiple .proto files in Windows
// by recursively finding all .proto files and constructing the appropriate
// protoc command with all files explicitly listed.
package protoc

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// ErrProtocNotFound is returned when the protoc command is not found in PATH.
var ErrProtocNotFound = errors.New("protoc command not found in PATH")

// ErrNoProtoFiles is returned when no .proto files are found in the specified directory.
var ErrNoProtoFiles = errors.New("no .proto files found")

// ErrPluginNotFound is returned when a required protoc plugin is not found.
type ErrPluginNotFound struct {
	Plugin string
}

func (e ErrPluginNotFound) Error() string {
	return fmt.Sprintf("protoc plugin %q not found in PATH", e.Plugin)
}

// Compiler provides a high-level API for compiling Protocol Buffer files.
type Compiler struct {
	// protoDir is the directory containing .proto files to compile.
	protoDir string

	// outputDir is the directory where generated files will be placed.
	outputDir string

	// protoPaths are additional include paths for protoc (-I flags).
	protoPaths []string

	// plugins specifies which protoc plugins to use.
	plugins []string

	// goOpts are options for the go plugin.
	goOpts []string

	// goGrpcOpts are options for the go-grpc plugin.
	goGrpcOpts []string

	// verbose enables verbose output to stdout.
	verbose bool

	// ctx is the context for cancellation and timeout.
	ctx context.Context

	// foundFiles caches the found .proto files.
	foundFiles []string

	// mu protects concurrent access to the compiler.
	mu sync.RWMutex

	// autoDetectImports enables automatic detection of import dependencies.
	autoDetectImports bool

	// additionalProtoPaths stores automatically detected include paths.
	additionalProtoPaths []string

	// smartFilter enables automatic filtering of imported-only files.
	smartFilter bool
}

// NewCompiler creates a new Compiler with default options.
func NewCompiler() *Compiler {
	return &Compiler{
		protoDir:          ".",
		outputDir:         ".",
		plugins:           []string{"go"},
		goOpts:            []string{"paths=source_relative"},
		goGrpcOpts:        []string{"paths=source_relative"},
		ctx:               context.Background(),
		autoDetectImports: true,
		smartFilter:       true,
	}
}

// WithProtoDir sets the directory containing .proto files.
func (c *Compiler) WithProtoDir(dir string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.protoDir = dir
	return c
}

// WithOutputDir sets the output directory for generated files.
func (c *Compiler) WithOutputDir(dir string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.outputDir = dir
	return c
}

// WithProtoPaths sets additional include paths for protoc.
func (c *Compiler) WithProtoPaths(paths ...string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.protoPaths = paths
	return c
}

// WithPlugins sets which protoc plugins to use.
func (c *Compiler) WithPlugins(plugins ...string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.plugins = plugins
	return c
}

// WithGoOpts sets options for the go plugin.
func (c *Compiler) WithGoOpts(opts ...string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.goOpts = opts
	return c
}

// WithGoGrpcOpts sets options for the go-grpc plugin.
func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.goGrpcOpts = opts
	return c
}

// WithVerbose enables verbose output.
func (c *Compiler) WithVerbose(verbose bool) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.verbose = verbose
	return c
}

// WithContext sets the context for cancellation and timeout.
func (c *Compiler) WithContext(ctx context.Context) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ctx = ctx
	return c
}

// WithAutoDetectImports enables or disables automatic import detection.
// When enabled (default), the compiler will automatically detect import
// dependencies and add necessary include paths.
func (c *Compiler) WithAutoDetectImports(enabled bool) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.autoDetectImports = enabled
	return c
}

// WithSmartFilter enables or disables smart file filtering.
// When enabled (default), the compiler will automatically filter out files
// that are only imported by other files, preventing duplicate compilation.
func (c *Compiler) WithSmartFilter(enabled bool) *Compiler {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.smartFilter = enabled
	return c
}

// FindFiles recursively finds all .proto files in the configured directory.
// This method can be called before Compile to inspect which files will be compiled.
func (c *Compiler) FindFiles() ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	absProtoDir, err := filepath.Abs(c.protoDir)
	if err != nil {
		return nil, fmt.Errorf("resolve proto directory: %w", err)
	}

	var files []string
	err = filepath.Walk(absProtoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), ".proto") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	c.foundFiles = files

	// If auto-detect is enabled, analyze imports
	if c.autoDetectImports && len(files) > 0 {
		if err := c.analyzeImports(files); err != nil {
			return files, fmt.Errorf("analyze imports: %w", err)
		}
	}

	return files, nil
}

// Compile compiles all found .proto files.
// If FindFiles hasn't been called, it will automatically find files first.
func (c *Compiler) Compile() (string, error) {
	c.mu.Lock()

	// Validate configuration
	if err := c.validate(); err != nil {
		c.mu.Unlock()
		return "", err
	}

	// Find files if not already found
	if len(c.foundFiles) == 0 {
		c.mu.Unlock()
		if _, err := c.FindFiles(); err != nil {
			return "", err
		}
		c.mu.Lock()
	}

	if len(c.foundFiles) == 0 {
		c.mu.Unlock()
		return "", ErrNoProtoFiles
	}

	// Apply smart filtering if enabled
	if c.smartFilter && len(c.foundFiles) > 1 {
		filteredFiles, err := c.filterImportedOnlyFiles()
		if err != nil {
			c.mu.Unlock()
			return "", fmt.Errorf("smart filter: %w", err)
		}

		if len(filteredFiles) > 0 {
			c.foundFiles = filteredFiles
		}
	}

	// Create output directory
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		c.mu.Unlock()
		return "", fmt.Errorf("create output directory: %w", err)
	}

	// Build command
	cmd := c.buildCommand()

	if c.verbose {
		fmt.Printf("Found %d .proto files:\n", len(c.foundFiles))
		for _, file := range c.foundFiles {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Printf("Executing: %s\n", strings.Join(cmd.Args, " "))
	}

	c.mu.Unlock()

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("protoc execution failed: %w", err)
	}

	if c.verbose && len(output) > 0 {
		fmt.Printf("protoc output: %s\n", output)
	}

	return string(output), nil
}

// validate validates the compiler configuration.
func (c *Compiler) validate() error {
	// Check protoc exists
	if _, err := exec.LookPath("protoc"); err != nil {
		return ErrProtocNotFound
	}

	// Check plugins exist
	for _, plugin := range c.plugins {
		switch plugin {
		case "go":
			if _, err := exec.LookPath("protoc-gen-go"); err != nil {
				return ErrPluginNotFound{Plugin: "protoc-gen-go"}
			}
		case "go-grpc":
			if _, err := exec.LookPath("protoc-gen-go-grpc"); err != nil {
				return ErrPluginNotFound{Plugin: "protoc-gen-go-grpc"}
			}
		}
	}

	// Resolve absolute paths
	if absPath, err := filepath.Abs(c.protoDir); err == nil {
		c.protoDir = absPath
	}

	if absPath, err := filepath.Abs(c.outputDir); err == nil {
		c.outputDir = absPath
	}

	return nil
}

// buildCommand builds the exec.Cmd for protoc.
func (c *Compiler) buildCommand() *exec.Cmd {
	args := []string{}

	// According to the optimization document best practice, we use only one -I parameter
	// The standard command format is: protoc -I <proto_root> --go_out=... <proto_files>
	// We use the proto directory as the single include path
	args = append(args, "-I", c.protoDir)

	// Add plugin outputs
	for _, plugin := range c.plugins {
		switch plugin {
		case "go":
			args = append(args, "--go_out="+buildPluginOpts("", c.goOpts, c.outputDir))
		case "go-grpc":
			args = append(args, "--go-grpc_out="+buildPluginOpts("", c.goGrpcOpts, c.outputDir))
		default:
			args = append(args, fmt.Sprintf("--%s_out=%s", plugin, c.outputDir))
		}
	}

	// Add all proto files with paths relative to the proto directory
	// This matches the standard command format from the optimization document
	for _, file := range c.foundFiles {
		// Get relative path from proto directory
		relPath, err := filepath.Rel(c.protoDir, file)
		if err != nil {
			// If we can't get relative path, the file is outside protoDir
			// In this case, we need to adjust our approach
			if c.verbose {
				fmt.Printf("Warning: File %s is outside proto directory %s\n", file, c.protoDir)
			}
			// For standard optimization, all files should be within protoDir
			// We'll use the absolute path as a fallback
			args = append(args, file)
		} else {
			// Use relative path as in the standard command format
			args = append(args, relPath)
		}
	}

	return exec.CommandContext(c.ctx, "protoc", args...)
}

// filterImportedOnlyFiles filters out files that are only imported by other files.
// This prevents duplicate compilation errors when a file is both imported
// and directly compiled.
func (c *Compiler) filterImportedOnlyFiles() ([]string, error) {
	if len(c.foundFiles) <= 1 {
		// No filtering needed for single file or empty
		return c.foundFiles, nil
	}

	// Parse imports from all files
	importMap := make(map[string][]string)  // file -> []imports
	fileImportCount := make(map[string]int) // file -> how many times it's imported

	for _, file := range c.foundFiles {
		imports, err := parseImportsFromFile(file)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		importMap[file] = imports

		// Count how many times each imported file is referenced
		for _, imp := range imports {
			// Convert import path to absolute path
			absImportPath, err := resolveImportPath(imp, file, c.protoDir)
			if err != nil {
				continue
			}
			fileImportCount[absImportPath]++
		}
	}

	// Identify files that should be kept (not filtered out)
	var filtered []string
	for _, file := range c.foundFiles {
		// Check if this file is imported by other files
		importCount := fileImportCount[file]

		// Always keep files that:
		// 1. Are not imported by any other file (importCount == 0)
		// 2. Have service definitions (likely main files)
		// 3. Have message definitions but are not imported

		hasService, _ := fileHasServiceDefinition(file)
		hasMessages, _ := fileHasMessageDefinitions(file)

		if importCount == 0 || hasService || (hasMessages && importCount == 0) {
			// This is likely a "main" file that should be compiled directly
			filtered = append(filtered, file)
		} else {
			// This file is only imported by others and doesn't have services
			// It will be compiled automatically through imports
			if c.verbose {
				fmt.Printf("Smart filter: Excluding %s (imported by %d other files)\n",
					filepath.Base(file), importCount)
			}
		}
	}

	// If filtering would remove all files, keep the original list
	if len(filtered) == 0 {
		return c.foundFiles, nil
	}

	return filtered, nil
}

// resolveImportPath resolves an import path to an absolute file path.
func resolveImportPath(importPath, sourceFile, protoDir string) (string, error) {
	// Try relative to source file directory first
	sourceDir := filepath.Dir(sourceFile)
	possiblePath := filepath.Join(sourceDir, importPath)

	if _, err := os.Stat(possiblePath); err == nil {
		absPath, err := filepath.Abs(possiblePath)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	// Try relative to proto directory
	possiblePath = filepath.Join(protoDir, importPath)
	if _, err := os.Stat(possiblePath); err == nil {
		absPath, err := filepath.Abs(possiblePath)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	return "", fmt.Errorf("cannot resolve import path: %s", importPath)
}

// fileHasServiceDefinition checks if a proto file contains service definitions.
func fileHasServiceDefinition(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Simple regex to check for service definitions
	serviceRegex := regexp.MustCompile(`service\s+\w+\s*{`)
	return serviceRegex.Match(content), nil
}

// fileHasMessageDefinitions checks if a proto file contains message definitions.
func fileHasMessageDefinitions(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Simple regex to check for message definitions
	messageRegex := regexp.MustCompile(`message\s+\w+\s*{`)
	return messageRegex.Match(content), nil
}

// buildPluginOpts builds the plugin options string.
func buildPluginOpts(prefix string, options []string, outputDir string) string {
	var opts []string

	if prefix != "" {
		opts = append(opts, prefix)
	}

	opts = append(opts, options...)

	optStr := ""
	if len(opts) > 0 {
		optStr = strings.Join(opts, ",") + ":"
	}

	return optStr + outputDir
}

// analyzeImports analyzes proto files to detect import dependencies.
// With the single -I parameter optimization, we don't add additional include paths.
// Instead, we validate that all imports can be resolved relative to the proto directory.
func (c *Compiler) analyzeImports(files []string) error {
	absProtoDir, err := filepath.Abs(c.protoDir)
	if err != nil {
		return fmt.Errorf("resolve proto directory: %w", err)
	}

	// Also check user-specified proto paths for import resolution
	allSearchPaths := []string{absProtoDir}
	for _, path := range c.protoPaths {
		absPath, err := filepath.Abs(path)
		if err == nil {
			allSearchPaths = append(allSearchPaths, absPath)
		}
	}

	// Collect all imports from all files
	allImports := make(map[string]bool)
	for _, file := range files {
		imports, err := parseImportsFromFile(file)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		for _, imp := range imports {
			allImports[imp] = true
		}
	}

	// For each import, try to find if it can be resolved
	importsFound := 0
	totalImports := len(allImports)

	for importPath := range allImports {
		found := false

		// Check all search paths
		for _, searchPath := range allSearchPaths {
			possiblePath := filepath.Join(searchPath, importPath)
			if _, err := os.Stat(possiblePath); err == nil {
				found = true
				importsFound++
				break
			}
		}

		if !found && c.verbose {
			fmt.Printf("Warning: Import %s not found in any search path\n", importPath)
		}
	}

	if c.verbose && totalImports > 0 {
		fmt.Printf("Import resolution: %d/%d imports found\n", importsFound, totalImports)
	}

	// Clear additionalProtoPaths since we're using single -I parameter approach
	c.additionalProtoPaths = []string{}

	return nil
}

// parseImportsFromFile parses import statements from a proto file.
func parseImportsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imports []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Remove comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}

		// Check for import statement
		if strings.Contains(line, "import") {
			// Try to match import pattern
			matches := regexp.MustCompile(`import\s+(?:"([^"]+)"|'([^']+)')`).FindStringSubmatch(line)
			if matches != nil {
				// matches[1] is for double quotes, matches[2] is for single quotes
				if matches[1] != "" {
					imports = append(imports, matches[1])
				} else if matches[2] != "" {
					imports = append(imports, matches[2])
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return imports, nil
}

// findImportDirectory tries to find the directory containing an imported file.
func findImportDirectory(importPath string, protoFiles []string, protoDir string) (string, error) {
	// First, check if the import path is relative to any of the proto files
	for _, protoFile := range protoFiles {
		protoFileDir := filepath.Dir(protoFile)
		possiblePath := filepath.Join(protoFileDir, importPath)

		// Check if the file exists
		if _, err := os.Stat(possiblePath); err == nil {
			return protoFileDir, nil
		}
	}

	// Check relative to proto directory
	possiblePath := filepath.Join(protoDir, importPath)
	if _, err := os.Stat(possiblePath); err == nil {
		return protoDir, nil
	}

	// If not found, try to find parent directories
	// For example, if import is "act/common.proto" and we're in "act7001",
	// we should check parent directories

	// Start from protoDir and go up
	currentDir := protoDir
	for {
		possiblePath := filepath.Join(currentDir, importPath)
		if _, err := os.Stat(possiblePath); err == nil {
			return currentDir, nil
		}

		// Go up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("cannot find directory for import: %s", importPath)
}
