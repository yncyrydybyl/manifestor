// Package grab finds and copies the most recent file from ~/Downloads.
package grab

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Latest copies the most recently modified file from ~/Downloads to dest,
// sanitizing the file name. Returns the path of the new file.
func Latest(dest string) (string, error) {
	dir, err := downloadsDir()
	if err != nil {
		return "", err
	}

	src, err := newestFile(dir)
	if err != nil {
		return "", err
	}

	name := Sanitize(filepath.Base(src))
	out := filepath.Join(dest, name)

	if err := copyFile(src, out); err != nil {
		return "", fmt.Errorf("copy %s -> %s: %w", src, out, err)
	}

	return out, nil
}

func downloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	dir := filepath.Join(home, "Downloads")
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("downloads directory not found: %w", err)
	}
	return dir, nil
}

func newestFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", dir, err)
	}

	var newest string
	var newestTime time.Time

	for _, e := range entries {
		if e.IsDir() || e.Name()[0] == '.' {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newest = filepath.Join(dir, e.Name())
		}
	}

	if newest == "" {
		return "", fmt.Errorf("no files found in %s", dir)
	}
	return newest, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
