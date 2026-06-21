package extractor

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/noddingSloth/asset-cat/internal/geom"
)

// Header represents the 12-byte GLB file header.
type Header struct {
	Magic   uint32 // Must be 0x46546C67 ("glTF")
	Version uint32 // Must be 2
	Length  uint32 // Total file length in bytes
}

// Chunk represents a structured data block in the GLB container.
type Chunk struct {
	Length uint32
	Type   uint32
	Data   []byte
}

// glTF JSON types (only what we need per ADR-5)
type glTF struct {
	Meshes      []glTFMesh       `json:"meshes"`
	Accessors   []glTFAccessor   `json:"accessors"`
	BufferViews []glTFBufferView `json:"bufferViews"`
	Buffers     []glTFBuffer     `json:"buffers"`
}

type glTFMesh struct {
	Primitives []glTFPrimitive `json:"primitives"`
}

type glTFPrimitive struct {
	Attributes map[string]int `json:"attributes"` // "POSITION" → accessor index
	Indices    *int           `json:"indices,omitempty"`
}

type glTFAccessor struct {
	BufferView    *int   `json:"bufferView,omitempty"`
	ComponentType int    `json:"componentType"`
	Count         int    `json:"count"`
	Type          string `json:"type"` // "VEC3", "SCALAR"
	ByteOffset    int    `json:"byteOffset,omitempty"`
}

type glTFBufferView struct {
	Buffer     int `json:"buffer"`
	ByteOffset int `json:"byteOffset,omitempty"`
	ByteLength int `json:"byteLength"`
}

type glTFBuffer struct {
	ByteLength int    `json:"byteLength"`
	URI        string `json:"uri,omitempty"`
}

// GLBExtractor handles reading and parsing GLB files.
type GLBExtractor struct {
	Reader io.Reader
}

// NewGLBExtractor creates a new extractor from a reader.
func NewGLBExtractor(r io.Reader) *GLBExtractor {
	return &GLBExtractor{Reader: r}
}

// ReadHeader parses the 12-byte header of the GLB container.
func (e *GLBExtractor) ReadHeader() (Header, error) {
	var header Header
	err := binary.Read(e.Reader, binary.LittleEndian, &header)
	if err != nil {
		return Header{}, fmt.Errorf("failed to read GLB header: %w", err)
	}

	// Validate magic bytes
	const expectedMagic = 0x46546C67
	if header.Magic != expectedMagic {
		return Header{}, fmt.Errorf("invalid GLB magic: expected 0x%X, got 0x%X", expectedMagic, header.Magic)
	}

	// Validate version
	const expectedVersion = 2
	if header.Version != expectedVersion {
		return Header{}, fmt.Errorf("unsupported GLB version: expected %d, got %d", expectedVersion, header.Version)
	}

	return header, nil
}

// ReadChunks parses all structured chunks remaining in the GLB reader.
func (e *GLBExtractor) ReadChunks() ([]Chunk, error) {
	var chunks []Chunk
	for {
		var chunkLength uint32
		err := binary.Read(e.Reader, binary.LittleEndian, &chunkLength)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk length: %w", err)
		}

		var chunkType uint32
		err = binary.Read(e.Reader, binary.LittleEndian, &chunkType)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk type: %w", err)
		}

		chunkData := make([]byte, chunkLength)
		_, err = io.ReadFull(e.Reader, chunkData)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk data of size %d: %w", chunkLength, err)
		}

		chunks = append(chunks, Chunk{
			Length: chunkLength,
			Type:   chunkType,
			Data:   chunkData,
		})
	}
	return chunks, nil
}

// ExtractMeshes reads and parses the GLB file, returning all static mesh geometry.
func (e *GLBExtractor) ExtractMeshes() ([]Mesh, error) {
	// 1. Read header (advances reader past 12 bytes)
	_, err := e.ReadHeader()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	// 2. Read chunks
	chunks, err := e.ReadChunks()
	if err != nil {
		return nil, fmt.Errorf("reading chunks: %w", err)
	}

	// 3. Find JSON and BIN chunks
	var jsonData, binData []byte
	for _, chunk := range chunks {
		switch chunk.Type {
		case 0x4E4F534A: // JSON
			jsonData = chunk.Data
		case 0x004E4942: // BIN
			binData = chunk.Data
		}
	}
	if jsonData == nil {
		return nil, fmt.Errorf("no JSON chunk found")
	}

	// 4. Parse glTF structure
	var gltf glTF
	if err := json.Unmarshal(jsonData, &gltf); err != nil {
		return nil, fmt.Errorf("parsing glTF JSON: %w", err)
	}

	// 5. Extract meshes
	var meshes []Mesh
	for _, m := range gltf.Meshes {
		for _, prim := range m.Primitives {
			mesh, err := extractPrimitive(prim, gltf, binData)
			if err != nil {
				return nil, fmt.Errorf("extracting primitive: %w", err)
			}
			meshes = append(meshes, *mesh)
		}
	}

	return meshes, nil
}

func extractPrimitive(prim glTFPrimitive, gltf glTF, binData []byte) (*Mesh, error) {
	// Get POSITION accessor
	posAccessorIdx, ok := prim.Attributes["POSITION"]
	if !ok {
		return nil, fmt.Errorf("primitive missing POSITION attribute")
	}

	// Read vertices
	vertices, err := readVec3Accessor(gltf.Accessors[posAccessorIdx], gltf.BufferViews, binData)
	if err != nil {
		return nil, fmt.Errorf("reading positions: %w", err)
	}

	mesh := &Mesh{
		Vertices: vertices,
	}

	// Read indices (optional — if absent, vertices form sequential triangles)
	if prim.Indices != nil {
		indices, err := readScalarAccessor(gltf.Accessors[*prim.Indices], gltf.BufferViews, binData)
		if err != nil {
			return nil, fmt.Errorf("reading indices: %w", err)
		}
		mesh.Faces = buildFaces(indices)
		mesh.Edges = buildEdges(indices)
	} else {
		// Non-indexed geometry: vertices are sequential triangles
		mesh.Faces = buildSequentialFaces(len(vertices))
		mesh.Edges = buildSequentialEdges(len(vertices))
	}

	return mesh, nil
}

func readVec3Accessor(acc glTFAccessor, views []glTFBufferView, bin []byte) ([]geom.Vector3, error) {
	if acc.Type != "VEC3" || acc.ComponentType != 5126 { // 5126 = FLOAT
		return nil, fmt.Errorf("expected VEC3 FLOAT accessor, got %s type %d", acc.Type, acc.ComponentType)
	}
	if acc.BufferView == nil {
		return nil, fmt.Errorf("accessor missing bufferView")
	}

	view := views[*acc.BufferView]
	offset := view.ByteOffset + acc.ByteOffset
	data := bin[offset : offset+view.ByteLength]

	vertices := make([]geom.Vector3, acc.Count)
	for i := 0; i < acc.Count; i++ {
		vertices[i].X = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[i*12:])))
		vertices[i].Y = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[i*12+4:])))
		vertices[i].Z = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[i*12+8:])))
	}
	return vertices, nil
}

func readScalarAccessor(acc glTFAccessor, views []glTFBufferView, bin []byte) ([]uint32, error) {
	if acc.BufferView == nil {
		return nil, fmt.Errorf("accessor missing bufferView")
	}

	view := views[*acc.BufferView]
	offset := view.ByteOffset + acc.ByteOffset
	data := bin[offset : offset+view.ByteLength]

	indices := make([]uint32, acc.Count)
	switch acc.ComponentType {
	case 5123: // UNSIGNED_SHORT
		for i := range indices {
			indices[i] = uint32(binary.LittleEndian.Uint16(data[i*2:]))
		}
	case 5125: // UNSIGNED_INT
		for i := range indices {
			indices[i] = binary.LittleEndian.Uint32(data[i*4:])
		}
	default:
		return nil, fmt.Errorf("unsupported index component type: %d", acc.ComponentType)
	}
	return indices, nil
}

func buildFaces(indices []uint32) [][3]int {
	faces := make([][3]int, len(indices)/3)
	for i := 0; i < len(faces); i++ {
		faces[i] = [3]int{
			int(indices[i*3]),
			int(indices[i*3+1]),
			int(indices[i*3+2]),
		}
	}
	return faces
}

func buildEdges(indices []uint32) [][2]int {
	edgeSet := make(map[[2]int]struct{})
	addEdge := func(a, b uint32) {
		if a < b {
			edgeSet[[2]int{int(a), int(b)}] = struct{}{}
		} else {
			edgeSet[[2]int{int(b), int(a)}] = struct{}{}
		}
	}

	for i := 0; i < len(indices); i += 3 {
		a, b, c := indices[i], indices[i+1], indices[i+2]
		addEdge(a, b)
		addEdge(b, c)
		addEdge(c, a)
	}

	edges := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		edges = append(edges, e)
	}
	return edges
}

func buildSequentialFaces(vertexCount int) [][3]int {
	faceCount := vertexCount / 3
	faces := make([][3]int, faceCount)
	for i := 0; i < faceCount; i++ {
		faces[i] = [3]int{i * 3, i*3 + 1, i*3 + 2}
	}
	return faces
}

func buildSequentialEdges(vertexCount int) [][2]int {
	edgeSet := make(map[[2]int]struct{})
	addEdge := func(a, b int) {
		if a < b {
			edgeSet[[2]int{a, b}] = struct{}{}
		} else {
			edgeSet[[2]int{b, a}] = struct{}{}
		}
	}

	faceCount := vertexCount / 3
	for i := 0; i < faceCount; i++ {
		a, b, c := i*3, i*3+1, i*3+2
		addEdge(a, b)
		addEdge(b, c)
		addEdge(c, a)
	}

	edges := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		edges = append(edges, e)
	}
	return edges
}
