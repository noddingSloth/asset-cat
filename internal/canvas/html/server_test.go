package html_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/html"
)

func TestServerInitialization(t *testing.T) {
	srv := &html.Server{Addr: ":8080"}
	if srv.Addr != ":8080" {
		t.Errorf("expected Server address to be :8080, got %s", srv.Addr)
	}
}
