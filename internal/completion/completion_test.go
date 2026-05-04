package completion

import (
	"strings"
	"testing"
)

func TestBashContainsAnimNames(t *testing.T) {
	out := Bash()
	if !strings.Contains(out, "rainbow-beam") {
		t.Error("bash completion missing animation names")
	}
	if !strings.Contains(out, "complete -o default -F _m_completions m") {
		t.Error("bash completion missing complete command")
	}
	if !strings.Contains(out, "complete -o default -F _m_completions mm") {
		t.Error("bash completion missing mm alias")
	}
	if !strings.Contains(out, "complete -o default -F _m_completions manifestor") {
		t.Error("bash completion missing manifestor canonical name")
	}
	if !strings.Contains(out, "--anim-size") {
		t.Error("bash completion missing --anim-size flag")
	}
	if !strings.Contains(out, `compgen -W "1 5 full"`) {
		t.Error("bash completion missing --anim-size value completion")
	}
}

func TestZshContainsAnimDescriptions(t *testing.T) {
	out := Zsh()
	if !strings.Contains(out, "#compdef m mm manifestor") {
		t.Error("zsh completion missing compdef header")
	}
	if !strings.Contains(out, "rainbow-beam") {
		t.Error("zsh completion missing animation names")
	}
	if !strings.Contains(out, "'--anim[play a specific animation]") {
		t.Error("zsh completion missing --anim flag")
	}
	if !strings.Contains(out, "completion[generate shell completions]") {
		t.Error("zsh completion missing completion subcommand")
	}
	if !strings.Contains(out, "--anim-size[animation size mode]") {
		t.Error("zsh completion missing --anim-size flag")
	}
}

func TestFishContainsBothCommands(t *testing.T) {
	out := Fish()
	if !strings.Contains(out, "complete -c m ") {
		t.Error("fish completion missing m command")
	}
	if !strings.Contains(out, "complete -c mm ") {
		t.Error("fish completion missing mm command")
	}
	if !strings.Contains(out, "complete -c manifestor ") {
		t.Error("fish completion missing manifestor command")
	}
	if !strings.Contains(out, "rainbow-beam") {
		t.Error("fish completion missing animation names")
	}
	if !strings.Contains(out, "'bash zsh fish'") {
		t.Error("fish completion missing shell subcommands")
	}
	if !strings.Contains(out, "anim-size") {
		t.Error("fish completion missing --anim-size flag")
	}
}

func TestAllShellsNonEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		fn   func() string
	}{
		{"bash", Bash},
		{"zsh", Zsh},
		{"fish", Fish},
	} {
		out := tc.fn()
		if len(out) < 100 {
			t.Errorf("%s completion suspiciously short: %d bytes", tc.name, len(out))
		}
	}
}
