import pathlib
import platform
import urllib.request
from pathlib import Path

import pytest

from dosh import external_commands as cmd
from dosh.base_commands import CommandException, normalize_path
from dosh.logger import get_logger, set_verbosity

logger = get_logger()


def test_copy(tmp_path):
    (tmp_path / "foo" / "bar" / "baz").mkdir(parents=True)
    (tmp_path / "foo" / "bar" / "meow.py").touch()
    (tmp_path / "foo" / "bar" / "baz" / "hell.txt").touch()

    # test with files "foo/*"
    src1 = str(tmp_path / "foo" / "*")
    dst1 = str(tmp_path / "dst1")

    cmd.copy(src1, dst1)
    assert (tmp_path / "dst1" / "bar" / "meow.py").exists()
    assert (tmp_path / "dst1" / "bar" / "baz" / "hell.txt").exists()

    # try overwrite, we don't expect to create a folder "foo" in dst1.
    cmd.copy(src1, dst1)
    assert not (tmp_path / "dst1" / "foo").exists()

    # test with a folder name "foo"
    src2 = str(tmp_path / "foo")
    dst2 = str(tmp_path / "dst2")

    cmd.copy(src2, dst2)
    assert (tmp_path / "dst2" / "bar" / "meow.py").exists()
    assert (tmp_path / "dst2" / "bar" / "baz" / "hell.txt").exists()

    # try overwrite, we expect to create a folder "foo" in dst2.
    cmd.copy(src2, dst2)
    assert (tmp_path / "dst2" / "foo" / "bar" / "meow.py").exists()
    assert (tmp_path / "dst2" / "foo" / "bar" / "baz" / "hell.txt").exists()


def test_echo(capsys):
    """Test echo command."""
    cmd.echo("Hello, world!")
    captured = capsys.readouterr()
    assert captured.out == "Hello, world!\n"


def test_run(caplog):
    set_verbosity(3)
    results = cmd.run("echo Hello!")
    assert results == [0]
    assert caplog.records[0].message == "[RUN] echo Hello!"
    assert caplog.records[1].message == "[RUN] Return code: 0"


def test_multiple_run(caplog):
    set_verbosity(3)
    results = cmd.run("echo Hello!", "echo World!")
    assert results == [0, 0]
    assert caplog.records[0].message == "[RUN] echo Hello!"
    assert caplog.records[1].message == "[RUN] Return code: 0"
    assert caplog.records[2].message == "[RUN] echo World!"
    assert caplog.records[3].message == "[RUN] Return code: 0"


def test_run_url(httpserver, caplog):
    set_verbosity(3)

    sh_content = "echo Hello!"
    httpserver.expect_request("/hello.sh").respond_with_data(
        sh_content, content_type="text/plain"
    )
    url = httpserver.url_for("/hello.sh")

    with urllib.request.urlopen(url) as response:
        content = response.read()

    assert content.decode("utf-8") == sh_content

    result = cmd.run_url(url)
    assert result == 0
    assert (
        caplog.records[2].message
        == f"[RUN_URL] http://{httpserver.host}:{httpserver.port}/hello.sh"
    )
    assert caplog.records[3].message == "[RUN_URL] Return code: 0"


def test_exists(tmp_path):
    src_path1 = tmp_path / "this-file-does-not-exist.txt"
    assert not cmd.exists(str(src_path1))

    src_path2 = tmp_path / "this-file-exists.txt"
    src_path2.touch()
    assert cmd.exists(str(src_path2))

    src_path3 = tmp_path / "this" / "folder" / "exists"
    src_path3.mkdir(parents=True)
    assert cmd.exists(str(src_path3))


def test_exists_command():
    if platform.system() == "Windows":
        assert cmd.exists_command("cmd.exe")
        assert cmd.exists_command("where.exe")
        assert not cmd.exists_command("dmc.exe")
    else:
        assert cmd.exists_command("bash")
        assert cmd.exists_command("which")
        assert not cmd.exists_command("hsab")


def test_normalize_path(monkeypatch):
    if platform.system() == "Windows":
        monkeypatch.setenv("USERPROFILE", "C:/Users/dosh")
        assert normalize_path("~/.config") == Path("C:/Users/dosh/.config")

    else:
        monkeypatch.setenv("HOME", "/home/dosh")  # for posix
        assert normalize_path("~/.config") == Path("/home/dosh/.config")

    assert normalize_path("/foo/bar/baz") == Path("/foo/bar/baz")
    assert normalize_path("current_dir") == Path.cwd() / "current_dir"


def test_scan_directory(monkeypatch):
    base_dir = pathlib.Path(__file__).parent.parent.resolve()
    monkeypatch.chdir(base_dir)

    # test current working directory
    result = cmd.scan_directory()
    assert result is not None
    assert str(base_dir / "README.md") in list(result.values())

    # test examples directory
    result = cmd.scan_directory("./tests")
    assert result is not None

    files = list(result.values())
    for file_name in ["test_commands.py", "test_logger.py", "test_task.py"]:
        assert str(base_dir / "tests" / file_name) in files

    # test invalid directory
    with pytest.raises(CommandException) as excinfo:
        cmd.scan_directory("./README.md")
        assert str(excinfo.value) == "Not a folder: ./README.md"
