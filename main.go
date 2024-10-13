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

func generateHelpOutput() string {
	// Common DOSH commands
	helpOutput := []string{
		"DOSH commands:",
		"  > help                 print this output",
		"  > init                 initialize a new config in current working directory",
		"  > version              print version of DOSH",
		"",
		"  -c, --config PATH      specify config path (default: dosh.lua)",
		"  -d, --directory PATH   change the working directory",
		"  -v|vv|vvv, --verbose   increase the verbosity of messages:",
		"                         1 - default, 2 - detailed, 3 - debug",
	}

	// TODO: prepare description

	// TODO: prepare tasks section
	return strings.Join(helpOutput, "\n")
}

func main() {
	flag.Parse()

	command := parseCommand(flag.Args())
	if command == CommandHelp {
		fmt.Println(generateHelpOutput())
		return
	}

	// l := lua.NewState()
	// lua.OpenLibraries(l)
	//
	//	if err := lua.DoFile(l, "hello.lua"); err != nil {
	//		panic(err)
	//	}
}
