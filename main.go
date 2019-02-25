package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
)

const (
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
	defaultFormat = "%g %b%a %m%d%u%t %s"
)

// Options defines command line args and options
type Options struct {
	NoColor              bool
	Verbose              bool
	Version              bool
	Dir                  string
	Timeout              int16
	Format               string
	NoGitTag             bool
	Output               string
	ShowVCS              bool
	ShowAheadBehind      bool
	ShowBranch           bool
	ShowRemote           bool
	ShowCommit           bool
	ShowUnstagedModified bool
	ShowStagedModified   bool
	ShowUnknown          bool
	ShowStash            bool
	ShowDiff             bool
}

var (
	cwd     string
	options Options
)

func init() {
	// Use cli args if present, else test args
	args := (func() []string {
		if len(os.Args) > 1 {
			return os.Args[1:]
		}
		// Insert test args here
		// To be used if no os.Args
		return []string{}
	})()

	flag.BoolVar(&options.NoColor, "n", false, "do not print color on prompt")
	flag.BoolVar(&options.Verbose, "v", false, "print verbose debug messages")
	flag.BoolVar(&options.Version, "version", false, "show version info and exit")
	flag.StringVar(&options.Dir, "d", "", "git repo location, if not cwd")
	flag.StringVar(&options.Format, "f", defaultFormat, "printf-style format string for git prompt")
	flag.StringVar(&options.Output, "o", "string", "output type: string, raw, {1,2,3...}")
	flag.BoolVar(&options.NoGitTag, "no-tag", false, "do not look for git tag if detached head")

	epilog := `
	Output Examples:

	[-o=s/string]
	  Prints based on [-f] FORMAT, which may contain:
	  %g  branch glyph ()
	  %n  VC name
	  %b  branch
	  %r  remote
	  %a  commits ahead/behind remote
	  %c  current commit hash
	  %m  unstaged changes (modified/added/removed)
	  %s  staged changes (modified/added/removed)
	  %u  untracked files
	  %d  diff lines, ex: "+20/-10"
	  %t  stashed files indicator

	[-o=r/raw]
	  Prints each value on a new line for easy parsing

	[-o={1,2,3...}]
	  Presets: sensible presets for ease of use
		
	  1: [%n:%b] (vcprompt default)
	  2: %b %c %a %u %m
	  3: %g %b@%c %a %u %m %s (similar to porcelain)
	`
	flag.Usage = func() {
		usageMsg := `
		Usage: gitprompt [-h] [-v] [-d DIR] [-f FORMAT]

		Git status for your prompt, similar to Greg Ward's vcprompt.

		Optional Arguments:
		`
		fmt.Fprintln(os.Stderr, detent(usageMsg))
		flag.PrintDefaults()
		fmt.Println(detent(epilog))
	}
	flag.Parse()

	// Discard logs unless --verbose is set
	logFile := ioutil.Discard

	if options.Verbose {
		logFile = os.Stderr
	}

	log.SetOutput(logFile)
	log.Printf("Raw args: %v", args)

	remaining := flag.Args()
	if len(remaining) > 0 {
		log.Printf("Remaining args: %+v", flag.Args())
	}

	if options.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	// Handle regular options
	if options.Dir != "" {
		cwd = options.Dir
	}

	if options.NoColor {
		color.NoColor = true
	}

	if options.Output == "r" {
		options.Output = "raw"
	}
	if options.Output == "s" {
		options.Output = "string"
	}

	presets := [3]string{
		"[%n:%b]",
		"%b %c %a %u %m",
		"%g %b@%c %a %u %m %s",
	}

	switch options.Output {
	case "raw":
		options.ShowAheadBehind = true
		options.ShowBranch = true
		options.ShowDiff = true
		options.ShowCommit = true
		options.ShowStagedModified = true
		options.ShowStash = true
		options.ShowUnknown = true
		options.ShowUnstagedModified = true
	case "1":
		options.Format = presets[0]
		options.Output = "string"
	case "2":
		options.Format = presets[1]
		options.Output = "string"
	case "3":
		options.Format = presets[2]
		options.Output = "string"
	case "string":
	default:
		fmt.Printf("error: invalid output format `%v'", options.Output)
		os.Exit(1)
	}

	if cwd == "" {
		var err error
		if cwd, err = os.Getwd(); err != nil {
			log.Printf("Error getting cwd: %s", err)
		}
	}
}

func parseFormatString() {
	format := options.Format
	for i := 0; i < len(format); i++ {
		if string(format[i]) == "%" {
			i++
			switch string(format[i]) {
			case "a":
				options.ShowAheadBehind = true
			case "n":
				options.ShowVCS = true
			case "b":
				options.ShowBranch = true
			case "c":
				options.ShowCommit = true
			case "u":
				options.ShowUnknown = true
			case "m":
				options.ShowUnstagedModified = true
			case "s":
				options.ShowStagedModified = true
			case "d":
				options.ShowDiff = true
			case "t":
				options.ShowStash = true
			case "g": // Show branch glyph
			case "%":
			default:
				fmt.Fprintf(os.Stderr, "error: invalid format string '%%%c'", format[i])
				os.Exit(1)
			}
		}
	}
}

func main() {
	log.Printf("Running gitprompt in directory %s", cwd)

	if options.Output == "string" {
		parseFormatString()
		fmt.Println(run().fmtString())
	} else {

		fmt.Println(run().FmtRaw())
	}

	// fmt.Print(run().Fmt())
	log.Printf("Options: %+v", options)
}
