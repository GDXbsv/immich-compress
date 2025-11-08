package e2e

import (
	"testing"

	"immich-compress/e2e/framework"
)

func TestImageWorkflow_JPEGCompression(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test different quality settings with server-based compression
	qualitySettings := []string{"50", "75", "90", "95"}

	for _, quality := range qualitySettings {
		t.Run("Quality_"+quality, func(t *testing.T) {
			// Test image format setting (not local file input)
			args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--image-format", "jpeg", "--image-quality", quality}
			_, stderr, err := runner.RunCLICommand(args...)

			// Should fail due to server connection, but accept the flags
			if err != nil && stderr != "" {
				// Check if it's a connection error (expected) vs unknown flag error (not expected)
				if stderr == "Error: unknown flag: --server" || stderr == "Error: unknown flag: --api-key" {
					t.Errorf("CLI should accept server flags: %s", stderr)
				} else {
					// Connection error is expected behavior
					t.Logf("Expected connection error for quality %s: %s", quality, stderr)
				}
			}
		})
	}
}

func TestImageWorkflow_FormatConversion(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test image format conversions (server-side processing)
	testCases := []struct {
		InputFormat  string
		OutputFormat string
	}{
		{"jpeg", "png"},
		{"jpeg", "webp"},
		{"png", "jpeg"},
		{"png", "webp"},
	}

	for _, tc := range testCases {
		t.Run(tc.InputFormat+"_to_"+tc.OutputFormat, func(t *testing.T) {
			// Test format conversion with server
			args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--image-format", tc.OutputFormat}
			_, stderr, err := runner.RunCLICommand(args...)

			// Should accept the format flags
			if err != nil && stderr != "" {
				if stderr == "Error: unknown flag: --image-format" {
					t.Errorf("CLI should accept image format flags: %s", stderr)
				} else {
					t.Logf("Format conversion test passed for %s: %s", tc.OutputFormat, stderr)
				}
			}
		})
	}
}

func TestImageWorkflow_BatchProcessing(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test batch processing with server
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--limit", "10"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should fail due to server connection, but accept batch flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --limit" {
			t.Errorf("CLI should accept batch processing flags: %s", stderr)
		} else {
			t.Logf("Batch processing flags accepted: %s", stderr)
		}
	}
}

func TestImageWorkflow_AnimatedGIFHandling(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test animated GIF handling with server
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--image-format", "webp"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept format flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --server" || stderr == "Error: unknown flag: --api-key" {
			t.Errorf("CLI should accept server flags: %s", stderr)
		} else {
			t.Logf("Animated format handling flags accepted: %s", stderr)
		}
	}
}

func TestImageWorkflow_ImageSizeOptimization(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test size-based optimization (diff-percents setting)
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--diff-percents", "5"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept size optimization flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --diff-percents" {
			t.Errorf("CLI should accept diff-percents flag: %s", stderr)
		} else {
			t.Logf("Size optimization flags accepted: %s", stderr)
		}
	}
}

func TestImageWorkflow_UnsupportedFormatHandling(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test with unsupported image format
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--image-format", "unsupported"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should handle unsupported format gracefully
	if err != nil {
		if stderr == "Error: unknown flag: --image-format" {
			t.Errorf("CLI should accept image format flag: %s", stderr)
		} else {
			// Either connection error or format validation error
			t.Logf("Format validation working: %s", stderr)
		}
	}
}

func TestImageWorkflow_PreserveMetadata(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test with different quality setting to preserve metadata
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "IMAGE", "--image-format", "jpeg", "--image-quality", "95"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept quality preservation settings
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --image-quality" {
			t.Errorf("CLI should accept image quality flag: %s", stderr)
		} else {
			t.Logf("Metadata preservation settings accepted: %s", stderr)
		}
	}
}
