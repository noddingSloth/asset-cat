package pipeline

import (
	"github.com/noddingSloth/asset-cat/internal/canvas"
	"github.com/noddingSloth/asset-cat/internal/extractor"
)

// Engine orchestrates the pipeline: GLB -> Projection -> Canvas2D
type Engine struct {
	Extractor *extractor.GLBExtractor
	Canvas    canvas.Canvas2D
}

func (e *Engine) Run() error {
	return nil
}
