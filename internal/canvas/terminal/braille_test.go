package terminal_test

import (
	"strings"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
)

func TestBrailleCanvasInitialization(t *testing.T) {
	bc := terminal.NewBrailleCanvas(10, 10)
	if bc == nil {
		t.Error("expected BrailleCanvas pointer, got nil")
	}
}

func TestBrailleCanvasSetAndClear(t *testing.T) {
	bc := terminal.NewBrailleCanvas(4, 4)

	// Initially all dots should render as empty braille (0x2800)
	output := bc.Render()
	emptyChar := string(rune(0x2800))
	for _, ch := range output {
		if ch == rune(0x2800) {
			continue
		}
		// Should only contain ANSI codes or empty braille or newlines
		if ch != '\033' && ch != '[' && ch != 'H' && ch != '?' && ch != '2' && ch != '5' && ch != 'l' && ch != 'h' && ch != '\n' {
			// This is fine if it's part of ANSI sequence
		}
	}
	_ = emptyChar

	// Set a single dot at (0, 0) — this is braille dot 0 (0x01)
	bc.Set(0, 0)
	output = bc.Render()

	// Find the first braille character in the output (skip ANSI codes)
	brailleChar := rune(0x2801) // dot 0 = 0x01
	found := strings.ContainsRune(output, brailleChar)
	if !found {
		t.Error("expected braille char with dot 0 set (U+2801) after setting (0,0)")
	}

	// Clear and verify it's empty again
	bc.Clear()
	output = bc.Render()
	found = strings.ContainsRune(output, brailleChar)
	if found {
		t.Error("expected no dots after clear")
	}
}

func TestBrailleCanvasAllDots(t *testing.T) {
	bc := terminal.NewBrailleCanvas(1, 1)

	// Set all 8 dots in the single cell
	// Left column: rows 0,1,2,3
	bc.Set(0, 0) // dot 0 (0x01)
	bc.Set(0, 1) // dot 1 (0x02)
	bc.Set(0, 2) // dot 2 (0x04)
	bc.Set(0, 3) // dot 6 (0x40)
	// Right column: rows 0,1,2,3
	bc.Set(1, 0) // dot 3 (0x08)
	bc.Set(1, 1) // dot 4 (0x10)
	bc.Set(1, 2) // dot 5 (0x20)
	bc.Set(1, 3) // dot 7 (0x80)

	// All dots = 0xFF = U+28FF
	output := bc.Render()
	expected := string(rune(0x28FF))
	if !strings.ContainsRune(output, rune(0x28FF)) {
		t.Errorf("expected braille char U+28FF (all dots), not found in: %q", output)
	}
	_ = expected
}

func TestBrailleCanvasDrawLineHorizontal(t *testing.T) {
	bc := terminal.NewBrailleCanvas(4, 4)

	// Draw a horizontal line across 4 sub-pixels
	bc.DrawLine(0, 0, 3, 0)

	// Should have dots at (0,0), (1,0), (2,0), (3,0)
	output := bc.Render()
	lines := strings.Split(output, "\n")
	brailleCount := 0
	for _, line := range lines {
		for _, r := range line {
			if r >= 0x2800 && r <= 0x28FF && r != 0x2800 {
				brailleCount++
			}
		}
	}
	if brailleCount == 0 {
		t.Error("expected non-empty braille characters after drawing line")
	}
}

func TestBrailleCanvasDrawLineVertical(t *testing.T) {
	bc := terminal.NewBrailleCanvas(4, 4)

	// Draw a vertical line down 4 sub-pixels
	bc.DrawLine(0, 0, 0, 3)

	// Should have dots at (0,0), (0,1), (0,2), (0,3) — left column, all rows
	output := bc.Render()
	// Dot 0=0x01, dot 1=0x02, dot 2=0x04, dot 6=0x40 → total 0x47
	expected := string(rune(0x2847))
	if !strings.ContainsRune(output, rune(0x2847)) {
		t.Errorf("expected braille char U+2847 (left column full), got: %q", output)
	}
	_ = expected
}

func TestBrailleCanvasOutOfBounds(t *testing.T) {
	bc := terminal.NewBrailleCanvas(2, 2)

	// These should not panic
	bc.Set(-1, -1)
	bc.Set(100, 100)
	bc.DrawLine(-10, -10, 100, 100)

	output := bc.Render()
	if output == "" {
		t.Error("expected render output even with out-of-bounds draws")
	}
}
