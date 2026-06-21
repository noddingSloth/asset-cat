package extractor_test

import (
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

func TestInspectTuxMesh(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	meshes, err := ext.ExtractMeshes()
	if err != nil {
		t.Fatalf("ExtractMeshes failed: %v", err)
	}

	// Write inspection output to a file
	outFile, err := os.Create("tux_inspection.txt")
	if err != nil {
		t.Fatalf("failed to create inspection file: %v", err)
	}
	defer outFile.Close()

	fmt.Fprintf(outFile, "=== Tux Model Inspection ===\n")
	fmt.Fprintf(outFile, "Total meshes: %d\n\n", len(meshes))

	for i, mesh := range meshes {
		fmt.Fprintf(outFile, "--- Mesh %d ---\n", i)
		fmt.Fprintf(outFile, "Vertices: %d\n", len(mesh.Vertices))
		fmt.Fprintf(outFile, "Faces:    %d\n", len(mesh.Faces))
		fmt.Fprintf(outFile, "Edges:    %d\n", len(mesh.Edges))
		fmt.Fprintf(outFile, "\n")

		// Print vertices with bounds
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

		// Print faces (with vertex count check)
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

		// Print edges (with adjacency info)
		if len(mesh.Edges) > 0 {
			fmt.Fprintf(outFile, "Edges:\n")

			// Count vertex valence (how many edges each vertex belongs to)
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

			// Print vertex valence statistics
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

			// Euler characteristic check (V - E + F should be ≈ 2 for closed manifold)
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
	for _, mesh := range meshes {
		totalVerts += len(mesh.Vertices)
		totalFaces += len(mesh.Faces)
		totalEdges += len(mesh.Edges)
	}
	fmt.Fprintf(outFile, "Total vertices: %d\n", totalVerts)
	fmt.Fprintf(outFile, "Total faces:    %d\n", totalFaces)
	fmt.Fprintf(outFile, "Total edges:    %d\n", totalEdges)

	t.Logf("Inspection written to tux_inspection.txt")
}
