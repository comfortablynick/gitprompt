package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

const (
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
)

var cwd string

// Options defines command line arguments
type Options struct {
	Verbose bool   `short:"v" long:"verbose" description:"see debug messages"`
	Version bool   `long:"version" description:"show version info and exit"`
	Dir     string `short:"d" long:"dir" description:"git repo location, if not cwd" value-name:"directory"`
	Timeout int16  `short:"t" long:"timeout" description:"timeout for git cmds in ms" value-name:"timeout_ms" default:"100"`
	Format  string `short:"f" long:"format" description:"printf-style format string for git prompt" value-name:"FORMAT" default:"[%n:%b]"`
}

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

	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	longDesc := `Git status for your prompt, similar to Greg Ward's vcprompt

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
	parser.LongDescription = longDesc
	extraArgs, err := parser.ParseArgs(args)

	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			typ := err.(*flags.Error).Type
			switch {
			case typ == flags.ErrHelp:
				break
			case typ == flags.ErrCommandRequired && len(extraArgs[0]) == 0:
				parser.WriteHelp(os.Stdout)
			default:
				log.Println(err.Error() + string(typ))
				parser.WriteHelp(os.Stdout)
			}
		} else {
			fmt.Printf("Exiting: %s", err.Error())
		}
		os.Exit(1)
	}

	// Discard logs unless --verbose is set
	logFile := ioutil.Discard

	if options.Verbose {
		logFile = os.Stderr
	}

	log.SetOutput(logFile)

	log.Printf("Raw args: %v", args)
	log.Printf("Parsed args: %+v", options)
	if len(extraArgs) > 0 {
		log.Printf("Remaining args: %v", extraArgs)
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
		if cwd, err = os.Getwd(); err != nil {
			log.Printf("Error getting cwd: %s", err)
		}
	}
}

func main() {
	log.Printf("Running gitprompt in directory %s", cwd)
	// TODO: switch here to determine whether to output prompt or raw data
	fmt.Fprintln(os.Stdout, run().Fmt())
}
