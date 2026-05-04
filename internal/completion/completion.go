// Package completion generates shell completion scripts for bash, zsh, and fish.
package completion

import (
	"fmt"
	"strings"

	"github.com/yncyrydybyl/manifestor/internal/anim"
)

// animNames returns a space-separated list of animation names.
func animNames() string {
	var names []string
	for _, a := range anim.List() {
		names = append(names, a.Name)
	}
	return strings.Join(names, " ")
}

// Bash generates a bash completion script.
func Bash() string {
	return fmt.Sprintf(`# bash completion for m (manifestor)
# Add to ~/.bashrc:  eval "$(m completion bash)"

_m_completions() {
    local cur prev opts anims
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    opts="--help --version --force --anim --anim-size --no-anim --list-anims"
    anims="%s"

    case "${prev}" in
        --anim)
            COMPREPLY=( $(compgen -W "${anims}" -- "${cur}") )
            return 0
            ;;
        --anim-size)
            COMPREPLY=( $(compgen -W "1 5 full" -- "${cur}") )
            return 0
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
        COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
        return 0
    fi

    # Default to directory completion for destination
    COMPREPLY=( $(compgen -d -- "${cur}") )
    return 0
}

complete -o default -F _m_completions m
complete -o default -F _m_completions mm
complete -o default -F _m_completions manifestor
`, animNames())
}

// Zsh generates a zsh completion script.
func Zsh() string {
	var animList strings.Builder
	for _, a := range anim.List() {
		animList.WriteString(fmt.Sprintf("        '%s:%s'\n", a.Name, a.Desc))
	}

	return fmt.Sprintf(`#compdef m mm manifestor
# zsh completion for m (manifestor)
# Add to ~/.zshrc:  eval "$(m completion zsh)"
# Or save to a file in your $fpath:
#   m completion zsh > "${fpath[1]}/_m"

_m() {
    local -a commands animations

    animations=(
%s    )

    _arguments \
        '(-h --help)'{-h,--help}'[show help]' \
        '(-v --version)'{-v,--version}'[show version]' \
        '(-f --force)'{-f,--force}'[skip staleness check]' \
        '--anim[play a specific animation]:animation:->anims' \
        '--anim-size[animation size mode]:size:(1 5 full)' \
        '--no-anim[skip animation]' \
        '--list-anims[list available animations]' \
        'completion[generate shell completions]:shell:(bash zsh fish)' \
        '*:destination:_directories'

    case "$state" in
        anims)
            _describe 'animation' animations
            ;;
    esac
}

_m "$@"
`, animList.String())
}

// Fish generates a fish completion script.
func Fish() string {
	var animCompletions strings.Builder
	for _, a := range anim.List() {
		animCompletions.WriteString(fmt.Sprintf(
			"complete -c m -l anim -xa '%s' -d '%s'\n", a.Name, a.Desc))
		animCompletions.WriteString(fmt.Sprintf(
			"complete -c mm -l anim -xa '%s' -d '%s'\n", a.Name, a.Desc))
		animCompletions.WriteString(fmt.Sprintf(
			"complete -c manifestor -l anim -xa '%s' -d '%s'\n", a.Name, a.Desc))
	}

	return fmt.Sprintf(`# fish completion for m (manifestor)
# Add to ~/.config/fish/completions/m.fish:
#   m completion fish > ~/.config/fish/completions/m.fish

complete -c m -s h -l help -d 'Show help'
complete -c m -s v -l version -d 'Show version'
complete -c m -s f -l force -d 'Skip staleness check'
complete -c m -l no-anim -d 'Skip animation'
complete -c m -l list-anims -d 'List available animations'
complete -c m -l anim-size -xa '1 5 full' -d 'Animation size mode'
complete -c m -a 'completion' -d 'Generate shell completions'
complete -c m -n '__fish_seen_subcommand_from completion' -xa 'bash zsh fish'

# Same for mm alias
complete -c mm -s h -l help -d 'Show help'
complete -c mm -s v -l version -d 'Show version'
complete -c mm -s f -l force -d 'Skip staleness check'
complete -c mm -l no-anim -d 'Skip animation'
complete -c mm -l list-anims -d 'List available animations'
complete -c mm -l anim-size -xa '1 5 full' -d 'Animation size mode'

# Same for manifestor canonical name
complete -c manifestor -s h -l help -d 'Show help'
complete -c manifestor -s v -l version -d 'Show version'
complete -c manifestor -s f -l force -d 'Skip staleness check'
complete -c manifestor -l no-anim -d 'Skip animation'
complete -c manifestor -l list-anims -d 'List available animations'
complete -c manifestor -l anim-size -xa '1 5 full' -d 'Animation size mode'
complete -c manifestor -a 'completion' -d 'Generate shell completions'
complete -c manifestor -n '__fish_seen_subcommand_from completion' -xa 'bash zsh fish'

# Animation names
%s`, animCompletions.String())
}
