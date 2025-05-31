"""DOSH config parser."""

import shutil
from dataclasses import fields
from typing import Any, Dict, List

from dosh.commands import COMMANDS
from dosh.base_commands import OperatingSystem, Task
from dosh.environments import DOSH_ENV, ENVIRONMENTS
from dosh.logger import get_logger
from dosh.lua_runtime import get_lua_environment

logger = get_logger()


class ConfigParser:
    """Dosh configuration parser."""

    tasks: List[Task] = []
    _vars: Dict[str, str] = {}

    def __init__(self, content: str) -> None:
        """Parse config first."""
        commands = COMMANDS.copy()
        commands["add_task"] = self.add_task
        self._vars = get_lua_environment(content, ENVIRONMENTS, commands)

    @property
    def description(self) -> str:
        """Get help description."""
        return self._vars["HELP_DESCRIPTION"]

    @property
    def epilog(self) -> str:
        """Get help epilog."""
        return self._vars["HELP_EPILOG"]

    def run_task(self, task_name: str, params: List[str]) -> None:
        """Find and run tasks with parameters."""
        task = None
        for current_task in self.tasks:
            if current_task.name == task_name:
                task = current_task
                break

        if task is None:
            return

        os_type = OperatingSystem.get_current().value
        if task.required_platforms and os_type not in task.required_platforms:
            logger.error(
                "The task `%s` works only on these operating system: %s (current: %s)",
                task.name,
                ", ".join(task.required_platforms),
                os_type,
            )
            return

        if task.environments and DOSH_ENV not in task.environments:
            logger.error(
                "The task `%s` works only in these environments: %s (current: %s)",
                task.name,
                ", ".join(task.environments),
                DOSH_ENV,
            )
            return

        for command in task.required_commands:
            if shutil.which(command) is None:
                logger.error("The command `%s` doesn't exist.", command)
                return

        try:
            task.command(*params)
        except KeyboardInterrupt:
            print("\r", end="")
            logger.error("keyboard interrupt...")

    def add_task(self, args: Dict[str, Any]) -> None:
        """Parse and add task to task list."""
        for field in filter(lambda field: field.type.startswith("List"), fields(Task)):
            val = args[field.name]
            args[field.name] = [] if val is None else list(val.values())

        task = Task.from_dict(args)
        self.tasks.append(task)
