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

	for _, disk := range list.WholeDisks {
		fmt.Println(disk)
	}
}
