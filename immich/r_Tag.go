package immich

import (
	"github.com/oapi-codegen/runtime/types"
)

const (
	TAG_ROOT                 = "__immich-compress__"
	TAG_COMPRESSED_AT        = "__compressed_at__"
	TAG_COMPRESSED_AT_FORMAT = "2006-01-02 15:04:05"
)

func (c *ClientSimple) tagCompressedAt() (types.UUID, error) {
	tagRootID, _, err := c.tagFindCreate(TAG_ROOT, nil)
	if err != nil {
		return tagRootID, err
	}
	tagCompressID, _, err := c.tagFindCreate(TAG_COMPRESSED_AT, &tagRootID)
	if err != nil {
		return tagCompressID, err
	}

	return tagCompressID, err
}

func (c *ClientSimple) tagFindCreate(name string, parent *types.UUID) (types.UUID, *TagResponseDto, error) {
	var uuid types.UUID
	tagExists := false
	var tagFound *TagResponseDto
	r, err := c.client.GetAllTagsWithResponse(c.ctx)
	if err != nil {
		return uuid, tagFound, err
	}
	for _, tagDto := range *r.JSON200 {
		if name == tagDto.Name {
			tagExists = true
			tagFound = &tagDto
		}
	}

	if !tagExists {
		rc, err := c.client.CreateTagWithResponse(c.ctx, CreateTagJSONRequestBody{
			Name:     name,
			ParentId: parent,
		})
		if err != nil {
			return uuid, tagFound, err
		}
		tagFound = rc.JSON201
	}

	uuid, err = UUUIDOfString(tagFound.Id)

	return uuid, tagFound, err

	// respParsed, err := ParseDownloadAssetResponse(r)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing download response: %w", err)
	// }
}
