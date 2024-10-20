package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type ConfigParser struct {
	configFile string
	tasks      []Task
	luaState   *lua.LState
	logger     *Logger
}

func NewConfigParser(configPath string) (*ConfigParser, error) {
	configFile := filepath.Join(configPath, "dosh.lua")
	logger := GetLogger()

	L := lua.NewState()
	L.PreloadModule("dosh_commands", DoshLuaLoader)
	if err := L.DoFile(configFile); err != nil {
		panic(err)
	}

	logger.logDefault(fmt.Sprintf("Using config file: %s", configFile))
	return &ConfigParser{configFile: configFile, tasks: globalTasks, logger: logger, luaState: L}, nil
}

func (cp *ConfigParser) close() {
	defer cp.luaState.Close()
}

func (cp *ConfigParser) getTasks() []Task {
	return cp.tasks
}

func (cp *ConfigParser) getDescription() string {
	description := os.Getenv("HELP_DESCRIPTION")
	if description == "" {
		description = "dosh - shell-independent task manager"
	}
	return description
}

func (cp *ConfigParser) getEpilog() string {
	return os.Getenv("HELP_EPILOG")
}

func (cp *ConfigParser) runTask(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task specified")
	}

	taskName := args[0]

	for _, task := range cp.tasks {
		if task.Name == taskName {
			cp.logger.logDefault("Running task: " + task.Name)
			task.run(cp.luaState)
			return nil
		}
	}

	return fmt.Errorf("Task not found: %s", taskName)
}

func (cp *ConfigParser) generateHelpOutput() string {
	helpOutput := []string{}
	description := cp.getDescription()

	if description != "" {
		helpOutput = append(helpOutput, description, "")
	}

	tasks := cp.getTasks()
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

	epilog := cp.getEpilog()
	if epilog != "" {
		helpOutput = append(helpOutput, "", epilog)
	}

	return strings.Join(helpOutput, "\n")
}
