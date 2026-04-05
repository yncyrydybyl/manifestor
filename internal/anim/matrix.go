package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register("matrix-manifest", "Matrix-style falling glyphs that form the filename", playMatrix)
}

func playMatrix(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	glyphs := []rune("ァイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン")
	greens := []int{22, 28, 34, 40, 46, 82, 118, 154, 190, 226}

	// Phase 1: matrix rain
	columns := make([]int, w) // current "head" position per column (cycle through brightness)
	for i := range columns {
		columns[i] = rand.Intn(len(greens))
	}

	for frame := 0; frame < 20; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder

		for i := 0; i < w; i++ {
			columns[i] = (columns[i] + 1) % len(greens)
			if rand.Float64() < 0.7 {
				g := glyphs[rand.Intn(len(glyphs))]
				brightness := greens[columns[i]]
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
				// Reveal real character
				line.WriteString(bold + fg(46) + string(msgRunes[i]))
			} else if rand.Float64() < (1.0-progress)*0.7 {
				g := glyphs[rand.Intn(len(glyphs))]
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
		c := greens[(i*3)%len(greens)]
		msg.WriteString(bold + fg(c) + string(r))
	}
	fmt.Fprintln(os.Stderr, center(msg.String()+reset, w+len(filename)*10))
	time.Sleep(400 * time.Millisecond)
}
