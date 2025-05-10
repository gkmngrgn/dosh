"""Logging support."""

from logging import DEBUG, ERROR, INFO, WARNING, Logger
from typing import Optional

import colorlog

__LOGGER: Optional[Logger] = None


def get_logger() -> Logger:
    """Get logger of dosh."""
    global __LOGGER  # pylint: disable=global-statement

    if __LOGGER is None:
        handler = colorlog.StreamHandler()
        handler.setFormatter(
            colorlog.ColoredFormatter("%(log_color)s%(name)s => %(message)s")
        )

        logger = colorlog.getLogger("DOSH")
        logger.addHandler(handler)

        __LOGGER = logger

    return __LOGGER


def set_verbosity(verbosity: int = 0) -> None:
    """Set verbosity level for logger."""
    if verbosity >= 3:
        level = DEBUG
    elif verbosity == 2:
        level = INFO
    elif verbosity == 1:
        level = WARNING
    else:
        level = ERROR

    get_logger().setLevel(level)
