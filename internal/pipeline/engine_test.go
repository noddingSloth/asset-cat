package pipeline_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/pipeline"
)

func TestEngineInitialization(t *testing.T) {
	eng := &pipeline.Engine{}
	if eng.Extractor != nil || eng.Canvas != nil {
		t.Error("expected unitialized Engine components to be nil")
	}
}
