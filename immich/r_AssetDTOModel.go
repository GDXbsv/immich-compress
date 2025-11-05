package immich

import (
	"fmt"
	"time"
)

func (a *AssetResponseDto) GetTag(tagName string) string {
	value := ""
	// check for nil before dereferencing!
	if a.Tags != nil {
		// Now it's safe to use *
		for _, tag := range *a.Tags {
			if tag.Name == tagName {
				value = tag.Value
				return value
			}
		}
	} else {
		return value
	}

	return value
}

func (a *AssetResponseDto) CompressedAfter(timestamp time.Time) bool {
	when := a.GetTag(TAG_COMPRESSED_AT)
	var timeCompressed time.Time
	var err error

	if when == "" {
		return false
	}
	timeCompressed, err = time.Parse(TAG_COMPRESSED_AT_FORMAT, when)
	if err != nil {
		fmt.Printf("✗ Tag '%s' with value '%s' parsed with error: %w\n", TAG_COMPRESSED_AT, timestamp, err)
	}

	if timeCompressed.Compare(timestamp) != 1 {
		return true
	}

	return false
}
