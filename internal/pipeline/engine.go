package pipeline

import (
	"fmt"
	"io"

	"github.com/noddingSloth/asset-cat/internal/canvas"
	"github.com/noddingSloth/asset-cat/internal/extractor"
	"github.com/noddingSloth/asset-cat/internal/geom"
)

// Engine orchestrates the pipeline: GLB → Projection → Canvas2D
type Engine struct {
	Model  *extractor.Model
	Camera *geom.Camera
	Canvas canvas.Canvas2D
	Width  int
	Height int
}

// NewEngineFromReader creates an Engine by extracting a model from a GLB reader.
func NewEngineFromReader(r io.Reader, canvas canvas.Canvas2D, width, height int) (*Engine, error) {
	ext := extractor.NewGLBExtractor(r)
	model, err := ext.ExtractModel()
	if err != nil {
		return nil, fmt.Errorf("extracting model: %w", err)
	}

	return &Engine{
		Model:  model,
		Camera: geom.DefaultCamera(),
		Canvas: canvas,
		Width:  width,
		Height: height,
	}, nil
}

// ViewportTransform converts normalized device coordinates (-1 to 1) to screen coordinates.
// Also flips Y so that +Y is up (NDC) vs down (screen).
func (e *Engine) ViewportTransform(ndc geom.Vector3) (float64, float64) {
	x := (ndc.X + 1.0) * 0.5 * float64(e.Width)
	y := (1.0 - ndc.Y) * 0.5 * float64(e.Height) // Flip Y
	return x, y
}

// RenderFrame clears the canvas and draws all meshes using the current camera.
func (e *Engine) RenderFrame() error {
	if e.Model == nil {
		return fmt.Errorf("no model loaded")
	}

	e.Canvas.Clear()

	aspect := float64(e.Width) / float64(e.Height)
	vpMatrix := e.Camera.ViewProjectionMatrix(aspect)

	for _, mesh := range e.Model.Meshes {
		// Project all vertices once
		projected := make([]geom.Vector3, len(mesh.Vertices))
		for i, v := range mesh.Vertices {
			// Transform vertex by combined view-projection matrix
			projected[i] = v.Transform(vpMatrix)
		}

		// Draw edges
		for _, edge := range mesh.Edges {
			p1 := projected[edge[0]]
			p2 := projected[edge[1]]

			// Skip edges where either vertex is behind the camera (w < 0 after perspective divide)
			// We check Z after projection — in NDC, visible Z is between -1 and 1
			if p1.Z < -1 || p1.Z > 1 || p2.Z < -1 || p2.Z > 1 {
				continue
			}

			x1, y1 := e.ViewportTransform(p1)
			x2, y2 := e.ViewportTransform(p2)

			e.Canvas.DrawLine(x1, y1, x2, y2)
		}
	}

	return e.Canvas.Render()
}

// RotateCamera rotates the camera around the Y axis by the given angle (radians).
func (e *Engine) RotateCamera(angle float64) {
	// Rotate camera position around target
	offset := e.Camera.Position.Sub(e.Camera.Target)
	rotMatrix := geom.RotateY(angle)
	rotatedOffset := offset.TransformDirection(rotMatrix)
	e.Camera.Position = e.Camera.Target.Add(rotatedOffset)
}

// ZoomCamera moves the camera closer to or further from the target.
func (e *Engine) ZoomCamera(delta float64) {
	direction := e.Camera.Target.Sub(e.Camera.Position).Normalize()
	e.Camera.Position = e.Camera.Position.Add(direction.Scale(delta))
}
