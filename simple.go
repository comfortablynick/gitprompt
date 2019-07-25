package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func parseSimple(status string) error {
	log.Println("Formatting simple output")
	lines := strings.Split(status, "\n")

	for i := 0; i < len(lines)-1; i++ {
		fmt.Println(lines[i])
	}
	return nil
}

func runSimple() error {
	log.Println("Running simple mode")
	status_cmd := []string{"status", "--porcelain", "--branch", "--untracked-files=normal"}
	cmd := exec.Command(gitExe, status_cmd...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Status:\n%s", string(out))
	parseSimple(string(out))
	return err
}
