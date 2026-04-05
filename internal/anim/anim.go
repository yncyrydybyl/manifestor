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

// Animation is a named terminal animation.
type Animation struct {
	Name string
	Desc string
	Play func(filename string)
}

var registry []Animation

func register(name, desc string, play func(string)) {
	registry = append(registry, Animation{Name: name, Desc: desc, Play: play})
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

// List returns all registered animation names and descriptions.
func List() []Animation {
	out := make([]Animation, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// termWidth returns the terminal width, defaulting to 60.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || w < 20 {
		return 60
	}
	return w
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

func fg(code int) string    { return fmt.Sprintf("\033[38;5;%dm", code) }
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
