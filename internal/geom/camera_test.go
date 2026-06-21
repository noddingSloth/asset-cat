package geom_test

import (
	"math"
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

func TestDefaultCamera(t *testing.T) {
	cam := geom.DefaultCamera()

	if cam.Position.Z != 3 {
		t.Errorf("expected Z=3, got %v", cam.Position.Z)
	}
	if cam.FOV != 60 {
		t.Errorf("expected FOV=60, got %v", cam.FOV)
	}
	if cam.Near != 0.1 {
		t.Errorf("expected Near=0.1, got %v", cam.Near)
	}
	if cam.Far != 100 {
		t.Errorf("expected Far=100, got %v", cam.Far)
	}
}

func TestCameraViewMatrix(t *testing.T) {
	cam := geom.DefaultCamera()
	view := cam.ViewMatrix()

	// The camera target (0,0,0) should be at distance 3 in front of camera
	v := geom.Vector3{X: 0, Y: 0, Z: 0}
	result := v.Transform(view)

	// Looking from (0,0,3) at (0,0,0), the target should be at (0,0,-3) in view space
	if math.Abs(result.Z+3) > 1e-9 {
		t.Errorf("expected Z=-3 in view space, got Z=%v", result.Z)
	}
}

func TestCameraProjectionMatrix(t *testing.T) {
	cam := geom.DefaultCamera()
	proj := cam.ProjectionMatrix(1.0) // square aspect

	if proj[0][0] == 0 {
		t.Error("projection matrix should not be zero")
	}
}

func TestCameraViewProjectionMatrix(t *testing.T) {
	cam := geom.DefaultCamera()
	vp := cam.ViewProjectionMatrix(1.0)

	// A point at origin should be visible through the camera
	v := geom.Vector3{X: 0, Y: 0, Z: 0}
	result := v.Transform(vp)

	// Should be in NDC range roughly
	if math.IsNaN(result.X) || math.IsNaN(result.Y) || math.IsNaN(result.Z) {
		t.Error("projected point should not be NaN")
	}
}
