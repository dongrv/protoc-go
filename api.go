// Package protoc provides a Go API for compiling Protocol Buffer files on Windows
// where wildcard patterns are not supported by the protoc command.
//
// This package solves the problem of compiling multiple .proto files in Windows
// by recursively finding all .proto files and constructing the appropriate
// protoc command with all files explicitly listed.
package protoc

import (
	"context"
	"fmt"
)

// Compile is a convenience function that compiles .proto files with default options.
// It's a shorthand for creating a Compiler with default options and calling Compile.
func Compile(protoDir, outputDir string) (string, error) {
	return NewCompiler().
		WithProtoDir(protoDir).
		WithOutputDir(outputDir).
		Compile()
}

// CompileWithOptions compiles .proto files with the specified options.
func CompileWithOptions(opts Options) (string, error) {
	compiler := NewCompiler().
		WithProtoDir(opts.ProtoDir).
		WithOutputDir(opts.OutputDir).
		WithProtoPaths(opts.ProtoPaths...).
		WithPlugins(opts.Plugins...).
		WithGoOpts(opts.GoOpts...).
		WithGoGrpcOpts(opts.GoGrpcOpts...).
		WithVerbose(opts.Verbose)

	if opts.Context != nil {
		compiler = compiler.WithContext(opts.Context)
	}

	return compiler.Compile()
}

// Options provides a functional-style configuration for protoc compilation.
type Options struct {
	// ProtoDir is the directory containing .proto files to compile.
	// If empty, the current directory is used.
	ProtoDir string

	// OutputDir is the directory where generated files will be placed.
	// If empty, files are generated in the same directory as the .proto files.
	OutputDir string

	// ProtoPaths are additional include paths for protoc (-I flags).
	// The ProtoDir is always included as the first include path.
	ProtoPaths []string

	// Plugins specifies which protoc plugins to use.
	// Common values: "go", "go-grpc".
	Plugins []string

	// GoOpts are options for the go plugin.
	// Common options: "paths=source_relative".
	GoOpts []string

	// GoGrpcOpts are options for the go-grpc plugin.
	// Common options: "paths=source_relative".
	GoGrpcOpts []string

	// Verbose enables verbose output to stdout.
	Verbose bool

	// Context for cancellation and timeout.
	// If nil, context.Background() is used.
	Context context.Context
}

// Option is a function that configures Options.
type Option func(*Options)

// NewOptions creates a new Options with default values.
func NewOptions(opts ...Option) Options {
	options := Options{
		ProtoDir:   ".",
		OutputDir:  ".",
		Plugins:    []string{"go"},
		GoOpts:     []string{"paths=source_relative"},
		GoGrpcOpts: []string{"paths=source_relative"},
		Context:    context.Background(),
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// WithProtoDir sets the proto directory.
func WithProtoDir(dir string) Option {
	return func(o *Options) {
		o.ProtoDir = dir
	}
}

// WithOutputDir sets the output directory.
func WithOutputDir(dir string) Option {
	return func(o *Options) {
		o.OutputDir = dir
	}
}

// WithProtoPaths sets additional proto include paths.
func WithProtoPaths(paths ...string) Option {
	return func(o *Options) {
		o.ProtoPaths = paths
	}
}

// WithPlugins sets which plugins to use.
func WithPlugins(plugins ...string) Option {
	return func(o *Options) {
		o.Plugins = plugins
	}
}

// WithGoOpts sets options for the go plugin.
func WithGoOpts(opts ...string) Option {
	return func(o *Options) {
		o.GoOpts = opts
	}
}

// WithGoGrpcOpts sets options for the go-grpc plugin.
func WithGoGrpcOpts(opts ...string) Option {
	return func(o *Options) {
		o.GoGrpcOpts = opts
	}
}

// WithVerbose enables verbose output.
func WithVerbose(verbose bool) Option {
	return func(o *Options) {
		o.Verbose = verbose
	}
}

// WithContext sets the context.
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// CompileWith is a functional-style API for compiling .proto files.
// Example:
//
//	output, err := protoc.CompileWith(
//		protoc.WithProtoDir("./proto"),
//		protoc.WithOutputDir("./generated"),
//		protoc.WithPlugins("go", "go-grpc"),
//		protoc.WithVerbose(true),
//	)
func CompileWith(opts ...Option) (string, error) {
	options := NewOptions(opts...)
	return CompileWithOptions(options)
}

// MustCompile is like Compile but panics on error.
// Useful for initialization in tests or examples.
func MustCompile(protoDir, outputDir string) string {
	output, err := Compile(protoDir, outputDir)
	if err != nil {
		panic(fmt.Sprintf("protoc.MustCompile: %v", err))
	}
	return output
}

// MustCompileWith is like CompileWith but panics on error.
func MustCompileWith(opts ...Option) string {
	output, err := CompileWith(opts...)
	if err != nil {
		panic(fmt.Sprintf("protoc.MustCompileWith: %v", err))
	}
	return output
}
