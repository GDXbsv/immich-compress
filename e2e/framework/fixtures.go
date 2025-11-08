package framework

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TestDataManager struct {
	BaseDir string
}

type MockAsset struct {
	ID           string `json:"id"`
	OriginalPath string `json:"original_path"`
	FileName     string `json:"file_name"`
	MIMEType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
	CreatedAt    string `json:"created_at"`
	AlbumID      string `json:"album_id,omitempty"`
}

type MockAlbum struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	AssetCount int         `json:"asset_count"`
	CreatedAt  string      `json:"created_at"`
	Assets     []MockAsset `json:"assets"`
}

func NewTestDataManager(baseDir string) *TestDataManager {
	return &TestDataManager{BaseDir: baseDir}
}

func (tm *TestDataManager) CreateTestImage(filename, format string, size int) (string, error) {
	imgDir := filepath.Join(tm.BaseDir, "images")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return "", err
	}

	imgPath := filepath.Join(imgDir, filename)

	// Create a minimal valid image file
	var data []byte
	switch format {
	case "jpeg", "jpg":
		data = createMinimalJPEG(size)
	case "png":
		data = createMinimalPNG(size)
	case "webp":
		data = createMinimalWebP(size)
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}

	return imgPath, os.WriteFile(imgPath, data, 0644)
}

func (tm *TestDataManager) CreateTestVideo(filename, format string, size int) (string, error) {
	vidDir := filepath.Join(tm.BaseDir, "videos")
	if err := os.MkdirAll(vidDir, 0755); err != nil {
		return "", err
	}

	vidPath := filepath.Join(vidDir, filename)

	// Create a minimal valid video file
	var data []byte
	switch format {
	case "mp4":
		data = createMinimalMP4(size)
	case "webm":
		data = createMinimalWebM(size)
	default:
		return "", fmt.Errorf("unsupported video format: %s", format)
	}

	return vidPath, os.WriteFile(vidPath, data, 0644)
}

func (tm *TestDataManager) CreateMockImmichServer() error {
	mockDir := filepath.Join(tm.BaseDir, "immich")
	if err := os.MkdirAll(mockDir, 0755); err != nil {
		return err
	}

	// Create mock API responses
	albums := []MockAlbum{
		{
			ID:         "album-1",
			Name:       "Test Album",
			AssetCount: 5,
			CreatedAt:  time.Now().Format(time.RFC3339),
			Assets: []MockAsset{
				{
					ID:           "asset-1",
					OriginalPath: "/photos/IMG_001.jpg",
					FileName:     "IMG_001.jpg",
					MIMEType:     "image/jpeg",
					FileSize:     2048000,
					CreatedAt:    time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				},
				{
					ID:           "asset-2",
					OriginalPath: "/photos/video.mp4",
					FileName:     "video.mp4",
					MIMEType:     "video/mp4",
					FileSize:     10485760,
					CreatedAt:    time.Now().Add(-12 * time.Hour).Format(time.RFC3339),
				},
			},
		},
	}

	// Save mock data
	albumsData, err := json.MarshalIndent(albums, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(mockDir, "mock-albums.json"), albumsData, 0644)
}

func (tm *TestDataManager) Cleanup() {
	os.RemoveAll(tm.BaseDir)
}

// Helper functions to create minimal valid files
func createMinimalJPEG(size int) []byte {
	// Create a minimal valid JPEG file
	if size < 100 {
		size = 100
	}
	data := make([]byte, size)
	// JPEG magic bytes
	data[0] = 0xFF
	data[1] = 0xD8
	// Add some basic JPEG structure
	data[2] = 0xFF
	data[3] = 0xE0
	// Fill with some data
	for i := 4; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}

func createMinimalPNG(size int) []byte {
	// Create a minimal valid PNG file
	if size < 50 {
		size = 50
	}
	data := make([]byte, size)
	// PNG magic bytes
	copy(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	// Fill with some data
	for i := 8; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}

func createMinimalWebP(size int) []byte {
	// Create a minimal valid WebP file
	if size < 30 {
		size = 30
	}
	data := make([]byte, size)
	// WebP magic bytes
	copy(data, []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50})
	// Fill with some data
	for i := 12; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}

func createMinimalMP4(size int) []byte {
	// Create a minimal valid MP4 file
	if size < 100 {
		size = 100
	}
	data := make([]byte, size)
	// MP4 magic bytes
	copy(data, []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x69, 0x73, 0x6F, 0x6D})
	// Fill with some data
	for i := 12; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}

func createMinimalWebM(size int) []byte {
	// Create a minimal valid WebM file
	if size < 50 {
		size = 50
	}
	data := make([]byte, size)
	// WebM magic bytes
	copy(data, []byte{0x1A, 0x45, 0xDF, 0xA3})
	// Fill with some data
	for i := 4; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}
