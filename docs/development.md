# Development

## Prerequisites

- Go 1.22+

## Build

```bash
make build
```

Binary goes to `bin/manifestor`.

## Test

```bash
make test
```

## Lint

```bash
make lint
```

## Release

1. Update `CHANGELOG.md`
2. Tag: `git tag v0.x.0`
3. Push: `git push origin v0.x.0`
4. GitHub Actions builds binaries and creates a release
5. Update the Brew formula SHA and URL
