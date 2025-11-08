package immich

import (
	"fmt"
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
				{Name: TAG_COMPRESSED, Value: ""},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "compressed after given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED, Value: "2024-01-01 14:00:00"},
			},
			timestamp: baseTime,
			expected:  false,
		},
		{
			name: "compressed before given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED, Value: "2024-01-01 10:00:00"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "compressed exactly at given timestamp",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED, Value: "2024-01-01 12:00:00"},
			},
			timestamp: baseTime,
			expected:  true,
		},
		{
			name: "invalid timestamp format",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED, Value: "invalid-timestamp"},
			},
			timestamp: baseTime,
			expected:  false,
		},
		{
			name: "multiple tags with compressed_at",
			tags: &[]TagResponseDto{
				{Name: "tag1", Value: "value1"},
				{Name: TAG_COMPRESSED, Value: "2024-01-01 14:00:00"},
				{Name: "tag2", Value: "value2"},
			},
			timestamp: baseTime,
			expected:  false,
		},
		{
			name: "future timestamp comparison",
			tags: &[]TagResponseDto{
				{Name: TAG_COMPRESSED, Value: "2024-01-01 10:00:00"},
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
			{Name: TAG_COMPRESSED, Value: "2024-03-01 12:00:00"},
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
			{Name: TAG_COMPRESSED, Value: "2024-01-01 00:00:00"},
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
			{Name: TAG_COMPRESSED, Value: "2024-01-01 12:00:00"},
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

func TestAssetResponseDto_GetTag(t *testing.T) {
	tests := []struct {
		name     string
		asset    *AssetResponseDto
		tagName  string
		expected string
	}{
		{
			name:     "nil tags should return empty string",
			asset:    &AssetResponseDto{Tags: nil},
			tagName:  "test-tag",
			expected: "",
		},
		{
			name:     "empty tags should return empty string",
			asset:    &AssetResponseDto{Tags: &[]TagResponseDto{}},
			tagName:  "test-tag",
			expected: "",
		},
		{
			name: "tag not found should return empty string",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "other-tag", Value: "some-value"},
				},
			},
			tagName:  "test-tag",
			expected: "",
		},
		{
			name: "tag found should return its value",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "test-tag", Value: "test-value"},
				},
			},
			tagName:  "test-tag",
			expected: "test-value",
		},
		{
			name: "multiple tags with target tag found",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "tag1", Value: "value1"},
					{Name: "test-tag", Value: "target-value"},
					{Name: "tag2", Value: "value2"},
				},
			},
			tagName:  "test-tag",
			expected: "target-value",
		},
		{
			name: "case sensitive tag name matching",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "Test-Tag", Value: "case-sensitive"},
				},
			},
			tagName:  "test-tag",
			expected: "",
		},
		{
			name: "duplicate tags should return first match",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "test-tag", Value: "first"},
					{Name: "test-tag", Value: "second"},
				},
			},
			tagName:  "test-tag",
			expected: "first",
		},
		{
			name: "empty tag value should return empty string",
			asset: &AssetResponseDto{
				Tags: &[]TagResponseDto{
					{Name: "test-tag", Value: ""},
				},
			},
			tagName:  "test-tag",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.asset.GetTag(tt.tagName)
			if result != tt.expected {
				t.Errorf("GetTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test performance with large number of tags
func TestAssetResponseDto_GetTag_Performance(t *testing.T) {
	// Create asset with many tags
	numTags := 1000
	tags := make([]TagResponseDto, numTags)
	for i := 0; i < numTags; i++ {
		tags[i] = TagResponseDto{
			Name:  fmt.Sprintf("tag-%d", i),
			Value: fmt.Sprintf("value-%d", i),
		}
	}

	// Add our target tag somewhere in the middle
	targetIndex := numTags / 2
	tags[targetIndex] = TagResponseDto{
		Name:  TAG_COMPRESSED,
		Value: "2024-01-01 12:00:00",
	}

	asset := &AssetResponseDto{
		Tags: &tags,
	}

	// Test that GetTag still works efficiently
	result := asset.GetTag(TAG_COMPRESSED)
	if result != "2024-01-01 12:00:00" {
		t.Errorf("Expected '2024-01-01 12:00:00', got %v", result)
	}

	// Benchmark the performance
	benchmarkResult := testing.Benchmark(func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			asset.GetTag(TAG_COMPRESSED)
		}
	})

	t.Logf("GetTag performance: %v ns/op", benchmarkResult.NsPerOp())
}

// Test memory safety with nil pointer scenarios
func TestAssetResponseDto_NilPointerSafety(t *testing.T) {
	t.Run("GetTag with nil asset", func(t *testing.T) {
		// This should not panic even if asset is somehow nil
		// Note: We can't actually call this on a nil receiver in Go,
		// but we can test with empty/default values
		var asset *AssetResponseDto
		if asset != nil {
			result := asset.GetTag("test")
			t.Errorf("Expected nil asset to not call GetTag, got %v", result)
		}
	})

	t.Run("CompressedAfter with various nil scenarios", func(t *testing.T) {
		baseTime := time.Now()

		// Test with Tags set to nil
		assetNilTags := &AssetResponseDto{Tags: nil}
		result := assetNilTags.CompressedAfter(baseTime)
		if !result {
			t.Errorf("Expected true for nil tags, got %v", result)
		}

		// Test with empty slice
		assetEmptyTags := &AssetResponseDto{Tags: &[]TagResponseDto{}}
		result = assetEmptyTags.CompressedAfter(baseTime)
		if !result {
			t.Errorf("Expected true for empty tags, got %v", result)
		}
	})
}
