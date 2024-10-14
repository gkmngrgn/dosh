package main

import (
	"flag"
	"fmt"
	"strings"
)

type Command int

const (
	CommandHelp Command = iota
	CommandInit
	CommandVersion
	CommandUnknown
)

type Task struct {
	Name        string
	Description string
	Command     string
}

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
		"  -c, --config PATH      specify config path (default: dosh.lua)",
		"  -d, --directory PATH   change the working directory",
		"  -v|vv|vvv, --verbose   increase the verbosity of messages:",
		"                         1 - default, 2 - detailed, 3 - debug",
	)

	if epilog != "" {
		helpOutput = append(helpOutput, "", epilog)
	}

	return strings.Join(helpOutput, "\n")
}

func main() {
	// parse arguments
	flag.Parse()
	args := flag.Args()

	// initialize config parser
	configParser := NewConfigParser()

	// run command looking at the first argument
	switch parseCommand(args) {
	case CommandHelp:
		tasks := configParser.getTasks()
		description := configParser.getDescription()
		epilog := configParser.getEpilog()
		fmt.Println(generateHelpOutput(tasks, description, epilog))
	case CommandVersion:
		fmt.Println(getVersion())
	}

	// l := lua.NewState()
	// lua.OpenLibraries(l)
	//
	//	if err := lua.DoFile(l, "hello.lua"); err != nil {
	//		panic(err)
	//	}
}
