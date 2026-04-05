# Changelog

## Unreleased

## 0.1.3.0 - 2026-04-05

### Added
- Shell completions for bash, zsh, and fish via `m completion <shell>`
- Tab-complete all flags and animation names
- Homebrew formula auto-installs completions

## 0.1.2.0 - 2026-04-05

### Added
- Manifestation animations: 9 unique terminal animations play after each copy
- `--anim <name>` to pick a specific animation (e.g., `m --anim fire-forge`)
- `--no-anim` to skip the animation
- `--list-anims` to see all available animations
- Animations: rainbow-beam, growing-rose, teleport, portal, crystallize, emoji-rain, matrix-manifest, fire-forge, starfield
- Random animation selected by default for variety

## 0.1.1.0 - 2026-04-05

### Added
- Short binary name: installs as `m` instead of `manifestor`
- `mm` symlink for force mode (skips all confirmations)
- Staleness check: prompts for confirmation if the newest download is older than 8 hours
- `--force` / `-f` flag to skip the staleness check
- Same-file guard: prevents data loss when running from `~/Downloads`

### Changed
- Refactored grab package into two-phase Find + Copy workflow
- Updated Makefile, Homebrew formula, and release workflow for new binary name

## 0.1.0 - 2026-04-05

### Added
- Initial release
- Grab latest file from ~/Downloads
- Sanitize file names (lowercase, no spaces, no special characters)
- Copy to current directory or specified destination
- `--help` and `--version` flags
