package main

import (
	"flag"
	"fmt"
	"strings"
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

func generateHelpOutput(tasks []Task, description string, epilog string) string {
	helpOutput := []string{}

	if description != "" {
		helpOutput = append(helpOutput, description, "")
	}

	if len(tasks) > 0 {
		helpOutput = append(helpOutput, "Tasks:")

		for _, task := range tasks {
			helpOutput = append(helpOutput, fmt.Sprintf("  > %-20s %s", task.Name, task.Description))
		}

		helpOutput = append(helpOutput, "")
	}

	helpOutput = append(helpOutput,
		"DOSH commands:",
		"  > help                 print this output",
		"  > init                 initialize a new config in current working directory",
		"  > version              print version of DOSH",
		"",
		"  -c string              specify config path (default: dosh.lua)",
		"  -d string              change the working directory",
		"  -v int                 increase the verbosity of messages:",
		"                         1 - default, 2 - detailed, 3 - debug",
	)

	if epilog != "" {
		helpOutput = append(helpOutput, "", epilog)
	}

	return strings.Join(helpOutput, "\n")
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

	logger.logDebug(fmt.Sprintf("Using config file: %s", configParser.configFile))

	// run command looking at the first argument
	switch parseCommand(args) {
	case CommandHelp:
		tasks := configParser.getTasks()
		description := configParser.getDescription()
		epilog := configParser.getEpilog()
		fmt.Println(generateHelpOutput(tasks, description, epilog))
	case CommandVersion:
		fmt.Println(getVersion())
	default:
		configParser.runTask(args)
	}
}
