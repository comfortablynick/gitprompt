package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
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
	// Branch
	branch   string
	commit   string
	remote   string
	upstream string
	ahead    int
	behind   int
	// Totals
	untracked int
	unmerged  int
	// Status for unstaged/staged
	Unstaged GitArea
	Staged   GitArea
}

func (ri *RepoInfo) hasUnmerged() bool {
	if ri.unmerged > 0 {
		return true
	}
	gitDir, err := PathToGitDir(cwd)
	if err != nil {
		log.Errorf("error calling PathToGitDir: %s", err)
		return false
	}
	// TODO figure out if output of MERGE_HEAD can be useful
	if _, err := ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Errorf("error reading MERGE_HEAD: %s", err)
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
	log.Infof("Formatting output: %s", ri.Debug())

	var (
		branchGlyph    = ""
		modifiedGlyph  = "Δ"
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
				if _, err := buf.WriteString(fmt.Sprintf(" %s%d ", aheadArrow, ri.ahead)); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			}
			if ri.behind > 0 {
				if _, err := buf.WriteString(fmt.Sprintf(" %s%d ", behindArrow, ri.behind)); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			}
			return buf.String()
		}(),
		func() string {
			var buf bytes.Buffer
			if ri.untracked > 0 {
				if _, err := buf.WriteString(untrackedGlyph); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			}
			if ri.hasUnmerged() {
				if _, err := buf.WriteString(unmergedGlyph); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			}
			if ri.hasModified() {
				if _, err := buf.WriteString(modifiedGlyph); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
			} else {
				if _, err := buf.WriteRune(' '); err != nil {
					log.Errorf("Buffer error: %s", err)
				}
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
	gitOut, err := GetGitStatusOutput(cwd)
	if err != nil {
		// Just log this as Info so that we don't return
		// any output by default when not in a repo
		log.Infoln(err)
		if err == ErrNotAGitRepo {
			os.Exit(1)
		}
		log.Errorf("error: %s", err)
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
