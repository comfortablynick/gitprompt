package main

import (
	"strings"
	"testing"
)

const expectedFmtOutput = ` [34mmaster[0m@[32;3m51c9c58[0m [2;43;30m ↑1 [0m[2;41;37m ↓10 [0m [2m?[0m[36m‼[0m[34mΔ[0m [31m✘[0m`

func TestFmtOutput(t *testing.T) {
	var ri = new(RepoInfo)
	if err := ri.ParseRepoInfo(strings.NewReader(gitoutput)); err != nil {
		t.Fatal(err)
	}

	if out := ri.Fmt(); out != expectedFmtOutput {
		t.Logf("\nexpected:\n%s\ngot:\n%s\n", expectedFmtOutput, out)
		t.FailNow()
	}
}
