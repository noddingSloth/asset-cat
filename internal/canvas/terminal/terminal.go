package terminal

// TerminalRenderer implements canvas.Canvas2D targeting stdout.
type TerminalRenderer struct {
	Width, Height int
}

func (tr *TerminalRenderer) Clear()                           {}
func (tr *TerminalRenderer) DrawLine(x1, y1, x2, y2 float64) {}
func (tr *TerminalRenderer) DrawPixel(x, y float64)           {}
func (tr *TerminalRenderer) Render() error                    { return nil }
