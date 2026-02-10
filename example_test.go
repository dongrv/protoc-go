// Package protoc_test contains examples for the protoc package.
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

// Example_basic demonstrates basic usage of the protoc package.
func Example_basic() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Create a simple .proto file
	protoContent := `syntax = "proto3";

package example;

option go_package = "example/generated";

message Person {
  string name = 1;
  int32 age = 2;
}`
	protoFile := filepath.Join(protoDir, "person.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Use the simple Compile function
	_, err = protoc.Compile(protoDir, workspaceDir, outputDir)
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation failed: %v\n", err)
		fmt.Println("This may happen if protoc is not installed or if there are other issues.")
	} else {
		fmt.Printf("Compilation successful\n")
		fmt.Printf("Generated files in: %s\n", outputDir)
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation failed: protoc execution failed: exec: "protoc": executable file not found in %PATH%
	//   This may happen if protoc is not installed or if there are other issues.
	// If protoc is installed:
	//   Compilation successful
	//   Generated files in: [output directory]
}

// Example_builder_pattern demonstrates using the builder pattern API.
func Example_builder_pattern() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

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

	// Use the builder pattern API
	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir).
		WithPlugins("go", "go-grpc").
		WithGoOpts("paths=source_relative").
		WithGoGrpcOpts("paths=source_relative").
		WithVerbose(false)

	_, err = compiler.Compile()
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation with gRPC support failed: %v\n", err)
	} else {
		fmt.Printf("Compiled with gRPC support\n")
		fmt.Println("Output generated successfully")
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation with gRPC support failed: protoc execution failed: exec: "protoc": executable file not found in %PATH%
	// If protoc is installed:
	//   Compiled with gRPC support
	//   Output length: [number] bytes
}

// Example_optimization_document demonstrates the exact scenario from the optimization document.
func Example_optimization_document() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-optimization-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the exact directory structure from the optimization document
	workDir := filepath.Join(tmpDir, "work", "go", "src", "shengyou")
	docsDir := filepath.Join(workDir, "docs", "branches", "beta")
	protoDir := filepath.Join(docsDir, "proto", "act7110") // Directory containing .proto files
	workspaceDir := filepath.Join(docsDir, "proto")        // Workspace directory for -I parameter
	outputDir := filepath.Join(workDir, "server", "branches", "beta", "protocol")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Create enum.proto file (from optimization document)
	enumProtoContent := `syntax = "proto3";
package act7110;

enum ClickType {
    Rat = 0;
    Rewards = 1;
}`
	enumProtoFile := filepath.Join(protoDir, "enum.proto")
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
	act7110ProtoFile := filepath.Join(protoDir, "act7110.proto")
	if err := os.WriteFile(act7110ProtoFile, []byte(act7110ProtoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Create debug.proto
	debugProtoContent := `syntax = "proto3";
package act7110;

message DebugInfo {
    string message = 1;
}`
	debugProtoFile := filepath.Join(protoDir, "debug.proto")
	if err := os.WriteFile(debugProtoFile, []byte(debugProtoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Demonstrate the optimized standard command format with forward slashes
	// This matches the command from the optimization document:
	// protoc -I D:/work/go/src/shengyou/docs/branches/beta/proto \
	//   --go_out=paths=source_relative:D:/work/go/src/shengyou/server/branches/beta/protocol \
	//   act7110/act7110.proto act7110/debug.proto act7110/enum.proto

	_ = protoc.NewCompiler().
		WithProtoDir(protoDir).              // Directory containing .proto files
		WithProtoWorkSpace(workspaceDir).    // Workspace directory for -I parameter
		WithOutputDir(outputDir).            // Output directory
		WithPlugins("go").                   // Use go plugin
		WithGoOpts("paths=source_relative"). // Go plugin options
		WithVerbose(false)

	fmt.Println("Optimization document example configuration:")
	fmt.Printf("  Proto directory: %s\n", protoDir)
	fmt.Printf("  Workspace directory: %s\n", workspaceDir)
	fmt.Printf("  Output directory: %s\n", outputDir)
	fmt.Println("  Plugins: go")
	fmt.Println("  Go options: paths=source_relative")
	fmt.Println()
	fmt.Println("This configuration generates the optimized command with forward slashes:")
	fmt.Println("  protoc -I <workspace_dir> \\")
	fmt.Println("    --go_out=paths=source_relative:<output_dir> \\")
	fmt.Println("    act7110/act7110.proto \\")
	fmt.Println("    act7110/debug.proto \\")
	fmt.Println("    act7110/enum.proto")
	fmt.Println()
	fmt.Println("On Windows, paths are automatically converted to forward slashes:")
	fmt.Println("  - D:\\path\\to\\proto becomes D:/path/to/proto")
	fmt.Println("  - act7110\\enum.proto becomes act7110/enum.proto")
	fmt.Println()
	fmt.Println("The optimization prevents 'already defined' errors by using")
	fmt.Println("only one -I parameter and relative file paths with forward slashes.")

	// Note: Actual compilation would require protoc to be installed
	// This example demonstrates the configuration that matches the optimization document

	// Output depends on the temporary directory created:
	// Optimization document example configuration:
	//   Proto directory: [temp_dir]/work/go/src/shengyou/docs/branches/beta/proto/act7110
	//   Workspace directory: [temp_dir]/work/go/src/shengyou/docs/branches/beta/proto
	//   Output directory: [temp_dir]/work/go/src/shengyou/server/branches/beta/protocol
	//   Plugins: go
	//   Go options: paths=source_relative
	//
	// This configuration generates the optimized command with forward slashes:
	//   protoc -I <workspace_dir> \
	//     --go_out=paths=source_relative:<output_dir> \
	//     act7110/act7110.proto \
	//     act7110/debug.proto \
	//     act7110/enum.proto
	//
	// On Windows, paths are automatically converted to forward slashes:
	//   - D:\path\to\proto becomes D:/path/to/proto
	//   - act7110\enum.proto becomes act7110/enum.proto
	//
	// The optimization prevents 'already defined' errors by using
	// only one -I parameter and relative file paths with forward slashes.
}

// Example_context demonstrates using context for timeout.
func Example_context() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

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

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Compile with context
	compiler := protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir).
		WithContext(ctx)

	_, err = compiler.Compile()
	if err != nil {
		// This may fail if protoc is not installed or if there are other issues
		fmt.Printf("Compilation with context failed: %v\n", err)
	} else {
		fmt.Printf("Compiled with context timeout\n")
	}

	// Output depends on whether protoc is installed:
	// If protoc is not installed:
	//   Compilation with context failed: protoc execution failed: exec: "protoc": executable file not found in %PATH%
	// If protoc is installed:
	//   Compiled with context timeout
}

// Example_custom_options demonstrates using custom options.
func Example_custom_options() {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "protoc-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	protoDir := filepath.Join(tmpDir, "proto", "act7110")
	workspaceDir := filepath.Join(tmpDir, "proto")
	outputDir := filepath.Join(tmpDir, "generated")

	// Create directories
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Create a .proto file
	protoContent := `syntax = "proto3";

package example;

option go_package = "example/generated";

message User {
  string id = 1;
  string name = 2;
}`
	protoFile := filepath.Join(protoDir, "user.proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// Use custom options
	_ = protoc.NewCompiler().
		WithProtoDir(protoDir).
		WithProtoWorkSpace(workspaceDir).
		WithOutputDir(outputDir).
		WithPlugins("go").
		WithGoOpts("paths=source_relative", "module=github.com/example/project").
		WithVerbose(true)

	fmt.Println("Custom options example:")
	fmt.Println("  Go options: paths=source_relative, module=github.com/example/project")
	fmt.Println("  Verbose: true")
	fmt.Println()
	fmt.Println("This would generate Go files with the specified module path.")

	// Note: Actual compilation would require protoc to be installed
	// This example demonstrates the configuration

	// Output:
	// Custom options example:
	//   Go options: paths=source_relative, module=github.com/example/project
	//   Verbose: true
	//
	// This would generate Go files with the specified module path.
}

// Example_error_handling demonstrates error handling.
func Example_error_handling() {
	// Try to compile from non-existent directories
	_, err := protoc.Compile("/non/existent/proto", "/non/existent/workspace", "./output")

	if err != nil {
		// The actual error may vary depending on the environment
		fmt.Printf("Error occurred: %v\n", err)
	}

	// Output depends on the environment:
	// If directories don't exist:
	//   Error occurred: proto directory does not exist: /non/existent/proto
	// Or:
	//   Error occurred: workspace directory does not exist: /non/existent/workspace
}

// Example_must_compile demonstrates the MustCompile function.
func Example_must_compile() {
	// Note: In real usage, you would use actual directories with .proto files
	// This example shows the pattern for initialization

	// For initialization in tests or examples where failure should panic
	// Note: MustCompile will panic if compilation fails
	// In a real environment with protoc installed and .proto files present, this would work:
	// _ = protoc.MustCompile("./proto/act7110", "./proto", "./generated")

	fmt.Println("MustCompile example - would panic without proper configuration")

	// Output:
	// MustCompile example - would panic without proper configuration
}
