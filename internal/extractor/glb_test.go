package extractor_test

import (
	"bytes"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/extractor"
)

func TestGLBExtractorInitialization(t *testing.T) {
	dummyData := []byte("dummy glb binary")
	ext := extractor.GLBExtractor{Reader: bytes.NewReader(dummyData)}
	if ext.Reader == nil {
		t.Error("expected GLBExtractor reader to be initialized")
	}
}
