package mcp

import (
	"testing"
)

func TestRequestIDCounter(t *testing.T) {
	// Test that request IDs are unique and incrementing
	ids := make(map[int64]bool)
	for i := 0; i < 100; i++ {
		req := NewRequest("test", nil)
		if req.ID != nil {
			if ids[*req.ID] {
				t.Errorf("Duplicate request ID: %d", *req.ID)
			}
			ids[*req.ID] = true
		}
	}

	if len(ids) != 100 {
		t.Errorf("Expected 100 unique IDs, got %d", len(ids))
	}
}
