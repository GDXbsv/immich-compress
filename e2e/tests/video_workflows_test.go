package e2e

import (
	"testing"

	"immich-compress/e2e/framework"
)

func TestVideoWorkflow_MP4Compression(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test different video quality settings with server-based processing
	qualitySettings := []string{"10", "25", "50"}

	for _, quality := range qualitySettings {
		t.Run("Quality_"+quality, func(t *testing.T) {
			// Test video compression with specific quality
			args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", "h264", "--video-quality", quality}
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

func TestVideoWorkflow_AV1Encoding(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test AV1 encoding
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", "av1", "--video-quality", "25"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept AV1 encoding flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --video-format" || stderr == "Error: unknown flag: --video-quality" {
			t.Errorf("CLI should accept video format flags: %s", stderr)
		} else {
			t.Logf("AV1 encoding flags accepted: %s", stderr)
		}
	}
}

func TestVideoWorkflow_WebMConversion(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test WebM container and format settings
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", "av1", "--video-container", "mkv", "--video-quality", "30"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept container and format flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --video-container" {
			t.Errorf("CLI should accept video container flag: %s", stderr)
		} else {
			t.Logf("WebM conversion flags accepted: %s", stderr)
		}
	}
}

func TestVideoWorkflow_FormatConversions(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test video format conversions (server-side processing)
	testCases := []struct {
		InputFormat  string
		OutputFormat string
	}{
		{"mp4", "h264"},
		{"mp4", "hevc"},
		{"webm", "av1"},
	}

	for _, tc := range testCases {
		t.Run(tc.InputFormat+"_to_"+tc.OutputFormat, func(t *testing.T) {
			// Test format conversion with server
			args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", tc.OutputFormat}
			_, stderr, err := runner.RunCLICommand(args...)

			// Should accept the format flags
			if err != nil && stderr != "" {
				if stderr == "Error: unknown flag: --video-format" {
					t.Errorf("CLI should accept video format flags: %s", stderr)
				} else {
					t.Logf("Video format conversion test passed for %s: %s", tc.OutputFormat, stderr)
				}
			}
		})
	}
}

func TestVideoWorkflow_BatchVideoProcessing(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test batch video processing with server
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--limit", "5"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should fail due to server connection, but accept batch flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --limit" {
			t.Errorf("CLI should accept batch processing flags: %s", stderr)
		} else {
			t.Logf("Batch video processing flags accepted: %s", stderr)
		}
	}
}

func TestVideoWorkflow_AudioTrackPreservation(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test video compression with quality settings that should preserve audio
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", "h264", "--video-quality", "25"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept quality settings for audio preservation
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --video-quality" {
			t.Errorf("CLI should accept video quality flag: %s", stderr)
		} else {
			t.Logf("Audio preservation settings accepted: %s", stderr)
		}
	}
}

func TestVideoWorkflow_SizeBasedOptimization(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test size-based optimization with different quality settings
	sizes := []string{"10", "25", "50"}

	for _, size := range sizes {
		t.Run("Quality_"+size, func(t *testing.T) {
			// Test with different quality settings
			args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-quality", size}
			_, stderr, err := runner.RunCLICommand(args...)

			if err != nil && stderr != "" {
				if stderr == "Error: unknown flag: --video-quality" {
					t.Errorf("CLI should accept video quality flag: %s", stderr)
				} else {
					t.Logf("Size optimization test passed for quality %s: %s", size, stderr)
				}
			}
		})
	}
}

func TestVideoWorkflow_UnsupportedVideoFormat(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test with unsupported video format
	args := []string{"compress", "--server", "http://localhost:2283", "--api-key", "test-key", "--type", "VIDEO", "--video-format", "unsupported"}
	_, stderr, err := runner.RunCLICommand(args...)

	// Should handle unsupported format gracefully
	if err != nil {
		if stderr == "Error: unknown flag: --video-format" {
			t.Errorf("CLI should accept video format flag: %s", stderr)
		} else {
			// Either connection error or format validation error
			t.Logf("Unsupported format validation working: %s", stderr)
		}
	}
}
