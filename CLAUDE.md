# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

savanna is a Go (1.24+) CLI tool that detects test smells in Go test files using AST analysis. Inspired by [kawasima/savanna-maven-plugin](https://github.com/kawasima/savanna-maven-plugin) (Java/Maven version). When smells are detected, a lion ASCII art roars with a @t_wada meme message. Zero external dependencies (stdlib only).

## Build & Test Commands

```bash
# Build
go build -o savanna ./cmd/savanna/

# Run all tests
go test ./...

# Run a single test
go test ./analyzer/ -run TestEmptyTest

# Run with verbose output
go test ./... -v

# Run the tool against a directory
go run ./cmd/savanna/ ./path/to/project
```

### CLI Flags

`-format console|json` `-output <dir>` (JSON output dir) `-fail` (exit 1 on smells) `-enable SMELL1,SMELL2` `-disable SMELL1,SMELL2` `-list` (show all smell types) `-version`

## Architecture

The tool parses `*_test.go` files into ASTs using `go/ast` and `go/parser`, then runs each function declaration through a pipeline of detectors. `vendor/` and dot-prefixed directories are automatically skipped.

**Flow:** `cmd/savanna/main.go` → `Scanner.ScanDir()` → parse files → run `Detector.Detect()` for each `*ast.FuncDecl` → collect `TestSmell` results → `Reporter.Report()`

### Key abstractions

- **`Detector` interface** (`analyzer/detector.go`): Each smell type implements `Detect(fset, file, fn) []TestSmell`. Detectors receive individual function declarations, not whole files.
- **`Scanner`** (`analyzer/scanner.go`): Orchestrates file walking, AST parsing, and detector execution. Supports enable/disable filtering by `SmellType`.
- **`Reporter` interface** (`reporter/reporter.go`): Output formatters (console with lion banner, JSON).

### Adding a new detector

1. Create `analyzer/<name>_detector.go` implementing `Detector` interface
2. Add the `SmellType` constant, display name, and lion message to `analyzer/smell.go`
3. Register in `analyzer/registry.go` → `AllDetectors()`
4. Add tests in `analyzer/scanner_test.go` using `scanner.ScanSource(src)`

### Conventions

- Detector files are named `<smell_name>_detector.go`
- `isTestFunc()` and `isHelperFunc()` in `detector.go` are shared helpers for identifying test/helper functions by signature
- `ScanSource(src string)` is the primary method for unit testing detectors — pass Go source as a string
- Messages and display names are in Japanese
