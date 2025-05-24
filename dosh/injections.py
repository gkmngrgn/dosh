"""Dosh script injections."""

import json
from inspect import isfunction
from typing import Any, Dict


def _parse_doc(desc: str) -> str:
    """Parse docstring of functions to use them in help output."""
    line_list = []
    for line in desc.split("\n"):
        line = line.strip()
        if not line:
            continue
        line_list.append(line)
    return "\n".join(line_list)


def print_commands(variables: Dict[str, Any]) -> None:
    """Inject a function to parse commands in user-defined dosh configuration."""
    prefix = "cmd_"
    commands = filter(
        lambda k: k.startswith(prefix) and isfunction(variables.get(k)),
        variables.keys(),
    )
    output = {}

    for cmd_func in commands:
        cmd_name = cmd_func[len(prefix) :]
        output[cmd_name] = _parse_doc(variables.get(cmd_func).__doc__ or "")

    print(json.dumps(output))
