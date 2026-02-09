// Command protoc-go-compiler is a command-line tool that wraps the github.com/dongrv/protoc-go
// package to provide a user-friendly interface for compiling Protocol Buffer files on Windows
// where wildcard patterns are not supported by the protoc command.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dongrv/protoc-go"
)

func main() {
	// Parse command line flags
	var (
		protoDir          string
		outputDir         string
		protoPaths        string
		plugins           string
		goOpts            string
		goGrpcOpts        string
		verbose           bool
		autoDetectImports bool
		showHelp          bool
		showVersion       bool
	)

	flag.StringVar(&protoDir, "proto-dir", ".", "Directory containing .proto files (default: current directory)")
	flag.StringVar(&protoDir, "p", ".", "Short form of -proto-dir")
	flag.StringVar(&outputDir, "output-dir", ".", "Output directory for generated files (default: current directory)")
	flag.StringVar(&outputDir, "o", ".", "Short form of -output-dir")
	flag.StringVar(&protoPaths, "proto-paths", "", "Additional proto include paths, comma-separated")
	flag.StringVar(&protoPaths, "I", "", "Short form of -proto-paths")
	flag.StringVar(&plugins, "plugins", "go", "Protoc plugins to use, comma-separated (e.g., 'go,go-grpc')")
	flag.StringVar(&goOpts, "go-opt", "paths=source_relative", "Options for go plugin, comma-separated")
	flag.StringVar(&goGrpcOpts, "go-grpc-opt", "paths=source_relative", "Options for go-grpc plugin, comma-separated")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Short form of -verbose")
	flag.BoolVar(&autoDetectImports, "auto-detect-imports", true, "Enable automatic import detection (default: true)")
	flag.BoolVar(&autoDetectImports, "a", true, "Short form of -auto-detect-imports")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Short form of -help")
	flag.BoolVar(&showVersion, "version", false, "Show version information")

	flag.Usage = func() {
		printUsage()
	}

	flag.Parse()

	// Handle help flag
	if showHelp {
		printUsage()
		return
	}

	// Handle version flag
	if showVersion {
		printVersion()
		return
	}

	// Build options from command line arguments
	opts := []protoc.Option{
		protoc.WithProtoDir(protoDir),
		protoc.WithOutputDir(outputDir),
		protoc.WithVerbose(verbose),
		protoc.WithAutoDetectImports(autoDetectImports),
	}

	// Parse comma-separated lists
	if protoPaths != "" {
		paths := splitCommaSeparated(protoPaths)
		opts = append(opts, protoc.WithProtoPaths(paths...))
	}

	if plugins != "" {
		pluginList := splitCommaSeparated(plugins)
		opts = append(opts, protoc.WithPlugins(pluginList...))
	}

	if goOpts != "" {
		goOptList := splitCommaSeparated(goOpts)
		opts = append(opts, protoc.WithGoOpts(goOptList...))
	}

	if goGrpcOpts != "" {
		goGrpcOptList := splitCommaSeparated(goGrpcOpts)
		opts = append(opts, protoc.WithGoGrpcOpts(goGrpcOptList...))
	}

	// Execute compilation
	output, err := protoc.CompileWith(opts...)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	// Print output if any
	if output != "" && verbose {
		fmt.Println(output)
	}

	if verbose {
		fmt.Println("✅ Compilation completed successfully")
	}
}

// splitCommaSeparated splits a comma-separated string into a slice of strings.
// It trims whitespace from each element and ignores empty elements.
func splitCommaSeparated(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// printUsage prints the usage information for the command.
func printUsage() {
	fmt.Fprintf(os.Stderr, `Protoc Go Compiler - A tool for compiling Protocol Buffer files on Windows

Usage:
  protoc-go-compiler [options]

Options:
  -p, -proto-dir string      Directory containing .proto files (default ".")
  -o, -output-dir string     Output directory for generated files (default ".")
  -I, -proto-paths string    Additional proto include paths, comma-separated
  -plugins string            Protoc plugins to use, comma-separated (default "go")
  -go-opt string             Options for go plugin, comma-separated (default "paths=source_relative")
  -go-grpc-opt string        Options for go-grpc plugin, comma-separated (default "paths=source_relative")
  -v, -verbose               Enable verbose output
  -a, -auto-detect-imports   Enable automatic import detection (default: true)
  -h, -help                  Show this help message
  -version                   Show version information

Examples:
  # Compile proto files in current directory
  protoc-go-compiler

  # Compile proto files in specific directory
  protoc-go-compiler -proto-dir=./proto -output-dir=./generated

  # Compile with gRPC support
  protoc-go-compiler -plugins=go,go-grpc

  # Compile with custom options and multiple include paths
  protoc-go-compiler \
    -proto-dir=./proto \
    -output-dir=./generated \
    -proto-paths=./proto,./vendor \
    -go-opt=paths=source_relative,module=github.com/example/project \
    -verbose

  # Disable auto import detection (use manual proto paths only)
  protoc-go-compiler -auto-detect-imports=false

  # Compile subdirectory that imports from parent directory
  protoc-go-compiler -proto-dir=./subdir -auto-detect-imports=true

Environment:
  This tool requires the following to be installed and available in PATH:
  - protoc (Protocol Buffers compiler)
  - protoc-gen-go (Go plugin for protoc)
  - protoc-gen-go-grpc (gRPC plugin for protoc, if using gRPC)

For more information, see: https://github.com/dongrv/protoc-go
`)
}

// printVersion prints version information.
func printVersion() {
	fmt.Println("protoc-go-compiler v1.0.0")
	fmt.Println("Using github.com/dongrv/protoc-go package")
	fmt.Println("Protocol Buffer compiler wrapper for Windows")
}
