package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const (
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
)

// Options defines command line arguments
// type Options struct {
//         Verbose bool   `short:"v" long:"verbose" description:"see debug messages"`
//         Version bool   `long:"version" description:"show version info and exit"`
//         Dir     string `short:"d" long:"dir" description:"git repo location, if not cwd" value-name:"directory"`
//         Timeout int16  `short:"t" long:"timeout" description:"timeout for git cmds in ms" value-name:"timeout_ms" default:"100"`
//         Format  string `short:"f" long:"format" description:"printf-style format string for git prompt" value-name:"FORMAT" default:"[%n:%b]"`
// }

// Options defines command line args
type Options struct {
	Verbose  bool
	Version  bool
	Dir      string
	Timeout  int16
	Format   string
	NoGitTag bool
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

	// var parser = flags.NewParser(&options, flags.Default)

	flag.BoolVar(&options.Verbose, "v", false, "print verbose debug messages")
	flag.BoolVar(&options.Version, "version", false, "show version info and exit")
	flag.StringVar(&options.Dir, "d", ".", "git repo location, if not cwd")
	flag.StringVar(&options.Format, "f", "[%n:%b]", "printf-style format string for git prompt")
	flag.BoolVar(&options.NoGitTag, "no-tag", false, "do not look for git tag if detached head")

	longDesc := `
	FORMATTING
	----------
	Prints according to FORMAT, which may contain:
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
	`

	flag.Usage = func() {
		usageMsg := `
		Usage: gitprompt [-h] [-v] [-d DIR] [-f FORMAT]

		Git status for your prompt, similar to Greg Ward's vcprompt.

		Arguments:
		`
		fmt.Fprintln(os.Stderr, detent(usageMsg))
		flag.PrintDefaults()
		fmt.Println(detent(longDesc))
	}
	flag.Parse()

	// parser.LongDescription = longDesc
	// extraArgs, err := parser.ParseArgs(args)
	//
	// if err != nil {
	//         if _, ok := err.(*flags.Error); ok {
	//                 typ := err.(*flags.Error).Type
	//                 switch {
	//                 case typ == flags.ErrHelp:
	//                         break
	//                 case typ == flags.ErrCommandRequired && len(extraArgs[0]) == 0:
	//                         parser.WriteHelp(os.Stdout)
	//                 default:
	//                         log.Println(err.Error() + string(typ))
	//                         parser.WriteHelp(os.Stdout)
	//                 }
	//         } else {
	//                 fmt.Printf("Exiting: %s", err.Error())
	//         }
	//         os.Exit(1)
	// }

	// Discard logs unless --verbose is set
	logFile := ioutil.Discard

	if options.Verbose {
		logFile = os.Stderr
	}

	log.SetOutput(logFile)

	log.Printf("Raw args: %v", args)
	log.Printf("Parsed args: %+v", options)

	log.Println(flag.Args())

	// if len(extraArgs) > 0 {
	//         log.Printf("Remaining args: %v", extraArgs)
	// }

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

func detent(s string) string {
	return regexp.MustCompile("(?m)^[\t]*").ReplaceAllString(s, "")
}

func main() {
	log.Printf("Running gitprompt in directory %s", cwd)
	// TODO: switch here to determine whether to output prompt or raw data
	fmt.Fprint(os.Stdout, run().Fmt())
}
