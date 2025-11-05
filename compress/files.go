package compress

import (
	"fmt"

	"immich-compress/immich"

	"github.com/cshum/vipsgen/vips"
	"github.com/google/uuid"
)

func compresFile(client immich.ClientSimple, asset immich.AssetResponseDto, imageConfig ImageConfig) error {
	switch asset.Type {
	case "IMAGE":
		fileBytes, err := compressImage(client, asset, imageConfig)
		if err != nil {
			return err
		}
		err = uploadBuffer(client, asset, fileBytes)
		if err != nil {
			return err
		}
		return nil

	case "VIDEO":
		// compressVideo(asset)
		return fmt.Errorf("VIDEO not yet supported")
	default:
		return fmt.Errorf("we do not support type: %s", asset.Type)
	}
}

func uploadBuffer(client immich.ClientSimple, asset immich.AssetResponseDto, fileBuffer []byte) error {
	// // Parse asset UUID
	// assetUUID, err := immich.UUUIDOfString(asset.Id)
	// if err != nil {
	// 	return err
	// }

	// // Get current timestamp for __compressed__ tag
	// timestamp := time.Now().Format("2006-01-02 15:04:05")

	// tagRootId, tagRoot, err := client.TagFindCreate(TAG_ROOT, nil)

	// tagResponse, err := client.CreateTagWithResponse(ctx, tagRequest)
	// if err != nil {
	// 	return fmt.Errorf("failed to create compressed tag: %w", err)
	// }

	// // Get the tag ID from the response
	// if tagResponse.JSON201 != nil {
	// 	tagID = tagResponse.JSON201.Id
	// } else {
	// 	// If tag creation didn't return a new tag, it might already exist
	// 	// Get all existing tags and find our compressed tag
	// 	tagsResponse, err := client.GetAllTagsWithResponse(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get existing tags: %w", err)
	// 	}

	// 	// Find the compressed tag
	// 	if tagsResponse.JSON200 != nil {
	// 		for _, tag := range *tagsResponse.JSON200 {
	// 			if tag.Name == compressedTag {
	// 				tagID = tag.Id
	// 				break
	// 			}
	// 		}
	// 	}

	// 	// If we still don't have a tag ID, try to create a new one
	// 	if tagID == "" {
	// 		return fmt.Errorf("failed to create or find compressed tag")
	// 	}
	// }

	// // Create the file reader for the compressed buffer
	// fileReader := bytes.NewReader(fileBuffer)

	// // Prepare replace asset parameters - preserve the original file creation time
	// params := &immich.ReplaceAssetParams{
	// 	// Keep the original file creation timestamp and metadata
	// 	Key:  &asset.DeviceAssetId,
	// 	Slug: &asset.OriginalFileName,
	// }

	// // Replace the asset with the compressed version
	// _, err = client.ReplaceAssetWithBody(ctx, assetUUID, params, "application/octet-stream", fileReader)
	// if err != nil {
	// 	return fmt.Errorf("failed to replace asset '%s': %w", asset.OriginalFileName, err)
	// }

	// // Tag the replaced asset with __compressed__ tag
	// if tagID != "" {
	// 	tagIDUUID, parseErr := uuid.Parse(tagID)
	// 	if parseErr != nil {
	// 		return fmt.Errorf("invalid tag UUID format: %w", parseErr)
	// 	}

	// 	assetUUIDs := []types.UUID{assetUUID}
	// 	tagUUIDs := []types.UUID{tagIDUUID}

	// 	bulkTagRequest := immich.BulkTagAssetsJSONRequestBody{
	// 		AssetIds: assetUUIDs,
	// 		TagIds:   tagUUIDs,
	// 	}

	// 	_, err = client.BulkTagAssetsWithResponse(ctx, bulkTagRequest)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to tag asset '%s' with compressed tag: %w", asset.OriginalFileName, err)
	// 	}
	// }

	return nil
}

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

var ImageAvailableFormats = []ImageFormat{JPG, JPEG, JXL, WEBP, HEIF}

func compressImage(client immich.ClientSimple, asset immich.AssetResponseDto, imageConfig ImageConfig) ([]byte, error) {
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
	// Use resp.Body (which is []byte) directly
	source := vips.NewSource(resp.Body)
	defer source.Close() // source needs to remain available during image lifetime

	image, err := vips.NewVipsloadSource(source, vips.DefaultVipsloadSourceOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer image.Close() // always close images to free memory

	var imageBytes []byte
	var exportErr error

	// 3. Use a switch to call the correct exporter
	switch imageConfig.Format {
	case JPEG, JPG:
		options := vips.DefaultJpegsaveBufferOptions()
		options.Q = imageConfig.Quality
		options.Keep = vips.KeepAll
		imageBytes, exportErr = image.JpegsaveBuffer(options)

	case JXL:
		options := vips.DefaultJxlsaveBufferOptions()
		options.Q = imageConfig.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.JxlsaveBuffer(options)

	case WEBP:
		options := vips.DefaultWebpsaveBufferOptions()
		options.Q = imageConfig.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.WebpsaveBuffer(options)

	case HEIF:
		options := vips.DefaultHeifsaveBufferOptions()
		options.Q = imageConfig.Quality
		options.Keep = vips.KeepAll
		options.Effort = 9
		imageBytes, exportErr = image.HeifsaveBuffer(options)

	default:
		return nil, fmt.Errorf("unsupported output format: %s", imageConfig.Format)
	}

	// Check for errors during the export
	if exportErr != nil {
		return nil, fmt.Errorf("failed to export image to %s: %w", imageConfig.Format, exportErr)
	}

	return imageBytes, nil
}
