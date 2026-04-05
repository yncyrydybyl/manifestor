package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register("starfield", "Warp-speed starfield with cosmic arrival", playStarfield)
}

func playStarfield(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	stars := []string{".", "·", "∗", "⋆", "✦", "★", "✧"}
	starColors := []int{255, 253, 251, 249, 247, 245, 243, 241, 239}

	// Phase 1: stars streak from center outward
	for frame := 0; frame < 20; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		speed := float64(frame) / 20.0

		mid := w / 2
		for i := 0; i < w; i++ {
			dist := float64(i-mid) / float64(mid)
			if dist < 0 {
				dist = -dist
			}
			// Stars appear more at edges (they've traveled further)
			if rand.Float64() < dist*speed*0.8 {
				si := int(dist * float64(len(stars)-1))
				if si >= len(stars) {
					si = len(stars) - 1
				}
				ci := int((1.0 - dist) * float64(len(starColors)-1))
				if ci < 0 {
					ci = 0
				}
				// At high speed, stars become streaks
				if speed > 0.5 && rand.Float64() < speed*0.3 {
					if i < mid {
						line.WriteString(fg(starColors[ci]) + "━")
					} else {
						line.WriteString(fg(starColors[ci]) + "━")
					}
				} else {
					line.WriteString(fg(starColors[ci]) + stars[si])
				}
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(40 * time.Millisecond)
	}

	// Phase 2: hyperspeed lines
	for frame := 0; frame < 6; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		for i := 0; i < w; i++ {
			if rand.Float64() < 0.7 {
				c := starColors[rand.Intn(3)]
				line.WriteString(bold + fg(c) + "━")
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(30 * time.Millisecond)
	}

	// Phase 3: arrival
	fmt.Fprint(os.Stderr, clearLn)
	msg := fmt.Sprintf("🚀 %s arrived from across the cosmos 🌌", filename)
	fmt.Fprintln(os.Stderr, center(bold+fg(255)+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}
