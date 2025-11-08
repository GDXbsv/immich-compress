package immich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (c *ClientSimple) AssetUploadCopy(asset AssetResponseDto, file *os.File) (*openapi_types.UUID, error) {
	origNameWithoutExt := strings.TrimSuffix(asset.OriginalFileName, filepath.Ext(asset.OriginalFileName))

	uuidOrig, err := uuid.Parse(asset.Id)
	if err != nil {
		return nil, err
	}
	// Get metadata from asset if it exists
	var metadata []AssetMetadataUpsertItemDto
	if asset.HasMetadata {
		resp, err := c.client.GetAssetMetadataWithResponse(c.ctx, uuidOrig)
		if err == nil && resp.JSON200 != nil {
			for _, item := range *resp.JSON200 {
				metadata = append(metadata, AssetMetadataUpsertItemDto{
					Key:   item.Key,
					Value: item.Value,
				})
			}
		}
	}
	params := &uploadAssetBody{
		DeviceAssetID:    asset.DeviceAssetId,
		DeviceID:         asset.DeviceId,
		Duration:         asset.Duration,
		FileCreatedAt:    asset.FileCreatedAt,
		FileModifiedAt:   time.Now(),
		Filename:         origNameWithoutExt + filepath.Ext(file.Name()),
		IsFavorite:       asset.IsFavorite,
		Metadata:         metadata,
		Visibility:       string(asset.Visibility),
		LivePhotoVideoID: *asset.LivePhotoVideoId,
	}

	// 3. Create the multipart body and content type
	body, contentType, err := assetUploadMultipartBody(params, file)
	if err != nil {
		return nil, err
	}

	// Stream the file directly without loading into memory
	rUp, err := c.client.UploadAssetWithBodyWithResponse(c.ctx, &UploadAssetParams{}, contentType, body)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	uuidNew, err := uuid.Parse(rUp.JSON201.Id)
	if err != nil {
		return nil, err
	}
	t := true
	// copy asset with API
	_, err = c.client.CopyAssetWithResponse(c.ctx, CopyAssetJSONRequestBody{
		Albums:      &t,
		Favorite:    &t,
		SharedLinks: &t,
		Sidecar:     &t,
		SourceId:    uuidOrig,
		Stack:       new(bool),
		TargetId:    uuidNew,
	})
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	// Copy tags from old asset to new one
	if asset.Tags != nil && len(*asset.Tags) > 0 {
		tagIds := make([]openapi_types.UUID, 0, len(*asset.Tags))
		for _, tag := range *asset.Tags {
			tagUUID, err := uuid.Parse(tag.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to parse tag UUID '%s': %w", tag.Id, err)
			}
			tagIds = append(tagIds, tagUUID)
		}

		_, err = c.client.BulkTagAssetsWithResponse(c.ctx, TagBulkAssetsDto{
			AssetIds: []openapi_types.UUID{uuidNew},
			TagIds:   tagIds,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to copy tags: %w", err)
		}
	}

	return &uuidNew, nil
}

type uploadAssetBody struct {
	DeviceAssetID    string                       `json:"deviceAssetId"`
	DeviceID         string                       `json:"deviceId"`
	Duration         string                       `json:"duration"`
	FileCreatedAt    time.Time                    `json:"fileCreatedAt"`
	FileModifiedAt   time.Time                    `json:"fileModifiedAt"`
	Filename         string                       `json:"filename"`
	IsFavorite       bool                         `json:"isFavorite"`
	LivePhotoVideoID string                       `json:"livePhotoVideoId"`
	Metadata         []AssetMetadataUpsertItemDto `json:"metadata"`
	// SidecarData is binary, would be handled like assetData if needed
	// Visibility is AssetVisibility, likely a string enum
	Visibility string `json:"visibility"`
}

// createMultipartBody is the core logic
func assetUploadMultipartBody(params *uploadAssetBody, assetFile *os.File) (io.Reader, string, error) {
	// Create a new buffer to write the multipart body to
	bodyBuf := &bytes.Buffer{}
	// Create a new multipart writer
	writer := multipart.NewWriter(bodyBuf)

	// --- 1. Stream the 'assetData' file ---
	// CreateFormFile returns an io.Writer for the file part
	part, err := writer.CreateFormFile("assetData", params.Filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}
	// Stream the file's content directly into the multipart part
	// This avoids loading the entire file into memory
	_, err = io.Copy(part, assetFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to copy file to multipart: %w", err)
	}

	// --- 2. Add all other metadata fields ---

	// Simple string fields
	_ = writer.WriteField("deviceAssetId", params.DeviceAssetID)
	_ = writer.WriteField("deviceId", params.DeviceID)
	_ = writer.WriteField("duration", params.Duration)
	_ = writer.WriteField("filename", params.Filename)
	_ = writer.WriteField("livePhotoVideoId", params.LivePhotoVideoID)
	_ = writer.WriteField("visibility", params.Visibility)

	// Boolean field
	_ = writer.WriteField("isFavorite", fmt.Sprintf("%t", params.IsFavorite))

	// DateTime fields (format as ISO 8601 string)
	_ = writer.WriteField("fileCreatedAt", params.FileCreatedAt.Format(time.RFC3339))
	_ = writer.WriteField("fileModifiedAt", params.FileModifiedAt.Format(time.RFC3339))

	// --- 3. Add the 'metadata' (JSON array) ---
	// Marshal the metadata struct to a JSON string
	metaJSON, err := json.Marshal(params.Metadata)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	_ = writer.WriteField("metadata", string(metaJSON))

	// --- 4. Add 'sidecarData' (if you have it) ---
	// If you had sidecar data (e.g., from another file), you would
	// use CreateFormFile again, just like with assetData.
	// Example:
	// sidecarPart, err := writer.CreateFormFile("sidecarData", "my-sidecar.xmp")
	// _, err = io.Copy(sidecarPart, sidecarFileReader)

	// --- 5. Finalize ---
	// Close the writer to write the final boundary
	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Return the buffer (which is an io.Reader) and the content type
	return bodyBuf, writer.FormDataContentType(), nil
}
