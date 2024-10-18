package main

import (
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Task struct {
	Name        string
	Description string
	Command     string
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
		Command:     params.RawGetString("command").String(),
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
