package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register("emoji-rain", "Cascading emoji rain with lightning reveal", playEmojiRain)
}

func playEmojiRain(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	drops := []string{"🌟", "⭐", "💫", "✨", "🔮", "🌈", "🦋", "🌸", "💜", "💙", "💚", "💛", "🧡", "❤️"}

	// Phase 1: emoji rain builds up
	for frame := 0; frame < 20; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		density := float64(frame) / 20.0

		for i := 0; i < w/2; i++ { // divide by 2 because emoji are double-width
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
