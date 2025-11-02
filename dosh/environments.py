"""Pre-defined environment variables."""

from __future__ import annotations

import getpass
import os
from shellingham import detect_shell, ShellDetectionFailure
from typing import Any, Final

from dosh.base_commands import OperatingSystem

__all__ = ["ENVIRONMENTS"]


def detect_shell_or_provide_default() -> str:
    try:
        return detect_shell()[0]  # type: ignore
    except ShellDetectionFailure:
        pass

    match os.name:
        case "posix":
            env_key = "SHELL"
        case "nt":
            env_key = "COMSPEC"
        case _:
            return ""

    return os.getenv(env_key) or ""


SHELL: Final[str] = detect_shell_or_provide_default()
CURRENT_OS: Final[OperatingSystem] = OperatingSystem.get_current()
DOSH_ENV: Final[str] = os.getenv("DOSH_ENV") or ""
ENVIRONMENTS: Final[dict[str, Any]] = {
    "USER": getpass.getuser(),
    "HELP_DESCRIPTION": "dosh - shell-independent task manager",
    "HELP_EPILOG": "",
    "DOSH_ENV": DOSH_ENV,
    # shell type
    "IS_ZSH": SHELL == "zsh",
    "IS_BASH": SHELL == "bash",
    "IS_PWSH": SHELL == "powershell",
    # os type
    "IS_MACOS": CURRENT_OS == OperatingSystem.MACOS,
    "IS_LINUX": CURRENT_OS == OperatingSystem.LINUX,
    "IS_WINDOWS": CURRENT_OS == OperatingSystem.WINDOWS,
}
