package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

func (c *ConfigParser) runTask(args []string) {
	taskName := args[0]

	for _, task := range c.tasks {
		if task.Name == taskName {
			c.logger.logDebug("Running task: " + task.Name)
			runCommand(task.Command)
			return
		}
	}

	c.logger.logDebug("Task not found: " + taskName)
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
