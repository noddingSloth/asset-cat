# 2. Use Go as the primary programming language

Date: 2026-06-21

## Status

Accepted

## Context

We need a programming language to implement the 3D-to-2D rendering pipeline and core math library.
The constraints and requirements are:
- Fast runtime performance for matrix manipulation and 3D vector transformations.
- Easy compilation into a single, static binary with zero external dependencies for terminal execution.
- A reliable, lightweight standard library for HTTP server and networking (needed for HTML Canvas visualizer).
- Strong typing to avoid runtime errors when performing vector geometry math.
- Low runtime footprint.

## Decision

We will use Go (Golang) as the primary backend language for the pipeline, parser, and terminal renderer.

## Consequences

- **Pros:**
  - Fast execution speed and low memory usage, which is key for real-time terminal rendering.
  - Zero-dependency compilation to run natively in various terminal environments without installation requirements.
  - Built-in, high-performance concurrency features (goroutines) if parallel processing of frames or assets becomes necessary.
  - Robust built-in HTTP server capabilities for serving the HTML canvas UI.
- **Cons:**
  - Go does not have a native interactive GUI layer, so we must rely on the web browser (HTML/JS Canvas) for rich graphical representation.
  - Standard math libraries in Go are simple, requiring us to implement custom 3D matrix and vector math tools.

