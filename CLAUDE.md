# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

`stressssd` is a macOS CLI utility to help keep external SSDs healthy and prevent data loss. The name is intentional.

## Commands

```bash
go build -o stressssd .   # Build binary
go run .                   # Run without building
go test ./...              # Run tests
go vet ./...               # Static analysis
```

## Architecture

Currently a single-file Go program (`main.go`). The program:

1. Invokes `diskutil list -plist external physical` to get machine-readable output
2. Parses the Apple plist XML using `howett.net/plist`
3. Extracts `WholeDisks` (top-level identifiers like `disk4`, `disk6`) and prints one per line to stdout
4. Writes errors to stderr only

The `vault/` directory is an Obsidian notes vault and is not part of the Go program. JSON files inside it are gitignored.

## Key Design Decisions

- **plist over text parsing**: `diskutil` text output is not stable; `-plist` is the intended machine-readable interface.
- **`WholeDisks` key**: Gives top-level identifiers without partition detail, appropriate for whole-disk operations.
- **No `/dev/` prefix**: Identifiers are printed as `diskN` (matching `diskutil` convention); callers prepend `/dev/` as needed.
- **stderr/stdout separation**: Errors go to stderr; results go to stdout so output can be piped cleanly.

## macOS Dependency

This tool is macOS-only — it relies on `diskutil`, which is not available on other platforms.
