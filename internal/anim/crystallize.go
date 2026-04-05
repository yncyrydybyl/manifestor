package anim

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func init() {
	register("crystallize", "Crystals form and shatter to reveal the file", playCrystallize)
}

func playCrystallize(filename string) {
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
			// Crystals grow from center outward
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
