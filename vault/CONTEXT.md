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

### Completed: Step 4 — Interactive disk selection and confirmation

The program (`main.go`) prints a columnar summary of all currently mounted external physical disks, then interacts with the user:

```
 #  Disk      Size        SMART Status
 1  disk4     2.0 GB      Not Supported
 2  disk6     4.0 TB      Not Supported

Enter line number to exercise (1-2), or 0 to quit: 2

Selected: disk6 (4.0 TB)
WARNING: All files on this disk must be closed before proceeding.
Ready to proceed? (yes/y to continue, no/n to quit): y
Preparing to fully read disk disk6...
```

If no external drives are detected, the program prints "No external physical drives found." and exits.

Input validation:
- Selection: non-integers and out-of-range values prompt a retry; 0 quits
- Confirmation: only yes/y/no/n accepted; anything else prompts a retry; no/n quits

The disk read itself is not yet implemented; the "Preparing to fully read..." message is a placeholder.

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

type diskEntry struct {
    name string
    info diskutilInfo
}
```

`diskEntry` collects results before printing so that user selections can be looked up by index after the list is displayed.

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

## Next Steps

- Implement the full sequential read of the selected disk (the core SSD exercise feature)
