// Package anim provides terminal animations for manifestation events.
package anim

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

// Size represents an animation's vertical footprint.
type Size int

const (
	OneLiner   Size = 1
	FiveLiner  Size = 5
	FullScreen Size = 0 // 0 means "use entire terminal height"
)

// ParseSize converts a user-facing string ("1", "5", "full") to a Size.
func ParseSize(s string) (Size, bool) {
	switch s {
	case "1":
		return OneLiner, true
	case "5":
		return FiveLiner, true
	case "full":
		return FullScreen, true
	default:
		return 0, false
	}
}

// String returns the user-facing representation.
func (s Size) String() string {
	switch s {
	case OneLiner:
		return "1"
	case FiveLiner:
		return "5"
	case FullScreen:
		return "full"
	default:
		return "1"
	}
}

// Animation is a named terminal animation.
type Animation struct {
	Name     string
	Desc     string
	Sizes    []Size // supported size modes
	MinWidth int    // minimum terminal width (0 = no minimum)
	HasEmoji bool   // uses emoji/wide unicode
	Play     func(filename string, size Size)
}

// SupportsSize reports whether the animation supports the given size.
func (a *Animation) SupportsSize(s Size) bool {
	for _, sz := range a.Sizes {
		if sz == s {
			return true
		}
	}
	return false
}

// BestSize returns the requested size if supported, otherwise falls back
// to the first supported size.
func (a *Animation) BestSize(requested Size) Size {
	if a.SupportsSize(requested) {
		return requested
	}
	if len(a.Sizes) > 0 {
		return a.Sizes[0]
	}
	return OneLiner
}

var registry []Animation

type regInfo struct {
	name     string
	desc     string
	sizes    []Size
	minWidth int
	hasEmoji bool
	play     func(string, Size)
}

func register(info regInfo) {
	registry = append(registry, Animation{
		Name:     info.name,
		Desc:     info.desc,
		Sizes:    info.sizes,
		MinWidth: info.minWidth,
		HasEmoji: info.hasEmoji,
		Play:     info.play,
	})
}

// Get returns an animation by name, or nil if not found.
func Get(name string) *Animation {
	for i := range registry {
		if registry[i].Name == name {
			return &registry[i]
		}
	}
	return nil
}

// Random returns a random animation.
func Random() *Animation {
	if len(registry) == 0 {
		return nil
	}
	return &registry[rand.Intn(len(registry))]
}

// RandomForSize returns a random animation that supports the given size.
// Falls back to any random animation if none match.
func RandomForSize(s Size) *Animation {
	var matching []*Animation
	for i := range registry {
		if registry[i].SupportsSize(s) {
			matching = append(matching, &registry[i])
		}
	}
	if len(matching) == 0 {
		return Random()
	}
	return matching[rand.Intn(len(matching))]
}

// List returns all registered animation names and descriptions.
func List() []Animation {
	out := make([]Animation, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// termWidth returns the terminal width, defaulting to 80.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || w < 20 {
		return 80
	}
	return w
}

// termHeight returns the terminal height, defaulting to 24.
func termHeight() int {
	_, h, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || h < 5 {
		return 24
	}
	return h
}

// linesForSize returns the number of output lines for a given Size.
func linesForSize(s Size) int {
	switch s {
	case OneLiner:
		return 1
	case FiveLiner:
		return 5
	case FullScreen:
		return termHeight() - 1 // leave room for prompt
	default:
		return 1
	}
}

// moveUp prints ANSI escape to move cursor up n lines.
func moveUp(n int) {
	if n > 0 {
		fmt.Fprintf(os.Stderr, "\033[%dA", n)
	}
}

// clearLines clears n lines above the cursor (inclusive of current line).
func clearLines(n int) {
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(os.Stderr, "\033[A")
		}
		fmt.Fprint(os.Stderr, clearLn)
	}
}

// ANSI helpers
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	hide    = "\033[?25l" // hide cursor
	show    = "\033[?25h" // show cursor
	clearLn = "\033[2K\r"
)

func fg(code int) string     { return fmt.Sprintf("\033[38;5;%dm", code) }
func rgb(r, g, b int) string { return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b) }

// center pads a string to center it in the terminal.
func center(s string, width int) string {
	// Approximate visible length (strips aren't perfect but good enough)
	visible := visibleLen(s)
	if visible >= width {
		return s
	}
	pad := (width - visible) / 2
	return strings.Repeat(" ", pad) + s
}

// visibleLen estimates the visible length of a string with ANSI codes.
func visibleLen(s string) int {
	inEsc := false
	n := 0
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}
