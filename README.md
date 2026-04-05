# manifestor

Grab the latest file from `~/Downloads`, sanitize its name, and drop it here.

Installs as `m`. One keystroke away from your last download.

## Why

You download a file. It lands in `~/Downloads` with a garbage name like
`Screenshot 2026-04-05 at 14.32.07.png`. You want it in your project with a
clean name. Instead of navigating, renaming, and copying:

```bash
m
```

## Install

### Homebrew

```bash
brew install yncyrydybyl/tap/manifestor
```

### From source

```bash
git clone https://github.com/yncyrydybyl/manifestor.git
cd manifestor
make install
```

## Usage

```bash
m [options] [destination]
```

## Examples

```bash
# Copy latest download to current directory
m

# Copy to a specific directory
m ./assets

# Force mode — skip the staleness check
m --force
mm              # same thing, shorter
```

### Staleness check

If the newest file in `~/Downloads` is older than 8 hours, `m` asks for
confirmation. This prevents accidentally grabbing a stale file.

```
  The newest file in ~/Downloads is 2 days old:
  budget-q3-2025.xlsx

  Hint: this might not be what you just downloaded.
  Use 'mm' or 'm --force' to skip this check.

  Grab it anyway? [y/N]
```

To skip the check: use `mm` or `m --force`.

### What it does

1. Finds the most recently modified file in `~/Downloads`
2. Sanitizes the file name (lowercase, no spaces, no special characters)
3. Copies it to the destination (default: `.`)

```
~/Downloads/Screenshot 2026-04-05 at 14.32.07.png
→ ./screenshot-2026-04-05-at-14.32.07.png
```

## Status

Early alpha. Works on Linux and macOS.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
