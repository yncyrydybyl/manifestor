package anim

import (
	"fmt"
	"os"
	"time"
)

func init() {
	register("portal", "A swirling portal opens and the file emerges", playPortal)
}

func playPortal(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	purple := rgb(147, 51, 234)
	cyan := rgb(6, 182, 212)
	white := rgb(255, 255, 255)

	// Portal frames — concentric rings expanding
	portalFrames := [][]string{
		{
			"        ·        ",
		},
		{
			"       ·°·       ",
			"      ·   ·      ",
			"       ·°·       ",
		},
		{
			"      ·°○°·      ",
			"     ·       ·   ",
			"    ○    ✦    ○  ",
			"     ·       ·   ",
			"      ·°○°·      ",
		},
		{
			"    ·°○◎○°·      ",
			"   ·    ✦    ·   ",
			"  ○   ◎   ◎   ○ ",
			" ·  ✦   ✧   ✦  ·",
			"  ○   ◎   ◎   ○ ",
			"   ·    ✦    ·   ",
			"    ·°○◎○°·      ",
		},
		{
			"  ╭─·°○◎●◎○°·─╮ ",
			"  ·    ◎✦◎    · ",
			" ○  ◎ ✧   ✧ ◎  ○",
			"·  ✦  ✧ ★ ✧  ✦  ·",
			" ○  ◎ ✧   ✧ ◎  ○",
			"  ·    ◎✦◎    · ",
			"  ╰─·°○◎●◎○°·─╯ ",
		},
	}

	// Phase 1: portal opens
	for fi, frame := range portalFrames {
		// Clear previous
		if fi > 0 {
			for range portalFrames[fi-1] {
				fmt.Fprintf(os.Stderr, "\033[A"+clearLn)
			}
		}
		for li, line := range frame {
			var colored string
			for _, r := range line {
				switch r {
				case '★', '✧':
					colored += white + bold + string(r)
				case '●', '◎', '◉', '✦':
					colored += cyan + string(r)
				case '○', '◇', '°', '·':
					colored += purple + string(r)
				case '╭', '╮', '╰', '╯', '─':
					colored += purple + dim + string(r)
				default:
					colored += string(r)
				}
			}
			_ = li
			fmt.Fprintln(os.Stderr, center(colored+reset, w+40))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Phase 2: flash and file emerges
	for range portalFrames[len(portalFrames)-1] {
		fmt.Fprintf(os.Stderr, "\033[A"+clearLn)
	}

	flash := []string{"⚡", "✴", "❇", "✳"}
	for _, f := range flash {
		fmt.Fprint(os.Stderr, clearLn)
		msg := fmt.Sprintf("%s%s %s %s%s", bold, cyan, f, f, reset)
		fmt.Fprintln(os.Stderr, center(msg, w+20))
		time.Sleep(60 * time.Millisecond)
		fmt.Fprintf(os.Stderr, "\033[A"+clearLn)
	}

	// Phase 3: file name materializes
	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🌀 %s%s emerged from the portal%s 🌀", bold+cyan, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+20))
	time.Sleep(400 * time.Millisecond)
}
