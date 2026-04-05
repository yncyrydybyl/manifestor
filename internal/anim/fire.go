package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register("fire-forge", "File is forged in flames", playFireForge)
}

func playFireForge(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	flames := []string{"▀", "▄", "█", "░", "▒", "▓", "^", "~"}
	fireColors := []int{52, 88, 124, 160, 196, 202, 208, 214, 220, 226, 228, 230} // dark red -> yellow -> white

	// Phase 1: flames build
	for frame := 0; frame < 18; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		intensity := float64(frame) / 18.0

		for i := 0; i < w; i++ {
			if rand.Float64() < intensity*0.8 {
				f := flames[rand.Intn(len(flames))]
				// Hotter in center
				distFromCenter := float64(i-w/2) / float64(w/2)
				if distFromCenter < 0 {
					distFromCenter = -distFromCenter
				}
				heat := int((1.0 - distFromCenter) * float64(len(fireColors)-1))
				if heat < 0 {
					heat = 0
				}
				// Add some randomness to heat
				heat += rand.Intn(3) - 1
				if heat < 0 {
					heat = 0
				}
				if heat >= len(fireColors) {
					heat = len(fireColors) - 1
				}
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
			f := flames[rand.Intn(len(flames))]
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
				// Cool down from white-hot to orange
				heat := int((1.0 - progress*0.5) * float64(len(fireColors)-1))
				line.WriteString(bold + fg(fireColors[heat]) + string(msgRunes[i]))
			} else if rand.Float64() < (1.0-progress)*0.5 {
				f := flames[rand.Intn(len(flames))]
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
