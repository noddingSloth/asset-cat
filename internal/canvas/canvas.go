package canvas

// Canvas2D defines methods required to render projected wireframe graphics.
type Canvas2D interface {
	Clear()
	DrawLine(x1, y1, x2, y2 float64)
	DrawPixel(x, y float64)
	Render() error
}
