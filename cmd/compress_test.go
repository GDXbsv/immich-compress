package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// TestCompressCommandExists verifies the compress command is properly defined
func TestCompressCommandExists(t *testing.T) {
	if compressCmd == nil {
		t.Fatal("compressCmd should not be nil")
	}

	if compressCmd.Use != "compress" {
		t.Errorf("Expected command use 'compress', got %q", compressCmd.Use)
	}

	if compressCmd.Short == "" {
		t.Error("compress command should have a short description")
	}

	if compressCmd.RunE == nil {
		t.Error("compress command should have a RunE function")
	}
}

// TestCompressCommandHasRequiredFlags verifies the required flags are present
func TestCompressCommandHasRequiredFlags(t *testing.T) {
	if compressCmd == nil {
		t.Fatal("compressCmd should not be nil")
	}

	// Test server flag exists
	serverFlag := compressCmd.PersistentFlags().Lookup("server")
	if serverFlag == nil {
		t.Error("server flag should be defined")
	} else {
		if serverFlag.Shorthand != "s" {
			t.Errorf("Expected server flag shorthand 's', got %q", serverFlag.Shorthand)
		}
	}

	// Test api-key flag exists
	apiKeyFlag := compressCmd.PersistentFlags().Lookup("api-key")
	if apiKeyFlag == nil {
		t.Error("api-key flag should be defined")
	} else {
		if apiKeyFlag.Shorthand != "a" {
			t.Errorf("Expected api-key flag shorthand 'a', got %q", apiKeyFlag.Shorthand)
		}
	}
}

// TestCompressCommandRequiredFlagsValidation tests that required flags work with Cobra
func TestCompressCommandRequiredFlagsValidation(t *testing.T) {
	if compressCmd == nil {
		t.Fatal("compressCmd should not be nil")
	}

	// Test missing both flags
	t.Run("missing both required flags", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		// Copy flag definitions from compressCmd
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
		if err := cmd.MarkPersistentFlagRequired("server"); err != nil {
			panic("Failed to mark server flag as required: " + err.Error())
		}
		if err := cmd.MarkPersistentFlagRequired("api-key"); err != nil {
			panic("Failed to mark api-key flag as required: " + err.Error())
		}

		cmd.SetArgs([]string{})

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when both required flags are missing")
		}
	})

	// Test missing server flag
	t.Run("missing server flag", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
		if err := cmd.MarkPersistentFlagRequired("server"); err != nil {
			panic("Failed to mark server flag as required: " + err.Error())
		}
		if err := cmd.MarkPersistentFlagRequired("api-key"); err != nil {
			panic("Failed to mark api-key flag as required: " + err.Error())
		}

		cmd.SetArgs([]string{"--api-key", "test-key"})

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when server flag is missing")
		}
	})

	// Test missing api-key flag
	t.Run("missing api-key flag", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test command",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
		if err := cmd.MarkPersistentFlagRequired("server"); err != nil {
			t.Fatalf("Failed to mark server flag as required: %v", err)
		}
		if err := cmd.MarkPersistentFlagRequired("api-key"); err != nil {
			t.Fatalf("Failed to mark api-key flag as required: %v", err)
		}

		cmd.SetArgs([]string{"--server", "https://test.com"})

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when api-key flag is missing")
		}
	})
}

// TestCompressCommandHelp tests help functionality
func TestCompressCommandHelp(t *testing.T) {
	if compressCmd == nil {
		t.Fatal("compressCmd should not be nil")
	}

	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "compress help",
			args:     []string{"compress", "--help"},
			wantHelp: true,
		},
		{
			name:     "compress -h",
			args:     []string{"compress", "-h"},
			wantHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:   "test",
				Short: "Test command",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}

			// Setup flags like in compressCmd
			cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
			cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
			cmd.MarkPersistentFlagRequired("server")
			cmd.MarkPersistentFlagRequired("api-key")

			cmd.SetArgs(tt.args)

			var output bytes.Buffer
			cmd.SetOut(&output)
			cmd.SetErr(&output)

			err := cmd.Execute()

			// Help commands should either succeed or return a help-related error
			if err != nil && !strings.Contains(err.Error(), "help") && !strings.Contains(err.Error(), "Usage") {
				t.Errorf("Unexpected error: %v", err)
			}

			outputStr := output.String()
			if tt.wantHelp && !strings.Contains(outputStr, "Usage:") {
				t.Errorf("Expected help output to contain 'Usage:', got: %s", outputStr)
			}
		})
	}
}

// TestCompressCommandFlagParsing tests that flags are properly parsed
func TestCompressCommandFlagParsing(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
	cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")

	tests := []struct {
		name  string
		args  []string
		check func(*testing.T)
	}{
		{
			name: "HTTPS server URL",
			args: []string{"--server", "https://example.com", "--api-key", "test-key"},
			check: func(t *testing.T) {
				if flagsCompress.flagServer != "https://example.com" {
					t.Errorf("Expected server 'https://example.com', got %q", flagsCompress.flagServer)
				}
				if flagsCompress.flagAPIKey != "test-key" {
					t.Errorf("Expected api-key 'test-key', got %q", flagsCompress.flagAPIKey)
				}
			},
		},
		{
			name: "HTTP localhost",
			args: []string{"--server", "http://localhost:3001", "--api-key", "abc123"},
			check: func(t *testing.T) {
				if flagsCompress.flagServer != "http://localhost:3001" {
					t.Errorf("Expected server 'http://localhost:3001', got %q", flagsCompress.flagServer)
				}
				if flagsCompress.flagAPIKey != "abc123" {
					t.Errorf("Expected api-key 'abc123', got %q", flagsCompress.flagAPIKey)
				}
			},
		},
		{
			name: "shorthand flags",
			args: []string{"-s", "https://test.com", "-a", "def456"},
			check: func(t *testing.T) {
				if flagsCompress.flagServer != "https://test.com" {
					t.Errorf("Expected server 'https://test.com', got %q", flagsCompress.flagServer)
				}
				if flagsCompress.flagAPIKey != "def456" {
					t.Errorf("Expected api-key 'def456', got %q", flagsCompress.flagAPIKey)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test
			flagsCompress.flagServer = ""
			flagsCompress.flagAPIKey = ""

			cmd.SetArgs(tt.args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			err := cmd.Execute()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			tt.check(t)
		})
	}
}

// TestCompressCommandContextHandling tests context handling
func TestCompressCommandContextHandling(t *testing.T) {
	// This test verifies that the command properly integrates with Cobra's context handling
	timeout := 1 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			select {
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			default:
				// Simulate quick execution
				return nil
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
	cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")

	// Test with short timeout
	cmd.SetArgs([]string{"--server", "https://test.com", "--api-key", "test"})

	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	err := cmd.ExecuteContext(ctx)

	// Either success or timeout error is acceptable
	if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("Unexpected error type: %v", err)
	}
}

// TestCompressCommandIntegration tests integration with the actual compress command
func TestCompressCommandIntegration(t *testing.T) {
	if compressCmd == nil {
		t.Fatal("compressCmd should not be nil")
	}

	// Test that the actual command can be instantiated without errors
	t.Run("command instantiation", func(t *testing.T) {
		// The compressCmd should be properly initialized
		if compressCmd.Use != "compress" {
			t.Errorf("Expected use 'compress', got %q", compressCmd.Use)
		}

		if compressCmd.RunE == nil {
			t.Error("compressCmd should have a RunE function")
		}

		// Test that flags are accessible
		serverFlag := compressCmd.PersistentFlags().Lookup("server")
		if serverFlag == nil {
			t.Error("server flag should be accessible")
		}

		apiKeyFlag := compressCmd.PersistentFlags().Lookup("api-key")
		if apiKeyFlag == nil {
			t.Error("api-key flag should be accessible")
		}
	})

	// Test command execution structure (without actual network calls)
	t.Run("command execution structure", func(t *testing.T) {
		// Create a test command that mimics the actual compress command structure
		cmd := &cobra.Command{
			Use:   "compress",
			Short: "Compress existing fotos/videos",
			RunE: func(cmd *cobra.Command, args []string) error {
				// This mimics the structure of the actual compress command
				// We can't test the actual compress.Compressing function without network dependencies
				// but we can test the command structure
				return nil
			},
		}

		cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")
		cmd.MarkPersistentFlagRequired("server")
		cmd.MarkPersistentFlagRequired("api-key")

		// Test valid execution
		cmd.SetArgs([]string{"--server", "https://test.com", "--api-key", "test-key"})

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error for valid command, got: %v", err)
		}
	})
}

// TestCompressCommandParallelAndAfterFlags tests integration with root command flags
func TestCompressCommandParallelAndAfterFlags(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	// Test that parallel and after flags from root are accessible
	parallelFlag := rootCmd.PersistentFlags().Lookup("parallel")
	if parallelFlag == nil {
		t.Error("parallel flag should be defined on root command")
	}

	afterFlag := rootCmd.PersistentFlags().Lookup("after")
	if afterFlag == nil {
		t.Error("after flag should be defined on root command")
	}

	// Test that these flags work with the compress command
	t.Run("root flags integration", func(t *testing.T) {
		cmd := &cobra.Command{
			Use:   "compress",
			Short: "Compress existing fotos/videos",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		// Add root command flags
		cmd.PersistentFlags().IntVarP(&flagsRoot.flagParallel, "parallel", "p", 4, "parallel")
		cmd.PersistentFlags().TimeVarP(&flagsRoot.flagAfter, "after", "t", time.Now(), []string{"2006-01-02 15:04:05"}, "after")

		// Add compress-specific flags
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
		cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")

		cmd.SetArgs([]string{
			"--server", "https://test.com",
			"--api-key", "test-key",
			"--parallel", "8",
			"--after", "2024-01-01 12:00:00",
		})

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify flags were parsed
		if flagsCompress.flagServer != "https://test.com" {
			t.Errorf("Expected server to be parsed correctly")
		}
		if flagsCompress.flagAPIKey != "test-key" {
			t.Errorf("Expected api-key to be parsed correctly")
		}
	})
}

// BenchmarkCompressCommandParsing benchmarks command parsing
func BenchmarkCompressCommandParsing(b *testing.B) {
	cmd := &cobra.Command{
		Use:   "compress",
		Short: "Compress existing fotos/videos",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&flagsCompress.flagServer, "server", "s", "", "The immich server address")
	cmd.PersistentFlags().StringVarP(&flagsCompress.flagAPIKey, "api-key", "a", "", "The immich server API key")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := []string{"--server", "https://test.com", "--api-key", "test-key"}
		cmd.SetArgs(args)
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})

		cmd.Execute()
	}
}
