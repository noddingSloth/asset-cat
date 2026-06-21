# 📄 Binary glTF (GLB) File Format Specification

This document describes the binary file format of glTF 2.0 (GLB), which is the standard format used to package 3D assets in a single self-contained binary file.

---

## 1. File Layout Overview

A GLB file consists of a **12-byte Header**, followed immediately by one or more structured **Chunks**. 

```text
┌──────────────────────────────────────────────────────────┐
│                      12-Byte Header                      │
│   magic (4B)   │   version (4B)   │     length (4B)      │
├──────────────────────────────────────────────────────────┤
│                     Chunk 0 (JSON)                       │
│ chunkLength (4B) │ chunkType (4B) │    JSON Data (N B)   │
├──────────────────────────────────────────────────────────┤
│                  Chunk 1 (BIN) [Optional]                │
│ chunkLength (4B) │ chunkType (4B) │  Binary Data (M B)   │
└──────────────────────────────────────────────────────────┘
```

---

## 2. 12-Byte Header

The header contains metadata verifying that the file is a valid GLB asset and declaring the total file length. All fields are in **little-endian** byte order.

| Offset (Bytes) | Type | Name | Description | Value / Range |
| :--- | :--- | :--- | :--- | :--- |
| **0** | `uint32` | `magic` | Identifies the file as a glTF binary. | `0x46546C67` (ASCII string `"glTF"`) |
| **4** | `uint32` | `version` | glTF container version number. | `2` |
| **8** | `uint32` | `length` | Total length of the GLB file in bytes. | File size on disk |

---

## 3. Chunk Structure

Following the header, the remainder of the file is divided into chunks. Each chunk has an 8-byte header prefix, followed by its payload.

| Offset (Relative) | Type | Name | Description |
| :--- | :--- | :--- | :--- |
| **0** | `uint32` | `chunkLength` | The length of the chunk data payload in bytes (does not include chunk headers). |
| **4** | `uint32` | `chunkType` | Indicates the chunk format. |
| **8** | `ubyte[]` | `chunkData` | The binary payload. |

### Chunk Types

There are two standard chunk types defined for glTF 2.0:

1. **JSON Chunk (`0x4E4F534A` / `"JSON"`)**
   - Contains a UTF-8/ASCII JSON string describing the scene graph, materials, buffers, cameras, meshes, accessors, and nodes.
   - **Must be the first chunk** in the GLB file.
   - Padded with trailing space characters (`0x20`) to align the chunk length to a 4-byte boundary.

2. **Binary Buffer Chunk (`0x004E4942` / `"BIN\x00"`)**
   - Contains raw binary buffers (vertices, indices, animation keyframes, or textures).
   - **Must be the second chunk** in the GLB file (if present).
   - Padded with trailing null bytes (`0x00`) to align the chunk length to a 4-byte boundary.

---

## 4. Reading Geometry Data (glTF Hierarchy)

To extract edges and faces for wireframe rendering, the Go parser must traverse the structured relationships declared inside the JSON chunk.

### Flow of References

```text
  [Mesh] 
     │
     └──> [Primitive] (Mode 1: Lines, Mode 4: Triangles)
             │
             ├──> indices ──────> [Accessor] ──> [BufferView] ──> [BIN Chunk]
             │
             └──> attributes.POSITION ─> [Accessor] ──> [BufferView] ──> [BIN Chunk]
```

### Key JSON Properties

1. **`buffers`**:
   - The first buffer inside a GLB file does not have a `uri` property. Instead, it refers directly to the binary chunk payload (`BIN\x00`).
   ```json
   "buffers": [
     {
       "byteLength": 102450
     }
   ]
   ```

2. **`bufferViews`**:
   - Divides the buffer into logical segments (e.g. vertices segment vs. indices segment).
   ```json
   "bufferViews": [
     {
       "buffer": 0,
       "byteOffset": 0,
       "byteLength": 4800,
       "target": 34962
     }
   ]
   ```

3. **`accessors`**:
   - Typed views over `bufferViews`. Explains how to parse elements (e.g. `FLOAT` type of size `VEC3` for 3D coordinates, or `UNSIGNED_SHORT` type of size `SCALAR` for vertex indices).
   ```json
   "accessors": [
     {
       "bufferView": 0,
       "byteOffset": 0,
       "componentType": 5126,
       "count": 400,
       "type": "VEC3"
     }
   ]
   ```
   - **Important Component Types (`componentType`)**:
     - `5123` (`UNSIGNED_SHORT` / 2 bytes) - Commonly used for vertex indices.
     - `5125` (`UNSIGNED_INT` / 4 bytes) - Used for larger index buffers.
     - `5126` (`FLOAT` / 4 bytes) - Used for vertex coordinates (`VEC3`).

4. **`meshes`**:
   - Defines geometry primitives. Check the `primitives` list to find the attributes and index mappings.
   - `attributes.POSITION` holds the indices of the vertex coordinate accessor.
   - `indices` holds the index accessor mapping vertices to triangles/faces.

---

## 5. Padding and Alignment

All chunk data sections must be aligned to a **4-byte boundary** so that the CPU can perform fast multi-byte memory access of numeric elements:
- The JSON chunk payload must be padded with **space characters (`0x20`)**.
- The BIN chunk payload must be padded with **null bytes (`0x00`)**.
- Read offsets from file headers directly, as trailing pads are included in the chunk size declaration.
