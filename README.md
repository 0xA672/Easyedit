# EasyEdit

[![Build Status](https://github.com/0xA672/Easyedit/actions/workflows/go.yml/badge.svg)](https://github.com/0xA672/Easyedit/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/0xA672/Easyedit)](go.mod)
[![License](https://img.shields.io/github/license/0xA672/Easyedit)](LICENSE)
[![Release](https://img.shields.io/github/v/release/0xA672/Easyedit)](https://github.com/0xA672/Easyedit/releases)

> **Vibe Coding Experiment** — This project is a pure AI/Agent playground.
> **No human-written code is allowed.** Every line is authored by an AI agent.

**A hybrid terminal text editor with stream processing capabilities.**

---

## 🚀 Features

- **TUI Editor**: Full-featured terminal UI with syntax highlighting (200+ languages), multi-cursor, and command mode.
- **Stream Mode (`sed` alternative)**: Process stdin/stdout non-interactively. Perfect for CI/CD or shell scripts.
  ```bash
  echo "hello foo" | easyedit -s -e 's/foo/bar/g'
  easyedit -s -i -e 's/old/new/g' config.txt
  ```
- **Gap Buffer**: Efficient text storage for large files.
- **Undo/Redo**: Full history stack (Ctrl+Z / Ctrl+Y).
- **Cross-Platform**: Pre-built binaries for **24+ platforms** (Linux, macOS, Windows, FreeBSD, ARM, RISC-V, etc.).

## 📦 Installation

### Download Prebuilt Binary (Recommended)
Visit [Releases](https://github.com/0xA672/Easyedit/releases) for binaries covering Linux, macOS, Windows, and more.

### Go Install
```bash
go install github.com/0xA672/Easyedit@latest
```

### Build from Source
```bash
git clone https://github.com/0xA672/Easyedit.git
cd Easyedit
go build -o easyedit .
```

## ⚡ Quick Start

### Interactive Mode
```bash
easyedit                # Open blank buffer
easyedit main.go        # Open existing file
```

### Stream Mode (Non-Interactive)
```bash
# Replace text
echo "hello world" | easyedit -s -e 's/world/Go/g'

# Delete lines matching pattern
cat log.txt | easyedit -s -e '/DEBUG/d'

# In-place edit
easyedit -s -i -e 's/foo/bar/g' file.txt
```

## 📚 Documentation

Detailed guides available in `/docs`:
- [Installation Guide](docs/installation.md)
- [User Manual](docs/usage.md)
- [Configuration](docs/config.md)
- [Contributing](docs/contributing.md)

## 🤝 Contributors Signature Wall

Every AI model that contributed to this project:

| Model | Provider | Contribution |
| :--- | :--- | :--- |
| **Qwen3.5** | Alibaba | Code review and maintenance |
| **GPT-4.1** | OpenAI | Architecture design |
| **Claude 4.7** | Anthropic | UI enhancements, Vix lang support, Stream mode, Multi-lang docs |

---
*Built with ❤️ using Go and Tcell*
