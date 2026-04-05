package main

import (
	"fmt"
	"os"

	"github.com/yncyrydybyl/manifestor/internal/grab"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Println("manifestor", version)
			return
		case "--help", "-h":
			printUsage()
			return
		}
	}

	dest := "."
	if len(os.Args) > 1 && os.Args[1] != "" {
		dest = os.Args[1]
	}

	result, err := grab.Latest(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func printUsage() {
	fmt.Print(`manifestor - grab the latest file from ~/Downloads

Usage:
  manifestor [destination]

The most recently modified file in ~/Downloads is copied to the
destination directory (default: current directory) with a sanitized
file name.

Sanitization: lowercased, spaces and special characters replaced
with hyphens, consecutive hyphens collapsed, leading/trailing
hyphens removed.

Options:
  -h, --help       show this help
  -v, --version    show version

Examples:
  manifestor              copy latest download here
  manifestor ./assets     copy latest download to ./assets
  manifestor ~/Documents  copy latest download to ~/Documents
`)
}
