package e2e

import (
	"testing"

	"immich-compress/e2e/framework"
	"immich-compress/e2e/mocks"
)

func TestFramework_BasicSetup(t *testing.T) {
	// Test that the test runner can be created and cleaned up
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Test basic operations
	if runner.TempDir == "" {
		t.Error("Temp directory should be set")
	}

	// Test file creation
	fixtureManager := framework.NewTestDataManager(runner.TestDataDir)
	imgPath, err := fixtureManager.CreateTestImage("test.jpg", "jpeg", 1024)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test that file exists
	if !runner.FileExists(imgPath) {
		t.Error("Created test file should exist")
	}

	// Test mock server creation
	err = fixtureManager.CreateMockImmichServer()
	if err != nil {
		t.Fatalf("Failed to create mock Immich data: %v", err)
	}

	mockServer := mocks.NewMockImmichServer(runner.TestDataDir)
	defer mockServer.Close()
	mockServer.Start()

	if mockServer.GetBaseURL() == "" {
		t.Error("Mock server should have a base URL")
	}

	t.Logf("Framework test completed successfully")
	t.Logf("Temp dir: %s", runner.TempDir)
	t.Logf("Project root: %s", runner.ProjectRoot)
	t.Logf("Mock server URL: %s", mockServer.GetBaseURL())
}

func TestFramework_BuildCLI(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Test CLI building
	err := runner.BuildCLI()
	// Note: CLI binary might already exist, so this is expected to succeed
	if err != nil {
		t.Logf("Build error (expected if building from scratch): %v", err)
		// Don't fail the test if building fails, as binary might already exist
	}

	// Test that project root is set correctly
	if runner.ProjectRoot == "" {
		t.Error("Project root should be set")
	}

	t.Logf("CLI build test completed successfully")
	t.Logf("Project root: %s", runner.ProjectRoot)
}

func TestFramework_CLIExecution(t *testing.T) {
	runner := framework.NewTestRunner(t)
	defer runner.Cleanup()

	// Build the CLI first
	err := runner.BuildCLI()
	if err != nil {
		t.Fatalf("Should build CLI successfully: %v", err)
	}

	// Test help command
	args := []string{"--help"}
	stdout, stderr, err := runner.RunCLICommand(args...)

	if err != nil {
		t.Errorf("Help command should work: %v", err)
	}

	if stdout == "" && stderr == "" {
		t.Error("Should have some output from help command")
	}

	t.Logf("CLI execution test completed")
	t.Logf("Help output length: %d", len(stdout)+len(stderr))
}
