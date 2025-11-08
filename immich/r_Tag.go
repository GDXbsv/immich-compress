package immich

import (
	"fmt"

	"github.com/oapi-codegen/runtime/types"
)

const (
	TAG_ROOT       = "__immich-compress__"
	TAG_COMPRESSED = "__compressed__"
)

func (c *ClientSimple) TagCompressedAdd(assetID types.UUID) error {
	_, err := c.client.BulkTagAssetsWithResponse(c.ctx, TagBulkAssetsDto{
		AssetIds: []types.UUID{assetID},
		TagIds:   []types.UUID{c.tags.compressedID},
	})
	if err != nil {
		return fmt.Errorf("failed to attach tags: %w", err)
	}

	return err
}

func (c *ClientSimple) tagCompressedAt() (types.UUID, error) {
	tagRootID, _, err := c.tagFindCreate(TAG_ROOT, nil)
	if err != nil {
		return tagRootID, err
	}
	tagCompressID, _, err := c.tagFindCreate(TAG_COMPRESSED, &tagRootID)
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
}
