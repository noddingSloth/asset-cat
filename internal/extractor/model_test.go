package extractor_test

import (
	"encoding/json"
	"fmt"
	"os"
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

func TestModelSerializationRoundTrip(t *testing.T) {
	// Create a simple model
	original := extractor.Model{
		Meshes: []extractor.Mesh{
			{
				Vertices: []geom.Vector3{
					{X: 0, Y: 0, Z: 0},
					{X: 1, Y: 0, Z: 0},
					{X: 0, Y: 1, Z: 0},
				},
				Faces: [][3]int{{0, 1, 2}},
				Edges: [][2]int{{0, 1}, {1, 2}, {2, 0}},
			},
		},
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal model: %v", err)
	}

	// Deserialize back
	var restored extractor.Model
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal model: %v", err)
	}

	// Verify round-trip integrity
	if len(restored.Meshes) != len(original.Meshes) {
		t.Fatalf("mesh count mismatch: expected %d, got %d", len(original.Meshes), len(restored.Meshes))
	}

	for i := range original.Meshes {
		orig := original.Meshes[i]
		rest := restored.Meshes[i]

		if len(rest.Vertices) != len(orig.Vertices) {
			t.Errorf("mesh %d: vertex count mismatch: expected %d, got %d", i, len(orig.Vertices), len(rest.Vertices))
		}
		if len(rest.Faces) != len(orig.Faces) {
			t.Errorf("mesh %d: face count mismatch: expected %d, got %d", i, len(orig.Faces), len(rest.Faces))
		}
		if len(rest.Edges) != len(orig.Edges) {
			t.Errorf("mesh %d: edge count mismatch: expected %d, got %d", i, len(orig.Edges), len(rest.Edges))
		}

		for j, v := range orig.Vertices {
			if rest.Vertices[j] != v {
				t.Errorf("mesh %d vertex %d: expected %+v, got %+v", i, j, v, rest.Vertices[j])
			}
		}
		for j, f := range orig.Faces {
			if rest.Faces[j] != f {
				t.Errorf("mesh %d face %d: expected %v, got %v", i, j, f, rest.Faces[j])
			}
		}
		for j, e := range orig.Edges {
			if rest.Edges[j] != e {
				t.Errorf("mesh %d edge %d: expected %v, got %v", i, j, e, rest.Edges[j])
			}
		}
	}
}

func TestTuxModelSerialization(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	model, err := ext.ExtractModel()
	if err != nil {
		t.Fatalf("ExtractModel failed: %v", err)
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal Tux model: %v", err)
	}

	// Write the serialized model to a file for inspection and reuse
	modelPath := "../../assets/tux/tux_model.json"
	if err := os.WriteFile(modelPath, data, 0644); err != nil {
		t.Fatalf("failed to write serialized model: %v", err)
	}
	t.Logf("Serialized Tux model written to %s (%d bytes)", modelPath, len(data))

	// Deserialize and verify
	var restored extractor.Model
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal Tux model: %v", err)
	}

	// Basic integrity checks
	if len(restored.Meshes) != len(model.Meshes) {
		t.Errorf("mesh count mismatch: expected %d, got %d", len(model.Meshes), len(restored.Meshes))
	}

	totalVerts := 0
	totalFaces := 0
	totalEdges := 0
	for i, mesh := range restored.Meshes {
		if len(mesh.Vertices) != len(model.Meshes[i].Vertices) {
			t.Errorf("mesh %d: vertex count mismatch", i)
		}
		if len(mesh.Faces) != len(model.Meshes[i].Faces) {
			t.Errorf("mesh %d: face count mismatch", i)
		}
		if len(mesh.Edges) != len(model.Meshes[i].Edges) {
			t.Errorf("mesh %d: edge count mismatch", i)
		}
		totalVerts += len(mesh.Vertices)
		totalFaces += len(mesh.Faces)
		totalEdges += len(mesh.Edges)
	}

	t.Logf("Tux model: %d meshes, %d vertices, %d faces, %d edges",
		len(restored.Meshes), totalVerts, totalFaces, totalEdges)
}

func TestDeserializeTuxModelFromFile(t *testing.T) {
	// This test verifies that the previously serialized model can be loaded
	modelPath := "../../assets/tux/tux_model.json"
	data, err := os.ReadFile(modelPath)
	if err != nil {
		t.Skipf("serialized model not found at %s (run TestTuxModelSerialization first): %v", modelPath, err)
	}

	var model extractor.Model
	if err := json.Unmarshal(data, &model); err != nil {
		t.Fatalf("failed to deserialize model from file: %v", err)
	}

	if len(model.Meshes) == 0 {
		t.Fatal("expected at least one mesh in deserialized model")
	}

	for i, mesh := range model.Meshes {
		if len(mesh.Vertices) == 0 {
			t.Errorf("mesh %d: no vertices", i)
		}
		// Verify face indices are in bounds
		for j, face := range mesh.Faces {
			for k, idx := range face {
				if idx < 0 || idx >= len(mesh.Vertices) {
					t.Errorf("mesh %d face %d: index %d out of range [0, %d)", i, j, k, len(mesh.Vertices))
				}
			}
		}
		// Verify edge indices are in bounds
		for j, edge := range mesh.Edges {
			if edge[0] < 0 || edge[0] >= len(mesh.Vertices) {
				t.Errorf("mesh %d edge %d: vertex 0 index %d out of range", i, j, edge[0])
			}
			if edge[1] < 0 || edge[1] >= len(mesh.Vertices) {
				t.Errorf("mesh %d edge %d: vertex 1 index %d out of range", i, j, edge[1])
			}
		}
	}

	t.Logf("Successfully deserialized model from %s: %d meshes", modelPath, len(model.Meshes))
}

func TestInspectTuxMesh(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	model, err := ext.ExtractModel()
	if err != nil {
		t.Fatalf("ExtractModel failed: %v", err)
	}

	// Write inspection output to a file
	outFile, err := os.Create("tux_inspection.txt")
	if err != nil {
		t.Fatalf("failed to create inspection file: %v", err)
	}
	defer outFile.Close()

	fmt.Fprintf(outFile, "=== Tux Model Inspection ===\n")
	fmt.Fprintf(outFile, "Total meshes: %d\n\n", len(model.Meshes))

	for i, mesh := range model.Meshes {
		fmt.Fprintf(outFile, "--- Mesh %d ---\n", i)
		fmt.Fprintf(outFile, "Vertices: %d\n", len(mesh.Vertices))
		fmt.Fprintf(outFile, "Faces:    %d\n", len(mesh.Faces))
		fmt.Fprintf(outFile, "Edges:    %d\n", len(mesh.Edges))
		fmt.Fprintf(outFile, "\n")

		if len(mesh.Vertices) > 0 {
			minX, minY, minZ := mesh.Vertices[0].X, mesh.Vertices[0].Y, mesh.Vertices[0].Z
			maxX, maxY, maxZ := minX, minY, minZ

			fmt.Fprintf(outFile, "Vertices:\n")
			for j, v := range mesh.Vertices {
				if j < 10 || j >= len(mesh.Vertices)-5 {
					fmt.Fprintf(outFile, "  [%4d] (% .4f, % .4f, % .4f)\n", j, v.X, v.Y, v.Z)
				} else if j == 10 {
					fmt.Fprintf(outFile, "  ... (%d vertices omitted) ...\n", len(mesh.Vertices)-15)
				}

				if v.X < minX {
					minX = v.X
				}
				if v.Y < minY {
					minY = v.Y
				}
				if v.Z < minZ {
					minZ = v.Z
				}
				if v.X > maxX {
					maxX = v.X
				}
				if v.Y > maxY {
					maxY = v.Y
				}
				if v.Z > maxZ {
					maxZ = v.Z
				}
			}
			fmt.Fprintf(outFile, "\n")
			fmt.Fprintf(outFile, "Bounding Box:\n")
			fmt.Fprintf(outFile, "  Min: (% .4f, % .4f, % .4f)\n", minX, minY, minZ)
			fmt.Fprintf(outFile, "  Max: (% .4f, % .4f, % .4f)\n", maxX, maxY, maxZ)
			fmt.Fprintf(outFile, "  Size: (% .4f, % .4f, % .4f)\n", maxX-minX, maxY-minY, maxZ-minZ)
			fmt.Fprintf(outFile, "\n")
		}

		if len(mesh.Faces) > 0 {
			fmt.Fprintf(outFile, "Faces:\n")
			for j, face := range mesh.Faces {
				if j < 5 || j >= len(mesh.Faces)-3 {
					fmt.Fprintf(outFile, "  [%4d] (%4d, %4d, %4d)\n", j, face[0], face[1], face[2])
				} else if j == 5 {
					fmt.Fprintf(outFile, "  ... (%d faces omitted) ...\n", len(mesh.Faces)-8)
				}
			}
			fmt.Fprintf(outFile, "\n")
		}

		if len(mesh.Edges) > 0 {
			fmt.Fprintf(outFile, "Edges:\n")

			valence := make([]int, len(mesh.Vertices))
			for _, edge := range mesh.Edges {
				valence[edge[0]]++
				valence[edge[1]]++
			}

			for j, edge := range mesh.Edges {
				if j < 5 || j >= len(mesh.Edges)-3 {
					fmt.Fprintf(outFile, "  [%4d] (%4d, %4d)  valence: [%d, %d]\n",
						j, edge[0], edge[1], valence[edge[0]], valence[edge[1]])
				} else if j == 5 {
					fmt.Fprintf(outFile, "  ... (%d edges omitted) ...\n", len(mesh.Edges)-8)
				}
			}
			fmt.Fprintf(outFile, "\n")

			fmt.Fprintf(outFile, "Vertex Valence Statistics:\n")
			minVal, maxVal := valence[0], valence[0]
			sumVal := 0
			for _, v := range valence {
				sumVal += v
				if v < minVal {
					minVal = v
				}
				if v > maxVal {
					maxVal = v
				}
			}
			avgVal := float64(sumVal) / float64(len(valence))
			fmt.Fprintf(outFile, "  Min: %d, Max: %d, Avg: %.2f\n", minVal, maxVal, avgVal)
			fmt.Fprintf(outFile, "  Total edges: %d, Total vertices: %d\n", len(mesh.Edges), len(mesh.Vertices))

			if len(mesh.Faces) > 0 {
				euler := len(mesh.Vertices) - len(mesh.Edges) + len(mesh.Faces)
				fmt.Fprintf(outFile, "  Euler characteristic (V - E + F): %d\n", euler)
			}
		}

		fmt.Fprintf(outFile, "\n")
	}

	fmt.Fprintf(outFile, "=== Summary ===\n")
	totalVerts := 0
	totalFaces := 0
	totalEdges := 0
	for _, mesh := range model.Meshes {
		totalVerts += len(mesh.Vertices)
		totalFaces += len(mesh.Faces)
		totalEdges += len(mesh.Edges)
	}
	fmt.Fprintf(outFile, "Total vertices: %d\n", totalVerts)
	fmt.Fprintf(outFile, "Total faces:    %d\n", totalFaces)
	fmt.Fprintf(outFile, "Total edges:    %d\n", totalEdges)

	t.Logf("Inspection written to tux_inspection.txt")
}
