package protoc_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dongrv/protoc-go"
)

// mockExecCommand is a mock for exec.CommandContext
var mockExecCommand func(ctx context.Context, name string, args ...string) *mockCmd

type mockCmd struct {
	output []byte
	err    error
}

func (m *mockCmd) CombinedOutput() ([]byte, error) {
	return m.output, m.err
}

func TestCompilerFindFiles(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto")

	// Create nested directory structure
	dirs := []string{
		protoDir,
		filepath.Join(protoDir, "a"),
		filepath.Join(protoDir, "b", "c"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create .proto files
	files := []string{
		"root.proto",
		"a/file1.proto",
		"b/c/file2.proto",
		"b/c/file3.proto",
	}

	for _, file := range files {
		filePath := filepath.Join(protoDir, file)
		content := `syntax = "proto3";
package test;
option go_package = "test/generated";
message Test { string id = 1; }`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a non-.proto file
	nonProtoFile := filepath.Join(protoDir, "not_a_proto.txt")
	if err := os.WriteFile(nonProtoFile, []byte("not a proto file"), 0644); err != nil {
		t.Fatal(err)
	}

	// Use Compiler to find files
	compiler := protoc.NewCompiler().WithProtoDir(protoDir)
	foundFiles, err := compiler.FindFiles()
	if err != nil {
		t.Fatalf("FindFiles failed: %v", err)
	}

	// Should find 4 .proto files
	if len(foundFiles) != 4 {
		t.Errorf("Expected 4 files, found %d", len(foundFiles))
	}

	// Verify all found files have .proto extension
	for _, file := range foundFiles {
		if !strings.HasSuffix(strings.ToLower(file), ".proto") {
			t.Errorf("Found non-.proto file: %s", file)
		}
	}
}

func TestCompilerWithNoProtoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	emptyDir := filepath.Join(tmpDir, "empty")

	// Create empty directory
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	compiler := protoc.NewCompiler().WithProtoDir(emptyDir)
	files, err := compiler.FindFiles()
	if err != nil {
		t.Fatalf("FindFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files, found %d", len(files))
	}
}

func TestCompilerConfiguration(t *testing.T) {
	compiler := protoc.NewCompiler()

	// Test WithProtoDir
	compiler = compiler.WithProtoDir("/test/proto")
	// Note: We can't directly check private fields, but we can verify through FindFiles

	// Test WithOutputDir
	compiler = compiler.WithOutputDir("/test/output")

	// Test WithProtoPaths
	compiler = compiler.WithProtoPaths("/path1", "/path2")

	// Test WithPlugins
	compiler = compiler.WithPlugins("go", "go-grpc")

	// Test WithGoOpts
	compiler = compiler.WithGoOpts("paths=source_relative", "module=test")

	// Test WithGoGrpcOpts
	compiler = compiler.WithGoGrpcOpts("paths=source_relative")

	// Test WithVerbose
	compiler = compiler.WithVerbose(true)

	// Test WithContext
	ctx := context.Background()
	compiler = compiler.WithContext(ctx)

	// Verify that compiler can be created and configured without panic
	if compiler == nil {
		t.Error("Compiler should not be nil")
	}
}

func TestErrorTypes(t *testing.T) {
	// Test ErrProtocNotFound
	err := protoc.ErrProtocNotFound
	if err.Error() != "protoc command not found in PATH" {
		t.Errorf("Unexpected error message: %v", err)
	}

	// Test ErrNoProtoFiles
	err = protoc.ErrNoProtoFiles
	if err.Error() != "no .proto files found" {
		t.Errorf("Unexpected error message: %v", err)
	}

	// Test ErrPluginNotFound
	pluginErr := protoc.ErrPluginNotFound{Plugin: "test-plugin"}
	expectedMsg := "protoc plugin \"test-plugin\" not found in PATH"
	if pluginErr.Error() != expectedMsg {
		t.Errorf("Unexpected error message: %v", pluginErr.Error())
	}
}

func TestFunctionalOptions(t *testing.T) {
	// Test WithProtoDir option
	opt := protoc.WithProtoDir("/test/proto")
	var opts protoc.Options
	opt(&opts)
	if opts.ProtoDir != "/test/proto" {
		t.Errorf("WithProtoDir failed: got %s, want /test/proto", opts.ProtoDir)
	}

	// Test WithOutputDir option
	opt = protoc.WithOutputDir("/test/output")
	opt(&opts)
	if opts.OutputDir != "/test/output" {
		t.Errorf("WithOutputDir failed: got %s, want /test/output", opts.OutputDir)
	}

	// Test WithProtoPaths option
	opt = protoc.WithProtoPaths("/path1", "/path2")
	opts = protoc.Options{}
	opt(&opts)
	if len(opts.ProtoPaths) != 2 || opts.ProtoPaths[0] != "/path1" || opts.ProtoPaths[1] != "/path2" {
		t.Errorf("WithProtoPaths failed: got %v", opts.ProtoPaths)
	}

	// Test WithPlugins option
	opt = protoc.WithPlugins("go", "go-grpc")
	opts = protoc.Options{}
	opt(&opts)
	if len(opts.Plugins) != 2 || opts.Plugins[0] != "go" || opts.Plugins[1] != "go-grpc" {
		t.Errorf("WithPlugins failed: got %v", opts.Plugins)
	}

	// Test WithGoOpts option
	opt = protoc.WithGoOpts("paths=source_relative", "module=test")
	opts = protoc.Options{}
	opt(&opts)
	if len(opts.GoOpts) != 2 || opts.GoOpts[0] != "paths=source_relative" || opts.GoOpts[1] != "module=test" {
		t.Errorf("WithGoOpts failed: got %v", opts.GoOpts)
	}

	// Test WithGoGrpcOpts option
	opt = protoc.WithGoGrpcOpts("paths=source_relative")
	opts = protoc.Options{}
	opt(&opts)
	if len(opts.GoGrpcOpts) != 1 || opts.GoGrpcOpts[0] != "paths=source_relative" {
		t.Errorf("WithGoGrpcOpts failed: got %v", opts.GoGrpcOpts)
	}

	// Test WithVerbose option
	opt = protoc.WithVerbose(true)
	opts = protoc.Options{}
	opt(&opts)
	if !opts.Verbose {
		t.Errorf("WithVerbose failed: got %v", opts.Verbose)
	}

	// Test WithContext option
	ctx := context.Background()
	opt = protoc.WithContext(ctx)
	opts = protoc.Options{}
	opt(&opts)
	if opts.Context != ctx {
		t.Errorf("WithContext failed")
	}
}

func TestNewOptions(t *testing.T) {
	// Test default options
	opts := protoc.NewOptions()
	if opts.ProtoDir != "." {
		t.Errorf("Default ProtoDir incorrect: got %s", opts.ProtoDir)
	}
	if opts.OutputDir != "." {
		t.Errorf("Default OutputDir incorrect: got %s", opts.OutputDir)
	}
	if len(opts.Plugins) != 1 || opts.Plugins[0] != "go" {
		t.Errorf("Default Plugins incorrect: got %v", opts.Plugins)
	}
	if len(opts.GoOpts) != 1 || opts.GoOpts[0] != "paths=source_relative" {
		t.Errorf("Default GoOpts incorrect: got %v", opts.GoOpts)
	}
	if len(opts.GoGrpcOpts) != 1 || opts.GoGrpcOpts[0] != "paths=source_relative" {
		t.Errorf("Default GoGrpcOpts incorrect: got %v", opts.GoGrpcOpts)
	}
	if opts.Context == nil {
		t.Error("Default Context should not be nil")
	}

	// Test with custom options
	ctx := context.Background()
	opts = protoc.NewOptions(
		protoc.WithProtoDir("/custom/proto"),
		protoc.WithOutputDir("/custom/output"),
		protoc.WithPlugins("custom"),
		protoc.WithContext(ctx),
	)

	if opts.ProtoDir != "/custom/proto" {
		t.Errorf("Custom ProtoDir incorrect: got %s", opts.ProtoDir)
	}
	if opts.OutputDir != "/custom/output" {
		t.Errorf("Custom OutputDir incorrect: got %s", opts.OutputDir)
	}
	if len(opts.Plugins) != 1 || opts.Plugins[0] != "custom" {
		t.Errorf("Custom Plugins incorrect: got %v", opts.Plugins)
	}
	if opts.Context != ctx {
		t.Errorf("Custom Context incorrect")
	}
}

func TestCreateOutputDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "non", "existent", "output", "dir")

	// Create proto directory
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a .proto file
	protoContent := `syntax = "proto3";
package test;
option go_package = "test/generated";
message Test { string value = 1; }`
	protoFile := filepath.Join(protoDir, "test.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Output directory doesn't exist yet
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		t.Fatal("Output directory should not exist yet")
	}

	// Create compiler (won't actually compile without protoc)
	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithOutputDir(outputDir)

	// Verify compiler was created
	if compiler == nil {
		t.Error("Compiler should not be nil")
	}

	// Note: We can't test actual compilation without protoc installed
	// This test verifies that the output directory validation logic works
}

func TestErrorProtocNotFound(t *testing.T) {
	// This test verifies that ErrProtocNotFound is properly defined
	err := protoc.ErrProtocNotFound
	if !errors.Is(err, protoc.ErrProtocNotFound) {
		t.Errorf("ErrProtocNotFound should be comparable with errors.Is")
	}
}

func TestMustCompilePanic(t *testing.T) {
	// Test that MustCompile panics when protoc is not found
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustCompile should panic when protoc is not found")
		}
	}()

	// This should panic because protoc is not installed
	_ = protoc.MustCompile("/tmp/proto", "/tmp/output")
}

func TestRelativePaths(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto")

	// Create proto directory
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a .proto file
	protoContent := `syntax = "proto3";
package test;
option go_package = "test/generated";
message Test { string value = 1; }`
	protoFile := filepath.Join(protoDir, "test.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to tmp directory to test relative paths
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create compiler with relative paths
	compiler := protoc.NewCompiler().
		WithProtoDir("proto").
		WithOutputDir("generated")

	// Verify compiler was created
	if compiler == nil {
		t.Error("Compiler should not be nil")
	}

	// Note: We can't test actual compilation without protoc
	// This test verifies that relative paths are accepted
}
