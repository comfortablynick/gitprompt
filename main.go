package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/subchen/go-log"
)

const (
	logloc  = "/tmp/porcelain.log"
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
)

var logLevels = []log.Level{
	log.WARN,
	log.INFO,
	log.DEBUG,
}

var cwd string

// Options defines command line arguments
var Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"see more debug messages"`
	Version bool   `long:"version" description:"show version info and exit"`
	Dir     string `short:"d" long:"dir" description:"git repo location" value-name:"directory" default:"."`
	Timeout int16  `short:"t" long:"timeout" description:"timeout for git cmds in ms" value-name:"timeout_ms" default:"100"`
	Format  string `short:"f" long:"format" description:"printf-style format string for git prompt" value-name:"FORMAT" default:"[%n:%b]"`
}

func init() {
	// flag.BoolVar(&debugFlag, "debug", false, "write logs to file ("+logloc+")")
	// flag.BoolVar(&fmtFlag, "fmt", true, "print formatted output (default)")
	// flag.BoolVar(&bashFmtFlag, "bash", false, "escape fmt output for bash")
	// flag.BoolVar(&zshFmtFlag, "zsh", false, "escape fmt output for zsh")
	// flag.BoolVar(&tmuxFmtFlag, "tmux", false, "escape fmt output for tmux")
	// flag.StringVar(&cwd, "path", "", "show output for path instead of the working directory")
	// flag.BoolVar(&versionFlag, "version", false, "print version and exit")

	// logtostderr := flag.Bool("logtostderr", false, "write logs to stderr")
	// flag.Parse()
	log.Default.Level = log.WARN

	// Use cli args if present, else test args
	args := (func() []string {
		if len(os.Args) > 1 {
			return os.Args[1:]
		}
		log.Infoln("Using test arguments")
		return []string{
			"-v",
		}
	})()

	var parser = flags.NewParser(&Options, flags.Default)
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
		if !flags.WroteHelp(err) {
			parser.WriteHelp(os.Stderr)
		}
		os.Exit(1)
	}

	// Get log level
	verbosity, maxLevels := len(Options.Verbose), len(logLevels)
	if verbosity > maxLevels-1 {
		verbosity = maxLevels - 1
	}

	log.Default.Level = logLevels[verbosity]

	log.Debugf("Raw args:\n%v", args)
	log.Debugf("Parsed args:\n%+v", Options)
	if len(extraArgs) > 0 {
		log.Debugf("Remaining args:\n%v", extraArgs)
	}

	if Options.Version {
		fmt.Println(version)
		os.Exit(0)
	}
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
}

func main() {
	log.Debugf("Running gitprompt in directory %s", cwd)

	var out string
	switch {
	case true:
		out = run().Fmt()
	default:
		flag.Usage()
		fmt.Println("\nOutside of a repository there will be no output.")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, out)
}
