// Package immich
package immich

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type ClientSimple struct {
	client   *Client
	ctx      context.Context
	parralel int
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

	return &ClientSimple{client: client, ctx: ctx, parralel: parralel}, nil
}

func (c *ClientSimple) GetAllAssets() <-chan struct {asset AssetResponseDto, err error} {
	ch := make(chan struct {
		asset AssetResponseDto
		err    error
	}, c.parralel)
	defer close(ch)
	resp, err := c.client.SearchAssets(c.ctx, SearchAssetsJSONRequestBody{})
	if err != nil {
		ch <- struct {
			asset AssetResponseDto
			er    error
		}{asset: nil, err: fmt.Errorf("Error getting assets: %w", err)}

		return ch
	}

	if resp.StatusCode != http.StatusOK {
		ch <- struct {
			asset AssetResponseDto
			er    error
		}{asset: nil, err: fmt.Errorf("Bad status code: %s", resp.Status)}

		return ch
	}

	parsedResp, err := ParseSearchAssetsResponse(resp)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
		ch <- struct {
			asset AssetResponseDto
			er    error
		}{asset: nil, err: fmt.Errorf("Error parsing response: %w", err)}

		return ch
	}
	for _, item := range parsedResp.JSON200.Assets.Items {
		ch <- struct {
			asset AssetResponseDto
			er    error
		}{asset: item, err: nil}
	}

	return ch
}
