package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/noddingSloth/asset-cat/internal/canvas/html"
	"github.com/noddingSloth/asset-cat/internal/canvas/terminal"
	"github.com/noddingSloth/asset-cat/internal/extractor"
	"github.com/noddingSloth/asset-cat/internal/pipeline"
	"golang.org/x/term"
)

func main() {
	cmd := "run"
	cmdArgs := os.Args[1:]

	// Check if first argument is a subcommand
	if len(cmdArgs) > 0 && !strings.HasPrefix(cmdArgs[0], "-") {
		switch cmdArgs[0] {
		case "run", "json", "serve":
			cmd = cmdArgs[0]
			cmdArgs = cmdArgs[1:]
		case "help", "-h", "--help":
			printUsage()
			os.Exit(0)
		}
	} else {
		// No subcommand — check for top-level help flags
		for _, arg := range cmdArgs {
			if arg == "-h" || arg == "--help" || arg == "help" {
				printUsage()
				os.Exit(0)
			}
		}
	}

	// Set up args for the subcommand's flag parser
	os.Args = append([]string{cmd}, cmdArgs...)

	switch cmd {
	case "run":
		runCmd()
	case "json":
		jsonCmd()
	case "serve":
		serveCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`asset-cat - 3D wireframe renderer for your terminal

Usage:
  asset-cat [command] [--input <file.glb>]

Commands:
  run     Render model in terminal (default)
  json    Output extracted model as JSON
  serve   Start web viewer server

If --input is omitted or set to "-", input is read from stdin.

Examples:
  asset-cat --input model.glb
  asset-cat run --input model.glb
  cat model.glb | asset-cat
  asset-cat json --input model.glb > model.json
  asset-cat serve --input model.glb
`)
}

func openInput(inputPath string) (io.ReadCloser, error) {
	if inputPath == "" || inputPath == "-" {
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(inputPath)
}

// ---- run command ----

func runCmd() {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	inputPath := flags.String("input", "", "Path to .glb file (or '-' for stdin)")
	flags.Parse(os.Args[1:])

	file, err := openInput(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Must read into memory for seeking (needed on resize)
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth, termHeight = 80, 24
	}

	filename := *inputPath
	if filename == "" || filename == "-" {
		filename = "stdin"
	}

	engine := createEngineFromData(data, termWidth, termHeight)

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[?25l")
	drawStatusLine(filename, termWidth, termHeight)

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

		if frameCount%15 == 0 {
			newWidth, newHeight, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil && (newWidth != termWidth || newHeight != termHeight) {
				termWidth = newWidth
				termHeight = newHeight
				engine = createEngineFromData(data, termWidth, termHeight)
				fmt.Print("\033[2J")
				drawStatusLine(filename, termWidth, termHeight)
			}
		}

		engine.RotateCamera(0.03)
		engine.RenderFrame()
	}
}

func createEngineFromData(data []byte, termWidth, termHeight int) *pipeline.Engine {
	statusRows := 1
	renderCols := termWidth
	renderRows := termHeight - statusRows
	if renderRows < 4 {
		renderRows = 4
	}

	viewportWidth := renderCols * 2
	viewportHeight := renderRows * 4

	canvas := terminal.NewTerminalRenderer(renderCols, renderRows)

	// Wrap data in a reader for each engine creation (resize support)
	reader := &byteReader{data: data}
	engine, err := pipeline.NewEngineFromReader(reader, canvas, viewportWidth, viewportHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating engine: %v\n", err)
		os.Exit(1)
	}

	engine.AutoPositionCamera()
	return engine
}

func drawStatusLine(filename string, width, height int) {
	fmt.Printf("\033[H\033[7m asset-cat | %s | %dx%d | Ctrl+C to exit \033[0m", filename, width, height)
}

// byteReader implements io.Reader over a byte slice with a position.
type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// ---- json command ----

func jsonCmd() {
	flags := flag.NewFlagSet("json", flag.ExitOnError)
	inputPath := flags.String("input", "", "Path to .glb file (or '-' for stdin)")
	pretty := flags.Bool("pretty", true, "Pretty-print JSON output")
	flags.Parse(os.Args[1:])

	file, err := openInput(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	model, err := ext.ExtractModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting model: %v\n", err)
		os.Exit(1)
	}

	var output []byte
	if *pretty {
		output, err = json.MarshalIndent(model, "", "  ")
	} else {
		output, err = json.Marshal(model)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

// ---- serve command ----

func serveCmd() {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	inputPath := flags.String("input", "", "Path to .glb file (or '-' for stdin)")
	addr := flags.String("addr", ":8080", "Address to listen on")
	staticDir := flags.String("static", "web", "Path to static files directory")
	flags.Parse(os.Args[1:])

	file, err := openInput(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	ext := extractor.NewGLBExtractor(&byteReader{data: data})
	model, err := ext.ExtractModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting model: %v\n", err)
		os.Exit(1)
	}

	// Create engine without a canvas — we'll project manually
	engine, err := pipeline.NewEngineFromReader(
		&byteReader{data: data},
		nil, // no canvas for web mode
		800, 600,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating engine: %v\n", err)
		os.Exit(1)
	}
	engine.AutoPositionCamera()

	srv := html.NewServer(*addr, *staticDir)

	// Start animation loop in background
	go func() {
		ticker := time.NewTicker(33 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			engine.RotateCamera(0.03)
			engine.Width = 800 // will be overridden by client
			engine.Height = 600

			frame := html.Frame{
				Clear: true,
				Lines: engine.ProjectLines(),
			}
			srv.Broadcast(frame)
		}
	}()

	fmt.Printf("Starting web server on %s\n", *addr)
	fmt.Printf("Open http://localhost%s in your browser\n", *addr)
	fmt.Printf("Model: %d meshes, %d vertices total\n",
		len(model.Meshes),
		func() int {
			count := 0
			for _, m := range model.Meshes {
				count += len(m.Vertices)
			}
			return count
		}(),
	)

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
