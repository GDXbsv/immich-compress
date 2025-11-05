package immich

import (
	"net/http"

	"github.com/oapi-codegen/runtime/types"
)

func (c *ClientSimple) AssetDownload(id types.UUID) (*http.Response, error) {
	r, err := c.clientRaw.DownloadAsset(c.ctx, id, nil)
	if err != nil {
		return nil, err
	}
	return r, nil

	// respParsed, err := ParseDownloadAssetResponse(r)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing download response: %w", err)
	// }
}
