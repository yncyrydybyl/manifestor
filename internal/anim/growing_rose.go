package anim

import (
	"fmt"
	"os"
	"time"
)

func init() {
	register("growing-rose", "A rose blooms in the terminal, petal by petal", playGrowingRose)
}

func playGrowingRose(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Rose frames — each frame adds more petals
	frames := []struct {
		lines []string
		color string
	}{
		{[]string{
			"    .",
			"    |",
		}, rgb(34, 139, 34)},
		{[]string{
			"   �‍🌱",
			"    |",
			"    |",
		}, rgb(34, 139, 34)},
		{[]string{
			"   🌿",
			"    |",
			"   �/|",
		}, rgb(34, 139, 34)},
		{[]string{
			"   ,",
			"  (@)",
			"   |",
			"  /|\\",
		}, rgb(255, 105, 180)},
		{[]string{
			"  .~.",
			" (@@@)",
			"  (@)",
			"   |",
			"  /|\\",
		}, rgb(255, 80, 150)},
		{[]string{
			" .~@~.",
			"(@@@@@)",
			" (@@@)",
			"  (@)",
			"   |",
			"  /|\\",
			"  |||",
		}, rgb(255, 50, 120)},
		{[]string{
			"  _____",
			" /~@@@~\\",
			"|@@@@@@@|",
			" \\@@@@@/",
			"  \\@@@/",
			"   (@)",
			"   |||",
			"  /|||\\ ",
			"   |||",
		}, rgb(255, 20, 100)},
	}

	green := rgb(34, 139, 34)
	pink := rgb(255, 20, 147)

	for i, f := range frames {
		// Clear previous frame
		if i > 0 {
			for range frames[i-1].lines {
				fmt.Fprintf(os.Stderr, "\033[A"+clearLn)
			}
		}
		for _, line := range f.lines {
			centered := center(line, w)
			// Color the stems green, petals pink
			var colored string
			for _, r := range centered {
				switch r {
				case '|', '/', '\\':
					colored += green + string(r)
				case '@', '~':
					colored += pink + string(r)
				default:
					colored += f.color + string(r)
				}
			}
			fmt.Fprintln(os.Stderr, colored+reset)
		}
		if i < len(frames)-1 {
			time.Sleep(250 * time.Millisecond)
		}
	}

	// Final message
	msg := fmt.Sprintf("🌹 %s has bloomed 🌹", filename)
	fmt.Fprintln(os.Stderr, center(bold+pink+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}
