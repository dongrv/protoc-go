// Package protoc provides Protocol Buffer compilation functionality.
package protoc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// compilerImpl is the internal implementation of the compiler.
type compilerImpl struct {
	protoDir     string
	workspaceDir string
	outputDir    string
	plugins      []string
	goOpts       []string
	goGrpcOpts   []string
	verbose      bool
	ctx          context.Context

	mu sync.Mutex
}

// compile implements the main compilation logic.
func (c *compilerImpl) compile() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Validate configuration
	if err := c.validate(); err != nil {
		return "", err
	}

	// Find all .proto files in the proto directory
	files, err := c.findProtoFiles()
	if err != nil {
		return "", fmt.Errorf("find proto files: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no .proto files found in %s", c.protoDir)
	}

	// Create output directory
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	// Build and execute protoc command
	cmd := c.buildCommand(files)

	if c.verbose {
		fmt.Printf("Found %d .proto files:\n", len(files))
		for _, file := range files {
			relPath, _ := filepath.Rel(c.workspaceDir, file)
			fmt.Printf("  - %s\n", relPath)
		}
		fmt.Printf("Executing: %s\n", strings.Join(cmd.Args, " "))
	}

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

// validate checks the compiler configuration.
func (c *compilerImpl) validate() error {
	if c.protoDir == "" {
		return fmt.Errorf("proto directory not specified")
	}

	if c.workspaceDir == "" {
		return fmt.Errorf("workspace directory not specified")
	}

	if c.outputDir == "" {
		return fmt.Errorf("output directory not specified")
	}

	// Check if proto directory exists
	if _, err := os.Stat(c.protoDir); os.IsNotExist(err) {
		return fmt.Errorf("proto directory does not exist: %s", c.protoDir)
	}

	// Check if workspace directory exists
	if _, err := os.Stat(c.workspaceDir); os.IsNotExist(err) {
		return fmt.Errorf("workspace directory does not exist: %s", c.workspaceDir)
	}

	// Verify proto directory is within workspace directory
	relPath, err := filepath.Rel(c.workspaceDir, c.protoDir)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("proto directory %s must be within workspace directory %s",
			c.protoDir, c.workspaceDir)
	}

	return nil
}

// findProtoFiles recursively finds all .proto files in the proto directory.
func (c *compilerImpl) findProtoFiles() ([]string, error) {
	var files []string

	absProtoDir, err := filepath.Abs(c.protoDir)
	if err != nil {
		return nil, fmt.Errorf("resolve proto directory: %w", err)
	}

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

	return files, nil
}

// buildCommand constructs the protoc command with the found files.
func (c *compilerImpl) buildCommand(files []string) *exec.Cmd {
	args := []string{}

	// Add workspace directory as single -I parameter
	// Use forward slashes for better cross-platform compatibility
	workspacePath := filepath.ToSlash(c.workspaceDir)
	args = append(args, "-I", workspacePath)

	// Add plugin outputs
	for _, plugin := range c.plugins {
		switch plugin {
		case "go":
			outputPath := filepath.ToSlash(c.outputDir)
			args = append(args, "--go_out="+buildPluginOpts("", c.goOpts, outputPath))
		case "go-grpc":
			outputPath := filepath.ToSlash(c.outputDir)
			args = append(args, "--go-grpc_out="+buildPluginOpts("", c.goGrpcOpts, outputPath))
		default:
			outputPath := filepath.ToSlash(c.outputDir)
			args = append(args, fmt.Sprintf("--%s_out=%s", plugin, outputPath))
		}
	}

	// Add all proto files with paths relative to workspace directory
	// Use forward slashes for better cross-platform compatibility
	for _, file := range files {
		relPath, err := filepath.Rel(c.workspaceDir, file)
		if err != nil {
			// This shouldn't happen since we validated the paths
			if c.verbose {
				fmt.Printf("Warning: cannot get relative path for %s: %v\n", file, err)
			}
			// Use forward slash for absolute paths too
			filePath := filepath.ToSlash(file)
			args = append(args, filePath)
		} else {
			// Convert relative path to use forward slashes
			relPath = filepath.ToSlash(relPath)
			args = append(args, relPath)
		}
	}

	return exec.CommandContext(c.ctx, "protoc", args...)
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
