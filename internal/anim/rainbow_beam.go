package anim

import (
	"fmt"
	"os"
	"time"
)

func init() {
	register("rainbow-beam", "Full-width rainbow bar with sparkle burst", playRainbowBeam)
}

func playRainbowBeam(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	blocks := []string{"░", "▒", "▓", "█"}
	rainbow := []int{196, 208, 220, 46, 51, 21, 93} // red orange yellow green cyan blue violet

	// Phase 1: beam fills left to right
	for i := 0; i <= w; i++ {
		fmt.Fprint(os.Stderr, clearLn)
		for j := 0; j < i; j++ {
			ci := (j * len(rainbow) / w) % len(rainbow)
			bi := 3 // full block for filled area
			if j == i-1 {
				bi = 0 // leading edge is light
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
