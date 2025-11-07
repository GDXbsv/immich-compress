package compress

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"immich-compress/immich"

	"github.com/cshum/vipsgen/vips"
	"github.com/google/uuid"
)

type ImageConfig struct {
	Format  ImageFormat
	Quality int
}

type ImageFormat string

const (
	JPG  ImageFormat = "jpg"
	JPEG ImageFormat = "jpeg"
	JXL  ImageFormat = "jxl"
	WEBP ImageFormat = "webp"
	HEIF ImageFormat = "heif"
)

var ImageFormatsAvailable = []ImageFormat{JPG, JPEG, JXL, WEBP, HEIF}

func (c *ImageConfig) compress(client *immich.ClientSimple, asset immich.AssetResponseDto) (*os.File, error) {
	uuid, err := uuid.Parse(asset.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid '%s': %w", asset.Id, err)
	}
	// Fetch an image from the client
	resp, err := client.AssetDownload(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body into a byte buffer
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Load image from buffer instead of source
	image, err := vips.NewImageFromBuffer(fileBytes, vips.DefaultLoadOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer image.Close() // always close images to free memory

	var imageBytes []byte
	var exportErr error

	// 3. Use a switch to call the correct exporter
	switch c.Format {
	case JPEG, JPG:
		options := vips.DefaultJpegsaveBufferOptions()
		options.Q = c.Quality
		options.Keep = vips.KeepAll
		imageBytes, exportErr = image.JpegsaveBuffer(options)

	case JXL:
		options := vips.DefaultJxlsaveBufferOptions()
		options.Q = c.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.JxlsaveBuffer(options)

	case WEBP:
		options := vips.DefaultWebpsaveBufferOptions()
		options.Q = c.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.WebpsaveBuffer(options)

	case HEIF:
		options := vips.DefaultHeifsaveBufferOptions()
		options.Q = c.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.HeifsaveBuffer(options)

	default:
		return nil, fmt.Errorf("unsupported output format: %s", c.Format)
	}

	// Check for errors during the export
	if exportErr != nil {
		return nil, fmt.Errorf("failed to export image to %s: %w", c.Format, exportErr)
	}

	// Create temporary output file
	fileOut, err := os.CreateTemp("", fmt.Sprintf("%s.%s", filepath.Base(asset.OriginalPath), c.Format))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer fileOut.Close()

	_, err = fileOut.Write(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to save image to temp file: %w", err)
	}
	fileOut.Close() // Close to ensure all data is written

	return fileOut, nil
}
