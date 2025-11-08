package immich

import (
	"context"
	"testing"

	"github.com/oapi-codegen/runtime/types"
)

func TestUUUIDOfString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected string
	}{
		{
			name:     "valid UUID",
			input:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr:  false,
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:     "invalid UUID format",
			input:    "not-a-valid-uuid",
			wantErr:  true,
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			wantErr:  true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UUUIDOfString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %q, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %q, got %v", tt.input, err)
				}

				if result.String() != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result.String())
				}
			}
		})
	}
}

func TestClientSimpleStruct(t *testing.T) {
	// Test that ClientSimple struct has the expected fields
	client := &ClientSimple{}

	// Test that the struct can be created
	if client == nil {
		t.Error("ClientSimple should not be nil")
	}

	// Test struct field access (even if zero values)
	_ = client.client
	_ = client.ctx
	_ = client.parallel

	// Test that the tags field exists and can be accessed
	_ = client.tags.compressedID

	t.Log("ClientSimple struct fields are accessible")
}

func TestClientSimpleWithMockContext(t *testing.T) {
	// Test with a mock context
	ctx := context.Background()

	client := &ClientSimple{
		ctx:      ctx,
		parallel: 4,
	}

	// Test that context is set correctly
	if client.ctx != ctx {
		t.Error("Context should be set correctly")
	}

	// Test that parallel is set correctly
	if client.parallel != 4 {
		t.Error("Parallel should be set correctly")
	}
}

func TestTagConstants(t *testing.T) {
	// Test that tag constants are defined correctly
	if TAG_ROOT != "__immich-compress__" {
		t.Errorf("Expected TAG_ROOT to be '__immich-compress__', got %q", TAG_ROOT)
	}

	if TAG_COMPRESSED != "__compressed__" {
		t.Errorf("Expected TAG_COMPRESSED to be '__compressed__', got %q", TAG_COMPRESSED)
	}
}

func TestClientSimpleParallelLimit(t *testing.T) {
	// Test parallel limit setting
	tests := []struct {
		parallel int
		expected bool
	}{
		{1, true},
		{4, true},
		{16, true},
		{0, true}, // Should still work with 0
	}

	for _, tt := range tests {
		t.Run("parallel_"+string(rune('0')+rune(tt.parallel)), func(t *testing.T) {
			client := &ClientSimple{parallel: tt.parallel}

			if client.parallel != tt.parallel {
				t.Errorf("Expected parallel %d, got %d", tt.parallel, client.parallel)
			}

			// All positive values (including 0) should be valid
			if tt.parallel >= 0 {
				t.Logf("Parallel %d is valid", tt.parallel)
			}
		})
	}
}

func TestGeneratedClientTypes(t *testing.T) {
	// Test that we can work with the generated types
	var uuid types.UUID
	_ = uuid

	// Test that we can create and use types.UUID
	t.Log("Types package is importable and usable")
}

func TestClientSimpleHttpClient(t *testing.T) {
	// Test that the client has the expected structure
	client := &ClientSimple{}

	// Test that the client interface fields exist
	if client.client == nil {
		t.Log("Client interface is nil (expected for new client)")
	}

	if client.clientRaw == nil {
		t.Log("Raw client interface is nil (expected for new client)")
	}
}
