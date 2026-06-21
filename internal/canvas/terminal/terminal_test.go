package terminal_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
)

func TestTerminalRendererInitialization(t *testing.T) {
	tr := &terminal.TerminalRenderer{Width: 80, Height: 24}
	if tr.Width != 80 || tr.Height != 24 {
		t.Errorf("expected 80x24 renderer, got %dx%d", tr.Width, tr.Height)
	}
}
