// Package compress
package compress

import (
	"context"
	"fmt"
	"sync"

	"immich-compress/immich"
)

func Compressing(ctx context.Context, parallel int, server string, apiKey string) error {
	ctxCancel, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	client, err := immich.NewClientSimple(ctxCancel, parallel, server, apiKey)
	if err != nil {
		return err
	}

	ch := client.GetAllAssets()

	// Create a worker pool to process assets in parallel
	var wg sync.WaitGroup

	// Create a channel to collect errors from workers
	errCh := make(chan error, parallel)

	// Start worker goroutines
	for range parallel {
		wg.Go(func() {
			for asset := range ch {
				if asset.Err != nil {
					select {
					case errCh <- asset.Err:
					default:
					}
					continue
				}

				// Process the asset here
				fmt.Printf("Processing file: %#v\n", asset.Asset)
				// TODO: Add actual compression logic here
			}
		})
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check for any errors
	select {
	case err := <-errCh:
		return err
	default:
	}

	return nil
}
