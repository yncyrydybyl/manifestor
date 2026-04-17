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

func TestCopyFileSameFile(t *testing.T) {
	dir := t.TempDir()

	f := filepath.Join(dir, "file.txt")
	os.WriteFile(f, []byte("precious data"), 0644)

	err := copyFile(f, f)
	if err == nil {
		t.Fatal("expected error when copying file onto itself, got nil")
	}

	// Verify the file wasn't truncated
	got, _ := os.ReadFile(f)
	if string(got) != "precious data" {
		t.Errorf("file was corrupted: %q", got)
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

func TestNewestNOrdering(t *testing.T) {
	dir := t.TempDir()

	names := []string{"a.txt", "b.txt", "c.txt", "d.txt"}
	for _, n := range names {
		os.WriteFile(filepath.Join(dir, n), []byte(n), 0644)
		time.Sleep(10 * time.Millisecond)
	}

	files, err := newestN(dir, 3)
	if err != nil {
		t.Fatalf("newestN error: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("got %d files, want 3", len(files))
	}
	wantOrder := []string{"d.txt", "c.txt", "b.txt"}
	for i, f := range files {
		if filepath.Base(f.path) != wantOrder[i] {
			t.Errorf("files[%d] = %q, want %q", i, filepath.Base(f.path), wantOrder[i])
		}
	}
}

func TestNewestNFewerAvailable(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)

	files, err := newestN(dir, 10)
	if err != nil {
		t.Fatalf("newestN error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("got %d files, want 2 (directory had fewer than requested)", len(files))
	}
}

func TestNewestNSkipsDirsAndDotfiles(t *testing.T) {
	dir := t.TempDir()

	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0644)
	time.Sleep(10 * time.Millisecond)
	os.WriteFile(filepath.Join(dir, "visible.txt"), []byte("v"), 0644)

	files, err := newestN(dir, 5)
	if err != nil {
		t.Fatalf("newestN error: %v", err)
	}
	if len(files) != 1 || filepath.Base(files[0].path) != "visible.txt" {
		t.Errorf("got %+v, want only visible.txt", files)
	}
}

func TestNewestNEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if _, err := newestN(dir, 3); err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestFindNValidation(t *testing.T) {
	cases := []struct {
		n       int
		wantErr bool
	}{
		{0, true},
		{-1, true},
		{MaxBatch + 1, true},
		{255, true},
	}
	for _, c := range cases {
		_, err := FindN(c.n, ".")
		if (err != nil) != c.wantErr {
			t.Errorf("FindN(%d) err=%v, wantErr=%v", c.n, err, c.wantErr)
		}
	}
}

func TestMaxBatch(t *testing.T) {
	if MaxBatch != 254 {
		t.Errorf("MaxBatch = %d, want 254", MaxBatch)
	}
}
