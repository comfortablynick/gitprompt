package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
)

// GitArea holds status info
type GitArea struct {
	modified int
	added    int
	deleted  int
	renamed  int
	copied   int
}

func (a *GitArea) hasChanged() bool {
	return a.added+a.deleted+a.modified+a.copied+a.renamed != 0
}

// RepoInfo holds data about the repo
type RepoInfo struct {
	workingDir string

	// Local branch data
	branch   string
	commit   string
	remote   string
	upstream string
	ahead    int
	behind   int

	// Branch totals
	untracked  int
	unmerged   int
	insertions int
	deletions  int

	// Status for staged/unstaged files
	Unstaged GitArea
	Staged   GitArea
}

func (ri *RepoInfo) hasUnmerged() bool {
	if ri.unmerged > 0 {
		return true
	}
	gitDir, err := PathToGitDir(cwd)
	if err != nil {
		log.Printf("error calling PathToGitDir: %s", err)
		return false
	}
	// TODO: figure out if output of MERGE_HEAD can be useful
	if _, err := os.Stat(path.Join(gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("error reading MERGE_HEAD: %s", err)
		return false
	}
	return true
}

func (ri *RepoInfo) hasModified() bool {
	return ri.Unstaged.hasChanged()
}

func (ri *RepoInfo) isDirty() bool {
	return ri.Staged.hasChanged()
}

// Debug prints repo info
func (ri *RepoInfo) Debug() string {
	return fmt.Sprintf("%#+v", ri)
}

// Fmt formats the output for the shell
func (ri *RepoInfo) Fmt() string {
	// TODO: make format user-configurable
	// Maybe with a TOML/ini file
	log.Printf("Formatting output: %s", ri.Debug())

	var (
		branchGlyph    = ""
		modifiedGlyph  = "Δ"
		dirtyGlyph     = "✘" // ✗
		cleanGlyph     = "✔" // ✓
		untrackedGlyph = "?"
		unmergedGlyph  = "‼"
		aheadArrow     = "↑"
		behindArrow    = "↓"
	)

	return fmt.Sprintf("%s %s@%s %s %s %s",
		branchGlyph,
		ri.branch,
		func() string {
			if ri.commit == "(initial)" {
				return ri.commit
			}
			return ri.commit[:7]
		}(),
		func() string {
			var buf bytes.Buffer
			if ri.ahead > 0 {
				if _, err := buf.WriteString(fmt.Sprintf(" %s%d ", aheadArrow, ri.ahead)); err != nil {
					log.Printf("Buffer error: %s", err)
				}
			}
			if ri.behind > 0 {
				if _, err := buf.WriteString(fmt.Sprintf(" %s%d ", behindArrow, ri.behind)); err != nil {
					log.Printf("Buffer error: %s", err)
				}
			}
			return buf.String()
		}(),
		func() string {
			var buf bytes.Buffer
			if ri.untracked > 0 {
				if _, err := buf.WriteString(untrackedGlyph); err != nil {
					log.Printf("Error writing untrackedGlyph: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Printf("Error writing rune: %s", err)
				}
			}
			if ri.hasUnmerged() {
				if _, err := buf.WriteString(unmergedGlyph); err != nil {
					log.Printf("Error writing unmergedGlyph: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Printf("Error writing rune: %s", err)
				}
			}
			if ri.hasModified() {
				if _, err := buf.WriteString(modifiedGlyph); err != nil {
					log.Printf("Error writing modifiedGlyph: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Printf("Error writing rune: %s", err)
				}
			}
			return buf.String()
		}(),
		// dirty/clean
		func() string {
			if ri.isDirty() {
				return dirtyGlyph
			}
			return cleanGlyph
		}(),
	)
}

func run() *RepoInfo {
	gitOut, err := GetGitStatusOutput(cwd)
	if err != nil {
		log.Printf("Git status error: %s", err)
		if err == ErrNotAGitRepo {
			// Expected if calling from prompt
			os.Exit(0)
		}
		// Some other error -- print to console
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	var repoInfo = new(RepoInfo)
	repoInfo.workingDir = cwd

	if err = repoInfo.ParseRepoInfo(gitOut); err != nil {
		log.Printf("Error parsing git repo: %s", err)
		os.Exit(1)
	}
	diffOut, err := GetGitNumstat(cwd)
	if err != nil {
		log.Printf("Git diff error: %s", err)
	}
	if err = repoInfo.parseDiffNumstat(diffOut); err != nil {
		log.Printf("Error parsing git diff: %v", err)
	}
	return repoInfo
}
