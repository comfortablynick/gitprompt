package main

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const notRepoStatus string = "exit status 128"

// ErrNotAGitRepo returned when no repo found
var ErrNotAGitRepo = errors.New("not a git repo")

// GetGitOutput returns a buffer of git status command output
func GetGitOutput(cwd string) (io.Reader, error) {
	if ok, err := IsInsideWorkTree(cwd); err != nil {
		if err == ErrNotAGitRepo {
			return nil, ErrNotAGitRepo
		}
		log.Debugf("error detecting work tree: %s", err)
		return nil, err
	} else if !ok {
		return nil, ErrNotAGitRepo
	}

	var buf = new(bytes.Buffer)
	cmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	cmd.Stdout = buf
	cmd.Dir = cwd
	log.Debugf("running %q", cmd.Args)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return buf, nil
}

// PathToGitDir returns parsed root of git repo
func PathToGitDir(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Dir = cwd
	log.Debugf("running %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsInsideWorkTree returns bool to indicate if path is inside git tree
func IsInsideWorkTree(cwd string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = cwd
	log.Debugf("running %q", cmd.Args)

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
