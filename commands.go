package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Task struct {
	Name        string
	Description string
	Command     *lua.LFunction
}

func (t Task) run(L *lua.LState) {
	if err := L.CallByParam(lua.P{
		Fn:      t.Command,
		NRet:    1,
		Protect: true,
	}); err != nil {
		fmt.Println("Error running task: " + t.Name)
		fmt.Println(err)
	}
}

var (
	globalTasks             = []Task{}
	editableEnvironmentKeys = []string{
		"HELP_DESCRIPTION",
		"HELP_EPILOG",
	}
	exports = map[string]lua.LGFunction{
		"add_task": addTask,
		"set_env":  setEnv,
		"run":      runCommand,
		"copy":     copyFiles,
	}
)

func DoshLuaLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "name", lua.LString("value"))
	L.Push(mod)
	return 1
}

func addTask(L *lua.LState) int {
	params := L.CheckTable(1)
	task := Task{
		Name:        params.RawGetString("name").String(),
		Description: params.RawGetString("description").String(),
		Command:     params.RawGetString("command").(*lua.LFunction),
	}
	globalTasks = append(globalTasks, task)
	return 1
}

func setEnv(L *lua.LState) int {
	params := L.CheckTable(1)
	params.ForEach(func(key lua.LValue, value lua.LValue) {
		k := trimSpaces(key.String())
		if !isEditableKey(k) {
			return
		}

		v := trimSpaces(value.String())
		if err := os.Setenv(k, v); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return
		}
	})

	L.Push(lua.LTrue)
	return 1
}

func runCommand(L *lua.LState) int {
	logger := GetLogger()
	command := trimSpaces(L.CheckString(1))
	logger.logDefault("Running command: " + command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1
}

func copyFiles(L *lua.LState) int {
	logger := GetLogger()
	srcPattern := L.CheckString(1)
	dstDir := L.CheckString(2)

	logger.logDefault(fmt.Sprintf("Copying files from %s to %s", srcPattern, dstDir))

	// Find files matching the glob pattern
	files, err := filepath.Glob(srcPattern)
	if err != nil {
		logger.logError(fmt.Sprintf("Failed to parse glob pattern: %v", err))
		return 2
	}

	// Expand tilde to home directory
	if strings.HasPrefix(dstDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.logError(fmt.Sprintf("Failed to get home directory: %v", err))
			return 2
		}
		dstDir = filepath.Join(homeDir, dstDir[1:])
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		logger.logError(fmt.Sprintf("Failed to create destination directory: %v", err))
		return 2
	}

	// Copy each file
	for _, src := range files {
		dst := filepath.Join(dstDir, filepath.Base(src))
		logger.logDetailed("Copying file " + src + " to " + dst)

		if err := copyFile(src, dst); err != nil {
			logger.logError(fmt.Sprintf("failed to copy file %s to %s: %v", src, dst, err))
			return 2
		}
	}

	return 1
}

func isEditableKey(k string) bool {
	for _, key := range editableEnvironmentKeys {
		if key == k {
			return true
		}
	}
	return false
}

func trimSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}
