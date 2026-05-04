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
		name:     "fire-forge",
		desc:     "File is forged in flames",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: true,
		play:     playFireForge,
	})
}

var fireFlames = []string{"▀", "▄", "█", "░", "▒", "▓", "^", "~"}
var fireColors = []int{52, 88, 124, 160, 196, 202, 208, 214, 220, 226, 228, 230} // dark red -> yellow -> white

func playFireForge(filename string, size Size) {
	switch size {
	case FiveLiner:
		playFireFive(filename)
	case FullScreen:
		playFireFull(filename)
	default:
		playFireOne(filename)
	}
}

// fireHeat returns a color index based on distance from center and randomness.
func fireHeat(i, w int) int {
	distFromCenter := float64(i-w/2) / float64(w/2)
	if distFromCenter < 0 {
		distFromCenter = -distFromCenter
	}
	heat := int((1.0 - distFromCenter) * float64(len(fireColors)-1))
	if heat < 0 {
		heat = 0
	}
	heat += rand.Intn(3) - 1
	if heat < 0 {
		heat = 0
	}
	if heat >= len(fireColors) {
		heat = len(fireColors) - 1
	}
	return heat
}

// playFireOne is the original single-line flame animation.
func playFireOne(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Phase 1: flames build
	for frame := 0; frame < 18; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		intensity := float64(frame) / 18.0

		for i := 0; i < w; i++ {
			if rand.Float64() < intensity*0.8 {
				f := fireFlames[rand.Intn(len(fireFlames))]
				heat := fireHeat(i, w)
				line.WriteString(fg(fireColors[heat]) + f)
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(40 * time.Millisecond)
	}

	// Phase 2: forge flash
	for frame := 0; frame < 4; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		for i := 0; i < w; i++ {
			c := fireColors[len(fireColors)-1-rand.Intn(3)]
			f := fireFlames[rand.Intn(len(fireFlames))]
			_ = i
			line.WriteString(bold + fg(c) + f)
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: cooling reveal
	for frame := 0; frame < 8; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		progress := float64(frame) / 8.0
		msgRunes := []rune(center(filename, w))

		for i := 0; i < w; i++ {
			if i < len(msgRunes) && msgRunes[i] != ' ' {
				heat := int((1.0 - progress*0.5) * float64(len(fireColors)-1))
				line.WriteString(bold + fg(fireColors[heat]) + string(msgRunes[i]))
			} else if rand.Float64() < (1.0-progress)*0.5 {
				f := fireFlames[rand.Intn(len(fireFlames))]
				heat := rand.Intn(len(fireColors) / 2)
				line.WriteString(dim + fg(fireColors[heat]) + f)
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🔥 %s forged 🔥", filename)
	fmt.Fprintln(os.Stderr, center(fg(208)+bold+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}

// playFireFive renders a 5-line flame wall forging the filename.
func playFireFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Print initial blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: flames build upward (hotter at bottom)
	for frame := 0; frame < 18; frame++ {
		intensity := float64(frame) / 18.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			// Bottom rows are hotter
			rowHeat := float64(h-1-r) / float64(h-1) // 1.0 at bottom, 0.0 at top
			for i := 0; i < w; i++ {
				if rand.Float64() < intensity*(0.3+rowHeat*0.6) {
					f := fireFlames[rand.Intn(len(fireFlames))]
					heat := fireHeat(i, w)
					// Reduce heat for upper rows
					heat = int(float64(heat) * (0.4 + rowHeat*0.6))
					if heat >= len(fireColors) {
						heat = len(fireColors) - 1
					}
					line.WriteString(fg(fireColors[heat]) + f)
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
		time.Sleep(45 * time.Millisecond)
	}

	// Phase 2: forge flash (all lines white-hot)
	for frame := 0; frame < 4; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				c := fireColors[len(fireColors)-1-rand.Intn(3)]
				f := fireFlames[rand.Intn(len(fireFlames))]
				_ = i
				line.WriteString(bold + fg(c) + f)
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: cooling reveal on center row
	msgRunes := []rune(center(filename, w))
	midRow := h / 2
	for frame := 0; frame < 8; frame++ {
		progress := float64(frame) / 8.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if r == midRow && i < len(msgRunes) && msgRunes[i] != ' ' {
					heat := int((1.0 - progress*0.5) * float64(len(fireColors)-1))
					line.WriteString(bold + fg(fireColors[heat]) + string(msgRunes[i]))
				} else if rand.Float64() < (1.0-progress)*0.5 {
					f := fireFlames[rand.Intn(len(fireFlames))]
					heat := rand.Intn(len(fireColors) / 2)
					line.WriteString(dim + fg(fireColors[heat]) + f)
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

	// Final
	clearLines(h)
	msg := fmt.Sprintf("🔥 %s forged 🔥", filename)
	fmt.Fprintln(os.Stderr, center(fg(208)+bold+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}

// playFireFull renders a full-screen inferno forging the filename.
func playFireFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Print initial blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: inferno builds (classic fire effect — hotter at bottom)
	for frame := 0; frame < 24; frame++ {
		intensity := float64(frame) / 24.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			rowHeat := float64(h-1-r) / float64(h-1)
			for i := 0; i < w; i++ {
				if rand.Float64() < intensity*(0.2+rowHeat*0.7) {
					f := fireFlames[rand.Intn(len(fireFlames))]
					heat := fireHeat(i, w)
					heat = int(float64(heat) * (0.3 + rowHeat*0.7))
					if heat >= len(fireColors) {
						heat = len(fireColors) - 1
					}
					line.WriteString(fg(fireColors[heat]) + f)
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
		time.Sleep(40 * time.Millisecond)
	}

	// Phase 2: forge flash
	for frame := 0; frame < 4; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				c := fireColors[len(fireColors)-1-rand.Intn(3)]
				f := fireFlames[rand.Intn(len(fireFlames))]
				_ = i
				line.WriteString(bold + fg(c) + f)
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: cooling reveal on center row
	msgRunes := []rune(center(filename, w))
	midRow := h / 2
	for frame := 0; frame < 12; frame++ {
		progress := float64(frame) / 12.0
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if r == midRow && i < len(msgRunes) && msgRunes[i] != ' ' {
					heat := int((1.0 - progress*0.5) * float64(len(fireColors)-1))
					line.WriteString(bold + fg(fireColors[heat]) + string(msgRunes[i]))
				} else if rand.Float64() < (1.0-progress)*0.4 {
					f := fireFlames[rand.Intn(len(fireFlames))]
					heat := rand.Intn(len(fireColors) / 2)
					line.WriteString(dim + fg(fireColors[heat]) + f)
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
		time.Sleep(50 * time.Millisecond)
	}

	// Final
	clearLines(h)
	msg := fmt.Sprintf("🔥 %s forged 🔥", filename)
	fmt.Fprintln(os.Stderr, center(fg(208)+bold+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}
