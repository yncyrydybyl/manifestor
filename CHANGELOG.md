# Changelog

## Unreleased

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
