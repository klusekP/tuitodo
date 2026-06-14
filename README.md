# tuitodo

**A fullscreen, keyboard-driven todo list in your terminal.**  
Built in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea) — no GUI, no web; tasks sync to a simple JSON file on your machine.

> **One-line summary for the GitHub repository “About” field:**  
> *Fullscreen terminal todo app in Go (Bubble Tea) with JSON persistence and a clean, responsive TUI.*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE) [![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

---

## Why this project

`tuitodo` is a small, [open source](https://opensource.org/) CLI tool: you can read the code, fork it, and adapt it. The codebase is split into clear packages (domain, persistence, application logic, TUI) so it stays easy to extend or test.

## Features

- Fullscreen, responsive UI — adapts to terminal resizing on the fly
- Scrollable list with position indicator (e.g. `1-20/150`)
- Add, edit, and delete tasks
- Toggle tasks as done or pending — completing a task moves it **below every still-pending** item (pending first, then done)
- Clear all completed tasks with a single key
- Automatic persistence to `~/.tuitodo.json`
- Clean, colorful interface ([lipgloss](https://github.com/charmbracelet/lipgloss))

## Requirements

- Go 1.21+

## Install

**From source (recommended for development):**

```bash
git clone https://github.com/klusekP/tuitodo.git
cd tuitodo
go build -o tuitodo ./cmd/tuitodo
# Optional: put the binary on your PATH
```

**Run without installing:**

```bash
go run ./cmd/tuitodo
```

## Project layout (packages)

| Package | Role |
| ------- | ---- |
| `cmd/tuitodo` | Composition root: config → JSON repo → app service → TUI |
| `internal/todo` | Domain: `Item` + `Repo` (persistence port) |
| `internal/storage/jsonfile` | JSON file adapter for `Repo` |
| `internal/app` | Task operations and save orchestration |
| `internal/tui` | Bubble Tea UI, layout, and theme |
| `internal/config` | Default data path |

## Keybindings

| Key                   | Action                           |
| --------------------- | -------------------------------- |
| `a` / `n` / `+`       | Add a new task                   |
| `e`                   | Edit the selected task           |
| `space` / `enter` / `x` | Toggle done / pending          |
| `d` / `delete`        | Delete the selected task         |
| `c`                   | Clear all completed tasks        |
| `j` / `↓`             | Move down                        |
| `k` / `↑`             | Move up                          |
| `PgUp` / `PgDn`       | Jump by a page                   |
| `g` / `home`          | Jump to top                      |
| `G` / `end`           | Jump to bottom                   |
| `q` / `ctrl+c`        | Quit                             |
| `esc` (in edit mode)  | Cancel                           |

## Data format

Tasks are stored as JSON in `~/.tuitodo.json`:

```json
[
  {
    "id": 1714000000000000000,
    "title": "Buy milk",
    "done": false,
    "created_at": "2026-04-24T08:33:00Z"
  }
]
```

## Contributing

Issues and pull requests are welcome. Please keep changes focused; match existing style and run:

```bash
go test ./...
go vet ./...
```

## License

This project is **open source** and released under the [MIT License](./LICENSE) — you may use, copy, modify, and distribute the software, subject to the terms in the license file.

**Suggested GitHub metadata**

| Field        | Suggested value |
| ------------ | --------------- |
| **About**    | `Fullscreen terminal todo in Go (Bubble Tea) + JSON persistence` |
| **Website**  | *(empty or your project page, if any)* |
| **Topics**   | `golang`, `tui`, `terminal`, `cli`, `bubble-tea`, `todo`, `open-source` |

> Repository: `github.com/klusekP/tuitodo` — rename the last segment if your repository name differs.
