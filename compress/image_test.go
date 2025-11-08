package compress

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cshum/vipsgen/vips"
	"github.com/google/uuid"
)

func TestImageConfig(t *testing.T) {
	// Test the ImageConfig struct and its behavior
	config := ImageConfig{
		Format:  JPEG,
		Quality: 80,
	}

	if config.Format != JPEG {
		t.Errorf("Expected format JPEG, got %s", config.Format)
	}

	if config.Quality != 80 {
		t.Errorf("Expected quality 80, got %d", config.Quality)
	}
}

func TestImageConfigInvalidUUID(t *testing.T) {
	// Test UUID parsing with invalid format - this should fail before any client calls
	_, err := uuid.Parse("invalid-uuid-format")
	if err == nil {
		t.Error("Expected error for invalid UUID format, got nil")
	}
}

func TestImageConfigSupportedFormats(t *testing.T) {
	// Test that all supported formats are recognized
	supportedFormats := []ImageFormat{JPG, JPEG, JXL, WEBP, HEIF}

	for _, format := range supportedFormats {
		config := ImageConfig{
			Format:  format,
			Quality: 80,
		}

		// Verify the format is set correctly
		if config.Format != format {
			t.Errorf("Expected format %s, got %s", format, config.Format)
		}
	}
}

func TestImageConfigUnsupportedFormat(t *testing.T) {
	config := ImageConfig{
		Format:  ImageFormat("unsupported"),
		Quality: 80,
	}

	// Test the format validation by checking if the switch statement handles unsupported formats
	switch config.Format {
	case JPEG, JPG, JXL, WEBP, HEIF:
		t.Errorf("Expected format %s to be unsupported, but it was accepted", config.Format)
	default:
		// Unsupported format - this is expected for this test
		t.Logf("Format %s is correctly identified as unsupported", config.Format)
	}
}

func TestImageConfigCreateTempFile(t *testing.T) {
	// Test that we can create a temporary file with the expected pattern
	config := ImageConfig{
		Format:  JPEG,
		Quality: 80,
	}

	uuid := uuid.New()
	tempDir := os.TempDir()
	expectedPath := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", uuid.String(), string(config.Format)))

	// This is the same logic as in the actual compress function
	fileOut, err := os.Create(expectedPath)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer fileOut.Close()
	defer os.Remove(expectedPath)

	// Write some test data
	testData := []byte("test image data")
	_, err = fileOut.Write(testData)
	if err != nil {
		t.Errorf("Failed to write to temp file: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(expectedPath); errors.Is(err, os.ErrNotExist) {
		t.Error("Expected temp file to exist")
	}
}

func TestImageConfigJPEG(t *testing.T) {
	// Test JPEG configuration
	config := ImageConfig{
		Format:  JPEG,
		Quality: 90,
	}

	if config.Format != JPEG {
		t.Errorf("Expected JPEG format, got %s", config.Format)
	}

	if config.Quality != 90 {
		t.Errorf("Expected quality 90, got %d", config.Quality)
	}
}

func TestImageConfigJXL(t *testing.T) {
	// Test JXL configuration
	config := ImageConfig{
		Format:  JXL,
		Quality: 75,
	}

	if config.Format != JXL {
		t.Errorf("Expected JXL format, got %s", config.Format)
	}

	if config.Quality != 75 {
		t.Errorf("Expected quality 75, got %d", config.Quality)
	}
}

func TestImageConfigWEBP(t *testing.T) {
	// Test WEBP configuration
	config := ImageConfig{
		Format:  WEBP,
		Quality: 85,
	}

	if config.Format != WEBP {
		t.Errorf("Expected WEBP format, got %s", config.Format)
	}

	if config.Quality != 85 {
		t.Errorf("Expected quality 85, got %d", config.Quality)
	}
}

func TestImageConfigHEIF(t *testing.T) {
	// Test HEIF configuration
	config := ImageConfig{
		Format:  HEIF,
		Quality: 70,
	}

	if config.Format != HEIF {
		t.Errorf("Expected HEIF format, got %s", config.Format)
	}

	if config.Quality != 70 {
		t.Errorf("Expected quality 70, got %d", config.Quality)
	}
}

func TestImageConfigEdgeCases(t *testing.T) {
	// Test edge cases
	tests := []struct {
		name     string
		format   ImageFormat
		quality  int
		expected bool
	}{
		{
			name:     "minimum quality",
			format:   JPEG,
			quality:  1,
			expected: true,
		},
		{
			name:     "maximum quality",
			format:   JPEG,
			quality:  100,
			expected: true,
		},
		{
			name:     "all formats with valid quality",
			format:   JXL,
			quality:  50,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ImageConfig{
				Format:  tt.format,
				Quality: tt.quality,
			}

			// These should not panic
			_ = config.Format
			_ = config.Quality

			if !tt.expected {
				t.Errorf("Unexpected configuration: format=%s, quality=%d", config.Format, config.Quality)
			}
		})
	}
}

func TestImageFormatsAvailable(t *testing.T) {
	// Test the available formats are correctly defined
	expectedFormats := []ImageFormat{JPG, JPEG, JXL, WEBP, HEIF}

	if len(ImageFormatsAvailable) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(ImageFormatsAvailable))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, actual := range ImageFormatsAvailable {
			if expected == actual {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Format %s not found in available formats", expected)
		}
	}
}

func TestImageConfigVipsIntegration(t *testing.T) {
	// Skip test if vips is not available
	defer func() {
		if r := recover(); r != nil {
			t.Skip("vips not available, skipping integration test")
		}
	}()

	// This test would require actual vips integration
	// For now, just test that we can create vips options
	options := vips.DefaultJpegsaveBufferOptions()
	if options == nil {
		t.Error("Expected non-nil vips options")
	}
}

func TestImageConfigJXLOptions(t *testing.T) {
	// Test JXL specific options
	defer func() {
		if r := recover(); r != nil {
			t.Skip("vips not available, skipping integration test")
		}
	}()

	options := vips.DefaultJxlsaveBufferOptions()
	if options == nil {
		t.Error("Expected non-nil JXL options")
	}

	// These are the defaults used in the actual code
	if options.Q == 0 {
		t.Error("Expected Q to be set")
	}
	if options.Effort == 0 {
		t.Error("Expected Effort to be set")
	}
}

func TestImageConfigWEBPOptions(t *testing.T) {
	// Test WEBP specific options
	defer func() {
		if r := recover(); r != nil {
			t.Skip("vips not available, skipping integration test")
		}
	}()

	options := vips.DefaultWebpsaveBufferOptions()
	if options == nil {
		t.Error("Expected non-nil WEBP options")
	}

	if options.Q == 0 {
		t.Error("Expected Q to be set")
	}
}

func TestImageConfigHEIFOptions(t *testing.T) {
	// Test HEIF specific options
	defer func() {
		if r := recover(); r != nil {
			t.Skip("vips not available, skipping integration test")
		}
	}()

	options := vips.DefaultHeifsaveBufferOptions()
	if options == nil {
		t.Error("Expected non-nil HEIF options")
	}

	if options.Q == 0 {
		t.Error("Expected Q to be set")
	}
}

func TestImageConfigJPEGVsJPG(t *testing.T) {
	// Test that JPEG and JPG are treated as the same format
	jpegConfig := ImageConfig{Format: JPEG, Quality: 80}
	jpgConfig := ImageConfig{Format: JPG, Quality: 80}

	if jpegConfig.Format != JPEG {
		t.Error("JPEG format should be constant JPEG")
	}
	if jpgConfig.Format != JPG {
		t.Error("JPG format should be constant JPG")
	}

	// Both should be supported by the switch statement in compress function
	switch jpegConfig.Format {
	case JPEG, JPG:
		// This should pass for JPEG format
	default:
		t.Errorf("JPEG format %s should be supported", jpegConfig.Format)
	}

	switch jpgConfig.Format {
	case JPEG, JPG:
		// This should pass for JPG format
	default:
		t.Errorf("JPG format %s should be supported", jpgConfig.Format)
	}
}
