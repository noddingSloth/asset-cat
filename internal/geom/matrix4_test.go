package geom_test

import (
	"math"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestMatrix4Initialization(t *testing.T) {
	var m geom.Matrix4
	if m[0][0] != 0 {
		t.Errorf("expected initialized matrix to be zero-filled, got %+v", m)
	}
}

func TestIdentity(t *testing.T) {
	m := geom.Identity()
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			if col == row {
				if m[col][row] != 1 {
					t.Errorf("expected m[%d][%d]=1, got %v", col, row, m[col][row])
				}
			} else {
				if m[col][row] != 0 {
					t.Errorf("expected m[%d][%d]=0, got %v", col, row, m[col][row])
				}
			}
		}
	}
}

func TestMultiplyIdentity(t *testing.T) {
	id := geom.Identity()
	translate := geom.Translate(10, 20, 30)

	result := id.Multiply(translate)
	if result != translate {
		t.Errorf("identity * translate should equal translate:\n%+v\n%+v", result, translate)
	}

	result = translate.Multiply(id)
	if result != translate {
		t.Errorf("translate * identity should equal translate:\n%+v\n%+v", result, translate)
	}
}

func TestTranslate(t *testing.T) {
	m := geom.Translate(5, 10, 15)
	v := geom.Vector3{X: 1, Y: 1, Z: 1}
	result := v.Transform(m)
	expected := geom.Vector3{X: 6, Y: 11, Z: 16}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestRotateY(t *testing.T) {
	// Rotate 90 degrees around Y: (1,0,0) → (0,0,-1)
	rot := geom.RotateY(math.Pi / 2)
	v := geom.Vector3{X: 1, Y: 0, Z: 0}
	result := v.TransformDirection(rot)

	if math.Abs(result.X-0) > 1e-9 {
		t.Errorf("expected X=0, got %v", result.X)
	}
	if math.Abs(result.Y-0) > 1e-9 {
		t.Errorf("expected Y=0, got %v", result.Y)
	}
	if math.Abs(result.Z+1) > 1e-9 {
		t.Errorf("expected Z=-1, got %v", result.Z)
	}
}

func TestRotateYFullCircle(t *testing.T) {
	rot := geom.RotateY(2 * math.Pi)
	v := geom.Vector3{X: 1, Y: 2, Z: 3}
	result := v.TransformDirection(rot)

	if math.Abs(result.X-1) > 1e-9 || math.Abs(result.Y-2) > 1e-9 || math.Abs(result.Z-3) > 1e-9 {
		t.Errorf("full rotation should return original: %+v", result)
	}
}

func TestRotateX(t *testing.T) {
	// Rotate 90 degrees around X: (0,1,0) → (0,0,1)
	rot := geom.RotateX(math.Pi / 2)
	v := geom.Vector3{X: 0, Y: 1, Z: 0}
	result := v.TransformDirection(rot)

	if math.Abs(result.Y-0) > 1e-9 {
		t.Errorf("expected Y=0, got %v", result.Y)
	}
	if math.Abs(result.Z-1) > 1e-9 {
		t.Errorf("expected Z=1, got %v", result.Z)
	}
}

func TestRotateZ(t *testing.T) {
	// Rotate 90 degrees around Z: (1,0,0) → (0,1,0)
	rot := geom.RotateZ(math.Pi / 2)
	v := geom.Vector3{X: 1, Y: 0, Z: 0}
	result := v.TransformDirection(rot)

	if math.Abs(result.X-0) > 1e-9 {
		t.Errorf("expected X=0, got %v", result.X)
	}
	if math.Abs(result.Y-1) > 1e-9 {
		t.Errorf("expected Y=1, got %v", result.Y)
	}
}

func TestScaleMatrix(t *testing.T) {
	m := geom.ScaleMatrix(2, 3, 4)
	v := geom.Vector3{X: 1, Y: 1, Z: 1}
	result := v.TransformDirection(m)
	expected := geom.Vector3{X: 2, Y: 3, Z: 4}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestMultiplyComposition(t *testing.T) {
	// Translate then rotate: first translate, then rotate
	translate := geom.Translate(10, 0, 0)
	rotate := geom.RotateY(math.Pi / 2) // 90° around Y

	// Combined: first translate, then rotate
	combined := rotate.Multiply(translate)

	v := geom.Vector3{X: 1, Y: 0, Z: 0}
	result := v.Transform(combined)

	// Translate first: (1,0,0) → (11,0,0)
	// Then rotate around Y 90°: (11,0,0) → (0,0,-11)
	if math.Abs(result.X-0) > 1e-9 || math.Abs(result.Z+11) > 1e-9 {
		t.Errorf("expected {0, 0, -11}, got %+v", result)
	}
}

func TestPerspective(t *testing.T) {
	persp := geom.Perspective(math.Pi/2, 1.0, 0.1, 100.0)

	// A point at (1, 1, -1) in view space (looking down -Z)
	// should map to roughly (1, 1) in NDC with FOV=90
	v := geom.Vector3{X: 1, Y: 1, Z: -1}
	result := v.Transform(persp)

	if math.Abs(result.X-1) > 1e-9 || math.Abs(result.Y-1) > 1e-9 {
		t.Errorf("expected {1, 1}, got {%v, %v}", result.X, result.Y)
	}
}

func TestLookAt(t *testing.T) {
	eye := geom.Vector3{X: 0, Y: 0, Z: 5}
	center := geom.Vector3{X: 0, Y: 0, Z: 0}
	up := geom.Vector3{X: 0, Y: 1, Z: 0}

	view := geom.LookAt(eye, center, up)

	// A point at the center should be at origin in view space
	v := geom.Vector3{X: 0, Y: 0, Z: 0}
	result := v.Transform(view)

	if math.Abs(result.X-0) > 1e-9 || math.Abs(result.Y-0) > 1e-9 || math.Abs(result.Z+5) > 1e-9 {
		t.Errorf("expected {0, 0, -5}, got %+v", result)
	}
}
