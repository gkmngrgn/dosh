"""DOSH CLI app."""

import importlib.metadata

try:
    __version__ = importlib.metadata.version(__package__ or __name__)
except importlib.metadata.PackageNotFoundError:
    __version__ = "0.0.0"


import sys
from enum import Enum
from pathlib import Path
from typing import List, Optional, Tuple

from dosh_core import DoshInitializer
from dosh_core.commands.internal import generate_help, init_config
from dosh_core.config import ConfigParser
from dosh_core.environments import ENVIRONMENTS
from dosh_core.logger import get_logger, set_verbosity

logger = get_logger()


class PredefinedTask(str, Enum):
    """Default DOSH tasks that run even if the configuration file is missing."""

    HELP = "help"
    INIT = "init"
    VERSION = "version"


class ArgumentParser:
    """CLI Argument parser."""

    def __get_arg_index(self, *args: str) -> Optional[int]:
        """Parse CLI arguments and return the index of expected arg."""
        index: Optional[int] = None
        for i, arg in enumerate(sys.argv):
            if arg in args:
                index = i
                break
        return index

    def __get_arg_value(self, *args: str) -> Optional[str]:
        """Return argument value."""
        index = self.__get_arg_index(*args)
        if (
            index is None
            or len(sys.argv) <= index
            or sys.argv[index + 1].startswith("-")
        ):
            return None
        return sys.argv[index + 1]

    def get_verbosity_level(self) -> int:
        """Get verbosity level."""
        index = self.__get_arg_index("-v", "-vv", "-vvv", "--verbose")

        if index is not None:
            arg = sys.argv[index]
            if arg.startswith("-v"):
                return arg.count("v")

            if arg == "--verbose" and len(sys.argv) > index + 1:
                level = sys.argv[index + 1]
                if level.isdigit():
                    return int(level)

        return 2  # default verbosity level.

    def get_config_path(self) -> Optional[Path]:
        """Return file path of dosh script."""
        filename = self.__get_arg_value("-c", "--config")
        return None if filename is None else Path.cwd() / filename

    def get_current_working_directory(self) -> Optional[Path]:
        """Return to get current working directory."""
        dir_str = self.__get_arg_value("-d", "--directory")
        return Path.cwd() if dir_str is None else Path(dir_str)

    def get_task_param(self, tasks: List[str]) -> Tuple[str, List[str]]:
        """Get task name with its parameters."""
        tasks += list(PredefinedTask)
        task_index = self.__get_arg_index(*tasks)
        if task_index is None:
            return (PredefinedTask.HELP, [])
        return sys.argv[task_index], sys.argv[task_index + 1 :]


class CLI:  # pylint: disable=too-few-public-methods
    """DOSH command line interface."""

    def __init__(self) -> None:
        """Initialize cli with config parser."""
        self.arg_parser = ArgumentParser()
        self.dosh_initializer = DoshInitializer()

        # update dosh initializer settings by cli arguments
        base_directory = self.arg_parser.get_current_working_directory()
        if base_directory is not None:
            self.dosh_initializer.base_directory = base_directory

        config_path = self.arg_parser.get_config_path()
        if config_path is not None:
            self.dosh_initializer.config_path = config_path

        # define config parser
        content = (
            self.dosh_initializer.config_path.read_text(encoding="utf-8")
            if self.dosh_initializer.config_path.exists()
            else ""
        )
        self.conf_parser = ConfigParser(content)

    def run(self) -> None:
        """Run cli reading the arguments."""
        set_verbosity(verbosity=self.arg_parser.get_verbosity_level())

        tasks = [t.name for t in self.conf_parser.tasks]
        task_name, task_params = self.arg_parser.get_task_param(tasks)

        if task_name == PredefinedTask.HELP:
            output = generate_help(
                tasks=self.conf_parser.tasks,
                description=self.conf_parser.description,
                epilog=self.conf_parser.epilog,
            )
            print(output)
        elif task_name == PredefinedTask.INIT:
            init_config(self.dosh_initializer.config_path)
        elif task_name == PredefinedTask.VERSION:
            print(__version__)
        else:
            dosh_env = ENVIRONMENTS["DOSH_ENV"] or "not specified."
            working_directory = self.arg_parser.get_current_working_directory()
            logger.debug("DOSH ENVIRONMENT : %s", dosh_env)
            logger.debug("CONFIG FILE PATH : %s", self.arg_parser.get_config_path())
            logger.debug("WORKING DIRECTORY: %s", working_directory)

            self.conf_parser.run_task(task_name, task_params)


if __name__ == "__main__":
    CLI().run()
