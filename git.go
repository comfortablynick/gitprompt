package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const notRepoStatus = "exit status 128"
const gitExe = "git"

// ErrNotAGitRepo returned when no repo found
var ErrNotAGitRepo = errors.New("not a git repo")

// GetGitStatusOutput returns a buffer of git status command output
func GetGitStatusOutput(cwd string) (io.Reader, error) {
	// if ok, err := IsInsideWorkTree(cwd); err != nil {
	//     if err == ErrNotAGitRepo {
	//         return nil, ErrNotAGitRepo
	//     }
	//     log.Printf("error detecting work tree: %s", err)
	//     return nil, err
	// } else if !ok {
	//     return nil, ErrNotAGitRepo
	// }

	var buf = new(bytes.Buffer)
	cmd := exec.Command(gitExe, "status", "--porcelain=v2", "--branch") // #nosec
	cmd.Stdout = buf
	cmd.Dir = cwd
	log.Printf("GetGitStatusOutput cmd: %q", cmd.Args)

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf, nil
}

// GetGitNumstat returns output of diff --numstat
func GetGitNumstat(cwd string) (string, error) {
	cmd := exec.Command(gitExe, "diff", "--numstat") // #nosec
	cmd.Dir = cwd
	log.Printf("GetGitNumstat cmd: %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GetGitTag returns tag name for detatched head
func GetGitTag(cwd string) (string, error) {
	cmd := exec.Command(gitExe, "describe", "--tags", "--exact-match") // #nosec
	cmd.Dir = cwd
	log.Printf("GetGitTag cmd: %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// PathToGitDir returns parsed root of git repo
func PathToGitDir(cwd string) (string, error) {
	cmd := exec.Command(gitExe, "rev-parse", "--absolute-git-dir") // #nosec
	cmd.Dir = cwd
	log.Printf("PathToGitDir cmd: %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsInsideWorkTree returns bool to indicate if path is inside git tree
func IsInsideWorkTree(cwd string) (bool, error) {
	cmd := exec.Command(gitExe, "rev-parse", "--is-inside-work-tree") // #nosec
	cmd.Dir = cwd
	log.Printf("IsInsideWorkTree cmd: %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 128 {
					return false, ErrNotAGitRepo
				}
			}
		}
		if cmd.ProcessState.String() == notRepoStatus {
			return false, ErrNotAGitRepo
		}
		return false, err
	}
	return strconv.ParseBool(strings.TrimSpace(string(out)))
}
