# stressssd — Project Context

## Project Goal

Build a Go command-line utility to help keep external SSDs healthy in a macOS environment, preventing data loss. The name `stressssd` is intentional.

## Repository

- GitHub: https://github.com/joelneely/stressssd
- Local: `/Users/joelneely/Documents/Projects/stressssd`
- Branch: `main` (tracks `origin/main`)

## Current State (as of session end)

### Completed: Step 1 — List external disk identifiers

The program (`main.go`) prints the macOS identifiers of all currently mounted external physical disks to stdout, one per line. Example output when two external drives are connected:

```
disk4
disk6
```

Identifiers are in the form `diskN` (without `/dev/` prefix), as returned by `diskutil`.

## Technical Approach

### Tool: `diskutil`

macOS's built-in `diskutil` CLI is used to query disk information. The specific invocation is:

```
diskutil list -plist external physical
```

- `-plist` produces machine-readable Apple Property List (XML) output
- `external physical` filters to externally connected physical disks only (excludes internal drives, disk images, virtual disks, containers)

### Plist Parsing: `howett.net/plist`

The output is parsed using the `howett.net/plist` library (v1.0.1). Only `plist.Unmarshal` is used — purely a read/decode operation, no writes.

The relevant key in the plist output is `WholeDisks`, an array of strings containing the disk identifiers (e.g., `["disk4", "disk6"]`). Only this key is extracted; partition details are ignored at this stage.

### Go struct used:

```go
type diskutilList struct {
    WholeDisks []string `plist:"WholeDisks"`
}
```

## Project Structure

```
stressssd/
  main.go       — entry point; all logic currently lives here
  go.mod        — module: stressssd, go 1.23
  go.sum        — dependency checksums
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
- **No `/dev/` prefix**: Identifiers are printed as `diskN` (matching `diskutil`'s own convention). Callers can prepend `/dev/` as needed.
- **stderr for errors, stdout for results**: Clean separation so output can be piped or processed by other tools.

## Next Steps (not yet defined)

The broader goal is SSD health monitoring and data-loss prevention. Potential next steps include:
- Retrieving detailed info (SMART status, filesystem health, capacity) for each identified disk
- Monitoring for unexpected unmounts or I/O errors
- Periodic health checks or alerting
