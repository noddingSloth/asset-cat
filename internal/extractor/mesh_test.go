package extractor_test

import (
	"testing"

	"github.com/noddingSloth/asset-cat/internal/extractor"
	"github.com/noddingSloth/asset-cat/internal/geom"
)

func TestMeshStructure(t *testing.T) {
	mesh := extractor.Mesh{
		Vertices: []geom.Vector3{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}},
		Edges:    [][2]int{{0, 1}},
	}
	if len(mesh.Vertices) != 2 || len(mesh.Edges) != 1 {
		t.Errorf("unexpected mesh structure: %+v", mesh)
	}
}
