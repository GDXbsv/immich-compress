// Package compress
package compress

import (
	"context"
	"fmt"
	"slices"
	"sync/atomic"
	"time"

	"immich-compress/immich"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Config holds configuration for compression
type Config struct {
	Parallel       int
	Limit          int
	AssetType      string
	AssetUUIDs     []string
	Server         string
	APIKey         string
	After          time.Time
	DiffPercent    int
	ImageFormat    ImageFormat
	ImageQuality   int
	VideoContainer VideoContainer
	VideoFormat    VideoFormat
	VideoQuality   int
}

func Compressing(ctx context.Context, config Config) error {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(config.Parallel)
	client, err := immich.NewClientSimple(gCtx, config.Parallel, config.Server, config.APIKey)
	if err != nil {
		return err
	}

	var counter int32 = 0

	searchOption := immich.SearchAssetsJSONRequestBody{}
	if config.AssetType != "ALL" {
		typeAsset := (immich.AssetTypeEnum)(config.AssetType)
		searchOption.Type = &typeAsset
	}
	if len(config.AssetUUIDs) == 1 {
		UUIDstring := config.AssetUUIDs[0]
		UUID, err := uuid.Parse(UUIDstring)
		if err != nil {
			return err
		}
		searchOption.Id = &UUID
	}
	ch := client.AssetSearch(config.Limit, searchOption)
	// Start the workers. Instead of 'for range parallel', we simply
	// read from the channel and run g.Go() for *each* element.
	// SetLimit(parallel) will take care of the limit.
	for asset := range ch {
		// Pass 'asset' to the closure to avoid race conditions
		asset := asset

		g.Go(func() error {
			// Check for cancellation. gCtx.Done() will trigger if:
			// 1. Parent 'ctx' is cancelled.
			// 2. Another goroutine in 'g' returned an error.
			select {
			case <-gCtx.Done():
				return gCtx.Err() // Return cancellation error
			default:
			}

			if asset.Err != nil {
				// Just return the error.
				// errgroup will automatically call cancel() for gCtx.
				return asset.Err
			}

			if len(config.AssetUUIDs) > 0 {
				if !slices.Contains(config.AssetUUIDs, asset.Asset.Id) {
					return nil
				}
			}

			if !asset.Asset.CompressedAfter(config.After) {
				return nil
			}
			// Process the asset here
			fmt.Printf("Processing file: %#v\n", asset.Asset.Id)
			err = compressFile(gCtx, client, asset.Asset, config.DiffPercent, ImageConfig{
				Format:  config.ImageFormat,
				Quality: config.ImageQuality,
			}, VideoConfig{
				Container: config.VideoContainer,
				Format:    config.VideoFormat,
				Quality:   config.VideoQuality,
			})
			if err != nil {
				return err
			}
			atomic.AddInt32(&counter, 1)

			return nil
		})
	}

	// g.Wait() waits for all goroutines to complete (like wg.Wait())
	// and returns the FIRST non-zero error returned by
	// any of the goroutines.
	if err := g.Wait(); err != nil {
		// If there was an error (including cancellation), we return it
		return err
	}

	fmt.Printf("Processed files: %d\n", counter)

	return nil
}
