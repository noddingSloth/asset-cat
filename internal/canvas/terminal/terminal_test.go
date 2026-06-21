package terminal_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
)

func TestTerminalRendererInitialization(t *testing.T) {
	tr := terminal.NewTerminalRenderer(80, 24)
	if tr.Width != 80 || tr.Height != 24 {
		t.Errorf("expected 80x24 renderer, got %dx%d", tr.Width, tr.Height)
	}
}

func TestTerminalRendererClear(t *testing.T) {
	tr := terminal.NewTerminalRenderer(10, 10)

	// Draw something, then clear
	tr.DrawLine(0, 0, 10, 10)
	tr.Clear()

	// Render should produce output without panic
	err := tr.Render()
	if err != nil {
		t.Errorf("Render returned error: %v", err)
	}
}

func TestTerminalRendererDrawLine(t *testing.T) {
	tr := terminal.NewTerminalRenderer(20, 10)

	// Draw a line across the canvas
	tr.DrawLine(0, 0, 20, 20)

	err := tr.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}

func TestTerminalRendererDrawPixel(t *testing.T) {
	tr := terminal.NewTerminalRenderer(5, 5)

	tr.DrawPixel(2, 2)
	err := tr.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}

func TestTerminalRendererMultipleLines(t *testing.T) {
	tr := terminal.NewTerminalRenderer(40, 20)

	// Draw a rectangle
	tr.DrawLine(10, 5, 30, 5)   // top
	tr.DrawLine(30, 5, 30, 15)  // right
	tr.DrawLine(30, 15, 10, 15) // bottom
	tr.DrawLine(10, 15, 10, 5)  // left

	err := tr.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
