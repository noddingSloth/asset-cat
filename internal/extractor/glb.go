package extractor

import "io"

// GLBExtractor handles reading and parsing GLB files.
type GLBExtractor struct {
	Reader io.Reader
}
