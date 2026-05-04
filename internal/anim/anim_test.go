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

func TestAllAnimationsHaveSizes(t *testing.T) {
	for _, a := range List() {
		if len(a.Sizes) == 0 {
			t.Errorf("animation %q has no Sizes declared", a.Name)
		}
	}
}

func TestAllAnimationsSupportOneLiner(t *testing.T) {
	// OneLiner is the default and should be universally supported.
	for _, a := range List() {
		if !a.SupportsSize(OneLiner) {
			t.Errorf("animation %q does not support OneLiner (the default)", a.Name)
		}
	}
}

func TestParseSize(t *testing.T) {
	cases := []struct {
		in   string
		want Size
		ok   bool
	}{
		{"1", OneLiner, true},
		{"5", FiveLiner, true},
		{"full", FullScreen, true},
		{"", 0, false},
		{"3", 0, false},
		{"FULL", 0, false},
		{"-1", 0, false},
	}
	for _, c := range cases {
		got, ok := ParseSize(c.in)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("ParseSize(%q) = (%v, %v), want (%v, %v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestSizeString(t *testing.T) {
	cases := []struct {
		s    Size
		want string
	}{
		{OneLiner, "1"},
		{FiveLiner, "5"},
		{FullScreen, "full"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Size(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

func TestBestSizeFallback(t *testing.T) {
	// All animations should yield a usable Size from BestSize, even when
	// requesting a size they don't support (fallback to first).
	for _, a := range List() {
		got := a.BestSize(FullScreen)
		if !a.SupportsSize(got) {
			t.Errorf("%s.BestSize(FullScreen) returned %v, which is not supported", a.Name, got)
		}
	}
}

func TestRandomForSize(t *testing.T) {
	for _, s := range []Size{OneLiner, FiveLiner, FullScreen} {
		a := RandomForSize(s)
		if a == nil {
			t.Errorf("RandomForSize(%v) returned nil", s)
		}
	}
}
