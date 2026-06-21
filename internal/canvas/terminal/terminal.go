package terminal

import (
	"fmt"
	"os"
	"strings"
)

// TerminalRenderer implements canvas.Canvas2D targeting stdout with braille sub-pixels.
type TerminalRenderer struct {
	Width  int
	Height int
	canvas *BrailleCanvas
}

// NewTerminalRenderer creates a new terminal renderer with the given character dimensions.
func NewTerminalRenderer(width, height int) *TerminalRenderer {
	return &TerminalRenderer{
		Width:  width,
		Height: height,
		canvas: NewBrailleCanvas(width, height),
	}
}

// Clear resets the canvas.
func (tr *TerminalRenderer) Clear() {
	tr.canvas.Clear()
}

// DrawLine draws a line between two points using sub-pixel braille resolution.
func (tr *TerminalRenderer) DrawLine(x1, y1, x2, y2 float64) {
	tr.canvas.DrawLine(int(x1), int(y1), int(x2), int(y2))
}

// DrawPixel draws a single point at sub-pixel resolution.
func (tr *TerminalRenderer) DrawPixel(x, y float64) {
	tr.canvas.Set(int(x), int(y))
}

// Render outputs the braille grid to stdout starting at row 2 (below status bar).
func (tr *TerminalRenderer) Render() error {
	grid := tr.canvas.Render()
	lines := strings.Split(grid, "\n")

	var output strings.Builder
	// Position cursor at row 2, column 1 (below the status bar)
	output.WriteString("\033[2;1H")

	for i := 0; i < tr.Height; i++ {
		if i < len(lines) {
			output.WriteString(lines[i])
		}
		output.WriteString("\033[K") // clear to end of line
		if i < tr.Height-1 {
			output.WriteString("\n")
		}
	}

	fmt.Fprint(os.Stdout, output.String())
	return nil
}
