package grab

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewestFile(t *testing.T) {
	dir := t.TempDir()

	older := filepath.Join(dir, "older.txt")
	newer := filepath.Join(dir, "newer.txt")

	os.WriteFile(older, []byte("old"), 0644)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(newer, []byte("new"), 0644)

	got, _, err := newestFile(dir)
	if err != nil {
		t.Fatalf("newestFile(%q) error: %v", dir, err)
	}
	if got != newer {
		t.Errorf("newestFile(%q) = %q, want %q", dir, got, newer)
	}
}

func TestNewestFileReturnsMtime(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "file.txt")
	os.WriteFile(f, []byte("content"), 0644)

	_, modTime, err := newestFile(dir)
	if err != nil {
		t.Fatalf("newestFile error: %v", err)
	}
	if time.Since(modTime) > 5*time.Second {
		t.Errorf("modTime too old: %v", modTime)
	}
}

func TestNewestFileSkipsDirs(t *testing.T) {
	dir := t.TempDir()

	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	time.Sleep(10 * time.Millisecond)

	f := filepath.Join(dir, "file.txt")
	os.WriteFile(f, []byte("content"), 0644)

	got, _, err := newestFile(dir)
	if err != nil {
		t.Fatalf("newestFile error: %v", err)
	}
	if got != f {
		t.Errorf("got %q, want %q", got, f)
	}
}

func TestNewestFileSkipsDotfiles(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("hidden"), 0644)
	time.Sleep(10 * time.Millisecond)
	visible := filepath.Join(dir, "visible.txt")
	os.WriteFile(visible, []byte("visible"), 0644)

	got, _, err := newestFile(dir)
	if err != nil {
		t.Fatalf("newestFile error: %v", err)
	}
	if got != visible {
		t.Errorf("got %q, want %q", got, visible)
	}
}

func TestNewestFileEmptyDir(t *testing.T) {
	dir := t.TempDir()

	_, _, err := newestFile(dir)
	if err == nil {
		t.Fatal("expected error for empty directory, got nil")
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()

	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")
	content := []byte("hello manifestor")

	os.WriteFile(src, content, 0644)

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dest: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content = %q, want %q", got, content)
	}
}

func TestCopyFileSourceNotFound(t *testing.T) {
	dir := t.TempDir()

	err := copyFile(filepath.Join(dir, "nonexistent"), filepath.Join(dir, "dst"))
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestStaleThreshold(t *testing.T) {
	if StaleThreshold != 8*time.Hour {
		t.Errorf("StaleThreshold = %v, want 8h", StaleThreshold)
	}
}

func TestResultIsStale(t *testing.T) {
	fresh := &Result{Age: 2 * time.Hour}
	if fresh.Age > StaleThreshold {
		t.Error("2h-old file should not be stale")
	}

	stale := &Result{Age: 10 * time.Hour}
	if stale.Age <= StaleThreshold {
		t.Error("10h-old file should be stale")
	}
}
