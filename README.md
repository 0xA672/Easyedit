# EasyEdit — Hybrid terminal text editor

> **Vibe Coding Experiment** — This project is a pure AI/Agent playground.  
> **No human-written code is allowed.** Every line is authored by an AI agent (vibe coding / agent-driven development).
> 
> Repository: [github.com/0xA672/easyedit](https://github.com/0xA672/easyedit)

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

## Shortcuts

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

## Command mode

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

## Config file

Windows: `%APPDATA%/easyedit/config.toml`
Linux/macOS: `~/.easyedit.toml`

## Build

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

## Usage

```bash
easyedit                # Open a blank file
easyedit filename       # Open specified file
```
