# DOSH - shell-independent task manager (CORE)

**STATUS: Archived.** I started to rewrite this project in Go. You can find the new project [here](https://github.com/gkmngrgn/dosh).

---

**DOSH-CORE** is a common library for parsing `dosh.lua` script file. The official CLI is [here](https://github.com/gkmngrgn/dosh-cli).

<img src="https://raw.githubusercontent.com/gkmngrgn/dosh-core/main/dosh-logo.svg"
     width="200"
     alt="DOSH logo" />


## CONTRIBUTION

Install these development dependencies manually:

- [poetry](https://python-poetry.org/)
- [poethepoet](https://github.com/nat-n/poethepoet)
- [pre-commit](https://pre-commit.com/)

```shell
$ poetry poe --help
[...]

CONFIGURED TASKS
  lint           Check code quality
  test           Run tests
    name         Filter tests by $name
```
