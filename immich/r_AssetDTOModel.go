package immich

import (
	"time"
)

func (a *AssetResponseDto) GetTag(tagName string) string {
	value := ""
	// check for nil before dereferencing!
	if a.Tags != nil {
		// Now it's safe to use *
		for _, tag := range *a.Tags {
			if tag.Name == tagName {
				value = tag.Id
				return value
			}
		}
	} else {
		return value
	}

	return value
}

func (a *AssetResponseDto) CompressedAfter(timestamp time.Time) bool {
	id := a.GetTag(TAG_COMPRESSED)
	if id == "" {
		return false
	}

	timeCompressed := a.FileModifiedAt

	return timeCompressed.Compare(timestamp) == 1
}
