package extractor_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/extractor"
)

const tuxGLBPath = "../../assets/tux/Linux mascot Tux.glb"

func openTuxAsset(t *testing.T) *os.File {
	t.Helper()
	file, err := os.Open(tuxGLBPath)
	if err != nil {
		t.Fatalf("failed to open Tux GLB asset: %v", err)
	}
	return file
}

func TestReadHeaderFromTuxAsset(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	header, err := ext.ReadHeader()
	if err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}

	const expectedMagic = 0x46546C67
	const expectedVersion = 2

	if header.Magic != expectedMagic {
		t.Errorf("expected magic 0x%X, got 0x%X", expectedMagic, header.Magic)
	}

	if header.Version != expectedVersion {
		t.Errorf("expected version %d, got %d", expectedVersion, header.Version)
	}

	info, err := file.Stat()
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if header.Length != uint32(info.Size()) {
		t.Errorf("expected length %d, got %d", info.Size(), header.Length)
	}
}

func TestReadChunksFromTuxAsset(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	header, err := ext.ReadHeader()
	if err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}

	chunks, err := ext.ReadChunks()
	if err != nil {
		t.Fatalf("ReadChunks failed: %v", err)
	}

	// Typical GLB file has exactly 2 chunks: JSON and BIN
	if len(chunks) < 1 {
		t.Fatalf("expected at least 1 chunk, got %d", len(chunks))
	}

	// Chunk 0 must be JSON
	const typeJSON = 0x4E4F534A
	if chunks[0].Type != typeJSON {
		t.Errorf("expected chunk 0 to be JSON (0x%X), got 0x%X", typeJSON, chunks[0].Type)
	}

	// If Chunk 1 exists, it must be BIN
	const typeBIN = 0x004E4942
	if len(chunks) > 1 {
		if chunks[1].Type != typeBIN {
			t.Errorf("expected chunk 1 to be BIN (0x%X), got 0x%X", typeBIN, chunks[1].Type)
		}
	}

	// Validate chunk sizing against total length
	// total = Header (12) + Sum of Chunks (8-byte header per chunk + data length)
	calculatedLength := uint32(12)
	for _, chunk := range chunks {
		calculatedLength += 8 + chunk.Length
	}

	if calculatedLength != header.Length {
		t.Errorf("expected calculated length %d to match header length %d", calculatedLength, header.Length)
	}
}

func TestExtractMeshesFromTuxAsset(t *testing.T) {
	file := openTuxAsset(t)
	defer file.Close()

	ext := extractor.NewGLBExtractor(file)
	meshes, err := ext.ExtractMeshes()
	if err != nil {
		t.Fatalf("ExtractMeshes failed: %v", err)
	}

	if len(meshes) == 0 {
		t.Fatal("expected at least one mesh, got none")
	}

	for i, mesh := range meshes {
		t.Run(fmt.Sprintf("mesh_%d", i), func(t *testing.T) {
			// Every mesh must have vertices
			if len(mesh.Vertices) == 0 {
				t.Errorf("mesh %d: expected vertices, got none", i)
			}

			// Vertices should be valid (not NaN, not Inf)
			for j, v := range mesh.Vertices {
				if math.IsNaN(v.X) || math.IsInf(v.X, 0) {
					t.Errorf("mesh %d vertex %d: X is NaN or Inf", i, j)
				}
				if math.IsNaN(v.Y) || math.IsInf(v.Y, 0) {
					t.Errorf("mesh %d vertex %d: Y is NaN or Inf", i, j)
				}
				if math.IsNaN(v.Z) || math.IsInf(v.Z, 0) {
					t.Errorf("mesh %d vertex %d: Z is NaN or Inf", i, j)
				}
			}

			// If faces exist, they must reference valid vertex indices
			for j, face := range mesh.Faces {
				for k, idx := range face {
					if idx < 0 || idx >= len(mesh.Vertices) {
						t.Errorf("mesh %d face %d: index %d out of range [0, %d)", i, j, k, len(mesh.Vertices))
					}
				}
			}

			// If edges exist, they must reference valid vertex indices
			for j, edge := range mesh.Edges {
				if edge[0] < 0 || edge[0] >= len(mesh.Vertices) {
					t.Errorf("mesh %d edge %d: vertex 0 index %d out of range [0, %d)", i, j, edge[0], len(mesh.Vertices))
				}
				if edge[1] < 0 || edge[1] >= len(mesh.Vertices) {
					t.Errorf("mesh %d edge %d: vertex 1 index %d out of range [0, %d)", i, j, edge[1], len(mesh.Vertices))
				}
				if edge[0] == edge[1] {
					t.Errorf("mesh %d edge %d: degenerate edge (same vertex %d)", i, j, edge[0])
				}
			}

			// Edges should be unique (no duplicates)
			edgeSet := make(map[[2]int]struct{}, len(mesh.Edges))
			for _, edge := range mesh.Edges {
				key := edge
				if key[0] > key[1] {
					key = [2]int{key[1], key[0]}
				}
				if _, exists := edgeSet[key]; exists {
					t.Errorf("mesh %d: duplicate edge [%d, %d]", i, edge[0], edge[1])
				}
				edgeSet[key] = struct{}{}
			}

			// Faces should be triangular
			for j, face := range mesh.Faces {
				if len(face) != 3 {
					t.Errorf("mesh %d face %d: expected 3 indices, got %d", i, j, len(face))
				}
			}

			// Edge-face consistency: every edge should belong to at least one face
			// (if faces exist)
			if len(mesh.Faces) > 0 {
				faceEdges := make(map[[2]int]bool)
				for _, face := range mesh.Faces {
					a, b, c := face[0], face[1], face[2]
					addFaceEdge(faceEdges, a, b)
					addFaceEdge(faceEdges, b, c)
					addFaceEdge(faceEdges, c, a)
				}
				for _, edge := range mesh.Edges {
					key := edge
					if key[0] > key[1] {
						key = [2]int{key[1], key[0]}
					}
					if !faceEdges[key] {
						t.Errorf("mesh %d: edge [%d, %d] not found in any face", i, edge[0], edge[1])
					}
				}
			}
		})
	}
}

func TestExtractMeshesInvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "empty data",
			data:        []byte{},
			expectError: true,
		},
		{
			name:        "invalid header magic",
			data:        make([]byte, 12), // all zeros
			expectError: true,
		},
		{
			name: "valid header but no JSON chunk",
			data: func() []byte {
				buf := make([]byte, 20)
				// Valid header
				binary.LittleEndian.PutUint32(buf[0:], 0x46546C67) // magic
				binary.LittleEndian.PutUint32(buf[4:], 2)          // version
				binary.LittleEndian.PutUint32(buf[8:], 20)         // length
				// Invalid chunk (not JSON)
				binary.LittleEndian.PutUint32(buf[12:], 0)          // chunk length
				binary.LittleEndian.PutUint32(buf[16:], 0x00000000) // chunk type
				return buf
			}(),
			expectError: true,
		},
		{
			name: "JSON chunk with no meshes",
			data: func() []byte {
				jsonData := []byte(`{"asset":{"version":"2.0"}}`)
				// Pad to 4-byte alignment
				for len(jsonData)%4 != 0 {
					jsonData = append(jsonData, 0x20) // space padding
				}

				totalLen := 12 + 8 + len(jsonData)
				buf := make([]byte, totalLen)
				binary.LittleEndian.PutUint32(buf[0:], 0x46546C67)             // magic
				binary.LittleEndian.PutUint32(buf[4:], 2)                      // version
				binary.LittleEndian.PutUint32(buf[8:], uint32(totalLen))       // length
				binary.LittleEndian.PutUint32(buf[12:], uint32(len(jsonData))) // chunk length
				binary.LittleEndian.PutUint32(buf[16:], 0x4E4F534A)            // JSON chunk type
				copy(buf[20:], jsonData)
				return buf
			}(),
			expectError: false, // No meshes is valid, just returns empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			ext := extractor.NewGLBExtractor(reader)
			meshes, err := ext.ExtractMeshes()

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(meshes) != 0 {
					t.Errorf("expected 0 meshes, got %d", len(meshes))
				}
			}
		})
	}
}
func TestReadVec3Accessor(t *testing.T) {
	// This tests the internal function through ExtractMeshes
	// but we can also test edge cases through malformed GLB data
	t.Run("missing buffer view", func(t *testing.T) {
		jsonData := []byte(`{
			"meshes": [{
				"primitives": [{
					"attributes": {"POSITION": 0}
				}]
			}],
			"accessors": [{
				"componentType": 5126,
				"type": "VEC3",
				"count": 3
			}]
		}`)

		buf := createGLBWithJSON(jsonData)
		reader := bytes.NewReader(buf)
		ext := extractor.NewGLBExtractor(reader)
		_, err := ext.ExtractMeshes()
		if err == nil {
			t.Error("expected error for missing buffer view, got nil")
		}
	})

	t.Run("wrong accessor type", func(t *testing.T) {
		jsonData := []byte(`{
			"meshes": [{
				"primitives": [{
					"attributes": {"POSITION": 0}
				}]
			}],
			"accessors": [{
				"bufferView": 0,
				"componentType": 5126,
				"type": "VEC4",
				"count": 3
			}],
			"bufferViews": [{
				"buffer": 0,
				"byteLength": 36
			}],
			"buffers": [{"byteLength": 36}]
		}`)

		buf := createGLBWithJSONAndBin(jsonData, make([]byte, 36))
		reader := bytes.NewReader(buf)
		ext := extractor.NewGLBExtractor(reader)
		_, err := ext.ExtractMeshes()
		if err == nil {
			t.Error("expected error for VEC4 accessor, got nil")
		}
	})
}

// Helper functions for test data creation

func addFaceEdge(edges map[[2]int]bool, a, b int) {
	if a < b {
		edges[[2]int{a, b}] = true
	} else {
		edges[[2]int{b, a}] = true
	}
}

func createGLBWithJSON(jsonData []byte) []byte {
	// Pad to 4-byte alignment
	for len(jsonData)%4 != 0 {
		jsonData = append(jsonData, 0x20)
	}

	totalLen := 12 + 8 + len(jsonData)
	buf := make([]byte, totalLen)
	binary.LittleEndian.PutUint32(buf[0:], 0x46546C67) // magic
	binary.LittleEndian.PutUint32(buf[4:], 2)          // version
	binary.LittleEndian.PutUint32(buf[8:], uint32(totalLen))
	binary.LittleEndian.PutUint32(buf[12:], uint32(len(jsonData)))
	binary.LittleEndian.PutUint32(buf[16:], 0x4E4F534A)
	copy(buf[20:], jsonData)
	return buf
}

func createGLBWithJSONAndBin(jsonData, binData []byte) []byte {
	// Pad JSON to 4-byte alignment
	for len(jsonData)%4 != 0 {
		jsonData = append(jsonData, 0x20)
	}

	totalLen := 12 + 8 + len(jsonData) + 8 + len(binData)
	buf := make([]byte, totalLen)

	// Header
	binary.LittleEndian.PutUint32(buf[0:], 0x46546C67)
	binary.LittleEndian.PutUint32(buf[4:], 2)
	binary.LittleEndian.PutUint32(buf[8:], uint32(totalLen))

	// JSON chunk
	offset := 12
	binary.LittleEndian.PutUint32(buf[offset:], uint32(len(jsonData)))
	binary.LittleEndian.PutUint32(buf[offset+4:], 0x4E4F534A)
	copy(buf[offset+8:], jsonData)

	// BIN chunk
	offset += 8 + len(jsonData)
	binary.LittleEndian.PutUint32(buf[offset:], uint32(len(binData)))
	binary.LittleEndian.PutUint32(buf[offset+4:], 0x004E4942)
	copy(buf[offset+8:], binData)

	return buf
}
