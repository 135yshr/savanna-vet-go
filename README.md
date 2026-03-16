# savanna

[![Go Version](https://img.shields.io/github/go-mod/go-version/135yshr/savanna-vet-go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/135yshr/savanna-vet-go)](https://goreportcard.com/report/github.com/135yshr/savanna-vet-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go CLI tool that detects test smells in Go test files.

## Overview

savanna analyzes Go test files using the `go/analysis` framework and detects common test smells — patterns that indicate potential issues with test quality and maintainability. Integrates with `go vet -vettool` for seamless CI/editor support.

- **`go vet` integration** — works as a standard `vettool`
- **Fast** — leverages `go/analysis` for static analysis
- **Configurable** — enable/disable individual smell detectors

## Supported Test Smells

| Smell Type | Description |
|---|---|
| `EMPTY_TEST` | Test function with no statements |
| `MISSING_ASSERTION` | Test function without any assertions |
| `SLEEPY_TEST` | Test using `time.Sleep` |
| `REDUNDANT_PRINT` | Test using `fmt.Println` or similar print statements |
| `CONDITIONAL_TEST_LOGIC` | Test containing `if`/`switch` conditional logic |
| `MAGIC_NUMBER_TEST` | Test using magic numbers in assertions |
| `MISSING_ERROR_CHECK` | Test ignoring returned errors |
| `MISSING_HELPER` | Helper function missing `t.Helper()` call |

## Installation

```bash
go install github.com/135yshr/savanna-vet-go/cmd/savanna@latest
go install github.com/135yshr/savanna-vet-go/cmd/savanna-vet@latest
```

### Build from source

```bash
git clone https://github.com/135yshr/savanna-vet-go.git
cd savanna-vet-go
go build -o savanna ./cmd/savanna/
go build -o savanna-vet ./cmd/savanna-vet/
```

## Usage

### go vet integration (recommended)

You can use `savanna-vet` as a `go vet -vettool` plugin.

```bash
# Run via go vet
go vet -vettool=$(which savanna-vet) ./...

# Run from a local build
go vet -vettool=./savanna-vet ./...
```

### Standalone (multichecker)

Both `savanna` and `savanna-vet` are built on `multichecker.Main` / `unitchecker.Main` from `golang.org/x/tools/go/analysis`. The following flags are provided by the analysis framework.

```bash
# Analyze packages
savanna ./...

# Output diagnostics as JSON
savanna -json ./...

# Disable a specific analyzer
savanna -emptytest=false ./...

# Show version
savanna -V
```

### CLI Flags

| Flag | Description |
|---|---|
| `-json` | Output diagnostics in JSON format |
| `-V` | Show version |
| `-c N` | Show N lines of context around each diagnostic |
| `-fix` | Apply suggested fixes |
| `-ANALYZER=false` | Disable a specific analyzer (e.g. `-emptytest=false`) |

### Console Output Example

When test smells are detected, savanna displays a lion ASCII art banner along with the detected smells:

```text
[EMPTY_TEST] example_test.go:10 in TestEmpty() - 空のテストとかお前それ@t_wadaの前でも同じ事言えんの？
[SLEEPY_TEST] example_test.go:20 in TestSlow() - テストでsleepとかお前それ@t_wadaの前でも同じ事言えんの？
```

## License

[MIT](LICENSE)

## Acknowledgments

This project is inspired by [savanna-maven-plugin](https://github.com/kawasima/savanna-maven-plugin) by [@kawasima](https://github.com/kawasima), a test smell detector for Java/Maven projects.
