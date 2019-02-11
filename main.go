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
	defaultFormat = "%g %b%a|%m%d%u%t|%s"
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
	ShowRevision         bool
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
	flag.StringVar(&options.Output, "o", "string", "output type: string, raw")
	flag.BoolVar(&options.NoGitTag, "no-tag", false, "do not look for git tag if detached head")

	epilog := `
	Output Examples:

	[-o=string]
	  Prints based on [-f] FORMAT, which may contain:
	  %g  branch glyph ()
	  %n  VC name
	  %b  commits ahead/behind remote
	  %b  branch
	  %r  current commit hash
	  %m  unstaged changes (modified/added/removed)
	  %s  staged changes (modified/added/removed)
	  %u  untracked files
	  %a  added files
	  %d  diff lines, ex: "+20/-10"
	  %t  stashed files indicator

	[-o=raw]
	  Prints each value on a new line for easy parsing
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

	if len(flag.Args()) > 0 {
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
			case "r":
				options.ShowRevision = true
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
		fmt.Print(run().fmtString())
	} else {
		fmt.Print(run().Fmt())
	}
	log.Printf("Options: %+v", options)
}
