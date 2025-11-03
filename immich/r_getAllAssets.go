package immich

import (
	"fmt"
	"io"
	"net/http"
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
		resp, err := c.client.SearchAssets(c.ctx, SearchAssetsJSONRequestBody{})
		if err != nil {
			ch <- struct {
				Asset AssetResponseDto
				Err   error
			}{Asset: AssetResponseDto{}, Err: fmt.Errorf("error getting assets: %w", err)}

			return
		}

		if resp.StatusCode != http.StatusOK {
			defer func() {
				_ = resp.Body.Close()
			}()
			bodyBytes, _ := io.ReadAll(resp.Body)
			ch <- struct {
				Asset AssetResponseDto
				Err   error
			}{Asset: AssetResponseDto{}, Err: fmt.Errorf("bad status code: %s, body: %s", resp.Status, string(bodyBytes))}

			return
		}

		parsedResp, err := ParseSearchAssetsResponse(resp)
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
				ch <- struct {
					Asset AssetResponseDto
					Err   error
				}{Asset: item, Err: nil}
			}
		}
	}()

	return ch
}
