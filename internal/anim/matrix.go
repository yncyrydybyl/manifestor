package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register(regInfo{
		name:     "matrix-manifest",
		desc:     "Matrix-style falling glyphs that form the filename",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: false,
		play:     playMatrix,
	})
}

var matrixGlyphs = []rune("ァイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン")
var matrixGreens = []int{22, 28, 34, 40, 46, 82, 118, 154, 190, 226}

func playMatrix(filename string, size Size) {
	switch size {
	case FiveLiner:
		playMatrixFive(filename)
	case FullScreen:
		playMatrixFull(filename)
	default:
		playMatrixOne(filename)
	}
}

// playMatrixOne is the original single-line animation.
func playMatrixOne(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	columns := make([]int, w)
	for i := range columns {
		columns[i] = rand.Intn(len(matrixGreens))
	}

	// Phase 1: matrix rain
	for frame := 0; frame < 20; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder

		for i := 0; i < w; i++ {
			columns[i] = (columns[i] + 1) % len(matrixGreens)
			if rand.Float64() < 0.7 {
				g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
				brightness := matrixGreens[columns[i]]
				line.WriteString(fg(brightness) + string(g))
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: glyphs converge into filename
	msgRunes := []rune(center(filename, w))
	for frame := 0; frame < 12; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		progress := float64(frame) / 12.0

		for i := 0; i < w; i++ {
			if i < len(msgRunes) && msgRunes[i] != ' ' && rand.Float64() < progress {
				line.WriteString(bold + fg(46) + string(msgRunes[i]))
			} else if rand.Float64() < (1.0-progress)*0.7 {
				g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
				line.WriteString(dim + fg(22) + string(g))
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	// Final
	fmt.Fprint(os.Stderr, clearLn)
	msg := strings.Builder{}
	for i, r := range filename {
		c := matrixGreens[(i*3)%len(matrixGreens)]
		msg.WriteString(bold + fg(c) + string(r))
	}
	fmt.Fprintln(os.Stderr, center(msg.String()+reset, w+len(filename)*10))
	time.Sleep(400 * time.Millisecond)
}

// playMatrixFive renders 5-line falling matrix columns converging to filename.
func playMatrixFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Each column has a "head" row that falls
	type col struct {
		head   int // current row of the bright head
		speed  int // rows per frame
		glyph  rune
		active bool
	}
	columns := make([]col, w)
	for i := range columns {
		columns[i] = col{
			head:   rand.Intn(h),
			speed:  1 + rand.Intn(2),
			glyph:  matrixGlyphs[rand.Intn(len(matrixGlyphs))],
			active: rand.Float64() < 0.6,
		}
	}

	// Print initial h blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: rain
	for frame := 0; frame < 20; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if !columns[i].active {
					line.WriteString(" ")
					continue
				}
				dist := columns[i].head - r
				if dist < 0 {
					dist = -dist
				}
				if dist == 0 {
					// bright head
					line.WriteString(bold + fg(46) + string(columns[i].glyph))
				} else if dist <= 2 && r <= columns[i].head {
					// trail
					ci := len(matrixGreens) - 1 - dist*3
					if ci < 0 {
						ci = 0
					}
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					line.WriteString(fg(matrixGreens[ci]) + string(g))
				} else if rand.Float64() < 0.15 {
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					line.WriteString(dim + fg(22) + string(g))
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)

		// Advance heads
		for i := range columns {
			columns[i].head = (columns[i].head + columns[i].speed) % (h + 3)
			if rand.Float64() < 0.05 {
				columns[i].active = !columns[i].active
			}
			if rand.Float64() < 0.1 {
				columns[i].glyph = matrixGlyphs[rand.Intn(len(matrixGlyphs))]
			}
		}
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 2: converge to filename on middle row
	msgRunes := []rune(center(filename, w))
	midRow := h / 2
	for frame := 0; frame < 12; frame++ {
		progress := float64(frame) / 12.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if r == midRow && i < len(msgRunes) && msgRunes[i] != ' ' && rand.Float64() < progress {
					line.WriteString(bold + fg(46) + string(msgRunes[i]))
				} else if rand.Float64() < (1.0-progress)*0.5 {
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					line.WriteString(dim + fg(22) + string(g))
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(60 * time.Millisecond)
	}

	// Final: clear the 5 lines, leave a single-line message
	clearLines(h)
	msg := strings.Builder{}
	for i, r := range filename {
		c := matrixGreens[(i*3)%len(matrixGreens)]
		msg.WriteString(bold + fg(c) + string(r))
	}
	fmt.Fprintln(os.Stderr, center(msg.String()+reset, w+len(filename)*10))
	time.Sleep(400 * time.Millisecond)
}

// playMatrixFull renders full-screen matrix rain converging to filename.
func playMatrixFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	type col struct {
		head  int
		speed int
		glyph rune
	}
	columns := make([]col, w)
	for i := range columns {
		columns[i] = col{
			head:  rand.Intn(h),
			speed: 1 + rand.Intn(3),
			glyph: matrixGlyphs[rand.Intn(len(matrixGlyphs))],
		}
	}

	// Print initial blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: classic matrix rain
	for frame := 0; frame < 30; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				dist := columns[i].head - r
				if dist < 0 {
					dist += h
				}
				if dist == 0 {
					line.WriteString(bold + fg(255) + string(columns[i].glyph))
				} else if dist > 0 && dist <= 4 {
					ci := len(matrixGreens) - 1 - dist*2
					if ci < 0 {
						ci = 0
					}
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					line.WriteString(fg(matrixGreens[ci]) + string(g))
				} else if rand.Float64() < 0.08 {
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					line.WriteString(dim + fg(22) + string(g))
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)

		for i := range columns {
			columns[i].head = (columns[i].head + columns[i].speed) % h
			if rand.Float64() < 0.08 {
				columns[i].glyph = matrixGlyphs[rand.Intn(len(matrixGlyphs))]
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: converge to filename at center
	msgRunes := []rune(center(filename, w))
	midRow := h / 2
	for frame := 0; frame < 15; frame++ {
		progress := float64(frame) / 15.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if r == midRow && i < len(msgRunes) && msgRunes[i] != ' ' && rand.Float64() < progress {
					line.WriteString(bold + fg(46) + string(msgRunes[i]))
				} else if rand.Float64() < (1.0-progress)*0.4 {
					g := matrixGlyphs[rand.Intn(len(matrixGlyphs))]
					ci := int((1.0 - progress) * float64(len(matrixGreens)-1))
					line.WriteString(dim + fg(matrixGreens[ci]) + string(g))
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(60 * time.Millisecond)
	}

	// Final: clear all lines, print single-line result
	clearLines(h)
	msg := strings.Builder{}
	for i, r := range filename {
		c := matrixGreens[(i*3)%len(matrixGreens)]
		msg.WriteString(bold + fg(c) + string(r))
	}
	fmt.Fprintln(os.Stderr, center(msg.String()+reset, w+len(filename)*10))
	time.Sleep(400 * time.Millisecond)
}
