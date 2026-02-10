// Package protoc provides a Go API for compiling Protocol Buffer files.
//
// This package solves the problem of compiling multiple .proto files by
// recursively finding all .proto files and constructing the appropriate
// protoc command with all files explicitly listed.
//
// The package implements the optimized standard command format:
//   - Single -I parameter using workspace directory
//   - Relative file paths from workspace to proto files
//   - Explicit listing of all .proto files to compile
package protoc

import (
	"context"
	"fmt"
)

// Compiler provides a high-level API for compiling Protocol Buffer files.
type Compiler struct {
	protoDir     string // Directory containing .proto files to compile
	workspaceDir string // Workspace directory for -I parameter
	outputDir    string // Output directory for generated files
	plugins      []string
	goOpts       []string
	goGrpcOpts   []string
	verbose      bool
	ctx          context.Context
}

// NewCompiler creates a new Compiler with default options.
func NewCompiler() *Compiler {
	return &Compiler{
		plugins:    []string{"go"},
		goOpts:     []string{"paths=source_relative"},
		goGrpcOpts: []string{"paths=source_relative"},
		ctx:        context.Background(),
	}
}

// WithProtoDir sets the directory containing .proto files to compile.
// The compiler will recursively find all .proto files in this directory.
func (c *Compiler) WithProtoDir(dir string) *Compiler {
	c.protoDir = dir
	return c
}

// WithProtoWorkSpace sets the workspace directory for the -I parameter.
// This should be the root directory containing all proto imports.
func (c *Compiler) WithProtoWorkSpace(dir string) *Compiler {
	c.workspaceDir = dir
	return c
}

// WithOutputDir sets the output directory for generated files.
func (c *Compiler) WithOutputDir(dir string) *Compiler {
	c.outputDir = dir
	return c
}

// WithPlugins sets which protoc plugins to use.
func (c *Compiler) WithPlugins(plugins ...string) *Compiler {
	c.plugins = plugins
	return c
}

// WithGoOpts sets options for the go plugin.
func (c *Compiler) WithGoOpts(opts ...string) *Compiler {
	c.goOpts = opts
	return c
}

// WithGoGrpcOpts sets options for the go-grpc plugin.
func (c *Compiler) WithGoGrpcOpts(opts ...string) *Compiler {
	c.goGrpcOpts = opts
	return c
}

// WithVerbose enables verbose output.
func (c *Compiler) WithVerbose(verbose bool) *Compiler {
	c.verbose = verbose
	return c
}

// WithContext sets the context for cancellation and timeout.
func (c *Compiler) WithContext(ctx context.Context) *Compiler {
	c.ctx = ctx
	return c
}

// Compile compiles all .proto files in the configured directory.
func (c *Compiler) Compile() (string, error) {
	if c.protoDir == "" {
		return "", fmt.Errorf("proto directory not specified")
	}
	if c.workspaceDir == "" {
		return "", fmt.Errorf("workspace directory not specified")
	}
	if c.outputDir == "" {
		return "", fmt.Errorf("output directory not specified")
	}

	// Create a new compiler instance to avoid mutating the original
	compiler := &compilerImpl{
		protoDir:     c.protoDir,
		workspaceDir: c.workspaceDir,
		outputDir:    c.outputDir,
		plugins:      c.plugins,
		goOpts:       c.goOpts,
		goGrpcOpts:   c.goGrpcOpts,
		verbose:      c.verbose,
		ctx:          c.ctx,
	}

	return compiler.compile()
}

// Compile is a convenience function that compiles .proto files with default options.
func Compile(protoDir, workspaceDir, outputDir string) (string, error) {
	return NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir).
		Compile()
}

// MustCompile is like Compile but panics on error.
func MustCompile(protoDir, workspaceDir, outputDir string) string {
	output, err := Compile(protoDir, workspaceDir, outputDir)
	if err != nil {
		panic(fmt.Sprintf("protoc.MustCompile: %v", err))
	}
	return output
}
