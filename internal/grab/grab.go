// Package grab finds and copies the most recent file from ~/Downloads.
package grab

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Result holds the newest file info before copying.
type Result struct {
	Source string
	Dest   string
	Age    time.Duration
}

// Find locates the most recently modified file in ~/Downloads and returns
// the source path, intended destination, and age of the file.
func Find(dest string) (*Result, error) {
	dir, err := downloadsDir()
	if err != nil {
		return nil, err
	}

	src, modTime, err := newestFile(dir)
	if err != nil {
		return nil, err
	}

	name := Sanitize(filepath.Base(src))
	out := filepath.Join(dest, name)

	return &Result{
		Source: src,
		Dest:   out,
		Age:    time.Since(modTime),
	}, nil
}

// Copy performs the actual file copy for a Result.
func (r *Result) Copy() error {
	if err := copyFile(r.Source, r.Dest); err != nil {
		return fmt.Errorf("copy %s -> %s: %w", r.Source, r.Dest, err)
	}
	return nil
}

// StaleThreshold is the age after which a file is considered stale.
const StaleThreshold = 8 * time.Hour

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

func newestFile(dir string) (string, time.Time, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("read %s: %w", dir, err)
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
		return "", time.Time{}, fmt.Errorf("no files found in %s", dir)
	}
	return newest, newestTime, nil
}

func copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	dstInfo, err := os.Stat(dst)
	if err == nil && os.SameFile(srcInfo, dstInfo) {
		return fmt.Errorf("source and destination are the same file: %s", src)
	}

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
