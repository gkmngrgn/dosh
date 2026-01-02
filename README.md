# DOSH - shell-independent task manager

**DOSH** is a command-line application to run your tasks on any platform, on any shell. Define your tasks, aliases, environments in a `dosh.lua` file and run `dosh`. DOSH will work like a CLI app reading your config file.

<img src="dosh-logo.svg" width="200" alt="DOSH logo" />

## INSTALLATION

Our GitHub Actions workflow builds app for these operating systems:

- Linux (aarch64, x86_64)
- MacOS (arm64)
- Windows (amd64)

### For Linux and MacOS

```shell
sh <(curl https://raw.githubusercontent.com/gkmngrgn/dosh/main/install.sh) --install-dir $HOME/.local/bin
```

### For Windows (PowerShell)

```powershell
iwr https://raw.githubusercontent.com/gkmngrgn/dosh/main/install.ps1 -useb | iex
```

## ANATOMY OF `dosh.lua`

```lua
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
```

When you run this command on MacOS, you will get a notification popup on the screen and see some logs in the console:

```shell
$ DOSH_ENV="development" dosh hello lua
DOSH => [RUN] osascript -e 'display notification "Hello, lua!" with title "Hi!"'
```

Take a look at the [`examples`](./examples) folder to find ready-in-use config files.

## ENVIRONMENT VARIABLES

#### HELP OUTPUT

Help outputs consist of four parts: **description**, **tasks**, **commands**, and **epilog**. The tasks will be generated getting task names and descriptions from your config file. The commands are including pre-defined dosh tasks and common task parameters. All help outputs start with a description and ends with an epilog if you have.

If you want to edit the default description and add an epilog to the help output, you can modify these variables:

- `env.HELP_DESCRIPTION`
- `env.HELP_EPILOG`

```
$ dosh help
dosh - shell-independent task manager                                           # HELP_DESCRIPTION HERE

Tasks:                                                                          # TASKS DEFINED BY THE USER
> hello                say hello

Dosh commands:
> help                 print this output
> init                 initialize a new config in current working directory
> version              print version of DOSH

-c, --config PATH      specify config path (default: dosh.lua)
-d, --directory PATH   change the working directory
-v|vv|vvv, --verbose   increase the verbosity of messages:
                       1 - default, 2 - detailed, 3 - debug

Wikipedia says that an epilog is a piece of writing at the end of a work of     # HELP_EPILOG HERE
literature, usually used to bring closure to the work.
```

#### OPERATING SYSTEM TYPE

All the following variables will return `true` or `false` depending on the operating system that you ran dosh:

- `env.IS_LINUX`
- `env.IS_MACOS`
- `env.IS_WINDOWS`

#### SHELL TYPE

It's like OS type checking. It's useful if you use shell-specific package like `ohmyzsh`.

- `env.IS_BASH`
- `env.IS_PWSH`
- `env.IS_ZSH`

#### DOSH-SPECIFIC ENVIRONMENTS

Consider you have some tasks that help you to test the project on your local and you want to restrict the task to prevent running it on the server by mistake. So the method `cmd.add_task` has an `environments` parameter and you can set your environment name for each target.

- `DOSH_ENV` (define it on your `~/.profile` file or CI/CD service)

_Check out the file [`dosh_environments.lua`](./examples/dosh_environments.lua) for example usage._

## COMMANDS

#### GENERAL PURPOSE

The main purpose of dosh to write one script that works on multiple operating systems and different shells. But it has to have a limit and it's nonsense to define functions for each cli command. So if you want to run a cli app (like `exa`, `bat`, `helix`, etc.), then you can use `cmd.run` for it.

`cmd.run` returns a list of results for each command executed. Each result object contains:
- `command_output`: The captured output of the command.
- `return_code`: The exit code of the command.

_Check out the file [`dosh_greet.lua`](./examples/dosh_greet.lua) for example usage._

#### FILE SYSTEM OPERATIONS

There are some ready-made functions both to keep the code readable and to make it work the same in all operating systems. You know Windows prefers backslash as a path separator but with dosh, use always `/` as in `/foo/bar/baz`, let dosh to find the path in a common way.

_Check out the file [`dosh_config.lua`](./examples/dosh_config.lua) for example usage._

#### PACKAGE MANAGERS

There are many package managers and I'm not sure if we need to implement all of them. But at least dosh supports these three of them mostly:

- `cmd.brew_install` (for MacOS and Linux)

- `packages`: list of strings, required.

- `cask`: boolean, default is `false`.

- `taps`: list of strings, optional.

- `cmd.apt_install` (for Debian based Linux distros)

- `packages`: list of strings, required.

- `cmd.winget_install` (for Windows)

- `packages`: list of strings, required.

_Check out the file [`dosh_config.lua`](./examples/dosh_config.lua) for example usage._

#### FILE, FOLDER, COMMAND EXISTENCY

To check if file or folder exists, use `cmd.exists`. And if you want to check if a command exists, use `cmd.exists_command`.

#### LOGGING

You can manage the command outputs by defining the verbosity level. It's still possible to use `print`, but if you want to hide the command outputs completely or print them by the verbosity level, you have to use these logging functions:

- `cmd.debug`
- `cmd.info`
- `cmd.warning`
- `cmd.error`

For more information about the verbosity parameter of dosh, type `dosh help`.

_Check out the file [`dosh_greet.lua`](./examples/dosh_greet.lua) for example usage._

## QUESTIONS

### CAN I TRUST THIS PROJECT?

No. Don't trust any project. The source code is open, trust yourself and read the code.

### BUT DO YOU USE THIS PROJECT YOURSELF?

Yes, of course. I use multiple operating systems with different shells, and I'm too tired to write my scripts in multiple languages. This is why I created this project.

### BUT YOU COULD USE MAKEFILE, CMAKE, OR ANOTHER SIMILAR TOOL.

They are typically used to build and package software for distribution and are more geared towards building and managing software projects, while Dosh is more focused on running tasks from the command line. They serve different purposes and are not directly comparable. I keep these rules in mind:

- If I need to add a paragraph to the `README.md` file to explain how to configure the development environment and need to run some commands on my local, write a DOSH task named `setup` instead, and then add just one sentence: "You can start development with a magic command: dosh setup." Or better yet, tell the contributors to type `dosh help` to see all available tasks.

- If I don't want to create a project or repository for my personal tasks, I create `dosh.lua` in my home folder and write my tasks directly. For example, I have a task named `git-sync` that pulls the latest changes from the remote server or warns me if there's a conflict in the repository.

- If I need a command alias but also need to run the command in Windows and Mac OS X, or in powershell and zsh, DOSH makes it simple.

### WHY DOESN'T THIS PROJECT HAVE `DOSH.LUA`?

It doesn't make sense to make this project dependent on itself.

### WHY DOESN'T DOSH HAVE ANY REMOVE COMMAND?

Because it's too dangerous! I don't use any remove command in my scripts indeed. If you really need a remove command, you can run it with `cmd.run`. But remember, contributors of this project don't guarantee anything.

## CONTRIBUTION

Install these development dependencies manually:

- [uv](https://github.com/astral-sh/uv)

```shell
$ uvx pre-commit run --all-files  # for linting
$ uv run pytest  # for testing
$ uv build # for building
$ uv pip install dist/dosh-*.whl  # for installing
$ uv run dosh  # testing the installation
$ uv run python -m dosh.cli  # testing the cli without installing
```
