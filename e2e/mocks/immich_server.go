package mocks

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MockImmichServer struct {
	Server       *httptest.Server
	URL          string
	MockDataPath string
}

type MockAsset struct {
	ID           string `json:"id"`
	OriginalPath string `json:"original_path"`
	FileName     string `json:"file_name"`
	MIMEType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
	CreatedAt    string `json:"created_at"`
}

type MockAlbum struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	AssetCount int         `json:"asset_count"`
	CreatedAt  string      `json:"created_at"`
	Assets     []MockAsset `json:"assets"`
}

type MockAPIResponse struct {
	Items []MockAlbum `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

func NewMockImmichServer(mockDataPath string) *MockImmichServer {
	mux := http.NewServeMux()

	server := &MockImmichServer{
		MockDataPath: mockDataPath,
		Server:       httptest.NewServer(mux),
	}

	// Set up routes
	mux.HandleFunc("/api/albums", server.handleAlbums)
	mux.HandleFunc("/api/assets/search", server.handleAssetSearch)
	mux.HandleFunc("/api/assets/", server.handleAssetOperations)
	mux.HandleFunc("/api/auth/login", server.handleLogin)
	mux.HandleFunc("/api/users/me", server.handleUserInfo)

	return server
}

func (m *MockImmichServer) handleAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// Load mock albums data
		albums, err := m.loadMockAlbums()
		if err != nil {
			http.Error(w, "Failed to load mock data", http.StatusInternalServerError)
			return
		}

		response := MockAPIResponse{
			Items: albums,
			Total: len(albums),
			Page:  1,
			Limit: 1000,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (m *MockImmichServer) handleAssetSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		// Return mock assets based on search criteria
		albums, err := m.loadMockAlbums()
		if err != nil {
			http.Error(w, "Failed to load mock data", http.StatusInternalServerError)
			return
		}

		var allAssets []MockAsset
		for _, album := range albums {
			allAssets = append(allAssets, album.Assets...)
		}

		response := struct {
			Total  int         `json:"total"`
			Assets []MockAsset `json:"assets"`
		}{
			Total:  len(allAssets),
			Assets: allAssets,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (m *MockImmichServer) handleAssetOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path

	if strings.Contains(path, "/download") {
		// Handle asset download
		if r.Method == http.MethodGet {
			// Return a minimal file response
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "mock file content")
			return
		}
	} else if strings.Contains(path, "/upload") {
		// Handle asset upload
		if r.Method == http.MethodPost {
			response := map[string]string{
				"id":     "uploaded-asset-id",
				"status": "success",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// Handle asset operations
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "deleted",
			})
			return
		}
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (m *MockImmichServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		// Mock successful login
		response := map[string]any{
			"accessToken": "mock-access-token",
			"userId":      "mock-user-id",
			"email":       "test@example.com",
			"firstName":   "Test",
			"lastName":    "User",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (m *MockImmichServer) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		response := map[string]any{
			"id":        "mock-user-id",
			"email":     "test@example.com",
			"firstName": "Test",
			"lastName":  "User",
			"createdAt": time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (m *MockImmichServer) loadMockAlbums() ([]MockAlbum, error) {
	mockFile := filepath.Join(m.MockDataPath, "mock-albums.json")

	if _, err := os.Stat(mockFile); os.IsNotExist(err) {
		// Create default mock data if file doesn't exist
		return m.createDefaultMockAlbums(), nil
	}

	data, err := os.ReadFile(mockFile)
	if err != nil {
		return nil, err
	}

	var albums []MockAlbum
	err = json.Unmarshal(data, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (m *MockImmichServer) createDefaultMockAlbums() []MockAlbum {
	return []MockAlbum{
		{
			ID:         "album-1",
			Name:       "Test Album",
			AssetCount: 2,
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
}

func (m *MockImmichServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

func (m *MockImmichServer) GetBaseURL() string {
	return m.Server.URL
}

func (m *MockImmichServer) Start() {
	log.Printf("Mock Immich server started at %s", m.Server.URL)
}
