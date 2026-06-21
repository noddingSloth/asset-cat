package pipeline_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas"
	"github.com/noddingSloth/asset-cat/internal/extractor"
	"github.com/noddingSloth/asset-cat/internal/geom"
	"github.com/noddingSloth/asset-cat/internal/pipeline"
)

// MockCanvas implements canvas.Canvas2D and records calls for assertions.
type MockCanvas struct {
	ClearCount  int
	RenderCount int
	LinesDrawn  [][4]float64 // x1, y1, x2, y2
	PixelsDrawn [][2]float64 // x, y
	RenderError error
}

func (m *MockCanvas) Clear() {
	m.ClearCount++
}

func (m *MockCanvas) DrawLine(x1, y1, x2, y2 float64) {
	m.LinesDrawn = append(m.LinesDrawn, [4]float64{x1, y1, x2, y2})
}

func (m *MockCanvas) DrawPixel(x, y float64) {
	m.PixelsDrawn = append(m.PixelsDrawn, [2]float64{x, y})
}

func (m *MockCanvas) Render() error {
	m.RenderCount++
	return m.RenderError
}

func TestEngineInitialization(t *testing.T) {
	eng := &pipeline.Engine{}
	if eng.Model != nil || eng.Canvas != nil {
		t.Error("expected uninitialized Engine components to be nil")
	}
}

func TestNewEngineFromReader_InvalidData(t *testing.T) {
	canvas := &MockCanvas{}
	reader := bytes.NewReader([]byte{})
	_, err := pipeline.NewEngineFromReader(reader, canvas, 80, 24)
	if err == nil {
		t.Error("expected error for empty reader, got nil")
	}
}

func TestNewEngineFromReader_ValidGLB(t *testing.T) {
	file, err := os.Open("../../assets/tux/Linux mascot Tux.glb")
	if err != nil {
		t.Fatalf("failed to open Tux GLB: %v", err)
	}
	defer file.Close()

	canvas := &MockCanvas{}
	engine, err := pipeline.NewEngineFromReader(file, canvas, 80, 24)
	if err != nil {
		t.Fatalf("NewEngineFromReader failed: %v", err)
	}

	if engine.Model == nil {
		t.Fatal("expected model to be extracted")
	}
	if len(engine.Model.Meshes) == 0 {
		t.Fatal("expected at least one mesh")
	}
	if engine.Camera == nil {
		t.Fatal("expected default camera")
	}
	if engine.Width != 80 || engine.Height != 24 {
		t.Errorf("expected 80x24 viewport, got %dx%d", engine.Width, engine.Height)
	}
}

func TestRenderFrame_ClearsAndRenders(t *testing.T) {
	engine := &pipeline.Engine{
		Model: &extractor.Model{
			Meshes: []extractor.Mesh{
				{
					Vertices: []geom.Vector3{
						{X: -1, Y: -1, Z: 0},
						{X: 1, Y: -1, Z: 0},
						{X: 0, Y: 1, Z: 0},
					},
					Faces: [][3]int{{0, 1, 2}},
					Edges: [][2]int{{0, 1}, {1, 2}, {2, 0}},
				},
			},
		},
		Camera: geom.DefaultCamera(),
		Canvas: &MockCanvas{},
		Width:  80,
		Height: 24,
	}

	err := engine.RenderFrame()
	if err != nil {
		t.Fatalf("RenderFrame failed: %v", err)
	}

	mock := engine.Canvas.(*MockCanvas)
	if mock.ClearCount != 1 {
		t.Errorf("expected 1 Clear call, got %d", mock.ClearCount)
	}
	if mock.RenderCount != 1 {
		t.Errorf("expected 1 Render call, got %d", mock.RenderCount)
	}
	if len(mock.LinesDrawn) != 3 {
		t.Errorf("expected 3 lines drawn (edges), got %d", len(mock.LinesDrawn))
	}
}

func TestRenderFrame_EmptyModel(t *testing.T) {
	engine := &pipeline.Engine{
		Model:  &extractor.Model{Meshes: []extractor.Mesh{}},
		Camera: geom.DefaultCamera(),
		Canvas: &MockCanvas{},
		Width:  80,
		Height: 24,
	}

	err := engine.RenderFrame()
	if err != nil {
		t.Fatalf("RenderFrame failed on empty model: %v", err)
	}

	mock := engine.Canvas.(*MockCanvas)
	if mock.ClearCount != 1 {
		t.Errorf("expected Clear call even on empty model")
	}
}

func TestRenderFrame_NilModel(t *testing.T) {
	engine := &pipeline.Engine{
		Model:  nil,
		Camera: geom.DefaultCamera(),
		Canvas: &MockCanvas{},
	}

	err := engine.RenderFrame()
	if err == nil {
		t.Error("expected error for nil model")
	}
}

func TestViewportTransform(t *testing.T) {
	engine := &pipeline.Engine{
		Width:  80,
		Height: 24,
		Scale:  1.0,
	}

	tests := []struct {
		name     string
		ndc      geom.Vector3
		expected [2]float64
	}{
		{
			name:     "origin",
			ndc:      geom.Vector3{X: 0, Y: 0, Z: 0},
			expected: [2]float64{40, 12},
		},
		{
			name:     "top-left",
			ndc:      geom.Vector3{X: -1, Y: 1, Z: 0},
			expected: [2]float64{0, 0},
		},
		{
			name:     "bottom-right",
			ndc:      geom.Vector3{X: 1, Y: -1, Z: 0},
			expected: [2]float64{80, 24},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := engine.ViewportTransform(tt.ndc)
			if x != tt.expected[0] || y != tt.expected[1] {
				t.Errorf("expected (%v, %v), got (%v, %v)", tt.expected[0], tt.expected[1], x, y)
			}
		})
	}
}

func TestRotateCamera(t *testing.T) {
	engine := &pipeline.Engine{
		Camera: geom.DefaultCamera(),
	}

	originalZ := engine.Camera.Position.Z
	engine.RotateCamera(0) // No rotation
	if engine.Camera.Position.Z != originalZ {
		t.Error("zero rotation should not change position")
	}

	// Rotate 90 degrees around Y
	engine.RotateCamera(3.14159 / 2)
	// Camera at (0, 0, 3) rotated 90° around Y should be at (~3, 0, ~0)
	if engine.Camera.Position.X < 2.9 || engine.Camera.Position.X > 3.1 {
		t.Errorf("expected X≈3 after 90° rotation, got %v", engine.Camera.Position.X)
	}
}

func TestZoomCamera(t *testing.T) {
	engine := &pipeline.Engine{
		Camera: geom.DefaultCamera(),
	}

	// Camera starts at Z=3, target at Z=0
	initialDist := engine.Camera.Position.Z // 3

	// Zoom in by 1
	engine.ZoomCamera(1)
	if engine.Camera.Position.Z != initialDist-1 {
		t.Errorf("expected Z=%v after zoom in, got %v", initialDist-1, engine.Camera.Position.Z)
	}

	// Zoom out by 2
	engine.ZoomCamera(-2)
	if engine.Camera.Position.Z != initialDist-1+2 {
		t.Errorf("expected Z=%v after zoom out, got %v", initialDist-1+2, engine.Camera.Position.Z)
	}
}

func TestRenderFrameWithTuxModel(t *testing.T) {
	file, err := os.Open("../../assets/tux/Linux mascot Tux.glb")
	if err != nil {
		t.Fatalf("failed to open Tux GLB: %v", err)
	}
	defer file.Close()

	canvas := &MockCanvas{}
	engine, err := pipeline.NewEngineFromReader(file, canvas, 80, 24)
	if err != nil {
		t.Fatalf("NewEngineFromReader failed: %v", err)
	}

	err = engine.RenderFrame()
	if err != nil {
		t.Fatalf("RenderFrame failed: %v", err)
	}

	// Should have drawn some lines (edges)
	if len(canvas.LinesDrawn) == 0 {
		t.Error("expected some lines to be drawn for Tux model")
	}
	t.Logf("Tux rendering: %d lines drawn across %d meshes", len(canvas.LinesDrawn), len(engine.Model.Meshes))

	// Verify all lines have valid screen coordinates
	for i, line := range canvas.LinesDrawn {
		x1, y1 := line[0], line[1]
		if x1 < -1000 || x1 > float64(engine.Width+1000) {
			t.Errorf("line %d: x1=%v out of reasonable range", i, x1)
		}
		if y1 < -1000 || y1 > float64(engine.Height+1000) {
			t.Errorf("line %d: y1=%v out of reasonable range", i, y1)
		}
	}

	// Rotate and render again
	engine.RotateCamera(0.5)
	canvas.LinesDrawn = nil // reset
	err = engine.RenderFrame()
	if err != nil {
		t.Fatalf("RenderFrame after rotation failed: %v", err)
	}
	if len(canvas.LinesDrawn) == 0 {
		t.Error("expected lines after rotation")
	}
	t.Logf("After rotation: %d lines drawn", len(canvas.LinesDrawn))
}

func TestCanvas2DInterfaceCompatibility(t *testing.T) {
	// Verify MockCanvas satisfies Canvas2D interface
	var c canvas.Canvas2D = &MockCanvas{}
	c.Clear()
	c.DrawLine(0, 0, 10, 10)
	c.DrawPixel(5, 5)
	err := c.Render()
	if err != nil {
		t.Errorf("unexpected Render error: %v", err)
	}
}
