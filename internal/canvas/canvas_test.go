package canvas_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas"
)

// MockCanvas implements Canvas2D for testing purposes.
type MockCanvas struct {
	linesDrawn int
}

func (m *MockCanvas) Clear()                           {}
func (m *MockCanvas) DrawLine(x1, y1, x2, y2 float64) { m.linesDrawn++ }
func (m *MockCanvas) DrawPixel(x, y float64)           {}
func (m *MockCanvas) Render() error                    { return nil }

func TestCanvas2DInterfaceUsage(t *testing.T) {
	var c canvas.Canvas2D = &MockCanvas{}
	c.DrawLine(0, 0, 10, 10)
	mock := c.(*MockCanvas)
	if mock.linesDrawn != 1 {
		t.Errorf("expected 1 line to be drawn, got %d", mock.linesDrawn)
	}
}
