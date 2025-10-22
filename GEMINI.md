# Project: DOSH

## Project Overview

DOSH is a shell-independent command-line task manager that allows users to define and run tasks, aliases, and environments using `dosh.lua` configuration files. Its primary goal is to provide a consistent task execution experience across various operating systems (Linux, macOS, Windows) and different shell environments. The project leverages Lua for task definition, offering flexibility and programmability.

## Technologies Used

*   **Primary Language:** Python
*   **Task Configuration:** Lua (`dosh.lua`)
*   **Dependency Management & Task Runner:** `uv`
*   **Logging:** `colorlog`
*   **Lua Integration:** `lupa`
*   **Testing:** `pytest`, `pytest-cov`, `pytest-httpserver`

## Building and Running

This project uses `uv` for dependency management, building, and testing.

*   **Linting:**
    ```bash
    uvx pre-commit run --all-files
    ```
*   **Testing:**
    ```bash
    uv run pytest
    ```
*   **Building:**
    ```bash
    uv build
    ```
*   **Installing:**
    ```bash
    uv pip install dist/dosh-*.whl
    ```
*   **Running the CLI (without installation):**
    ```bash
    uv run python -m dosh.cli
    ```

## Development Conventions

*   **Configuration:** Tasks and environments are defined in `dosh.lua` files. Examples can be found in the `examples/` directory.
*   **Linting:** Pre-commit hooks are used for linting, managed via `uvx pre-commit`.
*   **Testing:** Unit and integration tests are written using `pytest`. Test files are located in the `tests/` directory.
*   **Entry Point:** The main CLI entry point is `dosh.__main__:cli`.
