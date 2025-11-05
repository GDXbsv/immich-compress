package immich

import (
	"testing"
	"time"
)

func TestAssetResponseDto_CompressedAfter(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		tags      *[]TagResponseDto
		timestamp time.Time
		expected  bool
	}{
		{
			name:      "no tags",
			tags:      nil,
			timestamp: baseTime,
			expected:  true,
		},
		{
			name:      "empty tags",
			tags:      &[]TagResponseDto{},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "no compressed_at tag",
			tags: &[]TagResponseDto{
				{Name: "some_other_tag", Value: "some_value"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "empty compressed_at tag value",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: ""},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "compressed after given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 14:00:00"},
			},
			timestamp: baseTime,
			expected:  false,
		},
		{
			name: "compressed before given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 10:00:00"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "compressed exactly at given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 12:00:00"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "invalid timestamp format",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: "invalid-timestamp"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "multiple tags with compressed_at",
			tags: &[]TagResponseDto{
				{Name: "tag1", Value: "value1"},
				{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 14:00:00"},
				{Name: "tag2", Value: "value2"},
			},
			timestamp: baseTime,
			expected:  false,
		},
		{
			name: "future timestamp comparison",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 10:00:00"},
			},
			timestamp: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := &AssetResponseDto{
				Tags: tt.tags,
			}

			result := asset.CompressedAfter(tt.timestamp)

			if result != tt.expected {
				t.Errorf("CompressedAfter() = %v, expected %v for test case: %s", result, tt.expected, tt.name)
			}
		})
	}
}

func TestAssetResponseDto_CompressedAfter_EdgeCases(t *testing.T) {
	t.Run("leap year timestamp", func(t *testing.T) {
		leapDay := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC) // 2024 is a leap year

		tags := &[]TagResponseDto{
			{Name: TAG_COMPRESSED_AT, Value: "2024-03-01 12:00:00"},
		}

		asset := &AssetResponseDto{Tags: tags}
		result := asset.CompressedAfter(leapDay)

		if result {
			t.Errorf("Expected false for leap year test, got %v", result)
		}
	})

	t.Run("year boundary", func(t *testing.T) {
		endOfYear := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

		tags := &[]TagResponseDto{
			{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 00:00:00"},
		}

		asset := &AssetResponseDto{Tags: tags}
		result := asset.CompressedAfter(endOfYear)

		if result {
			t.Errorf("Expected false for year boundary test, got %v", result)
		}
	})

	t.Run("timezone handling", func(t *testing.T) {
		// Test with a different timezone (though the tag format doesn't include timezone info)
		estTime := time.Date(2024, 1, 1, 7, 0, 0, 0, time.FixedZone("EST", -5*60*60))

		tags := &[]TagResponseDto{
			{Name: TAG_COMPRESSED_AT, Value: "2024-01-01 12:00:00"},
		}

		asset := &AssetResponseDto{Tags: tags}
		result := asset.CompressedAfter(estTime)

		// Since the tag is parsed as UTC and compared with EST time (which is 5 hours behind),
		// the compressed time should be considered "after" the EST time
		if !result {
			t.Errorf("Expected true for timezone test, got %v", result)
		}
	})
}
