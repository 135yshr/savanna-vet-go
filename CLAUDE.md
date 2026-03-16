# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

savanna is a Go (1.25+) CLI tool that detects test smells in Go test files using the `go/analysis` framework. Inspired by [kawasima/savanna-maven-plugin](https://github.com/kawasima/savanna-maven-plugin) (Java/Maven version). Integrates with `go vet -vettool` for seamless CI/editor support.

### Dependencies

- `golang.org/x/tools/go/analysis` — required for `go vet -vettool` integration (`unitchecker`/`multichecker`)
- Transitive: `golang.org/x/mod`, `golang.org/x/sync` (required by `golang.org/x/tools`)

> **Note:** The original "zero external dependencies" policy was relaxed to adopt the `go/analysis` framework. Only `golang.org/x/*` modules (quasi-stdlib) are permitted.

## Build & Test Commands

```bash
# Build
go build -o savanna ./cmd/savanna/
go build -o savanna-vet ./cmd/savanna-vet/

# Run all tests
go test ./...

# Run a single analyzer test
go test ./analyzer/ -run TestAllAnalyzers/emptytest

# Run with verbose output
go test ./... -v

# Run the tool (standalone via multichecker)
go run ./cmd/savanna/ ./path/to/project

# Run via go vet
go vet -vettool=./savanna-vet ./...
```

## Architecture

The tool uses `golang.org/x/tools/go/analysis` framework. Each smell detector is an `*analysis.Analyzer` that receives parsed files via `analysis.Pass` and reports diagnostics with `pass.Reportf`.

**Flow:** `cmd/savanna/main.go` (`multichecker.Main`) or `cmd/savanna-vet/main.go` (`unitchecker.Main`) → `AllAnalyzers()` → each Analyzer's `Run` func processes `pass.Files` → `pass.Reportf` for diagnostics

### Key abstractions

- **`*analysis.Analyzer`** (`analyzer/*_analyzer.go`): Each smell type is an Analyzer with a `Run` function. Analyzers receive `*analysis.Pass` which provides file set, type info, and reporting.
- **`AllAnalyzers()`** (`analyzer/analyzers.go`): Registry returning all Analyzer instances.
- **Shared helpers** (`analyzer/detector.go`): `isTestFunc()`, `isHelperFunc()`, `testingParamName()`, `isPkgCall()`, assertion/print detection helpers.

### Entry points

- `cmd/savanna/main.go` — standalone execution via `multichecker.Main`
- `cmd/savanna-vet/main.go` — `go vet -vettool` integration via `unitchecker.Main`

### Adding a new analyzer

1. Create `analyzer/<name>_analyzer.go` with `var XxxAnalyzer = &analysis.Analyzer{...}`
2. Implement `Run` function: iterate `pass.Files`, filter `_test.go`, process `*ast.FuncDecl`
3. Register in `analyzer/analyzers.go` → `AllAnalyzers()`
4. Add testdata in `analyzer/testdata/src/<name>/<name>_test.go` with `// want` comments
5. Add package mapping in `analyzer/analyzer_test.go` → `analyzerPackages`

### Conventions

- Analyzer files are named `<smell_name>_analyzer.go`
- `isTestFunc()` and `isHelperFunc()` in `detector.go` are shared helpers for identifying test/helper functions by signature
- `isPkgCall()` uses type info to resolve package identity (handles aliases and shadowing)
- Use `testingParamName()` to get `*testing.T` parameter name for receiver verification
- Tests use `analysistest.Run` with testdata directories
- Messages and display names are in Japanese
