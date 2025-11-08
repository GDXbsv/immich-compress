package compress

import (
	"context"
	"fmt"
	"os"

	"immich-compress/immich"

	"github.com/oapi-codegen/runtime/types"
)

type compress interface {
	compress(ctx context.Context, client *immich.ClientSimple, asset immich.AssetResponseDto) (*os.File, error)
}

func compressFile(ctx context.Context, client *immich.ClientSimple, asset immich.AssetResponseDto, diffPercent int, imageConfig ImageConfig, videoConfig VideoConfig) error {
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

	file, err := compress.compress(ctx, client, asset)
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

	var uuidNew *types.UUID
	if sizeOrig-sizeNew > int64(float64(sizeOrig)*(float64(diffPercent)/100)) {
		uuidNew, err = uploadFile(client, asset, file)
		if err != nil {
			return err
		}
		err = client.TagCompressedAdd(*uuidNew)
		if err != nil {
			return err
		}
		client.AssetDelete(*uuidNew, false)
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

	return nil
}

func bytesToMB(bytes int64) float64 {
	return float64(bytes) / float64(1024*1024)
}

func uploadFile(client *immich.ClientSimple, asset immich.AssetResponseDto, file *os.File) (*types.UUID, error) {
	r, err := client.AssetUploadCopy(asset, file)
	if err != nil {
		return nil, fmt.Errorf("can not upload new file: %w", err)
	}
	return r, err
}
