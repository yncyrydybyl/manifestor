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
		name:     "starfield",
		desc:     "Warp-speed starfield with cosmic arrival",
		sizes:    []Size{OneLiner, FiveLiner, FullScreen},
		minWidth: 40,
		hasEmoji: true,
		play:     playStarfield,
	})
}

var starGlyphs = []string{".", "·", "∗", "⋆", "✦", "★", "✧"}
var starColors = []int{255, 253, 251, 249, 247, 245, 243, 241, 239}

func playStarfield(filename string, size Size) {
	switch size {
	case FiveLiner:
		playStarfieldFive(filename)
	case FullScreen:
		playStarfieldFull(filename)
	default:
		playStarfieldOne(filename)
	}
}

// playStarfieldOne is the original single-line warp streaks animation.
func playStarfieldOne(filename string) {
	w := termWidth()
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

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
			if rand.Float64() < dist*speed*0.8 {
				si := int(dist * float64(len(starGlyphs)-1))
				if si >= len(starGlyphs) {
					si = len(starGlyphs) - 1
				}
				ci := int((1.0 - dist) * float64(len(starColors)-1))
				if ci < 0 {
					ci = 0
				}
				if speed > 0.5 && rand.Float64() < speed*0.3 {
					line.WriteString(fg(starColors[ci]) + "━")
				} else {
					line.WriteString(fg(starColors[ci]) + starGlyphs[si])
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

// playStarfieldFive renders a 5-line starfield with depth layers.
func playStarfieldFive(filename string) {
	w := termWidth()
	h := 5
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	// Each row is a depth layer: row 0 = far (dim, small), row 4 = near (bright, big)
	type star struct {
		x     int
		glyph int // index into starGlyphs
		speed int // pixels per frame to move outward from center
	}

	// Seed stars per layer
	layers := make([][]star, h)
	for r := 0; r < h; r++ {
		depth := float64(r) / float64(h-1) // 0=far, 1=near
		count := 5 + int(depth*15)
		for s := 0; s < count; s++ {
			gi := int(depth * float64(len(starGlyphs)-1))
			layers[r] = append(layers[r], star{
				x:     rand.Intn(w),
				glyph: gi,
				speed: 1 + int(depth*3),
			})
		}
	}

	// Print initial blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	mid := w / 2

	// Phase 1: stars streak outward
	for frame := 0; frame < 20; frame++ {
		speed := float64(frame) / 20.0
		for r := 0; r < h; r++ {
			buf := make([]byte, w)
			for i := range buf {
				buf[i] = ' '
			}
			var line strings.Builder
			depth := float64(r) / float64(h-1)
			ci := int((1.0 - depth*0.5) * float64(len(starColors)-1))
			if ci >= len(starColors) {
				ci = len(starColors) - 1
			}
			// Place stars
			for si := range layers[r] {
				s := &layers[r][si]
				if s.x >= 0 && s.x < w {
					buf[s.x] = 1 // mark occupied
				}
				// Move away from center
				if s.x < mid {
					s.x -= s.speed
				} else {
					s.x += s.speed
				}
				// Respawn at center if out of bounds
				if s.x < 0 || s.x >= w {
					s.x = mid + rand.Intn(5) - 2
				}
			}
			for i := 0; i < w; i++ {
				if buf[i] == 1 {
					if speed > 0.5 && depth > 0.5 && rand.Float64() < speed*0.4 {
						line.WriteString(fg(starColors[ci]) + "━")
					} else {
						gi := int(depth * float64(len(starGlyphs)-1))
						line.WriteString(fg(starColors[ci]) + starGlyphs[gi])
					}
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(50 * time.Millisecond)
	}

	// Phase 2: hyperspeed (all lines become streaks)
	for frame := 0; frame < 6; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			density := 0.3 + float64(r)/float64(h-1)*0.5
			for i := 0; i < w; i++ {
				if rand.Float64() < density {
					c := starColors[rand.Intn(3)]
					line.WriteString(bold + fg(c) + "━")
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(30 * time.Millisecond)
	}

	// Final
	clearLines(h)
	msg := fmt.Sprintf("🚀 %s arrived from across the cosmos 🌌", filename)
	fmt.Fprintln(os.Stderr, center(bold+fg(255)+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}

// playStarfieldFull renders a full-screen warp tunnel.
func playStarfieldFull(filename string) {
	w := termWidth()
	h := linesForSize(FullScreen)
	fmt.Fprint(os.Stderr, hide)
	defer fmt.Fprint(os.Stderr, show)

	midX := w / 2
	midY := h / 2

	type star struct {
		x, y  float64 // position relative to center (-1..1)
		speed float64
	}

	numStars := w * h / 8
	if numStars > 500 {
		numStars = 500
	}
	stars := make([]star, numStars)
	for i := range stars {
		stars[i] = star{
			x:     (rand.Float64() - 0.5) * 2,
			y:     (rand.Float64() - 0.5) * 2,
			speed: 0.02 + rand.Float64()*0.06,
		}
	}

	// Print initial blank lines
	for r := 0; r < h; r++ {
		fmt.Fprintln(os.Stderr)
	}
	moveUp(h)

	// Phase 1: warp tunnel — stars fly outward from center
	for frame := 0; frame < 30; frame++ {
		warpSpeed := float64(frame) / 30.0
		// Build a 2D grid
		grid := make([][]byte, h)
		for r := range grid {
			grid[r] = make([]byte, w)
		}

		for si := range stars {
			s := &stars[si]
			// Move outward from center
			s.x *= 1.0 + s.speed*(1+warpSpeed*3)
			s.y *= 1.0 + s.speed*(1+warpSpeed*3)

			// Convert to screen coords
			sx := midX + int(s.x*float64(midX))
			sy := midY + int(s.y*float64(midY))

			if sx >= 0 && sx < w && sy >= 0 && sy < h {
				dist := s.x*s.x + s.y*s.y
				if dist < 0.1 {
					grid[sy][sx] = 1 // dim/small
				} else if dist < 0.5 {
					grid[sy][sx] = 2 // medium
				} else {
					grid[sy][sx] = 3 // bright/big
				}
			}

			// Respawn if out of view
			if sx < 0 || sx >= w || sy < 0 || sy >= h {
				s.x = (rand.Float64() - 0.5) * 0.3
				s.y = (rand.Float64() - 0.5) * 0.3
				s.speed = 0.02 + rand.Float64()*0.06
			}
		}

		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				switch grid[r][i] {
				case 1:
					line.WriteString(dim + fg(starColors[len(starColors)-1]) + ".")
				case 2:
					ci := rand.Intn(3) + 3
					line.WriteString(fg(starColors[ci]) + starGlyphs[2+rand.Intn(2)])
				case 3:
					if warpSpeed > 0.5 && rand.Float64() < warpSpeed*0.5 {
						line.WriteString(bold + fg(starColors[0]) + "━")
					} else {
						line.WriteString(bold + fg(starColors[rand.Intn(2)]) + starGlyphs[4+rand.Intn(3)])
					}
				default:
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(45 * time.Millisecond)
	}

	// Phase 2: hyperspeed flash
	for frame := 0; frame < 5; frame++ {
		for r := 0; r < h; r++ {
			var line strings.Builder
			for i := 0; i < w; i++ {
				if rand.Float64() < 0.7 {
					c := starColors[rand.Intn(3)]
					line.WriteString(bold + fg(c) + "━")
				} else {
					line.WriteString(" ")
				}
			}
			fmt.Fprint(os.Stderr, clearLn+line.String()+reset)
			if r < h-1 {
				fmt.Fprintln(os.Stderr)
			}
		}
		moveUp(h - 1)
		time.Sleep(30 * time.Millisecond)
	}

	// Final
	clearLines(h)
	msg := fmt.Sprintf("🚀 %s arrived from across the cosmos 🌌", filename)
	fmt.Fprintln(os.Stderr, center(bold+fg(255)+msg+reset, w+20))
	time.Sleep(400 * time.Millisecond)
}
