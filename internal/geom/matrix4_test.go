package geom_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestMatrix4Initialization(t *testing.T) {
	var m geom.Matrix4
	if m[0][0] != 0 {
		t.Errorf("expected initialized matrix to be zero-filled, got %+v", m)
	}
}
