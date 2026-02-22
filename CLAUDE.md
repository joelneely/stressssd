# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

`stressssd` is a macOS CLI utility to help keep external SSDs healthy and prevent data loss by fully reading a selected drive to surface latent errors. The name is intentional.

## Commands

```bash
go build -o stressssd .   # Build binary
go run .                   # Run without building (will sudo-escalate automatically)
go test ./...              # Run tests
go vet ./...               # Static analysis
```

## Architecture

Single-file Go program (`main.go`). Flow:

1. If not running as root, re-exec via `sudo` using `syscall.Exec` (replaces current process)
2. Invoke `diskutil list -plist external physical` → parse `WholeDisks` array
3. For each disk, invoke `diskutil info -plist /dev/diskN` → collect `TotalSize` and `SMARTStatus`
4. Print columnar list with index, disk name, size, SMART status
5. Prompt user to select a disk by line number (0 to quit; validates range)
6. Warn about open files; prompt for confirmation (yes/y/no/n)
7. Unmount disk (`diskutil unmountDisk`), read raw device (`/dev/rdiskN`) sequentially with live progress and elapsed time on completion, remount (`diskutil mountDisk`), print total run time

The `vault/` directory is an Obsidian notes vault and is not part of the Go program. JSON files inside it are gitignored.

## Key Design Decisions

- **plist over text parsing**: `diskutil` text output is not stable; `-plist` is the intended machine-readable interface.
- **`WholeDisks` key**: Gives top-level identifiers without partition detail, appropriate for whole-disk operations.
- **`/dev/rdiskN` for reads**: Raw character device bypasses the buffer cache for faster sequential reads.
- **Inline read loop**: Chosen over `dd`/`cat` to enable live progress display.
- **`sync.Once` for remount**: Guarantees remount runs exactly once on normal exit, error, or SIGINT/SIGTERM.
- **`syscall.Exec` for sudo**: Replaces the current process rather than spawning a child, so the elevated process owns the terminal cleanly.
- **stderr/stdout separation**: Errors go to stderr; results go to stdout so output can be piped cleanly.
- **Decimal sizes**: Matches `diskutil` convention (1 GB = 1,000,000,000 bytes).
- **Timing**: `programStart` captured at top of `main()`; `readStart` captured before the read loop. `formatDuration` formats as `Xh Ym Zs`, `Ym Zs`, or `Zs`. Read time shown on "Read complete."; total run time shown after successful remount.

## macOS Dependency

This tool is macOS-only — it relies on `diskutil` and macOS raw disk device paths.
