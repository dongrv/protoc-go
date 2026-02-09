# github.com/dongrv/protoc-go

A Go package for compiling Protocol Buffer files on Windows where wildcard patterns are not supported by the `protoc` command.

## Overview

This package solves the problem of compiling multiple `.proto` files in Windows by recursively finding all `.proto` files and constructing the appropriate `protoc` command with all files explicitly listed.

### Features

- ✅ **Recursive file discovery**: Automatically finds all `.proto` files in directory trees
- ✅ **Auto import detection**: Automatically detects import dependencies and adds necessary include paths
- ✅ **Multiple API styles**: Simple functions, functional options, and builder pattern
- ✅ **Plugin support**: Built-in support for `go` and `go-grpc` plugins
- ✅ **Custom options**: Flexible configuration for all protoc plugins
- ✅ **Context support**: Timeout and cancellation for long-running compilations
- ✅ **Error handling**: Comprehensive error types with clear messages
- ✅ **Cross-platform**: Works on Windows, Linux, and macOS
- ✅ **No external dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/dongrv/protoc-go
```

### Dependencies

This package requires the following tools to be installed and available in PATH:

1. **protoc** - Protocol Buffers compiler
   - Download from: https://github.com/protocolbuffers/protobuf/releases
   - Add to PATH environment variable

2. **Go plugins** (install with Go):
   ```bash
   # Install protoc-gen-go
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   
   # Install protoc-gen-go-grpc (for gRPC support)
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Quick Start

### Basic Usage

```go
import "github.com/dongrv/protoc-go"

func main() {
    // Simple API
    output, err := protoc.Compile("./proto", "./generated")
    if err != nil {
        log.Fatal(err)
    }
    
    // Functional options API
    output, err = protoc.CompileWith(
        protoc.WithProtoDir("./proto"),
        protoc.WithOutputDir("./generated"),
        protoc.WithPlugins("go", "go-grpc"),
        protoc.WithVerbose(true),
    )
    
    // Builder pattern API
    compiler := protoc.NewCompiler().
        WithProtoDir("./proto").
        WithOutputDir("./generated").
        WithPlugins("go", "go-grpc").
        WithVerbose(true)
    
    output, err = compiler.Compile()
}
```

## API Reference

### Simple Functions

```go
// Compile compiles .proto files with default options
func Compile(protoDir, outputDir string) (string, error)

// CompileWith compiles .proto files with functional options
func CompileWith(opts ...Option) (string, error)

// MustCompile compiles .proto files and panics on error
func MustCompile(protoDir, outputDir string) string

// MustCompileWith compiles with options and panics on error
func MustCompileWith(opts ...Option) string
```

### Functional Options

```go
// Option configures compilation options
type Option func(*Options)

// WithProtoDir sets the proto directory
func WithProtoDir(dir string) Option

// WithOutputDir sets the output directory
func WithOutputDir(dir string) Option

// WithProtoPaths sets additional proto include paths
func WithProtoPaths(paths ...string) Option

// WithPlugins sets which plugins to use
func WithPlugins(plugins ...string) Option

// WithGoOpts sets options for the go plugin
func WithGoOpts(opts ...string) Option

// WithGoGrpcOpts sets options for the go-grpc plugin
func WithGoGrpcOpts(opts ...string) Option

// WithVerbose enables verbose output
func WithVerbose(verbose bool) Option

// WithAutoDetectImports enables or disables automatic import detection
// When enabled (default), the compiler will automatically detect import
// dependencies and add necessary include paths.
func WithAutoDetectImports(enabled bool) Option

// WithContext sets the context for cancellation and timeout
func WithContext(ctx context.Context) Option
```

### Compiler Type (Builder Pattern)

```go
// Compiler provides a high-level API for compiling Protocol Buffer files
type Compiler struct { ... }

// NewCompiler creates a new Compiler with default options
func NewCompiler() *Compiler

// WithProtoDir sets the directory containing .proto files
func (c *Compiler) WithProtoDir(dir string) *Compiler

// WithOutputDir sets the output directory for generated files
func (c *Compiler) WithOutputDir(dir string) *Compiler

// WithProtoPaths sets additional include paths for protoc
func (c *Compiler) WithProtoPaths(paths ...string) *Compiler

// WithPlugins sets which protoc plugins to use
func (c *Compiler) WithPlugins(plugins ...string) *Compiler

// WithGoOpts sets options for the go plugin
func (c *Compiler) WithGoOpts(opts ...string) *Compiler

// WithGoGrpcOpts sets options for the go-grpc plugin
func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler

// WithVerbose enables verbose output
func (c *Compiler) WithVerbose(verbose bool) *Compiler

// WithAutoDetectImports enables or disables automatic import detection
// When enabled (default), the compiler will automatically detect import
// dependencies and add necessary include paths.
func (c *Compiler) WithAutoDetectImports(enabled bool) *Compiler

// WithContext sets the context for cancellation and timeout
func (c *Compiler) WithContext(ctx context.Context) *Compiler

// FindFiles recursively finds all .proto files in the configured directory
func (c *Compiler) FindFiles() ([]string, error)

// Compile compiles all found .proto files
func (c *Compiler) Compile() (string, error)
```

### Error Types

```go
// ErrProtocNotFound is returned when protoc command is not found
var ErrProtocNotFound = errors.New("protoc command not found in PATH")

// ErrNoProtoFiles is returned when no .proto files are found
var ErrNoProtoFiles = errors.New("no .proto files found")

// ErrPluginNotFound is returned when a required plugin is not found
type ErrPluginNotFound struct {
    Plugin string
}
```

## Examples

### Basic Compilation

```go
output, err := protoc.Compile("./proto", "./generated")
```

### With gRPC Support

```go
output, err := protoc.CompileWith(
    protoc.WithProtoDir("./proto"),
    protoc.WithOutputDir("./generated"),
    protoc.WithPlugins("go", "go-grpc"),
    protoc.WithGoOpts("paths=source_relative"),
    protoc.WithGoGrpcOpts("paths=source_relative"),
)
```

### With Custom Module Path

```go
output, err := protoc.CompileWith(
    protoc.WithProtoDir("./proto"),
    protoc.WithOutputDir("./generated"),
    protoc.WithGoOpts(
        "paths=source_relative",
        "module=github.com/yourusername/yourproject",
    ),
)
```

### With Auto Import Detection

```go
// Compile a subdirectory that imports files from parent directories
output, err := protoc.CompileWith(
    protoc.WithProtoDir("./subdir"),
    protoc.WithOutputDir("./generated"),
    protoc.WithAutoDetectImports(true), // Automatically find parent directories
)

// Disable auto import detection for manual control
output, err := protoc.CompileWith(
    protoc.WithProtoDir("./subdir"),
    protoc.WithOutputDir("./generated"),
    protoc.WithAutoDetectImports(false),
    protoc.WithProtoPaths(".."), // Manually add parent directory
)
```

### Using Builder Pattern

```go
compiler := protoc.NewCompiler().
    WithProtoDir("./proto").
    WithOutputDir("./generated").
    WithPlugins("go", "go-grpc").
    WithGoOpts("paths=source_relative").
    WithGoGrpcOpts("paths=source_relative").
    WithProtoPaths("./vendor/google/api").
    WithAutoDetectImports(true). // Enable auto import detection
    WithVerbose(true)

// Find files first
files, err := compiler.FindFiles()
if err != nil {
    log.Fatal(err)
}

// Then compile
output, err := compiler.Compile()
```

### With Context for Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

output, err := protoc.CompileWith(
    protoc.WithProtoDir("./proto"),
    protoc.WithOutputDir("./generated"),
    protoc.WithContext(ctx),
)
```

## Command Line Tool

The package includes a command-line tool in `cmd/protoc-go-compiler/`:

```bash
# Build the tool
go build -o protoc-go-compiler ./cmd/protoc-go-compiler

# Basic usage
protoc-go-compiler -proto-dir=./proto -output-dir=./generated

# With gRPC support
protoc-go-compiler -plugins=go,go-grpc

# With verbose output
protoc-go-compiler -verbose
```

### Integration Examples

#### In a Build Script

```go
// build.go
package main

import (
    "log"
    "github.com/dongrv/protoc-go"
)

func main() {
    if _, err := protoc.Compile("./proto", "./generated"); err != nil {
        log.Fatalf("Failed to compile proto files: %v", err)
    }
    log.Println("Proto files compiled successfully")
}
```

#### Compiling Subdirectories with Dependencies

```go
// Compile a specific subdirectory that imports from parent directories
output, err := protoc.CompileWith(
    protoc.WithProtoDir("./act/act7001"),
    protoc.WithOutputDir("./generated"),
    protoc.WithAutoDetectImports(true), // Automatically finds ../act directory
)
if err != nil {
    log.Fatalf("Failed to compile act7001: %v", err)
}
```

#### In a Makefile

```makefile
.PHONY: proto
proto:
    @echo "Compiling proto files..."
    @go run ./tools/build.go
    @echo "Proto compilation complete"

.PHONY: proto-subdir
proto-subdir:
    @echo "Compiling subdirectory with auto import detection..."
    @protoc-go-compiler -proto-dir=./act/act7001 -output-dir=./generated -auto-detect-imports=true
```

#### Using go:generate

```go
//go:generate go run github.com/dongrv/protoc-go/cmd/protoc-go-compiler -proto-dir=./proto -output-dir=./generated

// For subdirectories with imports
//go:generate go run github.com/dongrv/protoc-go/cmd/protoc-go-compiler -proto-dir=./act/act7001 -output-dir=./generated -auto-detect-imports=true
```

## Error Handling

```go
output, err := protoc.Compile("./proto", "./generated")
if err != nil {
    switch e := err.(type) {
    case *protoc.ErrPluginNotFound:
        log.Printf("Plugin not found: %s", e.Plugin)
    case *protoc.ErrProtocNotFound:
        log.Fatal("protoc not installed. Please install from: https://github.com/protocolbuffers/protobuf/releases")
    case *protoc.ErrNoProtoFiles:
        log.Fatal("No .proto files found in directory")
    default:
        log.Fatal(err)
    }
}
```

## Testing

Run the tests:

```bash
go test ./...
```

Run examples:

```bash
go run ./usage_example.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests: `go test ./...`
6. Submit a pull request

## License

MIT License

## See Also

- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf)
- [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)

## Why This Package?

On Windows, the `protoc` command doesn't support wildcard patterns like `*.proto` or `**/*.proto`. This makes it difficult to compile all `.proto` files in a directory tree. This package solves this problem by:

1. Recursively finding all `.proto` files
2. Building a `protoc` command with all files explicitly listed
3. Providing a clean Go API for integration into build systems
4. **Auto import detection**: Automatically finding and adding necessary include paths for imports

### Solving Import Dependencies

A common problem when compiling Protocol Buffer files is handling imports between directories. For example:
- Directory `act7001/` contains `act7001.proto` that imports `../act/common.proto`
- When compiling only `act7001/` directory, protoc cannot find the imported file

This package solves this with **auto import detection**:
- Automatically parses `.proto` files to find import statements
- Searches for imported files in parent and sibling directories
- Adds necessary include paths to the protoc command
- Works with nested directory structures

The package is designed to be:
- **Simple**: Easy to use with minimal configuration
- **Flexible**: Multiple API styles to suit different use cases
- **Smart**: Automatic import detection for complex dependency graphs
- **Robust**: Comprehensive error handling and validation
- **Standard**: Follows Go conventions and best practices