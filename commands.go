package main

import (
	lua "github.com/yuin/gopher-lua"
)

type Task struct {
	Name        string
	Description string
	Command     string
}

// Global variable to store tasks.
var globalTasks []Task

func DoshLuaLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "name", lua.LString("value"))
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"add_task": addTask,
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
