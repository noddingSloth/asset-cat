# 🗺️ Project Directory Tree Map

This document visually maps out the directory layout and package structure of the `asset-cat` project to help contributors navigate the codebase.

---

## 🌳 Visual Directory Tree

The following diagram illustrates the folder structure, Go files, tests, and web client assets:

```text
.
├── assets/
│   └── tux/
│       ├── source.txt
│       └── tux.glb
├── cmd/
│   └── asset-cat/
│       └── main.go
├── design/
│   ├── decisions/
│   │   └── architecture/
│   │       ├── 0001-record-architecture-decisions.md
│   │       ├── 0002-use-go-as-the-primary-programming-language.md
│   │       ├── 0003-pipeline-for-3d-wireframe-rendering.md
│   │       ├── 0004-terminal-and-html-canvas-rendering-interfaces.md
│   │       └── 0005-restrict-glb-extractor-to-static-geometry-only.md
│   └── documents/
│       └── project_tree.md
├── internal/
│   ├── canvas/
│   │   ├── html/
│   │   │   ├── client.go
│   │   │   ├── client_test.go
│   │   │   ├── server.go
│   │   │   └── server_test.go
│   │   ├── terminal/
│   │   │   ├── braille.go
│   │   │   ├── braille_test.go
│   │   │   ├── terminal.go
│   │   │   └── terminal_test.go
│   │   ├── canvas.go
│   │   └── canvas_test.go
│   ├── extractor/
│   │   ├── glb.go
│   │   ├── glb_test.go
│   │   ├── mesh.go
│   │   └── mesh_test.go
│   ├── geom/
│   │   ├── camera.go
│   │   ├── camera_test.go
│   │   ├── matrix4.go
│   │   ├── matrix4_test.go
│   │   ├── vector3.go
│   │   └── vector3_test.go
│   └── pipeline/
│       ├── engine.go
│       └── engine_test.go
├── web/
│   ├── app.js
│   ├── index.html
│   └── style.css
├── go.mod
├── LICENSE
└── README.md
```

---

## 📋 Directory Definitions & Responsibilities

### 1. `cmd/`
Contains the entrypoints for the application binaries. Subfolder `asset-cat/` builds the primary CLI executable.

### 2. `internal/`
Holds code private to the `asset-cat` project to enforce modular encapsulation.
- **`geom/`**: Vector math (3D coordinates, translation/rotation matrices, and camera viewing transforms).
- **`extractor/`**: Geometry extraction parser for extracting vertices, faces, and edges from `.glb` files.
- **`canvas/`**: Output rendering abstractions (`Canvas2D`).
  - **`canvas/terminal/`**: Translates lines to Unicode Braille characters and ANSI codes.
  - **`canvas/html/`**: Manages a WebSocket server streaming coordinates to browser clients.
- **`pipeline/`**: Orchestration logic running the extraction -> projection -> canvas cycle.

### 3. `web/`
Contains the visualizer UI code (HTML, CSS, JS) served by the Go server. Displays a canvas and controls to preview the wireframes at high framerates.

### 4. `assets/`
Storage for raw testing assets, starting with the `tux/` penguin model source and asset details.

### 5. `design/`
Architecture planning documents:
- **`decisions/architecture/`**: Active ADR files documenting tech stack choices (Go, custom pipeline, rendering outputs).
- **`documents/`**: Auxiliary architectural mapping documents (like this tree layout guide).
