# gopidof

A lightweight `pidof` implementation written in Go.

## Background

Homebrew's [`pidof`](https://formulae.brew.sh/formula/pidof) package has been deprecated. This provides a native Go alternative.

## Installation

### From releases

Download the latest binary from the [releases page](https://github.com/nusnewob/gopidof/releases).

### From source

```bash
go install github.com/nusnewob/gopidof@latest
```

## Usage

```bash
pidof <process-name>
```

## Development

### Running tests

```bash
go test -v ./...
```

### Building

```bash
go build -o pidof
```

## CI/CD

This project uses GitHub Actions for:

- **Automated testing**: Tests run on every push and pull request
- **Automated releases**: When a tag is pushed (e.g., `v1.0.0`), binaries are automatically built for multiple platforms and attached to the GitHub release
