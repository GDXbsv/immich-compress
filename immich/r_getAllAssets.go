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
		pageNum := float32(1)
		var respCurrent *http.Response
		var respNext *http.Response
		var err error
		var wg sync.WaitGroup
		for pageNum != 0 {
			wg.Wait()
			wg.Go(func() {
				respNext, err = c.client.SearchAssets(c.ctx, SearchAssetsJSONRequestBody{
					Page: &pageNum,
				})
			})
			if respNext == nil {
				wg.Wait()
			}
			respCurrent = respNext

			if err != nil {
				ch <- struct {
					Asset AssetResponseDto
					Err   error
				}{Asset: AssetResponseDto{}, Err: fmt.Errorf("error getting assets: %w", err)}

				return
			}

			if respCurrent.StatusCode != http.StatusOK {
				defer func() {
					_ = respCurrent.Body.Close()
				}()
				bodyBytes, _ := io.ReadAll(respCurrent.Body)
				ch <- struct {
					Asset AssetResponseDto
					Err   error
				}{Asset: AssetResponseDto{}, Err: fmt.Errorf("bad status code: %s, body: %s", respCurrent.Status, string(bodyBytes))}

				return
			}

			parsedResp, err := ParseSearchAssetsResponse(respCurrent)
			if err != nil {
				ch <- struct {
					Asset AssetResponseDto
					Err   error
				}{Asset: AssetResponseDto{}, Err: fmt.Errorf("error parsing response: %w", err)}
			}
			for _, item := range parsedResp.JSON200.Assets.Items {
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
			nextPage, _ := strconv.ParseFloat(*parsedResp.JSON200.Assets.NextPage, 32)
			pageNum = float32(nextPage)
		}
	}()

	return ch
}
