# 5. Restrict GLB Extractor to Static Geometry Only

Date: 2026-06-21

## Status

Accepted

## Context

We are building a lightweight 3D wireframe projection engine targeting terminal displays and browser visualizers. 
The glTF 2.0 specification is rich and complex, supporting features such as:
- PBR materials, textures, samplers, and embedded image buffers.
- Skeletal animations, armatures, skins, and morph targets.
- Embedded cameras, lighting systems, and extensions.

Implementing a full glTF loader and rendering system would require importing heavyweight rendering frameworks or writing thousands of lines of parsing and state-tracking code, defeating the goal of a lightweight, self-contained terminal tool.

## Decision

We will restrict our `GLBExtractor` to parse and extract only static structural geometry:
1. **Geometric Data:** Extract only mesh vertices (the `POSITION` attribute vectors) and edge/face connectivity arrays (the index arrays).
2. **Static Pose:** Render all meshes in their default bind pose.
3. **No Textures, Shading, or Material data:** Ignore colors, PBR parameters, materials, texture coordinates (`TEXCOORD`), normals (`NORMAL`), tangents, images, and samplers.
4. **No Animation/Skins:** Ignore animation nodes, channels, skin matrices, and morph targets.

## Consequences

- **Pros:**
  - Keep the parser code small, self-contained, and easy to maintain.
  - Low CPU overhead and minimal memory allocation footprint.
  - Decoupled design; the geometry math only deals with static 3D vector coordinates.
- **Cons:**
  - Animations, dynamic meshes, and character skeletons will not play; they will render only in their default static pose.
  - Models relying on materials or textures for their shape boundaries (e.g., alpha-masked textures) will render as full bounding shapes.

