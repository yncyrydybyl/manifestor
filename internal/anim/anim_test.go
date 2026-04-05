package anim

import "testing"

func TestRegistryHasAnimations(t *testing.T) {
	anims := List()
	if len(anims) < 7 {
		t.Errorf("expected at least 7 animations, got %d", len(anims))
	}
}

func TestGetKnown(t *testing.T) {
	names := []string{
		"rainbow-beam",
		"growing-rose",
		"teleport",
		"portal",
		"crystallize",
		"emoji-rain",
		"matrix-manifest",
		"fire-forge",
		"starfield",
	}
	for _, name := range names {
		a := Get(name)
		if a == nil {
			t.Errorf("Get(%q) returned nil", name)
		}
	}
}

func TestGetUnknown(t *testing.T) {
	a := Get("nonexistent-animation")
	if a != nil {
		t.Error("Get(unknown) should return nil")
	}
}

func TestRandom(t *testing.T) {
	a := Random()
	if a == nil {
		t.Error("Random() returned nil")
	}
}

func TestListSorted(t *testing.T) {
	anims := List()
	for i := 1; i < len(anims); i++ {
		if anims[i].Name < anims[i-1].Name {
			t.Errorf("List() not sorted: %q before %q", anims[i-1].Name, anims[i].Name)
		}
	}
}
