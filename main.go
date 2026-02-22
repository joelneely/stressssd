package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"howett.net/plist"
)

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

func getDiskInfo(disk string) (diskutilInfo, error) {
	cmd := exec.Command("diskutil", "info", "-plist", "/dev/"+disk)
	output, err := cmd.Output()
	if err != nil {
		return diskutilInfo{}, err
	}
	var info diskutilInfo
	if _, err = plist.Unmarshal(output, &info); err != nil {
		return diskutilInfo{}, err
	}
	return info, nil
}

func formatSize(bytes uint64) string {
	const (
		GB = 1_000_000_000
		TB = 1_000_000_000_000
	)
	if bytes >= TB {
		return fmt.Sprintf("%.1f TB", float64(bytes)/TB)
	}
	return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
}

func readDisk(name string, totalBytes uint64) error {
	// Unmount the disk to prevent concurrent filesystem access during the read.
	fmt.Printf("Unmounting %s...\n", name)
	out, err := exec.Command("diskutil", "unmountDisk", "/dev/"+name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unmount failed: %v: %s", err, strings.TrimSpace(string(out)))
	}
	fmt.Println(strings.TrimSpace(string(out)))

	// Remount exactly once, whether we exit normally, on error, or via signal.
	var once sync.Once
	remount := func() {
		once.Do(func() {
			fmt.Printf("Remounting %s...\n", name)
			out, err := exec.Command("diskutil", "mountDisk", "/dev/"+name).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: remount failed: %v: %s\n", err, strings.TrimSpace(string(out)))
			} else {
				fmt.Println(strings.TrimSpace(string(out)))
			}
		})
	}
	defer remount()

	// Handle Ctrl+C: remount before exiting.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nInterrupted.")
		remount()
		os.Exit(1)
	}()
	defer signal.Stop(sigCh)

	// Open the raw character device (/dev/rdiskN) for fast sequential reads.
	rawDev := "/dev/r" + name
	f, err := os.Open(rawDev)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", rawDev, err)
	}
	defer f.Close()

	const bufSize = 1024 * 1024 // 1 MB
	buf := make([]byte, bufSize)
	var bytesRead uint64

	fmt.Println("Reading disk...")
	for {
		n, err := f.Read(buf)
		bytesRead += uint64(n)
		if totalBytes > 0 {
			pct := float64(bytesRead) / float64(totalBytes) * 100
			fmt.Printf("\r  %5.1f%%  %s / %s", pct, formatSize(bytesRead), formatSize(totalBytes))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println()
			return fmt.Errorf("read error after %s: %v", formatSize(bytesRead), err)
		}
	}
	fmt.Println()
	fmt.Println("Read complete.")
	return nil
}

func main() {
	if os.Getuid() != 0 {
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error determining executable path: %v\n", err)
			os.Exit(1)
		}
		args := append([]string{"sudo", exe}, os.Args[1:]...)
		if err := syscall.Exec("/usr/bin/sudo", args, os.Environ()); err != nil {
			fmt.Fprintf(os.Stderr, "error re-launching with sudo: %v\n", err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("diskutil", "list", "-plist", "external", "physical")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running diskutil: %v\n", err)
		os.Exit(1)
	}

	var list diskutilList
	if _, err = plist.Unmarshal(output, &list); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing diskutil output: %v\n", err)
		os.Exit(1)
	}

	var disks []diskEntry
	for _, disk := range list.WholeDisks {
		info, err := getDiskInfo(disk)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting info for %s: %v\n", disk, err)
			continue
		}
		disks = append(disks, diskEntry{disk, info})
	}

	if len(disks) == 0 {
		fmt.Println("No external physical drives found.")
		return
	}

	fmt.Printf("%2s  %-8s  %-10s  %s\n", "#", "Disk", "Size", "SMART Status")
	for i, d := range disks {
		fmt.Printf("%2d  %-8s  %-10s  %s\n", i+1, d.name, formatSize(d.info.TotalSize), d.info.SMARTStatus)
	}

	reader := bufio.NewReader(os.Stdin)

	var selected int
	for {
		fmt.Printf("\nEnter line number to exercise (1-%d), or 0 to quit: ", len(disks))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || val < 0 || val > len(disks) {
			fmt.Printf("Please enter a number between 1 and %d, or 0 to quit.\n", len(disks))
			continue
		}
		if val == 0 {
			return
		}
		selected = val
		break
	}

	chosen := disks[selected-1]
	fmt.Printf("\nSelected: %s (%s)\n", chosen.name, formatSize(chosen.info.TotalSize))
	fmt.Println("WARNING: All files on this disk must be closed before proceeding.")

	for {
		fmt.Print("Ready to proceed? (yes/y to continue, no/n to quit): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
			os.Exit(1)
		}
		switch strings.TrimSpace(strings.ToLower(line)) {
		case "yes", "y":
			fmt.Printf("Preparing to fully read disk %s...\n", chosen.name)
			if err := readDisk(chosen.name, chosen.info.TotalSize); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			return
		case "no", "n":
			return
		default:
			fmt.Println("Please enter yes, y, no, or n.")
		}
	}
}
