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
}

type Command string

const (
	CommandHelp    Command = "help"
	CommandInit    Command = "init"
	CommandVersion Command = "version"
	CommandUnknown Command = "unknown"
)

func parseCommand(args []string) Command {
	if len(args) == 0 {
		return CommandHelp
	}

	switch Command(args[0]) {
	case CommandHelp, CommandInit, CommandVersion:
		return Command(args[0])
	default:
		return CommandUnknown
	}
}

func main() {
	flag.Usage = func() {} // disable default usage message
	flag.Parse()           // parse arguments
	args := flag.Args()    // get arguments
	logger := GetLogger()  // set verbosity level
	logger.setVerbosity(argVerbose)

	// initialize config parser
	configParser, err := NewConfigParser(argConfigPath)
	if err != nil {
		logger.logError(err.Error())
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

	configParser.close()
}
