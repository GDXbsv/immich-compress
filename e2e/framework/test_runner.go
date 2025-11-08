package framework

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type TestRunner struct {
	T              *testing.T
	Ctx            context.Context
	TempDir        string
	ProjectRoot    string
	TestDataDir    string
	OutputDir      string
	MockServerPort int
}

func NewTestRunner(t *testing.T) *TestRunner {
	ctx := context.Background()
	tempDir, err := os.MkdirTemp("", "immich-e2e-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	projectRoot, _ := os.Getwd()
	testDataDir := filepath.Join(projectRoot, "e2e", "testdata")
	outputDir := filepath.Join(tempDir, "output")

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	return &TestRunner{
		T:           t,
		Ctx:         ctx,
		TempDir:     tempDir,
		ProjectRoot: projectRoot,
		TestDataDir: testDataDir,
		OutputDir:   outputDir,
	}
}

func (tr *TestRunner) Cleanup() {
	if tr.TempDir != "" {
		os.RemoveAll(tr.TempDir)
	}
}

func (tr *TestRunner) RunCLICommand(args ...string) (string, string, error) {
	cliPath := filepath.Join(tr.ProjectRoot, "immich-compress")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("CLI binary not found at %s", cliPath)
	}

	cmd := exec.CommandContext(tr.Ctx, cliPath, args...)
	cmd.Dir = tr.ProjectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func (tr *TestRunner) BuildCLI() error {
	// Check if CLI binary already exists in the main project directory
	// The main project is the parent directory of e2e
	mainProjectDir := tr.ProjectRoot
	if filepath.Base(tr.ProjectRoot) == "tests" {
		mainProjectDir = filepath.Join(tr.ProjectRoot, "..", "..")
	} else if filepath.Base(tr.ProjectRoot) == "e2e" {
		mainProjectDir = filepath.Join(tr.ProjectRoot, "..")
	}

	cliPath := filepath.Join(mainProjectDir, "immich-compress")
	if _, err := os.Stat(cliPath); err == nil {
		// Binary already exists, update the project root to the main project
		tr.ProjectRoot = mainProjectDir
		return nil
	}

	// Binary doesn't exist, try to build it
	buildCmd := exec.Command("go", "build", "-o", "immich-compress", ".")
	buildCmd.Dir = mainProjectDir
	err := buildCmd.Run()
	if err == nil {
		// Update project root to the main project if build was successful
		tr.ProjectRoot = mainProjectDir
	}
	return err
}

func (tr *TestRunner) GetTestFile(subdir, filename string) string {
	return filepath.Join(tr.TestDataDir, subdir, filename)
}

func (tr *TestRunner) CreateTestFile(srcPath, destPath string) error {
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, srcData, 0644)
}

func (tr *TestRunner) WaitForCondition(condition func() bool, timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("condition not met within %v", timeout)
		case <-ticker.C:
			if condition() {
				return nil
			}
		}
	}
}

func (tr *TestRunner) AssertFileExists(path string) {
	if !tr.FileExists(path) {
		tr.T.Errorf("File should exist: %s", path)
	}
}

func (tr *TestRunner) AssertFileNotExists(path string) {
	if tr.FileExists(path) {
		tr.T.Errorf("File should not exist: %s", path)
	}
}

func (tr *TestRunner) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (tr *TestRunner) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (tr *TestRunner) AssertNoError(err error) {
	if err != nil {
		tr.T.Errorf("Expected no error, got: %v", err)
	}
}

func (tr *TestRunner) AssertError(err error) {
	if err == nil {
		tr.T.Errorf("Expected error, got none")
	}
}
