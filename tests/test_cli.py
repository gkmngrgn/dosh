import subprocess
import sys
from pathlib import Path

from dosh import __version__


def run_command(*params):
    return subprocess.run([sys.executable, "-m", "dosh.cli", *params], text=True)


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
