import platform
import textwrap

from dosh.base_commands import OperatingSystem
from dosh.config import ConfigParser
from dosh.logger import set_verbosity


def test_get_current_operating_system(monkeypatch):
    monkeypatch.setattr(platform, "system", lambda: "Darwin")
    assert OperatingSystem.get_current() == OperatingSystem.MACOS

    monkeypatch.setattr(platform, "system", lambda: "Linux")
    assert OperatingSystem.get_current() == OperatingSystem.LINUX

    monkeypatch.setattr(platform, "system", lambda: "Windows")
    assert OperatingSystem.get_current() == OperatingSystem.WINDOWS

    monkeypatch.setattr(platform, "system", lambda: "FreeBSD")
    assert OperatingSystem.get_current() == OperatingSystem.UNSUPPORTED


def test_task_if_available_on_linux(monkeypatch, caplog):
    set_verbosity(2)
    content = textwrap.dedent(
        """
        cmd.add_task{
            name="hello",
            required_platforms={ "windows", "macos" },
            command=function ()
                cmd.info("hello")
            end
        }
        """
    )
    config_parser = ConfigParser(content=content)

    # run task on unsupported operating system first.
    monkeypatch.setattr(platform, "system", lambda: "Linux")
    config_parser.run_task("hello", params=[])

    assert len(caplog.records) == 1
    assert (
        caplog.records[0].message
        == "The task `hello` works only on these operating system: windows, macos (current: linux)"
    )

    # run task on macos. that should be worked.
    monkeypatch.setattr(platform, "system", lambda: "Darwin")
    config_parser.run_task("hello", params=[])

    assert len(caplog.records) == 2
    assert caplog.records[1].message == "hello"
