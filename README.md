# stressssd

A macOS command-line utility to help keep external SSDs healthy and prevent data loss. The name is intentional.

## Current behavior

`stressssd` lists all external physical disks connected to the Mac, prompts the user to select one, then fully reads the selected disk to exercise it and surface any latent errors.

Example session:

```
 #  Disk      Size        SMART Status
 1  disk4     2.0 GB      Not Supported
 2  disk6     4.0 TB      Verified

Enter line number to exercise (1-2), or 0 to quit: 2

Selected: disk6 (4.0 TB)
WARNING: All files on this disk must be closed before proceeding.
Ready to proceed? (yes/y to continue, no/n to quit): y
Preparing to fully read disk disk6...
Unmounting disk6...
Unmount of all volumes on disk6 was successful
Reading disk...
   42.7%  1.7 TB / 4.0 TB
...
Read complete. (1h 12m 33s)
Remounting disk6...
Volume(s) mounted successfully
Total run time: 1h 13m 45s
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

The program requires root access to read raw disk devices. If not run as root, it will automatically re-execute itself via `sudo` and prompt for a password.
