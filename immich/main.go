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
