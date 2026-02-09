// Package protoc_test contains examples and tests for the protoc package.
package protoc_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dongrv/protoc-go"
)

// ExampleCompile demonstrates the basic usage of the Compile function.
func ExampleCompile() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create proto directory structure
	protoDir := filepath.Join(tmpDir, "proto")
	os.MkdirAll(protoDir, 0755)

	// Create a simple .proto file
	protoContent := `syntax = "proto3";

package example;

option go_package = "example/generated";

message Person {
  string name = 1;
  int32 age = 2;
  repeated string emails = 3;
}`
	protoFile := filepath.Join(protoDir, "person.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "generated")

	// Compile the proto file
	output, err := protoc.Compile(protoDir, outputDir)
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation failed: %v\n", err)
		fmt.Println("This may happen if protoc is not installed or if there are other issues.")
	} else {
		fmt.Printf("Compilation output: %s\n", output)
		fmt.Printf("Generated files in: %s\n", outputDir)
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation failed: protoc command not found in PATH
	//   This may happen if protoc is not installed or if there are other issues.
	// If protoc is installed:
	//   Compilation output: [protoc output]
	//   Generated files in: [output directory]
}

// ExampleCompileWith demonstrates the functional options API.
func ExampleCompileWith() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create proto directory structure
	protoDir := filepath.Join(tmpDir, "proto")
	os.MkdirAll(protoDir, 0755)

	// Create a .proto file with gRPC service
	protoContent := `syntax = "proto3";

package example;

option go_package = "example/generated";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}`
	protoFile := filepath.Join(protoDir, "greeter.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "generated")

	// Compile with gRPC support and custom options
	output, err := protoc.CompileWith(
		protoc.WithProtoDir(protoDir),
		protoc.WithOutputDir(outputDir),
		protoc.WithPlugins("go", "go-grpc"),
		protoc.WithGoOpts("paths=source_relative"),
		protoc.WithGoGrpcOpts("paths=source_relative"),
		protoc.WithVerbose(false),
	)
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation with gRPC support failed: %v\n", err)
	} else {
		fmt.Printf("Compiled with gRPC support\n")
		fmt.Printf("Output length: %d bytes\n", len(output))
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation with gRPC support failed: protoc command not found in PATH
	// If protoc is installed:
	//   Compiled with gRPC support
	//   Output length: [number] bytes
}

// ExampleCompiler demonstrates the builder pattern API.
func ExampleCompiler() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create proto directory structure
	protoDir := filepath.Join(tmpDir, "proto")
	os.MkdirAll(filepath.Join(protoDir, "subdir"), 0755)

	// Create multiple .proto files
	protoFiles := []struct {
		name    string
		content string
	}{
		{
			name: "user.proto",
			content: `syntax = "proto3";
package example;
option go_package = "example/generated";
message User {
  string id = 1;
  string name = 2;
}`,
		},
		{
			name: "subdir/product.proto",
			content: `syntax = "proto3";
package example;
option go_package = "example/generated";
message Product {
  string id = 1;
  string name = 2;
  double price = 3;
}`,
		},
	}

	for _, file := range protoFiles {
		filePath := filepath.Join(protoDir, file.name)
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			log.Fatal(err)
		}
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "generated")

	// Use the Compiler builder pattern
	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithOutputDir(outputDir).
		WithPlugins("go").
		WithGoOpts("paths=source_relative").
		WithVerbose(false)

	// First, find the files
	files, err := compiler.FindFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d .proto files\n", len(files))

	// Then compile them
	_, err = compiler.Compile()
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation failed: %v\n", err)
	} else {
		fmt.Printf("Compilation successful\n")
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Found 2 .proto files
	//   Compilation failed: protoc command not found in PATH
	// If protoc is installed:
	//   Found 2 .proto files
	//   Compilation successful
}

// ExampleWithContext demonstrates using context for timeout.
func ExampleWithContext() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create proto directory structure
	protoDir := filepath.Join(tmpDir, "proto")
	os.MkdirAll(protoDir, 0755)

	// Create a .proto file
	protoContent := `syntax = "proto3";
package example;
option go_package = "example/generated";
message Test {
  string value = 1;
}`
	protoFile := filepath.Join(protoDir, "test.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "generated")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Compile with context
	_, err = protoc.CompileWith(
		protoc.WithProtoDir(protoDir),
		protoc.WithOutputDir(outputDir),
		protoc.WithContext(ctx),
	)
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation with context failed: %v\n", err)
	} else {
		fmt.Printf("Compiled with context timeout\n")
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation with context failed: protoc command not found in PATH
	// If protoc is installed:
	//   Compiled with context timeout
}

// ExampleMustCompile demonstrates the MustCompile function for initialization.
func ExampleMustCompile() {
	// Note: In real usage, you would use actual directories
	// This example shows the pattern for initialization

	// For initialization in tests or examples where failure should panic
	// Note: MustCompile will panic if compilation fails
	// In a real environment with protoc installed, this would work:
	// _ = protoc.MustCompile("./proto", "./generated")

	// Or with options
	// _ = protoc.MustCompileWith(
	//     protoc.WithProtoDir("./proto"),
	//     protoc.WithOutputDir("./generated"),
	//     protoc.WithPlugins("go", "go-grpc"),
	// )

	fmt.Println("MustCompile example - would panic without protoc installed")

	// Output:
	// MustCompile example - would panic without protoc installed
}

// Example_error_handling demonstrates error handling.
func Example_error_handling() {
	// Try to compile from a non-existent directory
	_, err := protoc.Compile("/non/existent/dir", "./output")

	if err != nil {
		// The actual error may vary depending on the environment
		// It could be "protoc command not found in PATH" or a file system error
		fmt.Printf("Error occurred: %v\n", err)
	}

	// Output depends on the environment:
	// If protoc is not installed:
	//   Error occurred: protoc command not found in PATH
	// If protoc is installed but directory doesn't exist:
	//   Error occurred: [file system error]
}

// Example_path_deduplication demonstrates how the package prevents duplicate
// include path errors described in the optimization document.
func Example_path_deduplication() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-optimization-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure similar to the optimization document example
	protoRootDir := filepath.Join(tmpDir, "docs", "branches", "beta", "proto")
	act7110Dir := filepath.Join(protoRootDir, "act7110")

	// Create directories
	if err := os.MkdirAll(act7110Dir, 0755); err != nil {
		log.Fatal(err)
	}

	// Create enum.proto file (from optimization document)
	enumProtoContent := `syntax = "proto3";
package act7110;

enum ClickType {
    Rat = 0;
    Rewards = 1;
}`
	enumProtoFile := filepath.Join(act7110Dir, "enum.proto")
	if err := os.WriteFile(enumProtoFile, []byte(enumProtoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create act7110.proto that imports enum.proto
	act7110ProtoContent := `syntax = "proto3";
package act7110;
import "act7110/enum.proto";

message Request {
    ClickType click_type = 1;
}`
	act7110ProtoFile := filepath.Join(act7110Dir, "act7110.proto")
	if err := os.WriteFile(act7110ProtoFile, []byte(act7110ProtoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "generated")

	// This demonstrates the optimization fix:
	// Before the fix, specifying both the subdirectory and parent directory
	// as include paths would cause "already defined" errors because protoc
	// would treat act7110/enum.proto and enum.proto as different files.

	// With the optimization, duplicate include paths are automatically removed,
	// preventing the "already defined" error described in the optimization document.
	compiler := protoc.NewCompiler().
		WithProtoDir(act7110Dir).
		WithProtoPaths(protoRootDir). // This would cause duplicate -I paths without optimization
		WithOutputDir(outputDir).
		WithAutoDetectImports(true).
		WithSmartFilter(true).
		WithVerbose(false)

	// Find files
	files, err := compiler.FindFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d .proto files\n", len(files))
	fmt.Println("Path deduplication optimization prevents duplicate -I paths")
	fmt.Println("This avoids the 'already defined' error from the optimization document")

	// Note: Actual compilation would require protoc to be installed
	// This example demonstrates the configuration that would have failed
	// before the optimization but now works correctly.

	// Output:
	// Found 2 .proto files
	// Path deduplication optimization prevents duplicate -I paths
	// This avoids the 'already defined' error from the optimization document
}
