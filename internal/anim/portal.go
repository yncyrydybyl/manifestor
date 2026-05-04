package anim

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func init() {
	register(regInfo{
		name:     "portal",
		desc:     "A swirling portal opens and the file emerges",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 30,
		hasEmoji: true,
		play:     playPortal,
	})
}

func playPortal(filename string, size Size) {
	switch size {
	case OneLiner:
		playPortalOneLiner(filename)
	case FiveLiner:
		playPortalFiveLiner(filename)
	case FullScreen:
		playPortalFullScreen(filename)
	default:
		playPortalFiveLiner(filename)
	}
}

func playPortalOneLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	purple := rgb(147, 51, 234)
	cyan := rgb(6, 182, 212)
	white := rgb(255, 255, 255)

	frames := []string{
		purple + "·" + reset,
		purple + "·°·" + reset,
		purple + "·°" + cyan + "○" + purple + "°·" + reset,
		purple + "·°" + cyan + "○◎○" + purple + "°·" + reset,
		purple + "·°" + cyan + "○◎" + white + bold + "✦" + reset + cyan + "◎○" + purple + "°·" + reset,
		purple + "╭·°" + cyan + "○◎" + white + bold + "★" + reset + cyan + "◎○" + purple + "°·╮" + reset,
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
		time.Sleep(150 * time.Millisecond)
	}

	// Flash
	flashes := []string{"⚡", "✴", "❇", "✳"}
	for _, f := range flashes {
		fmt.Fprint(os.Stderr, clearLn)
		s := bold + cyan + f + " " + f + " " + f + reset
		vis := visibleLen(s)
		if vis < w {
			pad := (w - vis) / 2
			fmt.Fprint(os.Stderr, strings.Repeat(" ", pad)+s)
		} else {
			fmt.Fprint(os.Stderr, s)
		}
		time.Sleep(60 * time.Millisecond)
	}

	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🌀 %s%s emerged from the portal%s 🌀", bold+cyan, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+20))
	time.Sleep(400 * time.Millisecond)
}

func playPortalFiveLiner(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	purple := rgb(147, 51, 234)
	cyan := rgb(6, 182, 212)
	white := rgb(255, 255, 255)

	type frame5 [5]string

	frames := []frame5{
		{
			"",
			"",
			"        ·        ",
			"",
			"",
		},
		{
			"",
			"       ·°·       ",
			"      ·   ·      ",
			"       ·°·       ",
			"",
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
			"   ·    ✦    ·   ",
			"    ·°○◎○°·      ",
		},
		{
			"  ╭─·°○◎●◎○°·─╮ ",
			"  ·    ◎✦◎    · ",
			" ○  ◎ ✧   ✧ ◎  ○",
			"  ·    ◎✦◎    · ",
			"  ╰─·°○◎●◎○°·─╯ ",
		},
	}

	colorPortalLine := func(line string) string {
		var colored strings.Builder
		for _, r := range line {
			switch r {
			case '★', '✧':
				colored.WriteString(white + bold + string(r))
			case '●', '◎', '◉', '✦':
				colored.WriteString(cyan + string(r))
			case '○', '◇', '°', '·':
				colored.WriteString(purple + string(r))
			case '╭', '╮', '╰', '╯', '─':
				colored.WriteString(purple + dim + string(r))
			default:
				colored.WriteString(string(r))
			}
		}
		return colored.String()
	}

	for i, f := range frames {
		if i > 0 {
			moveUp(5)
		}
		for _, line := range f {
			fmt.Fprint(os.Stderr, clearLn)
			if line == "" {
				fmt.Fprintln(os.Stderr)
				continue
			}
			fmt.Fprintln(os.Stderr, center(colorPortalLine(line)+reset, w+40))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Flash
	moveUp(5)
	clearLines(5)
	flashes := []string{"⚡", "✴", "❇", "✳"}
	for _, f := range flashes {
		fmt.Fprint(os.Stderr, clearLn)
		s := fmt.Sprintf("%s%s %s %s%s", bold, cyan, f, f, reset)
		fmt.Fprintln(os.Stderr, center(s, w+20))
		time.Sleep(60 * time.Millisecond)
		moveUp(1)
	}

	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🌀 %s%s emerged from the portal%s 🌀", bold+cyan, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+20))
	time.Sleep(400 * time.Millisecond)
}

func playPortalFullScreen(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	purple := rgb(147, 51, 234)
	cyan := rgb(6, 182, 212)
	white := rgb(255, 255, 255)

	portalStages := [][]string{
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
		{
			" ╭──·°○◎●◎●◎○°·──╮ ",
			" · ○  ◎ ✧   ✧ ◎  ○·",
			"·  ✦  ✧  ◎  ✧  ✦  ·",
			"○ ◎ ✧  ✦  ★  ✦  ✧ ◎○",
			"·  ✦  ✧  ◎  ✧  ✦  ·",
			" · ○  ◎ ✧   ✧ ◎  ○·",
			" ╰──·°○◎●◎●◎○°·──╯ ",
			"       ·°○°·         ",
			"        ·°·          ",
		},
	}

	colorPortalLine := func(line string) string {
		var colored strings.Builder
		for _, r := range line {
			switch r {
			case '★', '✧':
				colored.WriteString(white + bold + string(r))
			case '●', '◎', '◉', '✦':
				colored.WriteString(cyan + string(r))
			case '○', '◇', '°', '·':
				colored.WriteString(purple + string(r))
			case '╭', '╮', '╰', '╯', '─':
				colored.WriteString(purple + dim + string(r))
			default:
				colored.WriteString(string(r))
			}
		}
		return colored.String()
	}

	for i, stage := range portalStages {
		if i > 0 {
			moveUp(h)
		}
		// Center the portal vertically
		blank := (h - len(stage)) / 2
		for j := 0; j < h; j++ {
			fmt.Fprint(os.Stderr, clearLn)
			idx := j - blank
			if idx >= 0 && idx < len(stage) {
				fmt.Fprintln(os.Stderr, center(colorPortalLine(stage[idx])+reset, w+40))
			} else {
				fmt.Fprintln(os.Stderr)
			}
		}
		time.Sleep(250 * time.Millisecond)
	}

	// Flash
	moveUp(h)
	for j := 0; j < h; j++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	flashes := []string{"⚡", "✴", "❇", "✳"}
	for _, f := range flashes {
		fmt.Fprint(os.Stderr, clearLn)
		s := fmt.Sprintf("%s%s %s %s%s", bold, cyan, f, f, reset)
		fmt.Fprintln(os.Stderr, center(s, w+20))
		time.Sleep(60 * time.Millisecond)
		moveUp(1)
	}

	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🌀 %s%s emerged from the portal%s 🌀", bold+cyan, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+20))
	time.Sleep(400 * time.Millisecond)
}
