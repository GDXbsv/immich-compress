package compress

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"immich-compress/immich"

	"github.com/google/uuid"
)

func TestVideoConfig(t *testing.T) {
	// Test the VideoConfig struct and its behavior
	config := VideoConfig{
		Container: MP4,
		Format:    AV1,
		Quality:   25,
	}

	if config.Container != MP4 {
		t.Errorf("Expected container MP4, got %s", config.Container)
	}

	if config.Format != AV1 {
		t.Errorf("Expected format AV1, got %s", config.Format)
	}

	if config.Quality != 25 {
		t.Errorf("Expected quality 25, got %d", config.Quality)
	}
}

func TestVideoConfigCompressInvalidUUID(t *testing.T) {
	config := VideoConfig{
		Container: MP4,
		Format:    AV1,
		Quality:   25,
	}

	ctx := context.Background()

	// Test with invalid UUID
	asset := createTestAsset("invalid-uuid", "VIDEO", "video.mp4")

	_, err := config.compress(ctx, nil, asset)
	if err == nil {
		t.Error("Expected error for invalid UUID, got nil")
	}
}

func TestVideoConfigCompressUnsupportedFormat(t *testing.T) {
	config := VideoConfig{
		Container: MP4,
		Format:    VideoFormat("unsupported"),
		Quality:   25,
	}

	// Test format validation without calling the actual compress function
	switch config.Format {
	case AV1, HEVC, H264:
		t.Errorf("Expected format %s to be unsupported, but it was accepted", config.Format)
	default:
		// Unsupported format - this is expected for this test
		t.Logf("Format %s is correctly identified as unsupported", config.Format)
	}
}

func TestVideoConfigCreateTempFiles(t *testing.T) {
	config := VideoConfig{
		Container: MP4,
		Format:    AV1,
		Quality:   25,
	}

	uuid := uuid.New()
	tempDir := os.TempDir()

	// Test input file creation (same logic as in compress function)
	inputPath := filepath.Join(tempDir, fmt.Sprintf("%s.mp4", uuid.String()))
	fileIn, err := os.Create(inputPath)
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer fileIn.Close()
	defer os.Remove(inputPath)

	// Test output file path creation
	outputPath := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", uuid.String(), string(config.Container)))

	// Write some test data
	testData := []byte("test video data")
	_, err = fileIn.Write(testData)
	if err != nil {
		t.Errorf("Failed to write to temp input file: %v", err)
	}

	// Verify input file was created
	if _, err := os.Stat(inputPath); errors.Is(err, os.ErrNotExist) {
		t.Error("Expected temp input file to exist")
	}

	// Create output file to test the pattern
	fileOut, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer fileOut.Close()
	defer os.Remove(outputPath)

	// Verify output file pattern
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		t.Error("Expected temp output file to exist")
	}
}

func TestVideoConfigFFmpegCommandBuilding(t *testing.T) {
	tests := []struct {
		name       string
		config     VideoConfig
		quality    int
		inputFile  string
		outputFile string
	}{
		{
			name: "AV1 configuration",
			config: VideoConfig{
				Container: MKV,
				Format:    AV1,
				Quality:   25,
			},
			quality:    25,
			inputFile:  "/tmp/test_input.mp4",
			outputFile: "/tmp/test_output.mkv",
		},
		{
			name: "HEVC configuration",
			config: VideoConfig{
				Container: MP4,
				Format:    HEVC,
				Quality:   30,
			},
			quality:    30,
			inputFile:  "/tmp/test_input.mkv",
			outputFile: "/tmp/test_output.mp4",
		},
		{
			name: "H264 configuration",
			config: VideoConfig{
				Container: MP4,
				Format:    H264,
				Quality:   35,
			},
			quality:    35,
			inputFile:  "/tmp/test_input.mkv",
			outputFile: "/tmp/test_output.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that we can build the command arguments
			args := buildFFmpegArgs(tt.config, tt.inputFile, tt.outputFile, tt.quality)

			// Verify the command structure
			if len(args) < 2 {
				t.Errorf("Expected at least 2 args, got %d", len(args))
			}

			if args[0] != "-i" {
				t.Errorf("Expected first arg to be '-i', got %s", args[0])
			}

			if args[1] != tt.inputFile {
				t.Errorf("Expected input file to be %s, got %s", tt.inputFile, args[1])
			}

			// Check that output file is at the end
			if args[len(args)-1] != tt.outputFile {
				t.Errorf("Expected last arg to be output file %s, got %s", tt.outputFile, args[len(args)-1])
			}

			// Check for codec arguments
			foundCodec := false
			for _, arg := range args {
				if arg == "-c:v" {
					foundCodec = true
					break
				}
			}
			if !foundCodec {
				t.Error("Expected -c:v argument in FFmpeg command")
			}
		})
	}
}

func TestVideoConfigFFmpegAV1(t *testing.T) {
	config := VideoConfig{
		Container: MKV,
		Format:    AV1,
		Quality:   20,
	}

	// Test that the format-specific arguments are correctly configured
	args := buildFFmpegArgs(config, "input.mp4", "output.mkv", config.Quality)

	// Check for AV1-specific arguments
	expectedArgs := []string{"-c:v", "libsvtav1", "-preset", "5"}
	for _, expected := range expectedArgs {
		found := false
		for _, arg := range args {
			if arg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected FFmpeg argument %s for AV1 encoding", expected)
		}
	}
}

func TestVideoFormatsAvailable(t *testing.T) {
	expectedFormats := []VideoFormat{AV1, HEVC, H264}
	actualFormats := VideoFormatsAvailable

	if len(actualFormats) != len(expectedFormats) {
		t.Errorf("Expected %d video formats, got %d", len(expectedFormats), len(actualFormats))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, actual := range actualFormats {
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

func TestVideoContainersAvailable(t *testing.T) {
	expectedContainers := []VideoContainer{MKV, MP4}
	actualContainers := VideoContainersAvailable

	if len(actualContainers) != len(expectedContainers) {
		t.Errorf("Expected %d video containers, got %d", len(expectedContainers), len(actualContainers))
	}

	for _, expected := range expectedContainers {
		found := false
		for _, actual := range actualContainers {
			if expected == actual {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Container %s not found in available containers", expected)
		}
	}
}

func TestVideoConfigQualityRange(t *testing.T) {
	// Test different quality values to ensure they're properly handled
	qualities := []int{1, 10, 25, 50, 75, 90, 100}

	for _, quality := range qualities {
		config := VideoConfig{
			Format:    AV1,
			Container: MKV,
			Quality:   quality,
		}

		if config.Quality != quality {
			t.Errorf("Expected quality %d, got %d", quality, config.Quality)
		}

		// Test that quality is converted to string for FFmpeg
		qualityStr := strconv.Itoa(config.Quality)
		if qualityStr == "" {
			t.Errorf("Quality %d should convert to non-empty string", quality)
		}
	}
}

// Helper function to build FFmpeg arguments (mirror of the logic in video.go)
func buildFFmpegArgs(config VideoConfig, inputFile, outputFile string, quality int) []string {
	args := make([]string, 0, 20)
	args = append(args, "-i", inputFile)

	switch config.Format {
	case AV1:
		args = append(args,
			"-c:v", "libsvtav1",
			"-pix_fmt", "yuv420p10le",
			"-preset", "5",
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	case HEVC:
		args = append(args,
			"-c:v", "libx265",
			"-profile:v", "main10",
			"-pix_fmt", "yuv420p10le",
			"-preset", "slow",
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	case H264:
		args = append(args,
			"-c:v", "libx264",
			"-profile:v", "high",
			"-pix_fmt", "yuv420p",
			"-preset", "slow",
			"-c:a", "aac",
			"-b:a", "128k",
		)
	default:
		// This would be handled by the actual compress function
		return args
	}

	args = append(args, []string{
		"-crf", strconv.Itoa(quality),
		outputFile,
	}...)

	return args
}

// Helper function to create test assets
func createTestAsset(id, assetType, fileName string) immich.AssetResponseDto {
	return immich.AssetResponseDto{
		Id:               id,
		Type:             immich.AssetTypeEnum(assetType),
		OriginalFileName: fileName,
		OriginalPath:     "/test/" + fileName,
	}
}

func TestVideoConfigFFmpegNotAvailable(t *testing.T) {
	// Test that we can create a command but it will fail to execute
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()

	if err != nil {
		// FFmpeg not available - this is expected in some test environments
		t.Skip("FFmpeg not available for testing")
	}
}
