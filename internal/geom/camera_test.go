package geom_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestCameraInitialization(t *testing.T) {
	cam := geom.Camera{
		Position: geom.Vector3{X: 0, Y: 0, Z: 5},
		Target:   geom.Vector3{X: 0, Y: 0, Z: 0},
		Up:       geom.Vector3{X: 0, Y: 1, Z: 0},
	}
	if cam.Position.Z != 5 {
		t.Errorf("expected Camera Position.Z to be 5, got %v", cam.Position.Z)
	}
}
