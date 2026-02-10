// Package protoc provides a Go API for compiling Protocol Buffer files.
//
// This package solves the problem of compiling multiple .proto files by
// recursively finding all .proto files in a directory and constructing
// the appropriate protoc command with all files explicitly listed.
//
// The package implements the optimized standard command format:
//   - Single -I parameter using workspace directory
//   - Relative file paths from workspace to proto files
//   - Explicit listing of all .proto files to compile
//
// # Overview
//
// The package provides a builder pattern API through the Compiler type:
//
//	compiler := protoc.NewCompiler().
//	    WithProtoDir("./proto/act7110").
//	    WithProtoWorkSpace("./proto").
//	    WithOutputDir("./generated").
//	    WithPlugins("go", "go-grpc").
//	    WithGoOpts("paths=source_relative").
//	    WithVerbose(true)
//
//	output, err := compiler.Compile()
//
// # Installation
//
//	go get github.com/dongrv/protoc-go
//
// # Dependencies
//
// This package requires the following tools to be installed and available in PATH:
//
//   - protoc: Protocol Buffers compiler
//   - protoc-gen-go: Go plugin for protoc
//   - protoc-gen-go-grpc: gRPC plugin for protoc (optional, for gRPC support)
//
// # Quick Start
//
// Basic usage:
//
//	import "github.com/dongrv/protoc-go"
//
//	func main() {
//	    // Simple API
//	    output, err := protoc.Compile(
//	        "./proto/act7110",    // proto directory
//	        "./proto",            // workspace directory
//	        "./generated",        // output directory
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Builder pattern API
//	    compiler := protoc.NewCompiler().
//	        WithProtoDir("./proto/act7110").
//	        WithProtoWorkSpace("./proto").
//	        WithOutputDir("./generated").
//	        WithPlugins("go", "go-grpc").
//	        WithGoOpts("paths=source_relative").
//	        WithGoGrpcOpts("paths=source_relative").
//	        WithVerbose(true)
//
//	    output, err = compiler.Compile()
//	}
//
// # API Reference
//
// ## Compiler Type
//
//	type Compiler struct { ... }
//
//	func NewCompiler() *Compiler
//	func (c *Compiler) WithProtoDir(dir string) *Compiler
//	func (c *Compiler) WithProtoWorkSpace(dir string) *Compiler
//	func (c *Compiler) WithOutputDir(dir string) *Compiler
//	func (c *Compiler) WithPlugins(plugins ...string) *Compiler
//	func (c *Compiler) WithGoOpts(opts ...string) *Compiler
//	func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler
//	func (c *Compiler) WithVerbose(verbose bool) *Compiler
//	func (c *Compiler) WithContext(ctx context.Context) *Compiler
//	func (c *Compiler) Compile() (string, error)
//
// ## Simple Functions
//
//	func Compile(protoDir, workspaceDir, outputDir string) (string, error)
//	func MustCompile(protoDir, workspaceDir, outputDir string) string
//
// # Examples
//
// ## Basic Compilation
//
//	output, err := protoc.Compile(
//	    "./proto/act7110",
//	    "./proto",
//	    "./generated",
//	)
//
// ## With gRPC Support
//
//	compiler := protoc.NewCompiler().
//	    WithProtoDir("./proto/act7110").
//	    WithProtoWorkSpace("./proto").
//	    WithOutputDir("./generated").
//	    WithPlugins("go", "go-grpc").
//	    WithGoOpts("paths=source_relative").
//	    WithGoGrpcOpts("paths=source_relative")
//
//	output, err := compiler.Compile()
//
// ## With Custom Options
//
//	compiler := protoc.NewCompiler().
//	    WithProtoDir("./proto/act7110").
//	    WithProtoWorkSpace("./proto").
//	    WithOutputDir("./generated").
//	    WithPlugins("go").
//	    WithGoOpts("paths=source_relative", "module=github.com/example/project")
//
//	output, err := compiler.Compile()
//
// ## Using Context for Timeout
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	compiler := protoc.NewCompiler().
//	    WithProtoDir("./proto/act7110").
//	    WithProtoWorkSpace("./proto").
//	    WithOutputDir("./generated").
//	    WithContext(ctx)
//
//	output, err := compiler.Compile()
//
// # Error Handling
//
// The package returns descriptive error messages for common issues:
//
//   - "proto directory not specified"
//   - "workspace directory not specified"
//   - "output directory not specified"
//   - "proto directory does not exist"
//   - "workspace directory does not exist"
//   - "proto directory must be within workspace directory"
//   - "no .proto files found in [directory]"
//   - "protoc execution failed: [error]"
//
// # Notes
//
//   - This package is particularly useful on Windows where protoc doesn't support
//     wildcard patterns like `*.proto`.
//   - The package automatically creates the output directory if it doesn't exist.
//   - Uses the optimized standard command format with single -I parameter to prevent
//     "already defined" errors.
//   - All .proto files are specified with paths relative to the workspace directory.
//   - The proto directory must be within the workspace directory.
//   - Context support allows for cancellation and timeout of long-running compilations.
//   - Implements the exact optimization recommended in best practices documentation.
//
// # See Also
//
//   - Protocol Buffers: https://developers.google.com/protocol-buffers
//   - protoc-gen-go: https://pkg.go.dev/google.golang.org/protobuf
//   - protoc-gen-go-grpc: https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
package protoc
