package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func consumeNext(s *bufio.Scanner) string {
	if s.Scan() {
		return s.Text()
	}
	return ""
}

// ParseRepoInfo begins parsing data returned from `git status`
func (ri *RepoInfo) ParseRepoInfo(r io.Reader) error {
	var s = bufio.NewScanner(r)

	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}
		ri.ParseLine(s.Text())
	}
	return nil
}

// ParseLine parses each line of `git status` porcelain v2 output
func (ri *RepoInfo) ParseLine(line string) error {
	s := bufio.NewScanner(strings.NewReader(line))
	// switch to a word based scanner
	s.Split(bufio.ScanWords)

	for s.Scan() {
		switch s.Text() {
		case "#":
			_ = ri.parseBranchInfo(s)
		case "1":
			ri.parseTrackedFile(s)
		case "2":
			ri.parseRenamedFile(s)
		case "u":
			ri.unmerged++
		case "?":
			ri.untracked++
		}
	}
	return nil
}

func (ri *RepoInfo) parseBranchInfo(s *bufio.Scanner) (err error) {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		switch s.Text() {
		case "branch.oid":
			ri.commit = consumeNext(s)
		case "branch.head":
			ri.branch = consumeNext(s)
			if ri.branch == "(detached)" {
				if tag, err := GetGitTag(cwd); err == nil {
					ri.branch = tag
				}
			}
		case "branch.upstream":
			ri.upstream = consumeNext(s)
		case "branch.ab":
			err = ri.parseAheadBehind(s)
		}
	}
	return err
}

func (ri *RepoInfo) parseAheadBehind(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		i, err := strconv.Atoi(s.Text()[1:])
		if err != nil {
			return err
		}

		switch s.Text()[:1] {
		case "+":
			ri.ahead = i
		case "-":
			ri.behind = i
		}
	}
	return nil
}

// parseTrackedFile parses the porcelain v2 output for tracked entries
// doc: https://git-scm.com/docs/git-status#_changed_tracked_entries
func (ri *RepoInfo) parseTrackedFile(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	var index int
	for s.Scan() {
		switch index {
		case 0: // xy
			ri.parseXY(s.Text())
		default:
			continue
		}
		index++
	}
	return nil
}

func (ri *RepoInfo) parseXY(xy string) error {
	switch xy[:1] { // parse staged
	case "M":
		ri.Staged.modified++
	case "A":
		ri.Staged.added++
	case "D":
		ri.Staged.deleted++
	case "R":
		ri.Staged.renamed++
	case "C":
		ri.Staged.copied++
	}

	switch xy[1:] { // parse unstaged
	case "M":
		ri.Unstaged.modified++
	case "A":
		ri.Unstaged.added++
	case "D":
		ri.Unstaged.deleted++
	case "R":
		ri.Unstaged.renamed++
	case "C":
		ri.Unstaged.copied++
	}
	return nil
}

func (ri *RepoInfo) parseRenamedFile(s *bufio.Scanner) error {
	return ri.parseTrackedFile(s)
}
