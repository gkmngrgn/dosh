# DOSH - Shell-Independent Task Manager

## Overview
DOSH is a command-line application for running tasks across different platforms and shells. It uses Lua configuration files to define tasks, aliases, and environments, making it a versatile task manager that works on Linux, macOS, and Windows.

## Project Structure
- **Language**: Python 3.13
- **Dependencies**: colorlog, lupa (Lua bridge)
- **Configuration**: Uses `dosh.lua` files for task definitions
- **Package Manager**: UV for Python dependency management

## Setup Status
- ✅ Python 3.13 environment configured
- ✅ Virtual environment created in `.pythonlibs/`
- ✅ Dependencies installed via UV
- ✅ CLI application tested and working
- ✅ Workflow configured for console output
- ✅ Deployment configured as VM target

## Usage
The application provides a CLI interface that reads Lua configuration files:

```bash
# Activate environment and use DOSH
source .pythonlibs/bin/activate

# Show help
python -m dosh help

# Initialize new config
python -m dosh init

# Show version
python -m dosh version

# Run custom tasks (defined in dosh.lua)
python -m dosh <task-name>
```

## Configuration
Tasks are defined in `dosh.lua` files using Lua syntax. The application supports:
- Custom task definitions with Lua functions
- Environment variables and platform detection
- Package manager integrations (brew, apt, winget)
- File system operations
- Command execution with logging

## Recent Changes
- 2025-09-14: Initial project import and setup completed
- Python 3.13 environment established
- All dependencies installed and tested
- Deployment configuration added

## Architecture
The project follows a modular structure:
- `dosh/` - Main package with CLI, config parsing, and command execution
- `examples/` - Sample Lua configuration files
- `tests/` - Test suite for core functionality
- Virtual environment in `.pythonlibs/` for dependency isolation