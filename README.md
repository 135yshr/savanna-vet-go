# savanna

[![Go Version](https://img.shields.io/github/go-mod/go-version/135yshr/savanna-vet-go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/135yshr/savanna-vet-go)](https://goreportcard.com/report/github.com/135yshr/savanna-vet-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go CLI tool that detects test smells in Go test files.

## Overview

savanna analyzes Go test files using AST parsing and detects common test smells — patterns that indicate potential issues with test quality and maintainability. When smells are detected, a lion ASCII art roars with a [@t_wada](https://github.com/t-wada) meme message.

- **Zero external dependencies** — stdlib only
- **Fast** — leverages `go/ast` and `go/parser` for static analysis
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
```

### Build from source

```bash
git clone https://github.com/135yshr/savanna-vet-go.git
cd savanna-vet-go
go build -o savanna ./cmd/savanna/
```

## Usage

```bash
# Scan the current directory
savanna .

# Scan specific directories
savanna ./pkg ./internal

# Output as JSON
savanna -format json ./...

# Exit with code 1 if smells are found (useful for CI)
savanna -fail ./...

# Enable only specific smells
savanna -enable EMPTY_TEST,SLEEPY_TEST ./...

# Disable specific smells
savanna -disable MAGIC_NUMBER_TEST ./...

# List all available smell types
savanna -list

# Show version
savanna -version

# Output JSON to custom directory
savanna -format json -output ./test-reports ./...
```

### CLI Flags

| Flag | Default | Description |
|---|---|---|
| `-format` | `console` | Report format (`console` or `json`) |
| `-output` | `savanna-reports` | Output directory for JSON reports |
| `-fail` | `false` | Exit with code 1 when smells are detected |
| `-enable` | *(all)* | Comma-separated list of smell types to enable |
| `-disable` | *(none)* | Comma-separated list of smell types to disable |
| `-list` | `false` | Show all available smell types |
| `-version` | `false` | Show version |

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
