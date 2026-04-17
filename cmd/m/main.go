package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yncyrydybyl/manifestor/internal/anim"
	"github.com/yncyrydybyl/manifestor/internal/completion"
	"github.com/yncyrydybyl/manifestor/internal/grab"
	"golang.org/x/term"
)

var version = "dev"

func main() {
	args := os.Args[1:]
	force := isForceMode()

	var dest string
	var animName string
	noAnim := false
	listAnims := false
	count := 1

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--version", "-v":
			fmt.Println("manifestor", version)
			return
		case "--help", "-h":
			printUsage()
			return
		case "--force", "-f":
			force = true
		case "--no-anim":
			noAnim = true
		case "--anim":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				fmt.Fprintln(os.Stderr, "error: --anim requires an animation name (use --list-anims to see available)")
				os.Exit(1)
			}
			i++
			animName = args[i]
		case "--list-anims":
			listAnims = true
		case "completion":
			if i+1 < len(args) {
				printCompletion(args[i+1])
			} else {
				fmt.Fprintln(os.Stderr, "usage: m completion <bash|zsh|fish>")
				os.Exit(1)
			}
			return
		default:
			if n, ok := parseCount(args[i]); ok {
				if err := validateCount(n); err != nil {
					fmt.Fprintf(os.Stderr, "error: %s\n", err)
					os.Exit(1)
				}
				count = n
				continue
			}
			if dest == "" {
				dest = args[i]
			}
		}
	}

	if listAnims {
		printAnims()
		return
	}

	if dest == "" {
		dest = "."
	}

	if count > 1 {
		runBatch(count, dest, force)
		return
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

	// Play animation after successful copy (only if stderr is a TTY)
	if !noAnim && term.IsTerminal(int(os.Stderr.Fd())) {
		name := filepath.Base(result.Dest)
		playAnimation(animName, name)
	}

	fmt.Println(result.Dest)
}

// parseCount recognises "-N" where N is a positive integer. Returns (n, true)
// on match, else (0, false). It does not validate range — callers should.
func parseCount(arg string) (int, bool) {
	if len(arg) < 2 || arg[0] != '-' {
		return 0, false
	}
	n, err := strconv.Atoi(arg[1:])
	if err != nil {
		return 0, false
	}
	return n, true
}

func validateCount(n int) error {
	if n < 1 {
		return fmt.Errorf("invalid count: -%d (must be >= 1)", n)
	}
	if n > grab.MaxBatch {
		return fmt.Errorf("-%d too large (max %d)", n, grab.MaxBatch)
	}
	return nil
}

func runBatch(n int, dest string, force bool) {
	results, err := grab.FindN(n, dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if !force && len(results) > 0 && results[0].Age > grab.StaleThreshold {
		if !confirmStale(results[0]) {
			fmt.Fprintln(os.Stderr, "aborted.")
			os.Exit(1)
		}
	}

	copied := 0
	for _, r := range results {
		if err := r.Copy(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
		fmt.Println(r.Dest)
		copied++
	}

	if copied < n {
		fmt.Fprintf(os.Stderr, "grabbed %d of %d requested (only %d available)\n", copied, n, len(results))
	}
}

func playAnimation(name, filename string) {
	var a *anim.Animation
	if name != "" {
		a = anim.Get(name)
		if a == nil {
			fmt.Fprintf(os.Stderr, "unknown animation: %s (use --list-anims to see available)\n", name)
			return
		}
	} else {
		a = anim.Random()
	}
	if a != nil {
		a.Play(filename)
	}
}

func printCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Print(completion.Bash())
	case "zsh":
		fmt.Print(completion.Zsh())
	case "fish":
		fmt.Print(completion.Fish())
	default:
		fmt.Fprintf(os.Stderr, "unknown shell: %s (supported: bash, zsh, fish)\n", shell)
		os.Exit(1)
	}
}

func printAnims() {
	fmt.Println("Available animations:")
	fmt.Println()
	for _, a := range anim.List() {
		fmt.Printf("  %-20s %s\n", a.Name, a.Desc)
	}
	fmt.Println()
	fmt.Println("Usage: m --anim <name>")
	fmt.Println("       m --no-anim       (skip animation)")
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
file name. A random manifestation animation plays on success.

If the newest file is older than 8 hours, you'll be asked to confirm.
To skip the check, use --force or invoke as 'mm'.

Options:
  -N                 grab the N newest files (1-254); e.g. -5, -200
  -f, --force        skip the staleness check
  --anim <name>      play a specific animation
  --no-anim          skip the animation
  --list-anims       show available animations
  -h, --help         show this help
  -v, --version      show version

Commands:
  completion <shell>   generate shell completions (bash, zsh, fish)

Examples:
  m                          copy latest download here
  m ./assets                 copy latest download to ./assets
  m -5                       copy the 5 newest downloads here
  m -200 ./img               copy the 200 newest downloads to ./img
  mm                         force mode (no confirmation)
  m --anim rainbow-beam      use a specific animation
  m --anim fire-forge .      forge it in flames
  m --no-anim                just copy, no flair

Shell completions:
  eval "$(m completion bash)"           # bash
  eval "$(m completion zsh)"            # zsh
  m completion fish > ~/.config/fish/completions/m.fish  # fish
`)
}
