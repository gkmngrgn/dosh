"""DOSH module."""

from dataclasses import dataclass, field
from pathlib import Path
import importlib.metadata

try:
    __version__ = importlib.metadata.version(__package__ or __name__)
except importlib.metadata.PackageNotFoundError:
    __version__ = "0.0.0"


@dataclass
class DoshInitializer:  # pylint: disable=too-few-public-methods
    """Pre-configured dosh initializer to store app-specific settings."""

    base_directory: Path = field(default_factory=Path.cwd)
    config_path: Path = field(default_factory=lambda: Path.cwd() / "dosh.lua")
