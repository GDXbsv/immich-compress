package immich

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func (c *ClientSimple) GetAllAssets() <-chan struct {
	Asset AssetResponseDto
	Err   error
} {
	ch := make(chan struct {
		Asset AssetResponseDto
		Err   error
	}, c.parallel)
	go func() {
		defer close(ch)
		respNext, nextPage, err := c.getAssets(1)
		if err != nil {
			ch <- struct {
				Asset AssetResponseDto
				Err   error
			}{Asset: AssetResponseDto{}, Err: err}

			return
		}
		var respCurrent *SearchAssetsResponse

		var wg sync.WaitGroup
		for {
			wg.Wait()
			wg.Go(func() {
				respNext, nextPage, err = c.getAssets(nextPage)
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
				select {
				case <-c.ctx.Done():
					return
				default:
					if !item.IsTrashed {
						ch <- struct {
							Asset AssetResponseDto
							Err   error
						}{Asset: item, Err: nil}
					}
				}
			}
		}
	}()

	return ch
}

func (c *ClientSimple) getAssets(nextPage float32) (*SearchAssetsResponse, float32, error) {
	if nextPage == 0 {
		return nil, 0, nil
	}
	r, err := c.client.SearchAssets(c.ctx, SearchAssetsJSONRequestBody{
		Page: &nextPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error getting assets: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		defer func() {
			_ = r.Body.Close()
		}()
		bodyBytes, _ := io.ReadAll(r.Body)

		return nil, 0, fmt.Errorf("bad status code: %s, body: %s", r.Status, string(bodyBytes))
	}

	respParsed, err := ParseSearchAssetsResponse(r)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing response: %w", err)
	}
	nextPage64, _ := strconv.ParseFloat(*respParsed.JSON200.Assets.NextPage, 32)
	nextPage = float32(nextPage64)

	return respParsed, nextPage, nil
}
