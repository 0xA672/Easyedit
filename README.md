# 🧩 EasyEdit

[![Build Status](https://github.com/0xA672/Easyedit/actions/workflows/go.yml/badge.svg)](https://github.com/0xA672/Easyedit/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/0xA672/Easyedit)](go.mod)
[![License](https://img.shields.io/github/license/0xA672/Easyedit)](LICENSE)
[![Release](https://img.shields.io/github/v/release/0xA672/Easyedit)](https://github.com/0xA672/Easyedit/releases)
[![Platforms](https://img.shields.io/badge/platforms-24+-blue)](#cross-platform)

> 🤖 **Vibe Coding Experiment** — *No human-written code allowed.*  
> Every line of this editor is authored by AI agents.  
> [See who contributed →](#-contributors-signature-wall)

**EasyEdit** is a hybrid terminal text editor that offers **stream processing capabilities** — functioning both as a full‑featured TUI editor and a non‑interactive `sed`‑style filter.

---

## ✨ Features

| Mode | Description |
|------|-------------|
| 🖥️ **TUI Editor** | Full terminal UI with syntax highlighting for 200+ languages, multi‑cursor editing, undo/redo, and command mode |
| ⚡ **Stream Mode** | Non‑interactive `stdin`/`stdout` processing — ideal for CI/CD, shell pipelines, and in‑place edits |
| 📦 **Gap Buffer** | Efficient text storage for large files |
| 🔁 **Undo/Redo** | Full history stack (`Ctrl+Z` / `Ctrl+Y`) |
| 🌍 **Cross‑Platform** | Pre‑built binaries available for **24+ platforms**, including Linux, macOS, Windows, FreeBSD, ARM, RISC‑V, and more |

---

## 🚀 Quick Start

### 📥 Installation

**Option 1: Download Prebuilt Binary (Recommended)**  
Visit the [Releases](https://github.com/0xA672/Easyedit/releases) page – binaries are available for Linux, macOS, Windows, and many more.

**Option 2: Go Install**  
```bash
go install github.com/0xA672/Easyedit@latest
```

**Option 3: Build from Source**  
```bash
git clone https://github.com/0xA672/Easyedit.git
cd Easyedit
go build -o easyedit .
```

---

### ▶️ Usage Examples

#### Interactive Mode
```bash
easyedit                # Open an empty buffer
easyedit main.go        # Open an existing file
```

#### Stream Mode (sed‑like)
```bash
# Replace text in a pipeline
echo "hello world" | easyedit -s -e 's/world/Go/g'

# Delete lines matching a pattern
cat log.txt | easyedit -s -e '/DEBUG/d'

# In-place file editing
easyedit -s -i -e 's/foo/bar/g' file.txt
```

---

## 📚 Documentation

Detailed guides are available in the [`/docs`](docs/) folder:

- [Installation Guide](docs/installation.md)
- [User Manual](docs/usage.md)
- [Configuration](docs/config.md)
- [Contributing](docs/contributing.md)

---

## 🤝 Contributors Signature Wall

Every AI model that contributed to this project:

| Model | Provider | Contribution |
| :--- | :--- | :--- |
| **DeepSeek-V4-Flash** | DeepSeek AI | Core contributor – syntax highlighting optimization, Vix language support |
| **Qwen3.5** | Alibaba | Code review & maintenance |
| **GPT-4.1** | OpenAI | Architecture design |
| **Claude 4.7** | Anthropic | UI enhancements, Vix language support, Stream mode, multi‑language docs |
| **Doubao** | ByteDance | Documentation annotations |

---

## 🧠 Why “EasyEdit”?

- **Easy** to use interactively, and **easy** to script in pipelines.  
- Built entirely by **AI** – a living experiment in vibe‑driven development.  
- Lightweight, fast, and **portable** — a single binary.

---

*Made with ❤️ using Go and [Tcell](https://github.com/gdamore/tcell)*
