package extractor

import "github.com/noddingSloth/asset-cat/internal/geom"

// Mesh represents 3D model geometry extracted from a file.
type Mesh struct {
	Vertices []geom.Vector3
	Edges    [][2]int // indices of vertices forming edges
	Faces    [][3]int // indices of vertices forming faces
}

// Model holds all meshes extracted from a single GLB file.
// This is the primary container for serialization and pipeline processing.
type Model struct {
	Meshes []Mesh `json:"meshes"`
}
