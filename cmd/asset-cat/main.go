package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
	"github.com/noddingSloth/asset-cat/internal/pipeline"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: asset-cat <file.glb>\n")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Detect terminal size
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 80, 24
	}

	engine := createEngine(file, termWidth, termHeight)

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[?25l")
	drawStatusLine(os.Args[1], termWidth, termHeight)

	cleanup := func() {
		fmt.Print("\033[?25h\033[2J\033[H")
	}
	defer cleanup()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cleanup()
		os.Exit(0)
	}()

	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	frameCount := 0
	for range ticker.C {
		frameCount++

		// Check for resize every 15 frames
		if frameCount%15 == 0 {
			newWidth, newHeight, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil && (newWidth != termWidth || newHeight != termHeight) {
				termWidth = newWidth
				termHeight = newHeight

				// Recreate engine with new viewport size
				file.Seek(0, 0)
				newEngine := createEngine(file, termWidth, termHeight)
				newEngine.Camera = engine.Camera // preserve camera state
				newEngine.Scale = engine.Scale
				engine = newEngine

				// Redraw status line
				fmt.Print("\033[2J") // clear screen on resize
				drawStatusLine(os.Args[1], termWidth, termHeight)
			}
		}

		engine.RotateCamera(0.03)
		engine.RenderFrame()
	}
}

func createEngine(file *os.File, termWidth, termHeight int) *pipeline.Engine {
	statusRows := 1
	renderCols := termWidth
	renderRows := termHeight - statusRows
	if renderRows < 4 {
		renderRows = 4
	}

	viewportWidth := renderCols * 2
	viewportHeight := renderRows * 4

	canvas := terminal.NewTerminalRenderer(renderCols, renderRows)
	engine, err := pipeline.NewEngineFromReader(file, canvas, viewportWidth, viewportHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating engine: %v\n", err)
		os.Exit(1)
	}

	// Auto-position camera based on model's bounding box
	engine.AutoPositionCamera()

	return engine
}

func drawStatusLine(filename string, width, height int) {
	fmt.Printf("\033[H\033[7m asset-cat | %s | %dx%d | Ctrl+C to exit \033[0m", filename, width, height)
}
