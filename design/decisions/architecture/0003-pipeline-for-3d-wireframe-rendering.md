# 3. Pipeline for 3d wireframe rendering

Date: 2026-06-21

## Status

Accepted

## Context

We need a structured, decoupled data processing pipeline to ingest 3D assets (`.glb` files) and project them onto 2D canvases. The pipeline must handle:
- Parsing binary glTF formats (`.glb`).
- Extracting geometric data, specifically mesh vertices, edges, and faces.
- Creating a standard vector math system in Go.
- Projecting these points onto a 2D viewport.

## Decision

We will design a modular 3D wireframe rendering pipeline with the following phases:
1. **GLB Extraction Engine:** Read `.glb` files and extract mesh data. We will write our own lightweight `.glb` parsing and edge/face extraction code (to avoid depending on heavyweight frameworks like Godot).
2. **Vector3D Engine:** A custom Go math library defining vector and matrix structures (`Vector3D`, `Matrix4`, etc.) and basic transformations (translation, rotation, scaling).
3. **Projection Engine:** Custom 3D to 2D projection logic that takes `Vector3D` coordinates and transforms them into 2D camera coordinates. This will support:
   - Perspective projection (with dynamic FOV, aspect ratio, near/far clipping).
   - Orthographic projection.
4. **2D Canvas Interface:** An abstraction (`Canvas2D`) that acts as the pipeline's target, allowing different backends to handle rasterization and draw/render calls.

## Consequences

- **Pros:**
  - Standardized interface boundaries make components (parser, math, projector, canvas) easily testable and swappable.
  - Custom math library minimizes external dependencies and provides exact control over transformations.
  - Easy extension to other asset types (e.g., `.obj` or `.fbx`) by changing only the extraction engine.
- **Cons:**
  - Hand-written GLB parsing requires supporting binary layout specifications (JSON header, chunk data) manually, which can be prone to compatibility issues if assets use non-standard extensions.
  - Custom matrix mathematics must be carefully optimized to prevent memory allocations on every frame render.
