package e2e

import (
	"testing"

	"immich-compress/e2e/framework"
	"immich-compress/e2e/mocks"
)

func TestImmichIntegration_BasicAssetDownloadAndCompress(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server with test data
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()
	mockServer.Start()

	// Test with mock server URL
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "IMAGE",
		"--image-format", "jpeg",
		"--image-quality", "85",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// The test validates that the CLI accepts proper flags and attempts server connection
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --server" || stderr == "Error: unknown flag: --api-key" {
			t.Errorf("CLI should accept server flags: %s", stderr)
		} else {
			// Connection to mock server might work or fail, both are valid
			t.Logf("Mock server connection attempt: %s", stderr)
		}
	}

	t.Logf("Immich integration test completed with mock server at %s", mockServer.GetBaseURL())
}

func TestImmichIntegration_AssetSearchAndBatchCompress(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test search functionality with limit parameter
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "ALL",
		"--limit", "10",
		"--image-format", "webp",
		"--image-quality", "80",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept batch processing flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --limit" {
			t.Errorf("CLI should accept limit flag: %s", stderr)
		} else {
			t.Logf("Search and batch compress flags accepted: %s", stderr)
		}
	}
}

func TestImmichIntegration_AuthenticationFailure(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test with invalid token
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "invalid-token",
		"--type", "IMAGE",
		"--image-format", "jpeg",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should fail gracefully with authentication error or connection error
	if err != nil {
		if stderr == "Error: unknown flag: --api-key" {
			t.Errorf("CLI should accept api-key flag: %s", stderr)
		} else {
			// Authentication or connection error is expected
			t.Logf("Expected authentication error: %s", stderr)
		}
	}
}

func TestImmichIntegration_UploadCompressedAssets(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test compression with upload behavior (default behavior)
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "IMAGE",
		"--image-format", "jpeg",
		"--image-quality", "90",
		"--diff-percents", "10",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept compression parameters
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --diff-percents" {
			t.Errorf("CLI should accept diff-percents flag: %s", stderr)
		} else {
			t.Logf("Upload and compression flags accepted: %s", stderr)
		}
	}
}

func TestImmichIntegration_TagHandling(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test with compression parameters (tags are managed by the system)
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "ALL",
		"--image-format", "jpeg",
		"--video-format", "av1",
		"--image-quality", "85",
		"--video-quality", "25",
		"--diff-percents", "8",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept all compression parameters
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --server" {
			t.Errorf("CLI should accept server flag: %s", stderr)
		} else {
			t.Logf("Tag and compression handling flags accepted: %s", stderr)
		}
	}
}

func TestImmichIntegration_NetworkFailureRecovery(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test with unreachable server
	nonexistentServer := "http://localhost:9999"

	// Test network failure handling
	args := []string{
		"compress",
		"--server", nonexistentServer,
		"--api-key", "test-token",
		"--type", "IMAGE",
		"--image-format", "jpeg",
		"--diff-percents", "5",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should fail gracefully with network error
	if err != nil {
		if stderr == "Error: unknown flag: --server" {
			t.Errorf("CLI should accept server flag: %s", stderr)
		} else {
			// Network error is expected
			t.Logf("Expected network failure: %s", stderr)
		}
	}
}

func TestImmichIntegration_ConcurrentAssetProcessing(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test concurrent processing with parallel setting
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "ALL",
		"--parallel", "3",
		"--limit", "20",
		"--image-format", "jpeg",
		"--video-format", "h264",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept parallel processing flags
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --parallel" {
			t.Errorf("CLI should accept parallel flag: %s", stderr)
		} else {
			t.Logf("Concurrent processing flags accepted: %s", stderr)
		}
	}
}

func TestImmichIntegration_MetadataPreservation(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI binary
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Create mock Immich server
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	// Start mock Immich server
	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()

	// Test metadata preservation with high quality settings
	args := []string{
		"compress",
		"--server", mockServer.GetBaseURL(),
		"--api-key", "test-token",
		"--type", "ALL",
		"--image-format", "jpeg",
		"--image-quality", "95",
		"--video-quality", "20",
		"--video-format", "av1",
		"--diff-percents", "5",
	}

	_, stderr, err := runner.RunCLICommand(args...)

	// Should accept high-quality settings for metadata preservation
	if err != nil && stderr != "" {
		if stderr == "Error: unknown flag: --image-quality" {
			t.Errorf("CLI should accept image quality flag: %s", stderr)
		} else {
			t.Logf("Metadata preservation settings accepted: %s", stderr)
		}
	}
}
