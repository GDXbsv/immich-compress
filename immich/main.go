// Package immich
package immich

import (
	"context"
	"fmt"
	"net/http"
)

type ClientSimple struct {
	client   *Client
	ctx      context.Context
	parallel int
}

func NewClientSimple(ctx context.Context, parralel int, baseURL string, apiKey string) (*ClientSimple, error) {
	// Create a new client.
	// You must provide an http.Client that adds the API key to every request.
	client, err := NewClient(baseURL, WithRequestEditorFn(
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("x-api-key", apiKey)
			return nil
		}))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Now you can call API functions.
	// Let's ping the server to see if it's running.
	resp, err := client.PingServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging server: %w", err)
	}

	// Always close the response body
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close response body: %v\n", closeErr)
		}
	}()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server ping failed with status: %v", resp.Status)
	}

	return &ClientSimple{client: client, ctx: ctx, parallel: parralel}, nil
}

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
		}

		if resp.StatusCode != http.StatusOK {
			ch <- struct {
				Asset AssetResponseDto
				Err   error
			}{Asset: AssetResponseDto{}, Err: fmt.Errorf("bad status code: %s", resp.Status)}
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
