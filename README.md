# github.com/dongrv/protoc-go

A Go package for compiling Protocol Buffer files using the optimized standard command format.

## Overview

This package provides a clean Go API for compiling Protocol Buffer files by implementing the optimized standard command format:

```bash
protoc -I <workspace_dir> --go_out=paths=source_relative:<output_dir> <relative_proto_files>
```

### Key Features

- ✅ **Standard command format**: Implements the optimized single `-I` parameter approach
- ✅ **Recursive file discovery**: Automatically finds all `.proto` files in a directory
- ✅ **Builder pattern API**: Clean, chainable configuration methods
- ✅ **Plugin support**: Built-in support for `go` and `go-grpc` plugins
- ✅ **Custom options**: Flexible configuration for all protoc plugins
- ✅ **Context support**: Timeout and cancellation for long-running compilations
- ✅ **Validation**: Comprehensive validation of paths and configuration
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
    // Simple function API
    output, err := protoc.Compile(
        "./proto/act7110",    // Directory containing .proto files
        "./proto",            // Workspace directory for -I parameter
        "./generated",        // Output directory
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Builder pattern API
    compiler := protoc.NewCompiler().
        WithProtoDir("./proto/act7110").
        WithProtoWorkSpace("./proto").
        WithOutputDir("./generated").
        WithPlugins("go", "go-grpc").
        WithGoOpts("paths=source_relative").
        WithGoGrpcOpts("paths=source_relative").
        WithVerbose(true)
    
    output, err = compiler.Compile()
}
```

## API Reference

### Compiler Type (Builder Pattern)

```go
// Compiler provides a high-level API for compiling Protocol Buffer files
type Compiler struct { ... }

// NewCompiler creates a new Compiler with default options
func NewCompiler() *Compiler

// WithProtoDir sets the directory containing .proto files to compile
func (c *Compiler) WithProtoDir(dir string) *Compiler

// WithProtoWorkSpace sets the workspace directory for the -I parameter
func (c *Compiler) WithProtoWorkSpace(dir string) *Compiler

// WithOutputDir sets the output directory for generated files
func (c *Compiler) WithOutputDir(dir string) *Compiler

// WithPlugins sets which protoc plugins to use
func (c *Compiler) WithPlugins(plugins ...string) *Compiler

// WithGoOpts sets options for the go plugin
func (c *Compiler) WithGoOpts(opts ...string) *Compiler

// WithGoGrpcOpts sets options for the go-grpc plugin
func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler

// WithVerbose enables verbose output
func (c *Compiler) WithVerbose(verbose bool) *Compiler

// WithContext sets the context for cancellation and timeout
func (c *Compiler) WithContext(ctx context.Context) *Compiler

// Compile compiles all .proto files in the configured directory
func (c *Compiler) Compile() (string, error)
```

### Simple Functions

```go
// Compile is a convenience function that compiles .proto files
func Compile(protoDir, workspaceDir, outputDir string) (string, error)

// MustCompile is like Compile but panics on error
func MustCompile(protoDir, workspaceDir, outputDir string) string
```

## Examples

### Basic Compilation

```go
output, err := protoc.Compile(
    "./proto/act7110",
    "./proto",
    "./generated",
)
```

### With gRPC Support

```go
compiler := protoc.NewCompiler().
    WithProtoDir("./proto/act7110").
    WithProtoWorkSpace("./proto").
    WithOutputDir("./generated").
    WithPlugins("go", "go-grpc").
    WithGoOpts("paths=source_relative").
    WithGoGrpcOpts("paths=source_relative")

output, err := compiler.Compile()
```

### With Custom Options

```go
compiler := protoc.NewCompiler().
    WithProtoDir("./proto/act7110").
    WithProtoWorkSpace("./proto").
    WithOutputDir("./generated").
    WithPlugins("go").
    WithGoOpts("paths=source_relative", "module=github.com/example/project").
    WithVerbose(true)

output, err := compiler.Compile()
```

### Using Context for Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

compiler := protoc.NewCompiler().
    WithProtoDir("./proto/act7110").
    WithProtoWorkSpace("./proto").
    WithOutputDir("./generated").
    WithContext(ctx)

output, err := compiler.Compile()
```

### Optimization Document Example

This example matches the exact scenario from the optimization document:

```go
// Directory structure:
// D:\work\go\src\shengyou\docs\branches\beta\proto\act7110\*.proto
// D:\work\go\src\shengyou\docs\branches\beta\proto\ (workspace)
// D:\work\go\src\shengyou\server\branches\beta\protocol\ (output)

compiler := protoc.NewCompiler().
    WithProtoDir("D:\\work\\go\\src\\shengyou\\docs\\branches\\beta\\proto\\act7110").
    WithProtoWorkSpace("D:\\work\\go\\src\\shengyou\\docs\\branches\\beta\\proto").
    WithOutputDir("D:\\work\\go\\src\\shengyou\\server\\branches\\beta\\protocol").
    WithPlugins("go").
    WithGoOpts("paths=source_relative")

output, err := compiler.Compile()

// Generates the optimized command:
// protoc -I D:\work\go\src\shengyou\docs\branches\beta\proto \
//   --go_out=paths=source_relative:D:\work\go\src\shengyou\server\branches\beta\protocol \
//   act7110/act7110.proto act7110/debug.proto act7110/enum.proto
```

## Error Handling

The package returns descriptive error messages for common issues:

```go
output, err := protoc.Compile("./proto/act7110", "./proto", "./generated")
if err != nil {
    // Common errors include:
    // - "proto directory not specified"
    // - "workspace directory not specified"
    // - "output directory not specified"
    // - "proto directory does not exist"
    // - "workspace directory does not exist"
    // - "proto directory must be within workspace directory"
    // - "no .proto files found in [directory]"
    // - "protoc execution failed: [error]"
    log.Fatal(err)
}
```

## Standard Command Format Optimization

### The Problem

When compiling Protocol Buffer files, a common issue is the "already defined" error that occurs when duplicate `-I` parameters are used:

```bash
# Problematic command (causes "already defined" errors):
protoc -I D:\proto\act7110 -I D:\proto --go_out=... D:\proto\act7110\enum.proto
```

### The Solution

This package implements the optimized standard command format with a single `-I` parameter:

```bash
# Optimized command (no duplicate -I parameters):
protoc -I D:\proto --go_out=... act7110/enum.proto
```

### Key Optimization Principles

1. **Single -I parameter**: Only the workspace directory is specified with `-I`
2. **Relative file paths**: All `.proto` files are listed with paths relative to the workspace directory
3. **Unified output**: All generated files go to a single output directory
4. **Path validation**: The proto directory must be within the workspace directory

## Design Philosophy

### Simplicity

The API is designed to be intuitive and easy to use. The builder pattern provides a clean, chainable interface that makes configuration straightforward.

### Standards Compliance

The package implements the exact command format recommended in best practices documentation, ensuring compatibility and preventing common compilation errors.

### Robustness

Comprehensive validation ensures that configuration errors are caught early with clear, descriptive error messages.

### Performance

By using the optimized single `-I` parameter approach, the package eliminates the performance overhead and compilation errors associated with duplicate include paths.

## Testing

Run the tests:

```bash
go test ./...
```

Run examples:

```bash
go test -v -run Example ./...
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