# stressssd

A macOS command-line utility to help keep external SSDs healthy and prevent data loss. The name is intentional.

> **Work in progress.** As of February 2026, this tool is in early development. The disk read feature is not yet implemented.

## Current behavior

When run, `stressssd` lists all external physical disks currently connected to the Mac, then prompts the user to select one for exercising:

```
 #  Disk      Size        SMART Status
 1  disk4     2.0 GB      Not Supported
 2  disk6     4.0 TB      Verified

Enter line number to exercise (1-2), or 0 to quit: 2

Selected: disk6 (4.0 TB)
WARNING: All files on this disk must be closed before proceeding.
Ready to proceed? (yes/y to continue, no/n to quit):
```

If no external drives are detected, the program prints a message and exits.

`Not Supported` for SMART status is common with disks connected via USB adapters that don't pass SMART commands through to the drive.

## Requirements

- macOS (depends on `diskutil`)
- Go 1.23 or later (to build from source)

## Build and run

```bash
go build -o stressssd .
./stressssd
```

Or without building:

```bash
go run .
```

## Planned direction

The tool will be enhanced to fully read a selected drive, helping to minimize the risk of data loss from SSD failure.
