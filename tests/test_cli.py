import subprocess
import sys

from dosh import __version__


def test_version(capfd):
    ret = subprocess.run([sys.executable, "-m", "dosh", "version"], text=True)
    assert ret.returncode == 0

    cap = capfd.readouterr()
    assert cap.out.strip() == __version__
    assert cap.err.strip() == ""
