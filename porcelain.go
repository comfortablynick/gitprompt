package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
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
	gitDir     string
	branch     string
	commit     string
	remote     string
	upstream   string
	stashed    bool
	ahead      int
	behind     int
	untracked  int
	unmerged   int
	insertions int
	deletions  int
	Unstaged   GitArea
	Staged     GitArea
}

func (ri *RepoInfo) hasUnmerged() bool {
	if ri.unmerged > 0 {
		return true
	}
	if ri.gitDir == "" {
		var err error
		ri.gitDir, err = PathToGitDir(cwd)
		if err != nil {
			log.Printf("error calling PathToGitDir: %s", err)
			return false
		}
	}
	// TODO: figure out if output of MERGE_HEAD can be useful
	if _, err := os.Stat(path.Join(ri.gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("error reading MERGE_HEAD: %s", err)
		return false
	}
	return true
}

func (ri *RepoInfo) hasStash() bool {
	if ri.gitDir == "" {
		var err error
		ri.gitDir, err = PathToGitDir(cwd)
		if err != nil {
			log.Printf("error calling PathToGitDir: %s", err)
			return false
		}
	}
	if _, err := os.Stat(path.Join(ri.gitDir, "logs/refs/stash")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("error reading stash: %s", err)
		return false
	}
	return true
}

// Debug prints repo info
func (ri *RepoInfo) Debug(pretty bool) string {
	if !pretty {
		return fmt.Sprintf("%+v", ri)
	}
	return detent(fmt.Sprintf(`
	RepoInfo
	========
	workingDir: %v
	gitDir:     %v
	branch:     %v
	commit:     %v
	remote:     %v
	upstream:   %v
	stashed:	%-v
	ahead:      %4d
	behind:     %4d
	untracked:  %4d
	unmerged:   %4d
	insertions: %4d
	deletions:  %4d

	Unstaged
	--------
	modified:   %4d
	added:      %4d
	deleted:    %4d
	renamed:    %4d
	copied:     %4d

	Staged
	--------
	modified:   %4d
	added:      %4d
	deleted:    %4d
	renamed:    %4d
	copied:     %4d`, ri.workingDir, ri.gitDir, ri.branch, ri.commit, ri.remote, ri.upstream,
		ri.stashed, ri.ahead, ri.behind, ri.untracked, ri.unmerged, ri.insertions, ri.deletions,
		ri.Unstaged.modified, ri.Unstaged.added, ri.Unstaged.deleted, ri.Unstaged.renamed,
		ri.Unstaged.copied, ri.Staged.modified, ri.Staged.added, ri.Staged.deleted,
		ri.Staged.renamed, ri.Staged.copied))
}

var (
	branchGlyph    = ""
	modifiedGlyph  = "Δ"
	dirtyGlyph     = "✘" // ✗
	cleanGlyph     = "✔" // ✓
	untrackedGlyph = "?"
	unmergedGlyph  = "‼"
	aheadArrow     = "↑"
	behindArrow    = "↓"
	stashGlyph     = "$"
)

// TODO: parse first, then format if called for by user

// fmtCleanDirty changes color depending on repo status
func (ri *RepoInfo) fmtCleanDirty(s string) string {
	if ri.Unstaged.hasChanged() {
		return color.HiRedString(s)
	}
	if ri.Staged.hasChanged() {
		return color.HiYellowString(s)
	}
	return color.HiGreenString(s)
}

// Fmt formats the output for the shell
func (ri *RepoInfo) Fmt() string {
	// TODO: make format user-configurable
	log.Println(ri.Debug(false))

	// Turn off color based on CLI option
	color.NoColor = options.NoColor

	cleanDirtyFmt := (func() func(...interface{}) string {
		if ri.Unstaged.hasChanged() {
			return color.New(color.FgHiRed).SprintFunc()
		}
		if ri.Staged.hasChanged() {
			return color.New(color.FgHiYellow).SprintFunc()
		}
		return color.New(color.FgHiGreen).SprintFunc()
	})()

	return fmt.Sprintf("%s %s@%s %s %s %s %s",
		branchGlyph,
		cleanDirtyFmt(ri.branch),
		cleanDirtyFmt(func() string {
			if ri.commit == "(initial)" {
				return ri.commit
			}
			return ri.commit[:7]
		}()),
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
			if ri.untracked == 0 {
				untrackedGlyph = " "
			}
			if !ri.hasUnmerged() {
				unmergedGlyph = " "
			}
			if !ri.Unstaged.hasChanged() {
				modifiedGlyph = " "
			}
			if _, err := buf.WriteString(untrackedGlyph + unmergedGlyph + modifiedGlyph); err != nil {
				log.Printf("Error writing glyphs: %s", err)
			}
			return buf.String()
		}(),
		func() string {
			if ri.Staged.hasChanged() {
				return dirtyGlyph
			}
			return cleanGlyph
		}(),
		func() string {
			var out string
			if ri.insertions > 0 {
				out += fmt.Sprintf("+%d", ri.insertions)
			}
			if ri.deletions > 0 {
				out += fmt.Sprintf(" -%d", ri.deletions)
			}
			return strings.TrimSpace(out)
		}(),
	)
}

// TODO: define custom format function that takes color param

// fmtString parses user-supplied format string
func (ri *RepoInfo) fmtString() string {
	log.Println(ri.Debug(false))

	format := options.Format
	var out string
	for i := 0; i < len(format); i++ {
		if string(format[i]) == "%" {
			i++
			switch string(format[i]) {
			case "n":
				out += "git"
			case "b":
				out += fmt.Sprintf("%s %s", branchGlyph, ri.fmtCleanDirty(ri.branch))
			case "r":
				out += ri.fmtCleanDirty(func() string {
					if ri.commit == "(initial)" {
						return ri.commit
					}
					return ri.commit[:7]
				}())
			case "u":
				if ri.untracked > 0 {
					out += color.HiYellowString(untrackedGlyph)
				}
			case "m":
				if ri.Unstaged.modified > 0 {
					out += color.HiBlueString(modifiedGlyph)
				}
			case "d":
				if ri.insertions+ri.deletions != 0 {
					out += color.HiBlueString(func() string {
						if ri.insertions != 0 && ri.deletions != 0 {
							return fmt.Sprintf("+%d/-%d", ri.insertions, ri.deletions)
						}
						if ri.insertions != 0 {
							return fmt.Sprintf("+%d", ri.insertions)
						}
						return fmt.Sprintf("-%d", ri.deletions)
					}())
				}
			case "s":
				if ri.stashed {
					out += stashGlyph
				}
			case "%":
				out += "%"
			default:
				out += string(format[i])
			}
			continue
		}
		out += string(format[i])
	}
	return strings.TrimSpace(out)
}

func run() *RepoInfo {
	gitOut, err := GetGitStatusOutput(cwd)
	if err != nil {
		log.Printf("Git status error: %s", err)
		os.Exit(1)
	}

	var repoInfo = new(RepoInfo)
	repoInfo.workingDir = cwd

	if err = repoInfo.ParseRepoInfo(gitOut); err != nil {
		log.Printf("Error parsing git repo: %s", err)
		os.Exit(1)
	}

	// Only get diff when there are changes
	if repoInfo.Unstaged.hasChanged() && options.ShowDiff {
		diffOut, err := GetGitNumstat(cwd)
		if err != nil {
			log.Printf("Git diff error: %s", err)
		}

		if err = repoInfo.parseDiffNumstat(diffOut); err != nil {
			log.Printf("Error parsing git diff: %v", err)
		}
	}

	if options.ShowStash {
		repoInfo.stashed = repoInfo.hasStash()
	}
	return repoInfo
}
