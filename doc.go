// Package protoc provides a Go API for compiling Protocol Buffer files on Windows
// where wildcard patterns are not supported by the protoc command.
//
// This package solves the problem of compiling multiple .proto files in Windows
// by recursively finding all .proto files and constructing the appropriate
// protoc command with all files explicitly listed.
//
// The package implements the optimized standard command format from best practices:
//   - Single -I parameter: Uses only one include path (proto root directory)
//   - Relative file paths: All .proto files are specified with paths relative to the -I directory
//   - Smart file filtering: Filters out imported-only files to prevent duplicate compilation
//   - Auto import validation: Validates that all imports can be resolved relative to the include directory
//
// # Overview
//
// The package provides three levels of API:
//
// 1. Simple function API: Compile, CompileWith, MustCompile
// 2. Functional options API: WithProtoDir, WithOutputDir, etc.
// 3. Builder pattern API: Compiler type with method chaining
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
//	    output, err := protoc.Compile("./proto", "./generated")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Functional options API
//	    output, err = protoc.CompileWith(
//	        protoc.WithProtoDir("./proto"),
//	        protoc.WithOutputDir("./generated"),
//	        protoc.WithPlugins("go", "go-grpc"),
//	        protoc.WithVerbose(true),
//	    )
//
//	    // Builder pattern API
//	    compiler := protoc.NewCompiler().
//	        WithProtoDir("./proto").
//	        WithOutputDir("./generated").
//	        WithPlugins("go", "go-grpc").
//	        WithVerbose(true)
//
//	    output, err = compiler.Compile()
//	}
//
// # API Reference
//
// ## Simple Functions
//
//	func Compile(protoDir, outputDir string) (string, error)
//	func CompileWith(opts ...Option) (string, error)
//	func MustCompile(protoDir, outputDir string) string
//	func MustCompileWith(opts ...Option) string
//
// ## Functional Options
//
//	func WithProtoDir(dir string) Option
//	func WithOutputDir(dir string) Option
//	func WithProtoPaths(paths ...string) Option
//	func WithPlugins(plugins ...string) Option
//	func WithGoOpts(opts ...string) Option
//	func WithGoGrpcOpts(opts ...string) Option
//	func WithVerbose(verbose bool) Option
//	func WithContext(ctx context.Context) Option
//
// ## Compiler Type
//
//	type Compiler struct { ... }
//
//	func NewCompiler() *Compiler
//	func (c *Compiler) WithProtoDir(dir string) *Compiler
//	func (c *Compiler) WithOutputDir(dir string) *Compiler
//	func (c *Compiler) WithProtoPaths(paths ...string) *Compiler
//	func (c *Compiler) WithPlugins(plugins ...string) *Compiler
//	func (c *Compiler) WithGoOpts(opts ...string) *Compiler
//	func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler
//	func (c *Compiler) WithVerbose(verbose bool) *Compiler
//	func (c *Compiler) WithContext(ctx context.Context) *Compiler
//	func (c *Compiler) FindFiles() ([]string, error)
//	func (c *Compiler) Compile() (string, error)
//
// # Examples
//
// ## Basic Compilation
//
//	output, err := protoc.Compile("./proto", "./generated")
//
// ## With gRPC Support
//
//	output, err := protoc.CompileWith(
//	    protoc.WithProtoDir("./proto"),
//	    protoc.WithOutputDir("./generated"),
//	    protoc.WithPlugins("go", "go-grpc"),
//	)
//
// ## With Custom Options
//
//	output, err := protoc.CompileWith(
//	    protoc.WithProtoDir("./proto"),
//	    protoc.WithOutputDir("./generated"),
//	    protoc.WithGoOpts("paths=source_relative", "module=github.com/example/project"),
//	    protoc.WithGoGrpcOpts("paths=source_relative"),
//	    protoc.WithProtoPaths("./proto", "./vendor"),
//	    protoc.WithVerbose(true),
//	)
//
// ## Using Builder Pattern
//
//	compiler := protoc.NewCompiler().
//	    WithProtoDir("./proto").
//	    WithOutputDir("./generated").
//	    WithPlugins("go", "go-grpc").
//	    WithGoOpts("paths=source_relative").
//	    WithGoGrpcOpts("paths=source_relative").
//	    WithProtoPaths("./vendor/google/api").
//	    WithVerbose(true)
//
//	files, err := compiler.FindFiles()
//	output, err := compiler.Compile()
//
// # Error Handling
//
// The package defines the following error types:
//
//	var ErrProtocNotFound = errors.New("protoc command not found in PATH")
//	var ErrNoProtoFiles = errors.New("no .proto files found")
//
//	type ErrPluginNotFound struct {
//	    Plugin string
//	}
//
// # Notes
//
//   - This package is particularly useful on Windows where protoc doesn't support
//     wildcard patterns like `*.proto`.
//   - The package automatically creates the output directory if it doesn't exist.
//   - Uses the optimized standard command format with single -I parameter to prevent
//     "already defined" errors.
//   - All .proto files are specified with paths relative to the include directory.
//   - The package validates that required tools (protoc, plugins) are in PATH.
//   - Context support allows for cancellation and timeout of long-running compilations.
//   - Implements the exact optimization recommended in best practices documentation.
//
// # See Also
//
//   - Protocol Buffers: https://developers.google.com/protocol-buffers
//   - protoc-gen-go: https://pkg.go.dev/google.golang.org/protobuf
//   - protoc-gen-go-grpc: https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
package protoc
