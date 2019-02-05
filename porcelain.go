package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"os"
	"path"

	"github.com/subchen/go-log"
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
	var changed bool
	if a.added != 0 {
		changed = true
	}
	if a.deleted != 0 {
		changed = true
	}
	if a.modified != 0 {
		changed = true
	}
	if a.copied != 0 {
		changed = true
	}
	if a.renamed != 0 {
		changed = true
	}
	return changed
}

// RepoInfo holds data about the repo
type RepoInfo struct {
	workingDir string

	branch   string
	commit   string
	remote   string
	upstream string
	ahead    int
	behind   int

	untracked int
	unmerged  int

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
	// TODO figure out if output of MERGE_HEAD can be useful
	if _, err := ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD")); err != nil {
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
// TODO should be configurable by the user
//
func (ri *RepoInfo) Fmt() string {
	log.Printf("formatting output: %s", ri.Debug())

	var (
		branchGlyph   = ""
		modifiedGlyph = "Δ"
		// deletedGlyph   string = "＊"
		dirtyGlyph     = "✘"
		cleanGlyph     = "✔"
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
				buf.WriteString(fmt.Sprintf(" %s%d ", aheadArrow, ri.ahead))
			}
			if ri.behind > 0 {
				buf.WriteString(fmt.Sprintf(" %s%d ", behindArrow, ri.behind))
			}
			return buf.String()
		}(),
		func() string {
			var buf bytes.Buffer
			if ri.untracked > 0 {
				buf.WriteString(untrackedGlyph)
			} else {
				buf.WriteRune(' ')
			}
			if ri.hasUnmerged() {
				buf.WriteString(unmergedGlyph)
			} else {
				buf.WriteRune(' ')
			}
			if ri.hasModified() {
				buf.WriteString(modifiedGlyph)
			} else {
				buf.WriteRune(' ')
			}
			// TODO star glyph
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
	gitOut, err := GetGitOutput(cwd)
	if err != nil {
		log.Printf("error: %s", err)
		if err == ErrNotAGitRepo {
			os.Exit(0)
		}
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	var repoInfo = new(RepoInfo)
	repoInfo.workingDir = cwd

	if err := repoInfo.ParseRepoInfo(gitOut); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	return repoInfo
}
