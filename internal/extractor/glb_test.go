package extractor_test

import (
	"os"
	"testing"

	"github.com/noddingSloth/asset-cat/internal/extractor"
)

func TestReadHeaderFromTuxAsset(t *testing.T) {
	file, err := os.Open("../../assets/tux/Linux mascot Tux.glb")
	if err != nil {
		t.Fatalf("failed to open Tux GLB asset: %v", err)
	}
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
	file, err := os.Open("../../assets/tux/Linux mascot Tux.glb")
	if err != nil {
		t.Fatalf("failed to open Tux GLB asset: %v", err)
	}
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
