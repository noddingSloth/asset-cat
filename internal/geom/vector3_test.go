package geom_test

import (
	"math"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestVector3Creation(t *testing.T) {
	v := geom.Vector3{X: 1, Y: 2, Z: 3}
	if v.X != 1 || v.Y != 2 || v.Z != 3 {
		t.Errorf("expected Vector3{1, 2, 3}, got %+v", v)
	}
}

func TestVector3Add(t *testing.T) {
	a := geom.Vector3{X: 1, Y: 2, Z: 3}
	b := geom.Vector3{X: 4, Y: 5, Z: 6}
	result := a.Add(b)
	expected := geom.Vector3{X: 5, Y: 7, Z: 9}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestVector3Sub(t *testing.T) {
	a := geom.Vector3{X: 5, Y: 7, Z: 9}
	b := geom.Vector3{X: 4, Y: 5, Z: 6}
	result := a.Sub(b)
	expected := geom.Vector3{X: 1, Y: 2, Z: 3}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestVector3Scale(t *testing.T) {
	v := geom.Vector3{X: 1, Y: 2, Z: 3}
	result := v.Scale(2)
	expected := geom.Vector3{X: 2, Y: 4, Z: 6}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestVector3Dot(t *testing.T) {
	a := geom.Vector3{X: 1, Y: 2, Z: 3}
	b := geom.Vector3{X: 4, Y: 5, Z: 6}
	result := a.Dot(b)
	if result != 32 { // 1*4 + 2*5 + 3*6 = 4 + 10 + 18
		t.Errorf("expected 32, got %v", result)
	}
}

func TestVector3Cross(t *testing.T) {
	a := geom.Vector3{X: 1, Y: 0, Z: 0}
	b := geom.Vector3{X: 0, Y: 1, Z: 0}
	result := a.Cross(b)
	expected := geom.Vector3{X: 0, Y: 0, Z: 1} // X × Y = Z
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestVector3Length(t *testing.T) {
	v := geom.Vector3{X: 3, Y: 4, Z: 0}
	result := v.Length()
	if result != 5 {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestVector3Normalize(t *testing.T) {
	v := geom.Vector3{X: 3, Y: 4, Z: 0}
	result := v.Normalize()
	expected := geom.Vector3{X: 0.6, Y: 0.8, Z: 0}
	if math.Abs(result.X-expected.X) > 1e-9 || math.Abs(result.Y-expected.Y) > 1e-9 || result.Z != 0 {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
	if math.Abs(result.Length()-1.0) > 1e-9 {
		t.Errorf("expected unit length, got %v", result.Length())
	}
}

func TestVector3NormalizeZeroVector(t *testing.T) {
	v := geom.Vector3{X: 0, Y: 0, Z: 0}
	result := v.Normalize()
	if result != (geom.Vector3{}) {
		t.Errorf("expected zero vector, got %+v", result)
	}
}

func TestVector3Transform(t *testing.T) {
	v := geom.Vector3{X: 1, Y: 2, Z: 3}
	// Translate by (10, 20, 30)
	translate := geom.Translate(10, 20, 30)
	result := v.Transform(translate)
	expected := geom.Vector3{X: 11, Y: 22, Z: 33}
	if result != expected {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestVector3TransformDirection(t *testing.T) {
	v := geom.Vector3{X: 1, Y: 0, Z: 0}
	// Rotate 90° around Y axis: X → Z
	rotate := geom.RotateY(math.Pi / 2)
	result := v.TransformDirection(rotate)
	if math.Abs(result.X-0) > 1e-9 || math.Abs(result.Z+1) > 1e-9 {
		t.Errorf("expected {0, 0, -1}, got %+v", result)
	}
}
