package main

import (
	"flag"
	"fmt"
)

var (
	argConfigPath string
	argVerbose    int
)

func init() {
	flag.StringVar(&argConfigPath, "c", ".", "Path to the configuration file")
	flag.IntVar(&argVerbose, "v", 0, "Enable verbose output")

	initLogger()
}

type Command int

const (
	CommandHelp Command = iota
	CommandInit
	CommandVersion
	CommandUnknown
)

func parseCommand(args []string) Command {
	if len(args) == 0 {
		// No command provided, print help
		return CommandHelp
	}

	switch args[0] {
	case "help":
		return CommandHelp
	case "init":
		return CommandInit
	case "version":
		return CommandVersion
	default:
		return CommandUnknown
	}
}

func main() {
	flag.Usage = func() {} // disable default usage message
	flag.Parse()           // parse arguments
	args := flag.Args()
	logger := NewLogger(argVerbose)

	// initialize config parser
	configParser, err := NewConfigParser(argConfigPath, argVerbose)
	if err != nil {
		logger.logDebug(err.Error())
		return
	}

	// run command looking at the first argument
	switch parseCommand(args) {
	case CommandHelp:
		fmt.Println(configParser.generateHelpOutput())
	case CommandVersion:
		fmt.Println(getVersion())
	default:
		if err := configParser.runTask(args); err != nil {
			fmt.Println(err)
			fmt.Println("Run 'dosh help' for usage information")
		}
	}
}
