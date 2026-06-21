package extractor

import "github.com/noddingSloth/asset-cat/internal/geom"

// Mesh represents 3D model geometry extracted from a file.
type Mesh struct {
	Vertices []geom.Vector3
	Edges    [][2]int // indices of vertices forming edges
	Faces    [][3]int // indices of vertices forming faces
}
