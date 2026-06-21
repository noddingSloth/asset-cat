package extractor

import "github.com/noddingSloth/asset-cat/internal/geom"

// Mesh represents 3D model geometry extracted from a file.
type Mesh struct {
	Vertices []geom.Vector3 `json:"vertices"`
	Edges    [][2]int       `json:"edges"`
	Faces    [][3]int       `json:"faces"`
}

// BoundingBox represents an axis-aligned bounding box.
type BoundingBox struct {
	Min geom.Vector3 `json:"min"`
	Max geom.Vector3 `json:"max"`
}

// Center returns the center point of the bounding box.
func (bb BoundingBox) Center() geom.Vector3 {
	return geom.Vector3{
		X: (bb.Min.X + bb.Max.X) / 2,
		Y: (bb.Min.Y + bb.Max.Y) / 2,
		Z: (bb.Min.Z + bb.Max.Z) / 2,
	}
}

// Size returns the dimensions (width, height, depth) of the bounding box.
func (bb BoundingBox) Size() geom.Vector3 {
	return geom.Vector3{
		X: bb.Max.X - bb.Min.X,
		Y: bb.Max.Y - bb.Min.Y,
		Z: bb.Max.Z - bb.Min.Z,
	}
}

// MaxDimension returns the largest dimension of the bounding box.
func (bb BoundingBox) MaxDimension() float64 {
	s := bb.Size()
	max := s.X
	if s.Y > max {
		max = s.Y
	}
	if s.Z > max {
		max = s.Z
	}
	return max
}

// Model holds all meshes extracted from a single GLB file.
type Model struct {
	Meshes      []Mesh      `json:"meshes"`
	BoundingBox BoundingBox `json:"boundingBox"`
}

// ComputeBoundingBox calculates the bounding box across all meshes in the model.
func (m *Model) ComputeBoundingBox() {
	if len(m.Meshes) == 0 || len(m.Meshes[0].Vertices) == 0 {
		return
	}

	// Initialize with first vertex
	first := m.Meshes[0].Vertices[0]
	bb := BoundingBox{
		Min: first,
		Max: first,
	}

	for _, mesh := range m.Meshes {
		for _, v := range mesh.Vertices {
			if v.X < bb.Min.X {
				bb.Min.X = v.X
			}
			if v.Y < bb.Min.Y {
				bb.Min.Y = v.Y
			}
			if v.Z < bb.Min.Z {
				bb.Min.Z = v.Z
			}
			if v.X > bb.Max.X {
				bb.Max.X = v.X
			}
			if v.Y > bb.Max.Y {
				bb.Max.Y = v.Y
			}
			if v.Z > bb.Max.Z {
				bb.Max.Z = v.Z
			}
		}
	}

	m.BoundingBox = bb
}
