package anim

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func init() {
	register(regInfo{
		name:     "growing-rose",
		desc:     "A rose blooms in the terminal, petal by petal",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 30,
		hasEmoji: true,
		play:     playGrowingRose,
	})
}

func playGrowingRose(filename string, size Size) {
	switch size {
	case OneLiner:
		playGrowingRoseOneLiner(filename)
	case FiveLiner:
		playGrowingRoseFiveLiner(filename)
	case FullScreen:
		playGrowingRoseFullScreen(filename)
	default:
		playGrowingRoseFiveLiner(filename)
	}
}

func playGrowingRoseOneLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	green := rgb(34, 139, 34)
	pink := rgb(255, 20, 147)

	frames := []string{
		green + "  .  " + reset,
		green + " .|  " + reset,
		green + " 🌱  " + reset,
		green + " 🌿  " + reset,
		pink + " ,  " + green + "|" + reset,
		pink + "(@)" + green + "|" + reset,
		pink + "(@@@)" + green + "|" + reset,
		pink + "🌹" + green + " --|--" + reset,
	}

	for _, f := range frames {
		fmt.Fprint(os.Stderr, clearLn)
		vis := visibleLen(f)
		if vis < w {
			pad := (w - vis) / 2
			fmt.Fprint(os.Stderr, strings.Repeat(" ", pad)+f)
		} else {
			fmt.Fprint(os.Stderr, f)
		}
		time.Sleep(180 * time.Millisecond)
	}

	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🌹 %s has bloomed 🌹", filename)
	fmt.Fprintln(os.Stderr, center(bold+pink+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}

func playGrowingRoseFiveLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	green := rgb(34, 139, 34)
	pink := rgb(255, 20, 147)

	type frame struct {
		lines [5]string
		color string
	}

	frames := []frame{
		{[5]string{
			"",
			"",
			"",
			"    .",
			"    |",
		}, green},
		{[5]string{
			"",
			"",
			"   🌱",
			"    |",
			"    |",
		}, green},
		{[5]string{
			"",
			"   🌿",
			"    |",
			"   /|",
			"    |",
		}, green},
		{[5]string{
			"   ,",
			"  (@)",
			"   |",
			"  /|\\",
			"   |",
		}, rgb(255, 105, 180)},
		{[5]string{
			"  .~.",
			" (@@@)",
			"  (@)",
			"   |",
			"  /|\\",
		}, rgb(255, 50, 120)},
		{[5]string{
			" .~@~.",
			"(@@@@@)",
			" \\@@@/",
			"  (@)",
			"  /|\\",
		}, rgb(255, 20, 100)},
	}

	for i, f := range frames {
		if i > 0 {
			moveUp(5)
		}
		for _, line := range f.lines {
			fmt.Fprint(os.Stderr, clearLn)
			if line == "" {
				fmt.Fprintln(os.Stderr)
				continue
			}
			var colored strings.Builder
			for _, r := range line {
				switch r {
				case '|', '/', '\\':
					colored.WriteString(green + string(r))
				case '@', '~':
					colored.WriteString(pink + string(r))
				default:
					colored.WriteString(f.color + string(r))
				}
			}
			fmt.Fprintln(os.Stderr, center(colored.String()+reset, w+40))
		}
		time.Sleep(250 * time.Millisecond)
	}

	moveUp(5)
	clearLines(5)
	msg := fmt.Sprintf("🌹 %s has bloomed 🌹", filename)
	fmt.Fprintln(os.Stderr, center(bold+pink+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}

func playGrowingRoseFullScreen(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	green := rgb(34, 139, 34)
	pink := rgb(255, 20, 147)

	// Build progressively taller rose frames, padded to h lines
	roseStages := [][]string{
		{
			"    .",
			"    |",
		},
		{
			"   🌱",
			"    |",
			"    |",
			"    |",
		},
		{
			"   🌿",
			"    |",
			"   /|",
			"    |",
			"    |",
		},
		{
			"   ,",
			"  (@)",
			"   |",
			"  /|\\",
			"   |",
			"   |",
		},
		{
			"  .~.",
			" (@@@)",
			"  (@)",
			"   |",
			"  /|\\",
			"   |",
			"   |",
			"   |",
		},
		{
			" .~@~.",
			"(@@@@@)",
			" (@@@)",
			"  (@)",
			"   |",
			"  /|\\",
			"  |||",
			"  |||",
			"  |||",
		},
		{
			"  _____",
			" /~@@@~\\",
			"|@@@@@@@|",
			" \\@@@@@/",
			"  \\@@@/",
			"   (@)",
			"   |||",
			"  /|||\\",
			"   |||",
			"   |||",
			"   |||",
			"   |||",
		},
		{
			"    ___",
			"  .~@@@~.",
			" /~@@@@@~\\",
			"|@@@@@@@@@|",
			"|@@@@@@@@@|",
			" \\@@@@@@@/",
			"  \\@@@@@/",
			"   \\@@@/",
			"    (@)",
			"    |||",
			"   /|||\\",
			"    |||",
			"    |||",
			"    |||",
			"    |||",
			"    |||",
		},
	}

	stageColors := []string{
		green,
		green,
		green,
		rgb(255, 105, 180),
		rgb(255, 80, 150),
		rgb(255, 50, 120),
		rgb(255, 20, 100),
		rgb(255, 0, 80),
	}

	for i, stage := range roseStages {
		if i > 0 {
			moveUp(h)
		}

		// Pad so the rose sits at the bottom of the frame
		blank := h - len(stage)
		for j := 0; j < h; j++ {
			fmt.Fprint(os.Stderr, clearLn)
			if j < blank {
				fmt.Fprintln(os.Stderr)
				continue
			}
			line := stage[j-blank]
			var colored strings.Builder
			for _, r := range line {
				switch r {
				case '|', '/', '\\':
					colored.WriteString(green + string(r))
				case '@', '~', '_':
					colored.WriteString(pink + string(r))
				default:
					colored.WriteString(stageColors[i] + string(r))
				}
			}
			fmt.Fprintln(os.Stderr, center(colored.String()+reset, w+40))
		}
		time.Sleep(300 * time.Millisecond)
	}

	moveUp(h)
	for i := 0; i < h; i++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	msg := fmt.Sprintf("🌹 %s has bloomed 🌹", filename)
	fmt.Fprintln(os.Stderr, center(bold+pink+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}
