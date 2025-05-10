"""Lua runtime for parsing configuration files."""

from typing import Any, Callable, Dict, Union

from lupa import LuaRuntime

LuaFunction = Callable[..., None]
LuaTable = Dict[Union[str, int], Any]
lua_runtime = LuaRuntime(unpack_returned_tuples=True)


def get_lua_environment(
    content: str, envs: Dict[str, Any], commands: Dict[str, Any]
) -> Dict[str, Any]:
    """Get lua environment variables."""
    lua_code = f"function (env, cmd) {content} return env end"
    lua_func = lua_runtime.eval(lua_code)
    result: Dict[str, Any] = lua_func(envs, commands)
    return result
