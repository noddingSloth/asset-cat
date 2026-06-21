# 4. Terminal and HTML Canvas Rendering Interfaces

Date: 2026-06-21

## Status

Accepted

## Context

The 3D-to-2D wireframe projection engine needs to support two output targets:
1. **Terminal Canvas:** Render 3D assets directly in the console for CLI usage.
2. **HTML Canvas:** Render 3D assets in a web browser for smooth visual presentation, high resolutions, and potential interactivity.

We need a flexible abstraction that handles both of these output mediums without coupling them to the core geometry projection engine.

## Decision

We will define an interface, `Canvas2D`, in Go. It will expose methods such as:
- `Clear()` to reset the viewport.
- `DrawLine(x1, y1, x2, y2 Color)` (or similar coordinates) to render lines/edges.
- `DrawPixel(x, y Color)` to set points.
- `Render()` to flush the buffer to the screen/socket.

We will implement two concrete types that satisfy this interface:
1. **TerminalRenderer:**
   - Targets stdout.
   - Uses ANSI escape sequences to clear the screen and set cursor position.
   - Employs text characters or Braille patterns (Unicode) to draw lines at sub-character resolution, simulating higher resolution graphics.
2. **HTMLCanvasRenderer:**
   - Starts a lightweight local HTTP server.
   - Uses WebSockets to stream projected wireframe coordinates (or pre-rendered SVG/Canvas draw commands) to a client browser.
   - The browser client runs a vanilla JS rendering loop that paints to an HTML5 `<canvas>` element.

## Consequences

- **Pros:**
  - High decoupled design; the projection math doesn't know (or care) whether it is rendering to stdout or a TCP socket.
  - Custom Terminal rendering provides immediate visual feedback directly inside the CLI.
  - HTML Canvas rendering allows rich aesthetics, color styling, scaling, and high frame rates that are impossible inside typical terminal environments.
- **Cons:**
  - Streaming raw draw commands over WebSockets introduces network overhead, though minimal for simple wireframes.
  - Sub-character rendering in the terminal (e.g., using Braille characters) has varying font and color support across different terminal emulators.

