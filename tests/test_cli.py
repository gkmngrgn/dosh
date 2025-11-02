import subprocess
import sys
from pathlib import Path

from dosh import __version__


def run_command(*params, cwd=None):
    return subprocess.run([sys.executable, "-m", "dosh", *params], text=True, cwd=cwd)


def test_version(capfd):
    ret = run_command("version")
    assert ret.returncode == 0

    cap = capfd.readouterr()
    assert cap.out.strip() == __version__
    assert cap.err.strip() == ""


def test_run_examples():
    examples_dir = Path(__file__).parent.parent / "examples"
    lua_files = list(examples_dir.glob("*.lua"))

    for lua_file in lua_files:
        ret = run_command("-c", str(lua_file))
        assert ret.returncode == 0


def test_init_command(tmp_path):
    """Test `init` command should not overwrite the existing file."""
    # when file exists
    (tmp_path / "dosh.lua").write_text("hello")
    run_command("init", cwd=tmp_path)
    assert (tmp_path / "dosh.lua").read_text() == "hello"

    # when file doesn't exist
    dosh_lua = tmp_path / "dosh.lua"
    dosh_lua.unlink()
    run_command("init", cwd=tmp_path)
    assert dosh_lua.exists()
