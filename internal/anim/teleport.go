package anim

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	register(regInfo{
		name:     "teleport",
		desc:     "Particles scatter and reassemble at destination",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: false,
		play:     playTeleport,
	})
}

func playTeleport(filename string, size Size) {
	switch size {
	case FiveLiner:
		playTeleportFive(filename)
	case FullScreen:
		playTeleportFull(filename)
	default:
		playTeleportOne(filename)
	}
}

func playTeleportOne(filename string) {
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
				if rand.Float64() < progress {
					p := particles[rand.Intn(len(particles))]
					line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
				} else if i < len(srcLabel) {
					line.WriteString(dim + fg(245) + string(srcLabel[i]))
				} else {
					line.WriteString(" ")
				}
			} else if pos < 0.7 {
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
				if rand.Float64() < (1.0-progress)*0.4 {
					g := glyphs[rand.Intn(len(glyphs))]
					c := colors[(i+frame)%len(colors)]
					line.WriteString(fg(c) + g)
				} else {
					line.WriteString(" ")
				}
			} else {
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

func playTeleportFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	particles := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	glyphs := []string{"◌", "◍", "◎", "●", "◉", "◈", "◆", "◇", "⬡", "⬢"}
	colors := []int{51, 45, 39, 33, 27, 21, 57, 93, 129, 165}

	// Reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: source dissolves — particles scatter vertically
	srcLabel := "~/Downloads"
	for frame := 0; frame < 12; frame++ {
		progress := float64(frame) / 12.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				pos := float64(i) / float64(w)
				if pos < 0.3 {
					// Source side — dissolving with vertical scatter
					if r == h/2 && !(rand.Float64() < progress) && i < len(srcLabel) {
						line.WriteString(dim + fg(245) + string(srcLabel[i]))
					} else if rand.Float64() < progress*0.5 {
						p := particles[rand.Intn(len(particles))]
						line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
					} else {
						line.WriteString(" ")
					}
				} else if pos < 0.7 {
					// Middle — vertical particle paths
					if rand.Float64() < progress*0.3 {
						g := glyphs[rand.Intn(len(glyphs))]
						c := colors[(i+frame+r)%len(colors)]
						line.WriteString(fg(c) + g)
					} else {
						line.WriteString(" ")
					}
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 2: particles converge to destination across all rows
	for frame := 0; frame < 12; frame++ {
		progress := float64(frame) / 12.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				pos := float64(i) / float64(w)
				if pos < 0.3 {
					line.WriteString(" ")
				} else if pos < 0.7 {
					if rand.Float64() < (1.0-progress)*0.3 {
						g := glyphs[rand.Intn(len(glyphs))]
						c := colors[(i+frame+r)%len(colors)]
						line.WriteString(fg(c) + g)
					} else {
						line.WriteString(" ")
					}
				} else {
					// Destination side — materializing with vertical convergence
					dstStart := int(float64(w) * 0.7)
					di := i - dstStart
					// Converge to middle row
					rowDist := float64(r-h/2) / float64(h/2)
					if rowDist < 0 {
						rowDist = -rowDist
					}
					if r == h/2 && di < len(filename) && rand.Float64() < progress {
						c := colors[(di+frame)%len(colors)]
						line.WriteString(bold + fg(c) + string(filename[di]))
					} else if rand.Float64() < (1.0-progress)*(1.0-rowDist)*0.4 {
						p := particles[rand.Intn(len(particles))]
						line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
					} else {
						line.WriteString(" ")
					}
				}
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(60 * time.Millisecond)
	}

	// Phase 3: clear and show final message on one line
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	fmt.Fprint(os.Stderr, clearLn)
	arrow := " ⚡ ~/Downloads  ━━━━━━━━▸  ./ "
	msg := fmt.Sprintf("%s%s%s%s%s", fg(51), arrow, bold, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+30))
	time.Sleep(400 * time.Millisecond)
}

func playTeleportFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	particles := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	glyphs := []string{"◌", "◍", "◎", "●", "◉", "◈", "◆", "◇", "⬡", "⬢"}
	colors := []int{51, 45, 39, 33, 27, 21, 57, 93, 129, 165}

	// Reserve space
	for i := 0; i < h; i++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: full-screen dissolve — source text scatters into particles
	srcLabel := "~/Downloads"
	midRow := h / 2
	for frame := 0; frame < 15; frame++ {
		progress := float64(frame) / 15.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				pos := float64(i) / float64(w)
				if pos < 0.3 {
					if r == midRow && !(rand.Float64() < progress) && i < len(srcLabel) {
						line.WriteString(dim + fg(245) + string(srcLabel[i]))
					} else if rand.Float64() < progress*0.3 {
						p := particles[rand.Intn(len(particles))]
						line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
					} else {
						line.WriteString(" ")
					}
				} else if pos < 0.7 {
					// Particle stream fills the middle zone
					density := progress * 0.25
					if rand.Float64() < density {
						g := glyphs[rand.Intn(len(glyphs))]
						c := colors[(i+frame+r)%len(colors)]
						line.WriteString(fg(c) + g)
					} else {
						line.WriteString(" ")
					}
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: particles converge and reassemble across the full screen
	for frame := 0; frame < 15; frame++ {
		progress := float64(frame) / 15.0
		for r := 0; r < h; r++ {
			fmt.Fprint(os.Stderr, clearLn)
			var line strings.Builder
			for i := 0; i < w; i++ {
				pos := float64(i) / float64(w)
				if pos < 0.3 {
					line.WriteString(" ")
				} else if pos < 0.7 {
					if rand.Float64() < (1.0-progress)*0.25 {
						g := glyphs[rand.Intn(len(glyphs))]
						c := colors[(i+frame+r)%len(colors)]
						line.WriteString(fg(c) + g)
					} else {
						line.WriteString(" ")
					}
				} else {
					// Destination — reassemble from edges to center
					dstStart := int(float64(w) * 0.7)
					di := i - dstStart
					rowDist := float64(r-midRow) / float64(h/2)
					if rowDist < 0 {
						rowDist = -rowDist
					}
					if r == midRow && di < len(filename) && rand.Float64() < progress {
						c := colors[(di+frame)%len(colors)]
						line.WriteString(bold + fg(c) + string(filename[di]))
					} else if rand.Float64() < (1.0-progress)*(1.0-rowDist*0.5)*0.3 {
						p := particles[rand.Intn(len(particles))]
						line.WriteString(fg(colors[rand.Intn(len(colors))]) + p)
					} else {
						line.WriteString(" ")
					}
				}
			}
			fmt.Fprint(os.Stderr, line.String()+reset)
			fmt.Fprintln(os.Stderr)
		}
		moveUp(h)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 3: clear and show final message
	for r := 0; r < h; r++ {
		fmt.Fprint(os.Stderr, clearLn)
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)
	fmt.Fprint(os.Stderr, clearLn)
	arrow := " ⚡ ~/Downloads  ━━━━━━━━▸  ./ "
	msg := fmt.Sprintf("%s%s%s%s%s", fg(51), arrow, bold, filename, reset)
	fmt.Fprintln(os.Stderr, center(msg, w+30))
	time.Sleep(400 * time.Millisecond)
}
