package extractor

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Header represents the 12-byte GLB file header.
type Header struct {
	Magic   uint32 // Must be 0x46546C67 ("glTF")
	Version uint32 // Must be 2
	Length  uint32 // Total file length in bytes
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
