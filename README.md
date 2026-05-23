# EasyEdit — Hybrid terminal text editor

> **Vibe Coding Experiment** — This project is a pure AI/Agent playground.  
> **No human-written code is allowed.** Every line is authored by an AI agent (vibe coding / agent-driven development).
>
> Repository: [github.com/0xA672/easyedit](https://github.com/0xA672/Easyedit)

## Features

- Hybrid input mode: insert (Vim-like) and command mode (`:` prefix)
- Syntax highlighting (powered by Chroma, supports 200+ languages)
- Gap Buffer text storage for efficient large-file editing
- Undo / Redo stack (Ctrl+Z / Ctrl+Y)
- Search (Ctrl+F)
- Cut / Copy / Paste (Ctrl+X/C/V)
- Select All (Ctrl+A)
- Soft wrap toggle (Alt+W)
- Mouse support (click to position cursor)
- Config via TOML
- Cross-platform (Windows / Linux / macOS)

## Install

### Build from source

Requires Go 1.21+.

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

If you are on Windows, this produces `easyedit.exe`.

### Download prebuilt binary

> Coming soon — releases page will host prebuilt binaries for Windows, Linux, and macOS.

## Usage

```bash
easyedit                # Open a blank buffer
easyedit filename       # Open an existing file
easyedit file1 file2    # Open multiple files (tab navigation, if implemented)
```

### Normal mode shortcuts

| Shortcut | Action |
|----------|--------|
| Arrow keys / Home / End / PgUp / PgDn | Cursor movement |
| Backspace / Delete | Delete character |
| Ctrl+S | Save |
| Ctrl+Q | Quit |
| Ctrl+F | Search (Enter to confirm, Esc to cancel) |
| Ctrl+A | Select All |
| Ctrl+Z | Undo |
| Ctrl+Y | Redo |
| Ctrl+X | Cut |
| Ctrl+C | Copy |
| Ctrl+V | Paste |
| Alt+W | Toggle soft wrap |
| `:` (colon) | Enter command mode |
| Mouse click | Move cursor to clicked position |

### Command mode

| Command | Action |
|---------|--------|
| `:q` | Quit |
| `:q!` | Force quit |
| `:w` | Save |
| `:wq` | Save and quit |
| `:w path` | Save as |
| `:e path` | Open file |
| `:10,20d` | Delete lines 10-20 |
| `:%s/old/new/g` | Replace all |
| `:3,5s/old/new` | Replace in range |
| `:set nu` | Show line numbers |
| `:set nonu` | Hide line numbers |
| `:set nowrap` | Disable soft wrap |
| `:set wrap` | Enable soft wrap |

## Configuration

Config is loaded from a TOML file:

- **Windows:** `%APPDATA%/easyedit/config.toml`
- **Linux / macOS:** `~/.easyedit.toml`

Example:

```toml
[theme]
background = "default"
foreground = "default"

[editor]
tab_width = 4
show_line_numbers = true
soft_wrap = false
```

## Project structure

```
easyedit/
├── main.go              # Entry point
├── command/
│   └── command.go       # Command mode parsing and execution (:d, :%s, :w, :q)
├── config/
│   └── config.go        # Config loading (TOML)
├── document/
│   ├── document.go      # Gap Buffer text storage + Document high-level wrapper
│   └── undo.go          # Undo/Redo stack
├── highlight/
│   └── highlight.go     # Syntax highlighting (Chroma)
└── ui/
    └── editor.go        # Main loop, event handling, rendering
```

## Contributing

This is a **Vibe Coding experiment** — contributions must follow these rules:

1. **No human-written code.** Every line must be authored by an AI agent.
2. Only configuration files (`.gitignore`, `README.md`, CI workflows, etc.) may be written or edited by humans.
3. If you open a PR, describe the problem clearly and let AI agents generate the fix.
4. All participants (human or agent) will be listed in the Model Signature Wall below.

## Model Signature Wall

Every AI model that has contributed code or design to this project is listed here.

| Model | Provider | Role |
|-------|----------|------|
| deepseek-v4-flash | DeepSeek | Core code, README, project setup |

*> Add your model! Open a PR or tell the maintainers which agent authored your contribution.*
