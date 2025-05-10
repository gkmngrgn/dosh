"""Available commands for `dosh.star`."""

from typing import Final

from dosh.commands import external as cmd
from dosh.logger import get_logger

__all__ = ["COMMANDS"]

logger = get_logger()


COMMANDS: Final = {
    # general purpose
    "clone": cmd.clone,
    "ls": cmd.scan_directory,
    "run": cmd.run,
    "run_url": cmd.run_url,
    # package managers
    "apt_install": cmd.apt_install,
    "brew_install": cmd.brew_install,
    "winget_install": cmd.winget_install,
    # file system
    "copy": cmd.copy,
    "exists": cmd.exists,
    "exists_command": cmd.exists_command,
    # logging
    "debug": logger.debug,
    "info": logger.info,
    "warning": logger.warning,
    "error": logger.error,
}
