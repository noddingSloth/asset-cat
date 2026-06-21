package terminal_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
)

func TestBrailleCanvasInitialization(t *testing.T) {
	bc := &terminal.BrailleCanvas{}
	if bc == nil {
		t.Error("expected BrailleCanvas pointer, got nil")
	}
}
