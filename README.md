# stressssd

A macOS command-line utility to help keep external SSDs healthy and prevent data loss. The name is intentional.

> **Work in progress.** As of February 2026, this tool is in early development. Current functionality is limited to listing attached external disks with basic health information.

## Current behavior

When run, `stressssd` lists all external physical disks currently connected to the Mac, one per line, with their total capacity and SMART status:

```
disk4     2.0 GB      Not Supported
disk6     4.0 TB      Verified
```

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
