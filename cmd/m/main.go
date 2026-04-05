package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yncyrydybyl/manifestor/internal/grab"
)

var version = "dev"

func main() {
	args := os.Args[1:]
	force := isForceMode()

	// Parse flags
	var dest string
	for _, a := range args {
		switch a {
		case "--version", "-v":
			fmt.Println("manifestor", version)
			return
		case "--help", "-h":
			printUsage()
			return
		case "--force", "-f":
			force = true
		default:
			if dest == "" {
				dest = a
			}
		}
	}

	if dest == "" {
		dest = "."
	}

	result, err := grab.Find(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if !force && result.Age > grab.StaleThreshold {
		if !confirmStale(result) {
			fmt.Fprintln(os.Stderr, "aborted.")
			os.Exit(1)
		}
	}

	if err := result.Copy(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(result.Dest)
}

// isForceMode returns true if the binary was invoked as "mm".
func isForceMode() bool {
	name := filepath.Base(os.Args[0])
	return name == "mm"
}

func confirmStale(r *grab.Result) bool {
	age := r.Age.Truncate(time.Minute)
	name := filepath.Base(r.Source)

	fmt.Fprintf(os.Stderr, "\n  The newest file in ~/Downloads is %s old:\n", formatAge(age))
	fmt.Fprintf(os.Stderr, "  %s\n\n", name)
	fmt.Fprintf(os.Stderr, "  Hint: this might not be what you just downloaded.\n")
	fmt.Fprintf(os.Stderr, "  Use 'mm' or 'm --force' to skip this check.\n\n")
	fmt.Fprintf(os.Stderr, "  Grab it anyway? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func formatAge(d time.Duration) string {
	hours := int(d.Hours())
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}
	days := hours / 24
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func printUsage() {
	fmt.Print(`manifestor (m) — grab the latest file from ~/Downloads

Usage:
  m [options] [destination]

The most recently modified file in ~/Downloads is copied to the
destination directory (default: current directory) with a sanitized
file name.

If the newest file is older than 8 hours, you'll be asked to confirm.
To skip the check, use --force or invoke as 'mm'.

Options:
  -f, --force      skip the staleness check
  -h, --help       show this help
  -v, --version    show version

Examples:
  m                 copy latest download here
  m ./assets        copy latest download to ./assets
  mm                copy latest download here (no confirmation)
  m --force ~/docs  copy to ~/docs, skip staleness check
`)
}
