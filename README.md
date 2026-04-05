# manifestor

Grab the latest file from `~/Downloads`, sanitize its name, and drop it here.

## Why

You download a file. It lands in `~/Downloads` with a garbage name like
`Screenshot 2026-04-05 at 14.32.07.png`. You want it in your project with a
clean name. Instead of navigating, renaming, and copying — just run `manifestor`.

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
manifestor [destination]
```

## Examples

```bash
# Copy latest download to current directory
manifestor

# Copy to a specific directory
manifestor ./assets

# Copy to ~/Documents
manifestor ~/Documents
```

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
