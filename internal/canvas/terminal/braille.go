package terminal

// BrailleCanvas provides high-resolution sub-pixel drawing using braille characters.
// Each braille character represents a 2×4 dot grid, giving 2x horizontal and 4x vertical
// sub-pixel resolution compared to standard terminal cells.
type BrailleCanvas struct {
	cols int
	rows int      // in braille cells (each cell is 2×4 dots)
	dots [][]bool // true = dot set; dimensions: [cols*2][rows*4]
}

// NewBrailleCanvas creates a new braille canvas with the given terminal dimensions
// in characters (cols, rows).
func NewBrailleCanvas(cols, rows int) *BrailleCanvas {
	dotCols := cols * 2
	dotRows := rows * 4
	dots := make([][]bool, dotCols)
	for i := range dots {
		dots[i] = make([]bool, dotRows)
	}
	return &BrailleCanvas{
		cols: cols,
		rows: rows,
		dots: dots,
	}
}

// Set sets a dot at the given sub-pixel coordinates.
func (bc *BrailleCanvas) Set(x, y int) {
	if x >= 0 && x < len(bc.dots) && y >= 0 && y < len(bc.dots[0]) {
		bc.dots[x][y] = true
	}
}

// Clear resets all dots to false.
func (bc *BrailleCanvas) Clear() {
	for col := range bc.dots {
		for row := range bc.dots[col] {
			bc.dots[col][row] = false
		}
	}
}

// Render returns the braille grid as a string with ANSI escape codes.
// Each 2×4 block of dots is converted to a Unicode braille character.
func (bc *BrailleCanvas) Render() string {
	// Braille dot positions within the 2×4 grid (Unicode standard):
	// 0 3
	// 1 4
	// 2 5
	// 6 7
	//
	// Bit values: dot0=0x01, dot1=0x02, dot2=0x04, dot3=0x08,
	//             dot4=0x10, dot5=0x20, dot6=0x40, dot7=0x80
	//
	// Mapping from our grid (col, row) to braille dot:
	// col 0, row 0 → dot 0 (0x01)
	// col 0, row 1 → dot 1 (0x02)
	// col 0, row 2 → dot 2 (0x04)
	// col 1, row 0 → dot 3 (0x08)
	// col 1, row 1 → dot 4 (0x10)
	// col 1, row 2 → dot 5 (0x20)
	// col 0, row 3 → dot 6 (0x40)
	// col 1, row 3 → dot 7 (0x80)

	dotMap := [2][4]uint8{
		{0x01, 0x02, 0x04, 0x40}, // left column
		{0x08, 0x10, 0x20, 0x80}, // right column
	}

	var result string
	for cellRow := 0; cellRow < bc.rows; cellRow++ {
		for cellCol := 0; cellCol < bc.cols; cellCol++ {
			var code uint8
			baseCol := cellCol * 2
			baseRow := cellRow * 4

			for dc := 0; dc < 2; dc++ {
				for dr := 0; dr < 4; dr++ {
					if bc.dots[baseCol+dc][baseRow+dr] {
						code |= dotMap[dc][dr]
					}
				}
			}
			result += string(rune(0x2800 + int(code)))
		}
		if cellRow < bc.rows-1 {
			result += "\n"
		}
	}
	return result
}

// DrawLine draws a line using Bresenham's algorithm at sub-pixel resolution.
func (bc *BrailleCanvas) DrawLine(x1, y1, x2, y2 int) {
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx := 1
	sy := 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx + dy

	x, y := x1, y1
	for {
		bc.Set(x, y)
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			if x == x2 {
				break
			}
			err += dy
			x += sx
		}
		if e2 <= dx {
			if y == y2 {
				break
			}
			err += dx
			y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
