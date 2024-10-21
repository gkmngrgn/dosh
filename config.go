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
	configFile := configPath
	if !strings.HasSuffix(configPath, ".lua") {
		configFile = filepath.Join(configPath, "dosh.lua")
	}

	logger := GetLogger()
	luaState := lua.NewState()
	luaState.PreloadModule("dosh_commands", DoshLuaLoader)

	if fileExists(configFile) {
		logger.logDefault(fmt.Sprintf("Using config file: %s", configFile))
		if err := luaState.DoFile(configFile); err != nil {
			panic(err)
		}
	}

	return &ConfigParser{configFile: configFile, tasks: globalTasks, logger: logger, luaState: luaState}, nil
}

func (cp *ConfigParser) close() {
	defer cp.luaState.Close()
}

func (cp *ConfigParser) init() error {
	if fileExists(cp.configFile) {
		return fmt.Errorf("Config file already exists: %s", cp.configFile)
	}

	file, err := os.Create(cp.configFile)
	if err != nil {
		return err
	}

	_, err = file.WriteString(`
local cmd = require("dosh_commands")          -- main module to run DOSH commands.
local name = "there"                          -- you can use all features of Lua programming language.

local function hello(there)                   -- even you can define your custom functions.
    there = there or name
    local message = "Hello, " .. there .. "!"
    cmd.run("osascript -e 'display notification \"" .. message .. "\" with title \"Hi!\"'")
end

cmd.add_task{                                 -- cmd comes from dosh.
    name="hello",                             -- task name, or subcommand for your cli.
    description="say hello",                  -- task description for the help output.
    required_commands={"osascript"},          -- check if the programs exist before running the task.
    required_platforms={"macos"},             -- check if the current operating system is available to run the task.
    environments={"development", "staging"},  -- DOSH_ENV variable should be either development or staging to run this task.
    command=hello                             -- run hello function with its parameters when the task ran.
}
`)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return !info.IsDir()
}
