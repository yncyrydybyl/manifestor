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
		name:     "emoji-rain",
		desc:     "Cascading emoji rain with lightning reveal",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 30,
		hasEmoji: true,
		play:     playEmojiRain,
	})
}

func playEmojiRain(filename string, size Size) {
	switch size {
	case OneLiner:
		playEmojiRainOneLiner(filename)
	case FiveLiner:
		playEmojiRainFiveLiner(filename)
	case FullScreen:
		playEmojiRainFullScreen(filename)
	default:
		playEmojiRainOneLiner(filename)
	}
}

func playEmojiRainOneLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	drops := []string{"🌟", "⭐", "💫", "✨", "🔮", "🌈", "🦋", "🌸", "💜", "💙", "💚", "💛", "🧡", "❤️"}

	// Phase 1: emoji rain builds up on one line
	for frame := 0; frame < 20; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		density := float64(frame) / 20.0
		slots := w / 2 // emoji are double-width
		for i := 0; i < slots; i++ {
			if rand.Float64() < density*0.6 {
				d := drops[rand.Intn(len(drops))]
				line.WriteString(d)
			} else {
				line.WriteString("  ")
			}
		}
		fmt.Fprint(os.Stderr, line.String())
		time.Sleep(40 * time.Millisecond)
	}

	// Phase 2: lightning flash
	zaps := []string{"⚡", "🌩️", "💥", "⚡"}
	for _, z := range zaps {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprint(os.Stderr, center(z+" "+z+" "+z, w))
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: reveal
	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🎆 %s manifested 🎆", filename)
	fmt.Fprintln(os.Stderr, center(bold+msg+reset, w+10))
	time.Sleep(400 * time.Millisecond)
}

func playEmojiRainFiveLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	drops := []string{"🌟", "⭐", "💫", "✨", "🔮", "🌈", "🦋", "🌸", "💜", "💙", "💚", "💛", "🧡", "❤️"}
	const rows = 5
	slots := w / 2

	// Build a grid of falling emoji
	grid := make([][]string, rows)
	for r := range grid {
		grid[r] = make([]string, slots)
		for c := range grid[r] {
			grid[r][c] = "  "
		}
	}

	// Print initial 5 blank lines
	for i := 0; i < rows; i++ {
		fmt.Fprintln(os.Stderr)
	}

	for frame := 0; frame < 25; frame++ {
		// Shift rows down: bottom row falls off, new row appears at top
		for r := rows - 1; r > 0; r-- {
			grid[r] = grid[r-1]
		}
		// Generate new top row
		density := float64(frame) / 25.0
		grid[0] = make([]string, slots)
		for c := range grid[0] {
			if rand.Float64() < density*0.5 {
				grid[0][c] = drops[rand.Intn(len(drops))]
			} else {
				grid[0][c] = "  "
			}
		}

		moveUp(rows)
		for r := 0; r < rows; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for c := 0; c < slots; c++ {
				line.WriteString(grid[r][c])
			}
			fmt.Fprintln(os.Stderr, line.String())
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Flash
	moveUp(rows)
	zaps := []string{"⚡", "🌩️", "💥", "⚡"}
	for _, z := range zaps {
		clearLines(rows)
		moveUp(rows - 1)
		midRow := rows / 2
		for r := 0; r < rows; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			if r == midRow {
				fmt.Fprintln(os.Stderr, center(z+" "+z+" "+z, w))
			} else {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(rows)
		time.Sleep(60 * time.Millisecond)
	}

	clearLines(rows)
	moveUp(rows - 1)
	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🎆 %s manifested 🎆", filename)
	fmt.Fprintln(os.Stderr, center(bold+msg+reset, w+10))
	time.Sleep(400 * time.Millisecond)
}

func playEmojiRainFullScreen(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	drops := []string{"🌟", "⭐", "💫", "✨", "🔮", "🌈", "🦋", "🌸", "💜", "💙", "💚", "💛", "🧡", "❤️"}
	slots := w / 2

	grid := make([][]string, h)
	for r := range grid {
		grid[r] = make([]string, slots)
		for c := range grid[r] {
			grid[r][c] = "  "
		}
	}

	// Print initial blank lines
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}

	totalFrames := 30
	for frame := 0; frame < totalFrames; frame++ {
		// Shift rows down
		for r := h - 1; r > 0; r-- {
			grid[r] = grid[r-1]
		}
		density := float64(frame) / float64(totalFrames)
		grid[0] = make([]string, slots)
		for c := range grid[0] {
			if rand.Float64() < density*0.4 {
				grid[0][c] = drops[rand.Intn(len(drops))]
			} else {
				grid[0][c] = "  "
			}
		}

		moveUp(h)
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for c := 0; c < slots; c++ {
				line.WriteString(grid[r][c])
			}
			fmt.Fprintln(os.Stderr, line.String())
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Flash
	moveUp(h)
	zaps := []string{"⚡", "🌩️", "💥", "⚡"}
	for _, z := range zaps {
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			if r == h/2 {
				fmt.Fprintln(os.Stderr, center(z+" "+z+" "+z, w))
			} else {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h)
		time.Sleep(60 * time.Millisecond)
	}

	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🎆 %s manifested 🎆", filename)
	fmt.Fprintln(os.Stderr, center(bold+msg+reset, w+10))
	time.Sleep(400 * time.Millisecond)
}
