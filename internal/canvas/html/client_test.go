package html_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/html"
)

func TestClientInitialization(t *testing.T) {
	c := &html.Client{ID: "test-client-1"}
	if c.ID != "test-client-1" {
		t.Errorf("expected ID test-client-1, got %s", c.ID)
	}
}
