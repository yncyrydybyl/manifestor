package anim

import (
	"fmt"
	"os"
	"time"
)

func init() {
	register(regInfo{
		name:     "rainbow-beam",
		desc:     "Full-width rainbow bar with sparkle burst",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: true,
		play:     playRainbowBeam,
	})
}

func playRainbowBeam(filename string, size Size) {
	switch size {
	case FiveLiner:
		playRainbowBeamFive(filename)
	case FullScreen:
		playRainbowBeamFull(filename)
	default:
		playRainbowBeamOne(filename)
	}
}

func playRainbowBeamOne(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	blocks := []string{"░", "▒", "▓", "█"}
	rainbow := []int{196, 208, 220, 46, 51, 21, 93}

	// Phase 1: beam fills left to right
	for i := 0; i <= w; i++ {
		fmt.Fprint(os.Stderr, clearLn)
		for j := 0; j < i; j++ {
			ci := (j * len(rainbow) / w) % len(rainbow)
			bi := 3
			if j == i-1 {
				bi = 0
			} else if j == i-2 {
				bi = 1
			} else if j == i-3 {
				bi = 2
			}
			fmt.Fprintf(os.Stderr, "%s%s", fg(rainbow[ci]), blocks[bi])
		}
		fmt.Fprint(os.Stderr, reset)
		time.Sleep(6 * time.Millisecond)
	}

	// Phase 2: sparkle burst
	sparkles := []string{"✦", "✧", "⟡", "◇", "⋆", "✺", "❋"}
	for frame := 0; frame < 8; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line string
		for j := 0; j < w; j++ {
			ci := (j * len(rainbow) / w) % len(rainbow)
			if (j+frame)%4 == 0 {
				s := sparkles[(j+frame)%len(sparkles)]
				line += fmt.Sprintf("%s%s%s", bold+fg(rainbow[ci]), s, reset)
			} else {
				line += fmt.Sprintf("%s█", fg(rainbow[ci]))
			}
		}
		fmt.Fprint(os.Stderr, line+reset)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: text reveal
	msg := fmt.Sprintf(" ✨ %s manifested ✨ ", filename)
	fmt.Fprint(os.Stderr, clearLn)
	centered := center(msg, w)
	for i, r := range centered {
		ci := (i * len(rainbow) / len(centered)) % len(rainbow)
		fmt.Fprintf(os.Stderr, "%s%s%c", bold, fg(rainbow[ci]), r)
	}
	fmt.Fprintln(os.Stderr, reset)
	time.Sleep(400 * time.Millisecond)
}

func playRainbowBeamFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	blocks := []string{"░", "▒", "▓", "█"}
	rainbow := []int{196, 208, 220, 46, 51, 21, 93}

	// Print initial blank lines to reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: beams build up row by row, each sweeping left to right
	for row := h - 1; row >= 0; row-- {
		steps := 12
		for s := 0; s <= steps; s++ {
			progress := s * w / steps
			// Redraw all h lines
			for r := 0; r < h; r++ {
				fmt.Fprint(os.Stderr, clearLn)
				if r < row {
					// not yet reached
					fmt.Fprintln(os.Stderr)
					continue
				}
				rowFill := w
				if r == row {
					rowFill = progress
				}
				for j := 0; j < rowFill; j++ {
					ci := (j*len(rainbow)/w + r) % len(rainbow)
					bi := 3
					if r == row && j >= progress-3 {
						bi = progress - j
						if bi < 0 {
							bi = 0
						}
						if bi > 3 {
							bi = 3
						}
					}
					fmt.Fprintf(os.Stderr, "%s%s", fg(rainbow[ci]), blocks[bi])
				}
				fmt.Fprint(os.Stderr, reset)
				fmt.Fprintln(os.Stderr)
			}
			moveUp(h)
			time.Sleep(8 * time.Millisecond)
		}
	}

	// Phase 2: sparkle burst across all 5 lines
	sparkles := []string{"✦", "✧", "⟡", "◇", "⋆", "✺", "❋"}
	for frame := 0; frame < 6; frame++ {
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			for j := 0; j < w; j++ {
				ci := (j*len(rainbow)/w + r) % len(rainbow)
				if (j+r+frame)%4 == 0 {
					s := sparkles[(j+r+frame)%len(sparkles)]
					fmt.Fprintf(os.Stderr, "%s%s%s", bold+fg(rainbow[ci]), s, reset)
				} else {
					fmt.Fprintf(os.Stderr, "%s█", fg(rainbow[ci]))
				}
			}
			fmt.Fprint(os.Stderr, reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: clear and show final message on one line
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	msg := fmt.Sprintf(" ✨ %s manifested ✨ ", filename)
	fmt.Fprint(os.Stderr, clearLn)
	centered := center(msg, w)
	for i, r := range centered {
		ci := (i * len(rainbow) / len(centered)) % len(rainbow)
		fmt.Fprintf(os.Stderr, "%s%s%c", bold, fg(rainbow[ci]), r)
	}
	fmt.Fprintln(os.Stderr, reset)
	time.Sleep(400 * time.Millisecond)
}

func playRainbowBeamFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	blocks := []string{"░", "▒", "▓", "█"}
	rainbow := []int{196, 208, 220, 46, 51, 21, 93}

	// Print initial blank lines to reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: rainbow fills top-to-bottom, each row sweeps left to right quickly
	rowsPerFrame := h / 10
	if rowsPerFrame < 1 {
		rowsPerFrame = 1
	}
	filledRows := 0
	for filledRows < h {
		filledRows += rowsPerFrame
		if filledRows > h {
			filledRows = h
		}
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			if r < filledRows {
				for j := 0; j < w; j++ {
					ci := (j*len(rainbow)/w + r) % len(rainbow)
					bi := 3
					if r >= filledRows-rowsPerFrame && r < filledRows {
						edge := (filledRows - r) * 3
						if edge < 3 {
							bi = edge
						}
					}
					fmt.Fprintf(os.Stderr, "%s%s", fg(rainbow[ci]), blocks[bi])
				}
				fmt.Fprint(os.Stderr, reset)
			}
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(30 * time.Millisecond)
	}

	// Phase 2: sparkle burst across all lines
	sparkles := []string{"✦", "✧", "⟡", "◇", "⋆", "✺", "❋"}
	for frame := 0; frame < 6; frame++ {
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			for j := 0; j < w; j++ {
				ci := (j*len(rainbow)/w + r) % len(rainbow)
				if (j+r+frame)%5 == 0 {
					s := sparkles[(j+r+frame)%len(sparkles)]
					fmt.Fprintf(os.Stderr, "%s%s%s", bold+fg(rainbow[ci]), s, reset)
				} else {
					fmt.Fprintf(os.Stderr, "%s█", fg(rainbow[ci]))
				}
			}
			fmt.Fprint(os.Stderr, reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(80 * time.Millisecond)
	}

	// Phase 3: clear screen and show final message
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	msg := fmt.Sprintf(" ✨ %s manifested ✨ ", filename)
	fmt.Fprint(os.Stderr, clearLn)
	centered := center(msg, w)
	for i, r := range centered {
		ci := (i * len(rainbow) / len(centered)) % len(rainbow)
		fmt.Fprintf(os.Stderr, "%s%s%c", bold, fg(rainbow[ci]), r)
	}
	fmt.Fprintln(os.Stderr, reset)
	time.Sleep(400 * time.Millisecond)
}
