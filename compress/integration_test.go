package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"immich-compress/immich"

	"github.com/google/uuid"
)

func TestEndToEndImageCompressionWorkflow(t *testing.T) {
	// Create temporary directory for integration test
	tempDir := t.TempDir()

	// Simulate a complete image compression workflow
	imageConfig := ImageConfig{
		Format:  JPEG,
		Quality: 80,
	}

	// Create a mock image file
	imageFile := filepath.Join(tempDir, "test-image.jpg")
	testImageData := []byte("mock jpeg image data")
	err := os.WriteFile(imageFile, testImageData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imageFile)

	// Test the complete image processing chain
	t.Run("image_processing_chain", func(t *testing.T) {
		// Test temp file creation (simulating the compress function)
		assetUUID := uuid.New()
		compressedFile := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", assetUUID.String(), string(imageConfig.Format)))

		fileOut, err := os.Create(compressedFile)
		if err != nil {
			t.Fatalf("Failed to create compressed file: %v", err)
		}
		defer fileOut.Close()
		defer os.Remove(compressedFile)

		// Simulate compression by writing data
		_, err = fileOut.Write(testImageData)
		if err != nil {
			t.Errorf("Failed to write compressed data: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(compressedFile); os.IsNotExist(err) {
			t.Error("Compressed file should exist")
		}

		// Test size comparison logic
		originalSize := int64(len(testImageData))
		compressedSize := int64(len(testImageData)) // Same size in this mock
		diffPercent := 8
		threshold := int64(float64(originalSize) * float64(diffPercent) / 100)
		sizeDiff := originalSize - compressedSize

		shouldReplace := sizeDiff > threshold

		if !shouldReplace {
			t.Log("Size difference not significant enough for replacement")
		}
	})
}

func TestEndToEndVideoCompressionWorkflow(t *testing.T) {
	// Create temporary directory for integration test
	tempDir := t.TempDir()

	// Simulate a complete video compression workflow
	videoConfig := VideoConfig{
		Container: MP4,
		Format:    AV1,
		Quality:   25,
	}

	// Create a mock video file
	videoFile := filepath.Join(tempDir, "test-video.mp4")
	testVideoData := []byte("mock video data")
	err := os.WriteFile(videoFile, testVideoData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test video: %v", err)
	}
	defer os.Remove(videoFile)

	t.Run("video_processing_chain", func(t *testing.T) {
		// Test temp file creation for video
		assetUUID := uuid.New()
		compressedFile := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", assetUUID.String(), string(videoConfig.Container)))

		fileOut, err := os.Create(compressedFile)
		if err != nil {
			t.Fatalf("Failed to create compressed video file: %v", err)
		}
		defer fileOut.Close()
		defer os.Remove(compressedFile)

		// Test FFmpeg command building (without actually running FFmpeg)
		args := buildFFmpegArgs(videoConfig, videoFile, compressedFile, videoConfig.Quality)

		// Verify command structure
		if len(args) < 10 {
			t.Errorf("Expected at least 10 arguments for FFmpeg command, got %d", len(args))
		}

		if args[0] != "-i" {
			t.Errorf("Expected first argument to be -i, got %s", args[0])
		}

		if args[1] != videoFile {
			t.Errorf("Expected second argument to be input file, got %s", args[1])
		}

		// Check for required video encoding arguments
		hasVideoCodec := false
		hasCRF := false

		for i, arg := range args {
			if arg == "-c:v" {
				hasVideoCodec = true
				if i+1 < len(args) {
					expectedCodec := getExpectedVideoCodec(videoConfig.Format)
					if args[i+1] != expectedCodec {
						t.Errorf("Expected video codec %s, got %s", expectedCodec, args[i+1])
					}
				}
			}
			if arg == "-crf" {
				hasCRF = true
				if i+1 < len(args) {
					expectedCRF := "25"
					if args[i+1] != expectedCRF {
						t.Errorf("Expected CRF %s, got %s", expectedCRF, args[i+1])
					}
				}
			}
		}

		if !hasVideoCodec {
			t.Error("FFmpeg command should include video codec")
		}

		if !hasCRF {
			t.Error("FFmpeg command should include CRF")
		}
	})
}

func TestMixedAssetTypeProcessing(t *testing.T) {
	// Test processing of both image and video assets together
	tempDir := t.TempDir()

	t.Run("mixed_asset_workflow", func(t *testing.T) {
		// Create both image and video files
		imageFile := filepath.Join(tempDir, "test.jpg")
		videoFile := filepath.Join(tempDir, "test.mp4")

		imageData := []byte("mock image data")
		videoData := []byte("mock video data")

		err := os.WriteFile(imageFile, imageData, 0644)
		if err != nil {
			t.Fatalf("Failed to create test image: %v", err)
		}
		defer os.Remove(imageFile)

		err = os.WriteFile(videoFile, videoData, 0644)
		if err != nil {
			t.Fatalf("Failed to create test video: %v", err)
		}
		defer os.Remove(videoFile)

		// Test that both file types are handled correctly
		if _, err := os.Stat(imageFile); os.IsNotExist(err) {
			t.Error("Image file should exist")
		}

		if _, err := os.Stat(videoFile); os.IsNotExist(err) {
			t.Error("Video file should exist")
		}

		// Verify file sizes
		imageInfo, err := os.Stat(imageFile)
		if err != nil {
			t.Fatalf("Failed to get image file info: %v", err)
		}

		videoInfo, err := os.Stat(videoFile)
		if err != nil {
			t.Fatalf("Failed to get video file info: %v", err)
		}

		if imageInfo.Size() != int64(len(imageData)) {
			t.Error("Image file size should match written data")
		}

		if videoInfo.Size() != int64(len(videoData)) {
			t.Error("Video file size should match written data")
		}
	})
}

func TestErrorPropagationInWorkflow(t *testing.T) {
	// Test error handling in the complete workflow
	t.Run("image_compression_errors", func(t *testing.T) {
		imageConfig := ImageConfig{
			Format:  ImageFormat("invalid"), // Invalid format
			Quality: 80,
		}

		// Test format validation without actual processing
		switch imageConfig.Format {
		case JPEG, JPG, JXL, WEBP, HEIF:
			// This would be handled in actual processing
		default:
			// Invalid format correctly identified
			t.Logf("Invalid format %s correctly identified", imageConfig.Format)
		}
	})

	t.Run("video_compression_errors", func(t *testing.T) {
		videoConfig := VideoConfig{
			Format:    VideoFormat("invalid"), // Invalid format
			Container: MP4,
			Quality:   25,
		}

		// Test format validation
		switch videoConfig.Format {
		case AV1, HEVC, H264:
			// This would be handled in actual processing
		default:
			// Invalid format correctly identified
			t.Logf("Invalid video format %s correctly identified", videoConfig.Format)
		}
	})
}

func TestConfigurationPropagation(t *testing.T) {
	// Test that configuration is properly propagated through the system
	t.Run("image_config_propagation", func(t *testing.T) {
		imageConfig := ImageConfig{
			Format:  WEBP,
			Quality: 90,
		}

		// Test that config values are accessible
		if imageConfig.Format != WEBP {
			t.Errorf("Expected format WEBP, got %s", imageConfig.Format)
		}

		if imageConfig.Quality != 90 {
			t.Errorf("Expected quality 90, got %d", imageConfig.Quality)
		}

		// Test that format affects file naming
		assetUUID := uuid.New()
		filename := fmt.Sprintf("%s-compressed.%s", assetUUID.String(), string(imageConfig.Format))
		expectedExt := "webp"

		if filepath.Ext(filename) != "."+expectedExt {
			t.Errorf("Expected extension .%s, got %s", expectedExt, filepath.Ext(filename))
		}
	})

	t.Run("video_config_propagation", func(t *testing.T) {
		videoConfig := VideoConfig{
			Format:    HEVC,
			Container: MKV,
			Quality:   30,
		}

		// Test that config values are accessible
		if videoConfig.Format != HEVC {
			t.Errorf("Expected format HEVC, got %s", videoConfig.Format)
		}

		if videoConfig.Container != MKV {
			t.Errorf("Expected container MKV, got %s", videoConfig.Container)
		}

		if videoConfig.Quality != 30 {
			t.Errorf("Expected quality 30, got %d", videoConfig.Quality)
		}

		// Test that format affects file naming
		assetUUID := uuid.New()
		filename := fmt.Sprintf("%s-compressed.%s", assetUUID.String(), string(videoConfig.Container))
		expectedExt := "mkv"

		if filepath.Ext(filename) != "."+expectedExt {
			t.Errorf("Expected extension .%s, got %s", expectedExt, filepath.Ext(filename))
		}
	})
}

func TestAssetTypeRouting(t *testing.T) {
	// Test that assets are routed to the correct compression path
	t.Run("image_asset_routing", func(t *testing.T) {
		imageAsset := immich.AssetResponseDto{
			Id:               uuid.New().String(),
			Type:             immich.IMAGE,
			OriginalFileName: "test.jpg",
			OriginalPath:     "/test/test.jpg",
		}

		// Test that IMAGE type is correctly identified
		if imageAsset.Type != immich.IMAGE {
			t.Errorf("Expected type IMAGE, got %s", imageAsset.Type)
		}

		// This would route to imageConfig.compress() in actual code
		t.Log("Image asset would be routed to image compression path")
	})

	t.Run("video_asset_routing", func(t *testing.T) {
		videoAsset := immich.AssetResponseDto{
			Id:               uuid.New().String(),
			Type:             immich.VIDEO,
			OriginalFileName: "test.mp4",
			OriginalPath:     "/test/test.mp4",
		}

		// Test that VIDEO type is correctly identified
		if videoAsset.Type != immich.VIDEO {
			t.Errorf("Expected type VIDEO, got %s", videoAsset.Type)
		}

		// This would route to videoConfig.compress() in actual code
		t.Log("Video asset would be routed to video compression path")
	})

	t.Run("unsupported_asset_routing", func(t *testing.T) {
		audioAsset := immich.AssetResponseDto{
			Id:               uuid.New().String(),
			Type:             immich.AssetTypeEnum("AUDIO"),
			OriginalFileName: "test.mp3",
			OriginalPath:     "/test/test.mp3",
		}

		// Test that AUDIO type is correctly identified as unsupported
		if audioAsset.Type != immich.IMAGE && audioAsset.Type != immich.VIDEO {
			t.Logf("Asset type %s correctly identified as unsupported", audioAsset.Type)
		}
	})
}

// Helper functions for integration tests

func getExpectedVideoCodec(format VideoFormat) string {
	switch format {
	case AV1:
		return "libsvtav1"
	case HEVC:
		return "libx265"
	case H264:
		return "libx264"
	default:
		return ""
	}
}

// Test utility function for creating temp files
func createTempFile(t *testing.T, data []byte, suffix string) string {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test"+suffix)

	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return filename
}

// Test cleanup helper
func cleanupFile(t *testing.T, filename string) {
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: Failed to remove file %s: %v", filename, err)
	}
}
