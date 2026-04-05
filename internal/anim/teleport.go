package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register("teleport", "Particles scatter and reassemble at destination", playTeleport)
}

func playTeleport(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	particles := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	glyphs := []string{"◌", "◍", "◎", "●", "◉", "◈", "◆", "◇", "⬡", "⬢"}
	colors := []int{51, 45, 39, 33, 27, 21, 57, 93, 129, 165}

	// Phase 1: source dissolves into particles
	srcLabel := "~/Downloads"
	for frame := 0; frame < 12; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		progress := float64(frame) / 12.0

		for i := 0; i < w; i++ {
			pos := float64(i) / float64(w)
			if pos < 0.3 {
				// Source side — dissolving
				if rand.Float64() < progress {
					p := particles[rand.Intn(len(particles))]
					line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
				} else if i < len(srcLabel) {
					line.WriteString(dim + fg(245) + string(srcLabel[i]))
				} else {
					line.WriteString(" ")
				}
			} else if pos < 0.7 {
				// Middle — particle stream
				if rand.Float64() < progress*0.4 {
					g := glyphs[rand.Intn(len(glyphs))]
					c := colors[(i+frame)%len(colors)]
					line.WriteString(fg(c) + g)
				} else {
					line.WriteString(" ")
				}
			} else {
				line.WriteString(" ")
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 2: particles converge to destination
	for frame := 0; frame < 12; frame++ {
		fmt.Fprint(os.Stderr, clearLn)
		var line strings.Builder
		progress := float64(frame) / 12.0

		for i := 0; i < w; i++ {
			pos := float64(i) / float64(w)
			if pos < 0.3 {
				line.WriteString(" ")
			} else if pos < 0.7 {
				// Middle — particles moving right
				if rand.Float64() < (1.0-progress)*0.4 {
					g := glyphs[rand.Intn(len(glyphs))]
					c := colors[(i+frame)%len(colors)]
					line.WriteString(fg(c) + g)
				} else {
					line.WriteString(" ")
				}
			} else {
				// Destination side — materializing
				dstStart := int(float64(w) * 0.7)
				di := i - dstStart
				if di < len(filename) && rand.Float64() < progress {
					c := colors[(di+frame)%len(colors)]
					line.WriteString(bold + fg(c) + string(filename[di]))
				} else if rand.Float64() < (1.0-progress)*0.3 {
					p := particles[rand.Intn(len(particles))]
					line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
				} else {
					line.WriteString(" ")
				}
			}
		}
		fmt.Fprint(os.Stderr, line.String()+reset)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: final reveal
	fmt.Fprint(os.Stderr, clearLn)
	arrow := " ⚡ ~/Downloads  ━━━━━━━━▸  ./ "
	msg := fmt.Sprintf("%s%s%s%s%s", fg(51), arrow, bold, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+30))
	time.Sleep(400 * time.Millisecond)
}
