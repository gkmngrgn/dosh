"""Available commands for `dosh.star`."""

from __future__ import annotations

import os
import shutil
import urllib.request
from pathlib import Path
from typing import List, Optional

from dosh.commands.base import (
    CommandException,
    check_command,
    copy_tree,
    is_url_valid,
    normalize_path,
    run_command_and_return_result,
)
from dosh.logger import get_logger
from dosh.lua_runtime import LuaTable, lua_runtime

logger = get_logger()


@check_command("apt")
def apt_install(packages: List[str]) -> None:
    """Install packages with apt."""
    command = "apt install"

    run(f"{command} {' '.join(packages)}")


@check_command("brew")
def brew_install(packages: LuaTable, options: Optional[LuaTable] = None) -> None:
    """Install packages with brew."""
    if options is None:
        options = lua_runtime.table()

    for tap_path in list((options["taps"] or lua_runtime.table()).values()):
        run(f"brew tap {tap_path}")

    command = "brew install"

    if options["cask"] is True:
        command = f"{command} --cask"

    run(f"{command} {' '.join(list(packages.values()))}")


@check_command("winget")
def winget_install(packages: List[str]) -> None:
    """Install packages with winget."""
    command = "winget install -e --id"
    run(f"{command} {' '.join(packages)}")


def copy(src: str, dst: str) -> None:
    """Copy files from source to destination. It works like `cp` command."""
    src_path = normalize_path(src)
    dst_path = normalize_path(dst)
    glob_index = -1
    src_splitted = src_path.as_posix().split("/")

    for index, value in enumerate(src_splitted):
        if "*" in value:
            glob_index = index
            break

    if glob_index >= 0:
        src_path = normalize_path("/".join(src_splitted[:glob_index]))
        for path in src_path.glob("/".join(src_splitted[glob_index:])):
            copy_tree(path, dst_path / path.name)
    else:
        copy_tree(src_path, dst_path / src_path.name if dst_path.exists() else dst_path)


@check_command("git")
def clone(url: str, options: Optional[LuaTable] = None) -> None:
    """Clone repository from VCS."""
    if options is None:
        options = lua_runtime.table()

    destination = normalize_path(
        options["destination"] or url.rsplit("/", 1)[-1].rsplit(".git", 1)[0]
    )

    if destination.exists():
        logger.info("[CLONE] the folder `%s` already exists. skipped.", destination)
    else:
        command = f"git clone {url} {destination}"
        run(command)


def scan_directory(parent_dir: str = ".", opts: Optional[LuaTable] = None) -> LuaTable:
    """
    List files and directories.

    Optional parameters:
        include_files: boolean (default: true)
        include_dirs: boolean (default: true)
        excludes: list[str] (default: [])
    """
    parent = normalize_path(parent_dir)
    options = opts or lua_runtime.table()

    if not parent.is_dir():
        raise CommandException(f"Not a folder: {parent_dir}")

    files, directories = [], []

    for item in parent.iterdir():
        if item.is_dir():
            directories.append(str(item))
        elif item.is_file():
            files.append(str(item))

    items = []

    if options["include_files"] is not False:
        items.extend(files)

    if options["include_dirs"] is not False:
        items.extend(directories)

    excludes = list((options["excludes"] or lua_runtime.table()).values())
    if excludes:
        items = list(filter(lambda i: i.rsplit(os.sep, 1)[-1] not in excludes, items))

    response: LuaTable = lua_runtime.table_from(sorted(items))
    return response


def run(command: str) -> int:
    """Run a shell command using subprocess."""
    log_prefix = "[RUN]"
    logger.info("%s %s", log_prefix, command)
    return run_command_and_return_result(command, log_prefix)


def run_url(url: str) -> int:
    """Run a remote shell script directly."""
    if not is_url_valid(url):
        raise CommandException(f"URL is not valid: {url}")

    with urllib.request.urlopen(url) as response:
        content = response.read().decode("utf-8")

    log_prefix = "[RUN_URL]"
    logger.info("%s %s", log_prefix, url)
    return run_command_and_return_result(content, log_prefix)


def exists(path: str) -> bool:
    """Check if the path exists in the file system."""
    return Path(path).exists()


def exists_command(command: str) -> bool:
    """Check if the command exists."""
    return shutil.which(command) is not None
