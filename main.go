package main

import (
	"fmt"
	"os"
	"os/exec"

	"howett.net/plist"
)

type diskutilList struct {
	WholeDisks []string `plist:"WholeDisks"`
}

type diskutilInfo struct {
	SMARTStatus string `plist:"SMARTStatus"`
	TotalSize   uint64 `plist:"TotalSize"`
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

func main() {
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

	fmt.Printf("%2s  %-8s  %-10s  %s\n", "#", "Disk", "Size", "SMART Status")
	n := 0
	for _, disk := range list.WholeDisks {
		info, err := getDiskInfo(disk)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting info for %s: %v\n", disk, err)
			continue
		}
		n++
		fmt.Printf("%2d  %-8s  %-10s  %s\n", n, disk, formatSize(info.TotalSize), info.SMARTStatus)
	}
}
