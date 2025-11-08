// Package immich
package immich

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type ClientSimple struct {
	client    ClientWithResponsesInterface
	clientRaw ClientInterface
	ctx       context.Context
	parallel  int
	tags      struct {
		compressedID types.UUID
	}
}

func NewClientSimple(ctx context.Context, parralel int, baseURL string, apiKey string) (*ClientSimple, error) {
	// Create a new client.
	// You must provide an http.Client that adds the API key to every request.
	client, err := NewClientWithResponses(baseURL, WithRequestEditorFn(
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("x-api-key", apiKey)
			return nil
		}))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	clientSimple := &ClientSimple{client: client, clientRaw: client.ClientInterface, ctx: ctx, parallel: parralel}

	tagCompressedAtID, err := clientSimple.tagCompressedAt()
	if err != nil {
		return nil, fmt.Errorf("can not get/create tags: %w", err)
	}

	clientSimple.tags = struct{ compressedID types.UUID }{compressedID: tagCompressedAtID}

	return clientSimple, nil
}

func UUUIDOfString(id string) (types.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		err = fmt.Errorf("failed to parse UUID '%s': %w", id, err)
	}
	return uuid, err
}
