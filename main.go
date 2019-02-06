package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	Log "github.com/sirupsen/logrus"
)

// Log default logger
// var Log = new(log.Logger)

// Log Global logger
// var Log = &log.Logger{
//     Level:     log.WARN,
//     Formatter: new(formatters.TextFormatter),
//     Out:       os.Stderr,
// }

const (
	logloc  = "/tmp/porcelain.log"
	version = `
gitprompt version 0.0.1

© 2019 Nicholas Murphy
(github.com/comfortablynick)
`
)

var logLevels = []Log.Level{
	Log.WarnLevel,
	Log.InfoLevel,
	Log.DebugLevel,
}

var cwd string

// Options defines command line arguments
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"see more debug messages"`
	Version bool   `long:"version" description:"show version info and exit"`
	Dir     string `short:"d" long:"dir" description:"git repo location" value-name:"directory" default:"."`
	Timeout int16  `short:"t" long:"timeout" description:"timeout for git cmds in ms" value-name:"timeout_ms" default:"100"`
	Format  string `short:"f" long:"format" description:"printf-style format string for git prompt" value-name:"FORMAT" default:"[%n:%b]"`
}

func init() {
	// log.Default.Level = log.WARN
	// Log.Level = Log.WarnLevel
	Log.SetLevel(Log.WarnLevel)

	// Use cli args if present, else test args
	args := (func() []string {
		if len(os.Args) > 1 {
			return os.Args[1:]
		}
		Log.Infoln("Using test arguments")
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
				Log.Info(err.Error() + string(typ))
				parser.WriteHelp(os.Stdout)
			}
		} else {
			Log.Fatalf("Exiting: %s", err.Error())
		}
		os.Exit(1)
	}

	// Get log level
	verbosity, maxLevels := len(options.Verbose), len(logLevels)
	if verbosity > maxLevels-1 {
		verbosity = maxLevels - 1
	}

	Log.SetLevel(logLevels[verbosity])

	Log.Debugf("Raw args:\n%v", args)
	Log.Debugf("Parsed args:\n%+v", options)
	if len(extraArgs) > 0 {
		Log.Debugf("Remaining args:\n%v", extraArgs)
	}

	if options.Version {
		fmt.Println(version)
		os.Exit(0)
	}
	if cwd == "" {
		cwd, _ = os.Getwd() // #nosec
	}
}

func main() {
	Log.Debugf("Running gitprompt in directory %s", cwd)

	// var out string
	// switch {
	// case true:
	//         out = run().Fmt()
	// default:
	//         flag.Usage()
	//         fmt.Println("\nOutside of a repository there will be no output.")
	//         os.Exit(1)
	// }
	fmt.Fprint(os.Stdout, run().Fmt())
}
