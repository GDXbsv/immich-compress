package compress

import (
	"fmt"
	"os"

	"immich-compress/immich"
)

type compress interface {
	compress(client *immich.ClientSimple, asset immich.AssetResponseDto) (*os.File, error)
}

func compressFile(client *immich.ClientSimple, asset immich.AssetResponseDto, diffPercent int, imageConfig ImageConfig, videoConfig VideoConfig) error {
	skipped := true
	sizeOrig := *asset.ExifInfo.FileSizeInByte
	var sizeNew int64
	var compress compress
	switch asset.Type {
	case "IMAGE":
		compress = &imageConfig
	case "VIDEO":
		compress = &videoConfig
	default:
		return fmt.Errorf("we do not support type: %s", asset.Type)
	}

	file, err := compress.compress(client, asset)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file stats: %w", err)
	}

	// 3. Get the size from the FileInfo
	sizeNew = fileInfo.Size()
	if sizeNew == 0 {
		return fmt.Errorf("compressed size is 0 most likely we have an error")
	}

	if sizeOrig-sizeNew > int64(float64(sizeOrig)*(float64(diffPercent)/100)) {
		err = uploadFile(client, asset, file)
		if err != nil {
			return err
		}
		skipped = false
	}

	sizeOrigMB := bytesToMB(sizeOrig)
	sizeNewMB := bytesToMB(sizeNew)
	sizeSavedMB := bytesToMB(sizeOrig - sizeNew)

	if skipped {
		fmt.Printf("✗ Skipped: %s (Original: %.2f MB, Converted: %.2f MB, No size reduction)\n", asset.OriginalFileName, sizeOrigMB, sizeNewMB)
		return nil
	}
	fmt.Printf("✓ Replaced: %s (Original: %.2f MB, Converted: %.2f MB, Saved: %.2f MB)\n", asset.OriginalFileName, sizeOrigMB, sizeNewMB, sizeSavedMB)

	// todo upload to immich
	// todo update tags
	// todo remove orig to trash

	return nil
}

func bytesToMB(bytes int64) float64 {
	return float64(bytes) / float64(1024*1024)
}

func uploadFile(client *immich.ClientSimple, asset immich.AssetResponseDto, file *os.File) error {
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
