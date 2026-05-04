package anim

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func init() {
	register(regInfo{
		name:     "crystallize",
		desc:     "Crystals form and shatter to reveal the file",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: true,
		play:     playCrystallize,
	})
}

func playCrystallize(filename string, size Size) {
	switch size {
	case FiveLiner:
		playCrystallizeFive(filename)
	case FullScreen:
		playCrystallizeFull(filename)
	default:
		playCrystallizeOne(filename)
	}
}

func playCrystallizeOne(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	ice := []int{159, 153, 147, 141, 135, 129, 123, 117, 111, 105}
	crystals := []string{"❄", "❅", "❆", "✶", "✸", "✹", "◆", "◇", "⬡", "⬢"}

	// Phase 1: crystals grow across the terminal
	for frame := 0; frame < 15; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		density := float64(frame) / 15.0

		for i := 0; i < w; i++ {
			distFromCenter := float64(i-w/2) / float64(w/2)
			if distFromCenter < 0 {
				distFromCenter = -distFromCenter
			}
			if distFromCenter < density {
				ci := (i + frame) % len(ice)
				cr := (i + frame*3) % len(crystals)
				line.WriteString(fg(ice[ci]) + crystals[cr])
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: crystals shimmer
	for frame := 0; frame < 6; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		for i := 0; i < w; i++ {
			ci := (i + frame*2) % len(ice)
			cr := (i + frame*5) % len(crystals)
			line.WriteString(bold + fg(ice[ci]) + crystals[cr])
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: shatter — crystals break apart revealing text
	msg := fmt.Sprintf(" 💎 %s crystallized 💎 ", filename)
	msgRunes := []rune(center(msg, w))

	for frame := 0; frame < 10; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		progress := float64(frame) / 10.0

		for i := 0; i < w; i++ {
			distFromCenter := float64(i-w/2) / float64(w/2)
			if distFromCenter < 0 {
				distFromCenter = -distFromCenter
			}
			if distFromCenter < progress && i < len(msgRunes) {
				ci := (i + frame) % len(ice)
				line.WriteString(bold + fg(ice[ci]) + string(msgRunes[i]))
			} else {
				ci := (i + frame*3) % len(ice)
				cr := (i + frame*7) % len(crystals)
				line.WriteString(dim + fg(ice[ci]) + crystals[cr])
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	fmt.Fprintln(os.Stderr)
	time.Sleep(300 * time.Millisecond)
}

func playCrystallizeFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	ice := []int{159, 153, 147, 141, 135, 129, 123, 117, 111, 105}
	crystals := []string{"❄", "❅", "❆", "✶", "✸", "✹", "◆", "◇", "⬡", "⬢"}

	// Reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: crystals grow upward from the bottom row
	for frame := 0; frame < 15; frame++ {
		density := float64(frame) / 15.0
		// How many rows are filled (bottom-up)
		activeRows := int(density*float64(h)) + 1
		if activeRows > h {
			activeRows = h
		}
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			rowFromBottom := h - 1 - r
			if rowFromBottom < activeRows {
				var line strings.Builder
				rowDensity := density
				if rowFromBottom == activeRows-1 {
					rowDensity = density * 0.5
				}
				for i := 0; i < w; i++ {
					distFromCenter := float64(i-w/2) / float64(w/2)
					if distFromCenter < 0 {
						distFromCenter = -distFromCenter
					}
					if distFromCenter < rowDensity {
						ci := (i + frame + r) % len(ice)
						cr := (i + frame*3 + r*2) % len(crystals)
						line.WriteString(fg(ice[ci]) + crystals[cr])
					} else {
						line.WriteString(" ")
					}
				}
				fmt.Fprint(os.Stderr, line.String()+reset)
			}
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: shimmer across all 5 lines
	for frame := 0; frame < 6; frame++ {
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				ci := (i + frame*2 + r) % len(ice)
				cr := (i + frame*5 + r*3) % len(crystals)
				line.WriteString(bold + fg(ice[ci]) + crystals[cr])
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: shatter to reveal text
	msg := fmt.Sprintf(" 💎 %s crystallized 💎 ", filename)
	msgRunes := []rune(center(msg, w))

	for frame := 0; frame < 10; frame++ {
		progress := float64(frame) / 10.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			if r == h/2 {
				// Middle row: text reveal
				var line strings.Builder
				for i := 0; i < w; i++ {
					distFromCenter := float64(i-w/2) / float64(w/2)
					if distFromCenter < 0 {
						distFromCenter = -distFromCenter
					}
					if distFromCenter < progress && i < len(msgRunes) {
						ci := (i + frame) % len(ice)
						line.WriteString(bold + fg(ice[ci]) + string(msgRunes[i]))
					} else {
						ci := (i + frame*3 + r) % len(ice)
						cr := (i + frame*7 + r*2) % len(crystals)
						line.WriteString(dim + fg(ice[ci]) + crystals[cr])
					}
				}
				fmt.Fprint(os.Stderr, line.String()+reset)
			} else {
				// Other rows: crystals fade
				var line strings.Builder
				for i := 0; i < w; i++ {
					if float64(r)/float64(h) < progress || float64(h-1-r)/float64(h) < progress {
						line.WriteString(" ")
					} else {
						ci := (i + frame*3 + r) % len(ice)
						cr := (i + frame*7 + r*2) % len(crystals)
						line.WriteString(dim + fg(ice[ci]) + crystals[cr])
					}
				}
				fmt.Fprint(os.Stderr, line.String()+reset)
			}
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(60 * time.Millisecond)
	}

	// Final: clear and show message on one line
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	fmt.Fprint(os.Stderr, clearLn)
	for i, r := range msgRunes {
		ci := (i) % len(ice)
		fmt.Fprintf(os.Stderr, "%s%s%c", bold, fg(ice[ci]), r)
	}
	fmt.Fprintln(os.Stderr, reset)
	time.Sleep(300 * time.Millisecond)
}

func playCrystallizeFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	ice := []int{159, 153, 147, 141, 135, 129, 123, 117, 111, 105}
	crystals := []string{"❄", "❅", "❆", "✶", "✸", "✹", "◆", "◇", "⬡", "⬢"}

	// Reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: crystal cave grows from center outward and bottom up
	for frame := 0; frame < 20; frame++ {
		density := float64(frame) / 20.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			rowFromBottom := h - 1 - r
			rowDensity := density * (1.0 - float64(rowFromBottom)/float64(h)*0.7)
			if rowDensity < 0 {
				rowDensity = 0
			}
			for i := 0; i < w; i++ {
				distFromCenter := float64(i-w/2) / float64(w/2)
				if distFromCenter < 0 {
					distFromCenter = -distFromCenter
				}
				if distFromCenter < rowDensity {
					ci := (i + frame + r) % len(ice)
					cr := (i + frame*3 + r*2) % len(crystals)
					line.WriteString(fg(ice[ci]) + crystals[cr])
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(40 * time.Millisecond)
	}

	// Phase 2: full cave shimmer
	for frame := 0; frame < 6; frame++ {
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				ci := (i + frame*2 + r) % len(ice)
				cr := (i + frame*5 + r*3) % len(crystals)
				line.WriteString(bold + fg(ice[ci]) + crystals[cr])
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: shatter from center revealing text on middle row
	msg := fmt.Sprintf(" 💎 %s crystallized 💎 ", filename)
	msgRunes := []rune(center(msg, w))
	midRow := h / 2

	for frame := 0; frame < 12; frame++ {
		progress := float64(frame) / 12.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			rowDist := float64(r-midRow) / float64(h/2)
			if rowDist < 0 {
				rowDist = -rowDist
			}
			if r == midRow {
				var line strings.Builder
				for i := 0; i < w; i++ {
					distFromCenter := float64(i-w/2) / float64(w/2)
					if distFromCenter < 0 {
						distFromCenter = -distFromCenter
					}
					if distFromCenter < progress && i < len(msgRunes) {
						ci := (i + frame) % len(ice)
						line.WriteString(bold + fg(ice[ci]) + string(msgRunes[i]))
					} else {
						ci := (i + frame*3 + r) % len(ice)
						cr := (i + frame*7 + r*2) % len(crystals)
						line.WriteString(dim + fg(ice[ci]) + crystals[cr])
					}
				}
				fmt.Fprint(os.Stderr, line.String()+reset)
			} else if rowDist < progress {
				// Row has shattered away
			} else {
				var line strings.Builder
				for i := 0; i < w; i++ {
					ci := (i + frame*3 + r) % len(ice)
					cr := (i + frame*7 + r*2) % len(crystals)
					line.WriteString(dim + fg(ice[ci]) + crystals[cr])
				}
				fmt.Fprint(os.Stderr, line.String()+reset)
			}
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(50 * time.Millisecond)
	}

	// Final: clear and show message
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	fmt.Fprint(os.Stderr, clearLn)
	for i, r := range msgRunes {
		ci := (i) % len(ice)
		fmt.Fprintf(os.Stderr, "%s%s%c", bold, fg(ice[ci]), r)
	}
	fmt.Fprintln(os.Stderr, reset)
	time.Sleep(300 * time.Millisecond)
}
