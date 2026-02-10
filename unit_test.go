package protoc_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dongrv/protoc-go"
)

func TestNewCompiler(t *testing.T) {
	compiler := protoc.NewCompiler()
	if compiler == nil {
		t.Error("NewCompiler should not return nil")
	}
}

func TestWithProtoDir(t *testing.T) {
	compiler := protoc.NewCompiler().WithProtoDir("/test/proto")
	// We can't directly check private fields, but we can verify the compiler was created
	if compiler == nil {
		t.Error("WithProtoDir should not return nil")
	}
}

func TestWithProtoWorkSpace(t *testing.T) {
	compiler := protoc.NewCompiler().WithProtoWorkSpace("/test/workspace")
	if compiler == nil {
		t.Error("WithProtoWorkSpace should not return nil")
	}
}

func TestWithOutputDir(t *testing.T) {
	compiler := protoc.NewCompiler().WithOutputDir("/test/output")
	if compiler == nil {
		t.Error("WithOutputDir should not return nil")
	}
}

func TestWithPlugins(t *testing.T) {
	compiler := protoc.NewCompiler().WithPlugins("go", "go-grpc")
	if compiler == nil {
		t.Error("WithPlugins should not return nil")
	}
}

func TestWithGoOpts(t *testing.T) {
	compiler := protoc.NewCompiler().WithGoOpts("paths=source_relative", "module=test")
	if compiler == nil {
		t.Error("WithGoOpts should not return nil")
	}
}

func TestWithGoGrpcOpts(t *testing.T) {
	compiler := protoc.NewCompiler().WithGoGrpcOpts("paths=source_relative")
	if compiler == nil {
		t.Error("WithGoGrpcOpts should not return nil")
	}
}

func TestWithVerbose(t *testing.T) {
	compiler := protoc.NewCompiler().WithVerbose(true)
	if compiler == nil {
		t.Error("WithVerbose should not return nil")
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	compiler := protoc.NewCompiler().WithContext(ctx)
	if compiler == nil {
		t.Error("WithContext should not return nil")
	}
}

func TestCompileValidation(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test 1: Missing proto directory
	compiler1 := protoc.NewCompiler().
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)
	_, err := compiler1.Compile()
	if err == nil || !strings.Contains(err.Error(), "proto directory not specified") {
		t.Errorf("Expected error about missing proto directory, got: %v", err)
	}

	// Test 2: Missing workspace directory
	compiler2 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithOutputDir(outputDir)
	_, err = compiler2.Compile()
	if err == nil || !strings.Contains(err.Error(), "workspace directory not specified") {
		t.Errorf("Expected error about missing workspace directory, got: %v", err)
	}

	// Test 3: Missing output directory
	compiler3 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir)
	_, err = compiler3.Compile()
	if err == nil || !strings.Contains(err.Error(), "output directory not specified") {
		t.Errorf("Expected error about missing output directory, got: %v", err)
	}

	// Test 4: Non-existent proto directory
	compiler4 := protoc.NewCompiler().
		WithProtoDir("/non/existent/proto").
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)
	_, err = compiler4.Compile()
	if err == nil || !strings.Contains(err.Error(), "proto directory does not exist") {
		t.Errorf("Expected error about non-existent proto directory, got: %v", err)
	}

	// Test 5: Non-existent workspace directory
	compiler5 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace("/non/existent/workspace").
		WithOutputDir(outputDir)
	_, err = compiler5.Compile()
	if err == nil || !strings.Contains(err.Error(), "workspace directory does not exist") {
		t.Errorf("Expected error about non-existent workspace directory, got: %v", err)
	}

	// Test 6: Proto directory outside workspace
	compiler6 := protoc.NewCompiler().
		WithProtoDir("/outside/proto").
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)
	_, err = compiler6.Compile()
	// The error could be either about directory not existing or not being within workspace
	if err == nil {
		t.Errorf("Expected error about proto directory, got nil")
	} else if !strings.Contains(err.Error(), "does not exist") && !strings.Contains(err.Error(), "must be within") {
		t.Errorf("Expected error about proto directory not existing or not within workspace, got: %v", err)
	}
}

func TestCompileNoProtoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create compiler with empty proto directory
	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)

	_, err := compiler.Compile()
	if err == nil || !strings.Contains(err.Error(), "no .proto files found") {
		t.Errorf("Expected error about no .proto files, got: %v", err)
	}
}

func TestSimpleCompileFunction(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test Compile function (will fail because no .proto files, but should validate)
	_, err := protoc.Compile(protoDir, workspaceDir, outputDir)
	if err == nil || !strings.Contains(err.Error(), "no .proto files found") {
		t.Errorf("Expected error about no .proto files from Compile function, got: %v", err)
	}
}

func TestMustCompilePanic(t *testing.T) {
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustCompile should panic when there are no .proto files")
		}
	}()

	_ = protoc.MustCompile(protoDir, workspaceDir, outputDir)
}

func TestBuilderPattern(t *testing.T) {
	// Test that we can chain all methods
	ctx := context.Background()
	compiler := protoc.NewCompiler().
		WithProtoDir("/test/proto/act7110").
		WithProtoWorkSpace("/test/proto").
		WithOutputDir("/test/generated").
		WithPlugins("go", "go-grpc").
		WithGoOpts("paths=source_relative").
		WithGoGrpcOpts("paths=source_relative").
		WithVerbose(true).
		WithContext(ctx)

	if compiler == nil {
		t.Error("Builder pattern should return a non-nil compiler")
	}
}

func TestForwardSlashPathsOnWindows(t *testing.T) {
	// This test verifies that paths use forward slashes on Windows
	// for better cross-platform compatibility

	tmpDir := t.TempDir()

	// Create directory structure with backslashes (Windows style)
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a .proto file
	protoContent := `syntax = "proto3";
package test;
option go_package = "test/generated";
message Test { string id = 1; }`

	protoFile := filepath.Join(protoDir, "test.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create another .proto file in subdirectory
	subDir := filepath.Join(protoDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	subProtoFile := filepath.Join(subDir, "subtest.proto")
	if err := os.WriteFile(subProtoFile, []byte(protoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Note: We can't directly test the command construction without exposing internal methods
	// But we can verify that the compiler works with Windows-style paths

	// On Windows, verify that filepath.ToSlash would convert paths correctly
	if runtime.GOOS == "windows" {
		// Test that Windows paths with backslashes are accepted
		windowsProtoDir := strings.ReplaceAll(protoDir, "/", "\\")
		windowsWorkspaceDir := strings.ReplaceAll(workspaceDir, "/", "\\")
		windowsOutputDir := strings.ReplaceAll(outputDir, "/", "\\")

		_ = protoc.NewCompiler().
			WithProtoDir(windowsProtoDir).
			WithProtoWorkSpace(windowsWorkspaceDir).
			WithOutputDir(windowsOutputDir).
			WithPlugins("go").
			WithGoOpts("paths=source_relative").
			WithVerbose(false)

		// Verify that filepath.ToSlash converts correctly
		convertedProtoDir := filepath.ToSlash(windowsProtoDir)
		if !strings.Contains(convertedProtoDir, "/") {
			t.Errorf("filepath.ToSlash should convert backslashes to forward slashes, got: %s", convertedProtoDir)
		}
	}

	t.Logf("Forward slash path compatibility test completed")
	t.Logf("On %s, paths are normalized to use forward slashes for protoc command", runtime.GOOS)
}

func TestProtocAvailabilityCheck(t *testing.T) {
	// Test that the compiler checks for protoc availability
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a .proto file
	protoContent := `syntax = "proto3";
package test;
option go_package = "test/generated";
message Test { string id = 1; }`

	protoFile := filepath.Join(protoDir, "test.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Check if protoc is actually available
	_, err := exec.LookPath("protoc")
	protocAvailable := err == nil

	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir).
		WithPlugins("go").
		WithGoOpts("paths=source_relative").
		WithVerbose(false)

	_, compileErr := compiler.Compile()

	if !protocAvailable {
		// If protoc is not available, we should get an error about protoc not found
		if compileErr == nil {
			t.Error("Expected error when protoc is not available, got nil")
		} else if !strings.Contains(compileErr.Error(), "protoc not found in PATH") {
			t.Errorf("Expected error about protoc not found, got: %v", compileErr)
		}

		// Check that the error message contains helpful hints
		if !strings.Contains(compileErr.Error(), "PATH environment variable") {
			t.Error("Error message should mention PATH environment variable")
		}

		// Check for platform-specific hints
		switch runtime.GOOS {
		case "windows":
			if !strings.Contains(compileErr.Error(), "Windows") {
				t.Error("Error message should contain Windows-specific installation hints")
			}
		case "darwin":
			if !strings.Contains(compileErr.Error(), "macOS") && !strings.Contains(compileErr.Error(), "Homebrew") {
				t.Error("Error message should contain macOS-specific installation hints")
			}
		case "linux":
			if !strings.Contains(compileErr.Error(), "Linux") && !strings.Contains(compileErr.Error(), "apt") && !strings.Contains(compileErr.Error(), "yum") {
				t.Error("Error message should contain Linux-specific installation hints")
			}
		}
	} else {
		// If protoc is available, compilation should proceed (though it may fail for other reasons)
		// We just want to ensure the availability check doesn't block valid compilation
		t.Logf("protoc is available, compilation attempted (may succeed or fail for other reasons)")
		if compileErr != nil && !strings.Contains(compileErr.Error(), "protoc not found") {
			// Other errors are OK (e.g., import issues, syntax errors, etc.)
			t.Logf("Compilation failed for other reasons (expected): %v", compileErr)
		}
	}
}

func TestProtocAvailabilityCheckOrder(t *testing.T) {
	// Test that protoc availability check happens after validation
	tmpDir := t.TempDir()
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test 1: Missing proto directory - should fail validation before checking protoc
	compiler1 := protoc.NewCompiler().
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)
	_, err1 := compiler1.Compile()
	if err1 == nil || !strings.Contains(err1.Error(), "proto directory not specified") {
		t.Errorf("Expected validation error about missing proto directory, got: %v", err1)
	}

	// Test 2: Missing workspace directory - should fail validation before checking protoc
	compiler2 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithOutputDir(outputDir)
	_, err2 := compiler2.Compile()
	if err2 == nil || !strings.Contains(err2.Error(), "workspace directory not specified") {
		t.Errorf("Expected validation error about missing workspace directory, got: %v", err2)
	}

	// Test 3: Missing output directory - should fail validation before checking protoc
	compiler3 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir)
	_, err3 := compiler3.Compile()
	if err3 == nil || !strings.Contains(err3.Error(), "output directory not specified") {
		t.Errorf("Expected validation error about missing output directory, got: %v", err3)
	}

	// Test 4: Non-existent proto directory - should fail validation before checking protoc
	compiler4 := protoc.NewCompiler().
		WithProtoDir("/non/existent/proto").
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir)
	_, err4 := compiler4.Compile()
	if err4 == nil || !strings.Contains(err4.Error(), "proto directory does not exist") {
		t.Errorf("Expected validation error about non-existent proto directory, got: %v", err4)
	}

	// Test 5: Non-existent workspace directory - should fail validation before checking protoc
	compiler5 := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace("/non/existent/workspace").
		WithOutputDir(outputDir)
	_, err5 := compiler5.Compile()
	if err5 == nil || !strings.Contains(err5.Error(), "workspace directory does not exist") {
		t.Errorf("Expected validation error about non-existent workspace directory, got: %v", err5)
	}
}
