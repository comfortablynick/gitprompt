package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
)

// Options defines command line args
type Options struct {
	Verbose  bool
	Version  bool
	Dir      string
	Timeout  int16
	Format   string
	NoGitTag bool
	Output   string
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

	flag.BoolVar(&options.Verbose, "v", false, "print verbose debug messages")
	flag.BoolVar(&options.Version, "version", false, "show version info and exit")
	flag.StringVar(&options.Dir, "d", ".", "git repo location, if not cwd")
	flag.StringVar(&options.Format, "f", "[%n:%b]", "printf-style format string for git prompt")
	flag.StringVar(&options.Output, "o", "string", "output type: string, raw")
	flag.BoolVar(&options.NoGitTag, "no-tag", false, "do not look for git tag if detached head")

	epilog := `
	Output Examples:

	[-o=string]
	  Prints based on [-f] FORMAT, which may contain:
	  %n  show VC name
	  %b  show branch
	  %r  show remote (default: ".")
	  %m  indicate uncomitted changes (modified/added/removed)
	  %u  show untracked file count
	  %a  show added file count
	  %d  show deleted file count
	  %s  show stash count
	  %x  show insertion count
	  %y  show deletion count

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
	log.Printf("Parsed args: %+v", options)

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

	if cwd == "" {
		var err error
		if cwd, err = os.Getwd(); err != nil {
			log.Printf("Error getting cwd: %s", err)
		}
	}
}

func main() {
	log.Printf("Running gitprompt in directory %s", cwd)
	fmt.Print(run().Fmt())
}
