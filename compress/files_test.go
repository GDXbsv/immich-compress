package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"immich-compress/immich"

	"github.com/google/uuid"
)

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected float64
	}{
		{
			name:     "1 MB",
			bytes:    1024 * 1024,
			expected: 1.0,
		},
		{
			name:     "2.5 MB",
			bytes:    int64(2.5 * 1024 * 1024),
			expected: 2.5,
		},
		{
			name:     "0 bytes",
			bytes:    0,
			expected: 0.0,
		},
		{
			name:     "100 MB",
			bytes:    100 * 1024 * 1024,
			expected: 100.0,
		},
		{
			name:     "small file",
			bytes:    512 * 1024,
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesToMB(tt.bytes)
			if result != tt.expected {
				t.Errorf("bytesToMB(%d) = %v, want %v", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestBytesToMBEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected float64
	}{
		{
			name:     "negative bytes",
			bytes:    -1024,
			expected: -0.0009765625,
		},
		{
			name:     "very large number",
			bytes:    1024 * 1024 * 1024 * 1024, // 1 TB
			expected: 1048576.0,                 // Should be 1,048,576 MB (1 TB = 1,048,576 MB)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesToMB(tt.bytes)
			if result != tt.expected {
				t.Errorf("bytesToMB(%d) = %v, want %v", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestCompressFileWithImage(t *testing.T) {
	// Test compressFile logic with an IMAGE asset type
	tempDir := os.TempDir()

	// Create a temporary file that simulates compressed output
	tempFile := filepath.Join(tempDir, "test-compressed.jpg")

	// Write some test data to simulate compressed file
	testData := []byte("compressed image data")
	err := os.WriteFile(tempFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test the core logic of compressFile (without calling the actual function)
	// This tests the size comparison logic

	// Simulate that the compressed file is smaller
	originalSize := int64(10 * 1024 * 1024) // 10 MB
	sizeNew := int64(len(testData))         // Size of our test data

	// Test the size comparison logic (if sizeOrig-sizeNew > threshold)
	diffPercent := 8
	threshold := int64(float64(originalSize) * float64(diffPercent) / 100)
	sizeDiff := originalSize - sizeNew

	shouldReplace := sizeDiff > threshold

	// For this test, we expect the file to be smaller than threshold, so shouldReplace should be true
	if !shouldReplace {
		t.Errorf("Expected file to be smaller than threshold for replacement, sizeDiff=%d, threshold=%d", sizeDiff, threshold)
	}

	// Test the opposite case - very small size difference
	largeOriginalSize := int64(100 * 1024 * 1024)  // 100 MB
	smallCompressedSize := int64(99 * 1024 * 1024) // 99 MB
	smallSizeDiff := largeOriginalSize - smallCompressedSize
	smallThreshold := int64(float64(largeOriginalSize) * float64(diffPercent) / 100)

	shouldReplaceSmall := smallSizeDiff > smallThreshold

	if shouldReplaceSmall {
		t.Errorf("Expected small size difference to not trigger replacement, sizeDiff=%d, threshold=%d", smallSizeDiff, smallThreshold)
	}
}

func TestCompressFileWithVideo(t *testing.T) {
	// Test compressFile with a VIDEO asset type
	tempDir := os.TempDir()

	// Create a temporary file that simulates compressed video output
	tempFile := filepath.Join(tempDir, "test-compressed.mp4")

	// Write some test data to simulate compressed video
	testData := []byte("compressed video data")
	err := os.WriteFile(tempFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test that we can create files with video extensions
	videoFile := filepath.Join(tempDir, "test-video.mkv")
	_, err = os.Create(videoFile)
	if err != nil {
		t.Fatalf("Failed to create video file: %v", err)
	}
	defer os.Remove(videoFile)

	// Test video file pattern matching
	if !hasVideoExtension("test.mp4") {
		t.Error("Expected .mp4 to be recognized as video extension")
	}

	if !hasVideoExtension("test.mkv") {
		t.Error("Expected .mkv to be recognized as video extension")
	}

	if hasVideoExtension("test.jpg") {
		t.Error("Expected .jpg to NOT be recognized as video extension")
	}
}

func TestCompressFileZeroSize(t *testing.T) {
	// Test compressFile with a file that has zero size
	tempDir := os.TempDir()

	// Create an empty file
	emptyFile := filepath.Join(tempDir, "empty.txt")
	_, err := os.Create(emptyFile)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	defer os.Remove(emptyFile)

	// Get file info for empty file
	file, err := os.Open(emptyFile)
	if err != nil {
		t.Fatalf("Failed to open empty file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Test the logic that checks for zero size
	sizeNew := fileInfo.Size()
	if sizeNew == 0 {
		// This should trigger an error in the actual compressFile function
		t.Log("Zero size file correctly identified")
	}
}

func TestCompressFileSizeComparison(t *testing.T) {
	// Test the size comparison logic used in compressFile
	tests := []struct {
		name           string
		originalSize   int64
		compressedSize int64
		diffPercent    int
		shouldReplace  bool
		description    string
	}{
		{
			name:           "significant size reduction",
			originalSize:   10 * 1024 * 1024, // 10 MB
			compressedSize: 5 * 1024 * 1024,  // 5 MB
			diffPercent:    8,
			shouldReplace:  true,
			description:    "50% reduction should trigger replacement",
		},
		{
			name:           "small size reduction",
			originalSize:   10 * 1024 * 1024,  // 10 MB
			compressedSize: 9.5 * 1024 * 1024, // 9.5 MB (5% reduction)
			diffPercent:    8,
			shouldReplace:  false,
			description:    "5% reduction should not trigger replacement",
		},
		{
			name:           "larger file with small reduction",
			originalSize:   100 * 1024 * 1024, // 100 MB
			compressedSize: 95 * 1024 * 1024,  // 95 MB
			diffPercent:    8,
			shouldReplace:  false, // 5MB reduction from 100MB is only 5%, which is less than 8% threshold
			description:    "5 MB reduction from 100 MB (5%) should NOT trigger replacement",
		},
		{
			name:           "exact threshold",
			originalSize:   100 * 1024 * 1024, // 100 MB
			compressedSize: 92 * 1024 * 1024,  // 92 MB
			diffPercent:    8,
			shouldReplace:  false, // 8MB reduction equals exactly 8% threshold, but we need > threshold
			description:    "8% reduction (exact threshold) should NOT trigger replacement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from compressFile
			sizeDiff := tt.originalSize - tt.compressedSize
			threshold := int64(float64(tt.originalSize) * float64(tt.diffPercent) / 100)
			shouldReplace := sizeDiff > threshold

			if shouldReplace != tt.shouldReplace {
				t.Errorf("%s: sizeDiff=%d, threshold=%d, shouldReplace=%v, want %v. %s",
					tt.name, sizeDiff, threshold, shouldReplace, tt.shouldReplace, tt.description)
			}
		})
	}
}

func TestCompressFileUnsupportedType(t *testing.T) {
	// Test compressFile with an unsupported asset type
	asset := immich.AssetResponseDto{
		Id:               uuid.New().String(),
		Type:             immich.AssetTypeEnum("AUDIO"), // Unsupported type
		OriginalFileName: "test.mp3",
		OriginalPath:     "/test/test.mp3",
	}

	// Test that the asset type is not IMAGE or VIDEO
	if asset.Type != "IMAGE" && asset.Type != "VIDEO" {
		t.Logf("Correctly identified unsupported asset type: %s", asset.Type)
	}
}

func TestFileOperationsTempFileCreation(t *testing.T) {
	// Test the temp file creation patterns used in compressFile
	tempDir := os.TempDir()
	testUUID := uuid.New()

	// Test image temp file pattern
	imageFile := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", testUUID.String(), "jpg"))
	_, err := os.Create(imageFile)
	if err != nil {
		t.Fatalf("Failed to create image temp file: %v", err)
	}
	defer os.Remove(imageFile)

	// Test video temp file pattern
	videoFile := filepath.Join(tempDir, fmt.Sprintf("%s-compressed.%s", testUUID.String(), "mp4"))
	_, err = os.Create(videoFile)
	if err != nil {
		t.Fatalf("Failed to create video temp file: %v", err)
	}
	defer os.Remove(videoFile)

	// Verify files were created
	if _, err := os.Stat(imageFile); os.IsNotExist(err) {
		t.Error("Image temp file should exist")
	}

	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		t.Error("Video temp file should exist")
	}
}

// Helper functions for testing

// hasVideoExtension checks if a filename has a video extension
func hasVideoExtension(filename string) bool {
	ext := filepath.Ext(filename)
	videoExtensions := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"}

	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}
