package geom_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestVector3Creation(t *testing.T) {
	v := geom.Vector3{X: 1, Y: 2, Z: 3}
	if v.X != 1 || v.Y != 2 || v.Z != 3 {
		t.Errorf("expected Vector3{1, 2, 3}, got %+v", v)
	}
}
