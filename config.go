package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type ConfigParser struct {
	configFile string
	tasks      []Task
	logger     *Logger
}

func NewConfigParser(configPath string, verboseLevel int) (*ConfigParser, error) {
	configFile := filepath.Join(configPath, "dosh.lua")
	logger := NewLogger(argVerbose)

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("dosh_commands", DoshLuaLoader)
	if err := L.DoFile(configFile); err != nil {
		panic(err)
	}

	logger.logDebug(fmt.Sprintf("Using config file: %s", configFile))
	return &ConfigParser{configFile: configFile, tasks: globalTasks, logger: logger}, nil
}

func (c *ConfigParser) getTasks() []Task {
	return c.tasks
}

func (c *ConfigParser) getDescription() string {
	description := os.Getenv("HELP_DESCRIPTION")
	if description == "" {
		description = "dosh - shell-independent task manager"
	}
	return description
}

func (c *ConfigParser) getEpilog() string {
	return os.Getenv("HELP_EPILOG")
}

func (c *ConfigParser) runTask(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task specified")
	}

	taskName := args[0]

	for _, task := range c.tasks {
		if task.Name == taskName {
			c.logger.logDebug("Running task: " + task.Name)
			runCommand(task.Command)
			return nil
		}
	}

	return fmt.Errorf("Task not found: %s", taskName)
}

func (c *ConfigParser) generateHelpOutput() string {
	helpOutput := []string{}
	description := c.getDescription()

	if description != "" {
		helpOutput = append(helpOutput, description, "")
	}

	tasks := c.getTasks()
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

	epilog := c.getEpilog()
	if epilog != "" {
		helpOutput = append(helpOutput, "", epilog)
	}

	return strings.Join(helpOutput, "\n")
}

func runCommand(s string) {
	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
