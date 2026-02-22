# stressssd — Project Context

## Project Goal

Build a Go command-line utility to help keep external SSDs healthy in a macOS environment, preventing data loss. The name `stressssd` is intentional.

## Repository

- GitHub: https://github.com/joelneely/stressssd
- Local: `/Users/joelneely/Documents/Projects/stressssd`
- Branch: `main` (tracks `origin/main`)

## Current State (as of session end)

### Completed: Step 1 — List external disk identifiers

### Completed: Step 2 — Add SMART status and disk size

### Completed: Step 3 — Add index column and header row

The program (`main.go`) prints a columnar summary of all currently mounted external physical disks to stdout. Example output:

```
 #  Disk      Size        SMART Status
 1  disk4     2.0 GB      Not Supported
 2  disk6     4.0 TB      Not Supported
```

Columns: sequential index, disk identifier, total size (decimal GB/TB), SMART status string. The index only increments for successfully retrieved disks.

`Not Supported` for SMART status is typical for disks connected via USB adapters that don't pass SMART commands through.

## Technical Approach

### Tool: `diskutil`

macOS's built-in `diskutil` CLI is used in two stages:

**Stage 1 — enumerate disks:**
```
diskutil list -plist external physical
```
- `-plist` produces machine-readable Apple Property List (XML) output
- `external physical` filters to externally connected physical disks only (excludes internal drives, disk images, virtual disks, containers)
- The `WholeDisks` key yields an array of top-level disk identifiers (e.g., `["disk4", "disk6"]`)

**Stage 2 — per-disk info:**
```
diskutil info -plist /dev/diskN
```
- Called once per disk identifier from Stage 1
- Keys used: `TotalSize` (uint64, bytes) and `SMARTStatus` (string)

### Plist Parsing: `howett.net/plist`

The output is parsed using the `howett.net/plist` library (v1.0.1). Only `plist.Unmarshal` is used — purely a read/decode operation, no writes.

### Go structs used:

```go
type diskutilList struct {
    WholeDisks []string `plist:"WholeDisks"`
}

type diskutilInfo struct {
    SMARTStatus string `plist:"SMARTStatus"`
    TotalSize   uint64 `plist:"TotalSize"`
}
```

### Size formatting

Decimal units matching `diskutil` convention: GB below 1 TB, TB at or above. One decimal place.

## Project Structure

```
stressssd/
  main.go       — entry point; all logic currently lives here
  go.mod        — module: stressssd, go 1.23
  go.sum        — dependency checksums
  CLAUDE.md     — guidance for Claude Code
  README.md     — project overview and build instructions
  vault/
    CONTEXT.md  — this file
```

## Dependencies

| Module             | Version | Purpose                        |
|--------------------|---------|--------------------------------|
| howett.net/plist   | v1.0.1  | Decode Apple plist output      |

## Design Notes & Decisions

- **plist over text parsing**: `diskutil` text output is human-readable but not stable; plist is the intended machine-readable format.
- **`WholeDisks` over `AllDisksAndPartitions`**: `WholeDisks` gives just the top-level disk identifiers without partition detail, which is what's needed for whole-disk operations.
- **No `/dev/` prefix in output**: Identifiers are printed as `diskN` (matching `diskutil`'s own convention). The `/dev/` prefix is used internally when calling `diskutil info`.
- **stderr for errors, stdout for results**: Clean separation so output can be piped or processed by other tools. Per-disk info errors are non-fatal; processing continues to the next disk.
- **Decimal sizes**: `diskutil` reports sizes in decimal (1 GB = 1,000,000,000 bytes), so the formatter uses the same convention.

## Next Steps (not yet defined)

The broader goal is SSD health monitoring and data-loss prevention. Potential next steps include:
- Monitoring for unexpected unmounts or I/O errors
- Periodic health checks or alerting
- Filesystem health checks
