package immich

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func (c *ClientSimple) AssetSearch(limit int, search SearchAssetsJSONRequestBody) <-chan struct {
	Asset AssetResponseDto
	Err   error
} {
	ch := make(chan struct {
		Asset AssetResponseDto
		Err   error
	}, c.parallel)
	go func() {
		defer close(ch)
		var page float32 = 1
		search.Page = &page
		respNext, nextPage, err := c.getAssets(search)
		if err != nil {
			ch <- struct {
				Asset AssetResponseDto
				Err   error
			}{Asset: AssetResponseDto{}, Err: err}

			return
		}
		var respCurrent *SearchAssetsResponse
		processedCount := 0

		var wg sync.WaitGroup
		for {
			wg.Wait()
			wg.Go(func() {
				search.Page = &nextPage
				respNext, nextPage, err = c.getAssets(search)
				if err != nil {
					ch <- struct {
						Asset AssetResponseDto
						Err   error
					}{Asset: AssetResponseDto{}, Err: err}
				}
			})
			if err != nil {
				return
			}
			if respNext == nil {
				return
			}
			respCurrent = respNext
			for _, item := range respCurrent.JSON200.Assets.Items {
				// Check if limit is reached (limit 0 means no limit)
				if limit > 0 && processedCount >= limit {
					return
				}

				select {
				case <-c.ctx.Done():
					return
				default:
					if item.IsTrashed {
						continue
					}

					ch <- struct {
						Asset AssetResponseDto
						Err   error
					}{Asset: item, Err: nil}
					processedCount++
				}
			}
		}
	}()

	return ch
}

func (c *ClientSimple) getAssets(search SearchAssetsJSONRequestBody) (*SearchAssetsResponse, float32, error) {
	if search.Page == nil {
		return nil, 0, nil
	}
	nextPage := *search.Page
	if nextPage == 0 {
		return nil, 0, nil
	}
	r, err := c.client.SearchAssetsWithResponse(c.ctx, search)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting assets: %w", err)
	}

	if r.StatusCode() != http.StatusOK {
		return nil, 0, fmt.Errorf("bad status code: %s, body: %s", r.Status(), string(r.Body))
	}
	var nextPage32 float32
	if r.JSON200.Assets.NextPage != nil {
		nextPage64, _ := strconv.ParseFloat(*r.JSON200.Assets.NextPage, 32)
		nextPage32 = float32(nextPage64)
	} else {
		nextPage32 = 0
	}

	return r, nextPage32, nil
}
