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

	// Constants for GLB specification
	const expectedMagic = 0x46546C67 // "glTF"
	const expectedVersion = 2

	if header.Magic != expectedMagic {
		t.Errorf("expected magic 0x%X (glTF), got 0x%X", expectedMagic, header.Magic)
	}

	if header.Version != expectedVersion {
		t.Errorf("expected version %d, got %d", expectedVersion, header.Version)
	}

	// The header length must match the actual file size
	info, err := file.Stat()
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if header.Length != uint32(info.Size()) {
		t.Errorf("expected length to match file size %d, got %d", info.Size(), header.Length)
	}
}
