# EasyEdit — Hybrid terminal text editor

[![Build Status](https://github.com/0xA672/Easyedit/actions/workflows/ci.yml/badge.svg)](https://github.com/0xA672/Easyedit/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/0xA672/Easyedit)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> **Vibe Coding Experiment** — This project is a pure AI/Agent playground.  
> **No human-written code is allowed.** Every line is authored by an AI agent (vibe coding / agent-driven development).
>
> Repository: [github.com/0xA672/Easyedit](https://github.com/0xA672/Easyedit)

---

**Languages:** [English](#english) | [中文](#中文)

---

## English

### Features

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

### Install

#### Go install (recommended)

Requires Go 1.21+.

```bash
go install github.com/0xA672/easyedit@latest
```

This installs the `easyedit` binary to your `$GOPATH/bin` (or `$HOME/go/bin`). Make sure this directory is in your `PATH`.

#### Build from source

Requires Go 1.21+.

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

If you are on Windows, this produces `easyedit.exe`.

#### Download prebuilt binary

> Coming soon — releases page will host prebuilt binaries for Windows, Linux, and macOS.

### Usage

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
| `:uninstall` | Remove config and delete the editor binary |

### Configuration

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

### Project structure

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

### Contributing

This is a **Vibe Coding experiment** — contributions must follow these rules:

1. **No human-written code.** Every line must be authored by an AI agent.
2. Only configuration files (`.gitignore`, `README.md`, CI workflows, etc.) may be written or edited by humans.
3. If you open a PR, describe the problem clearly and let AI agents generate the fix.
4. All participants (human or agent) will be listed in the Model Signature Wall below.

### Documentation

For detailed documentation, see the [/docs](./docs) directory.

### Quick Help

Run `easyedit --help` for quick usage information.

### Model Signature Wall

Every AI model that has contributed code or design to this project is listed here.

| Model             | Provider  | Role           |
|-------------------|-----------|----------------|
| deepseek-v4-flash | DeepSeek  | Core code, README, project setup |
| Qwen3.5           | Alibaba   | Code review and maintenance |
| GPT-4.1           | OpenAI    | Architecture design |
| Claude 4.7        | Anthropic | UI enhancements and signatures |
| Doubao            | ByteDance | Documentation annotation generation, README updates |

> *Add your model! Open a PR or tell the maintainers which agent authored your contribution.*

---

## 中文

### 核心特性

- 混合输入模式：插入模式（类似 Vim）和命令模式（`:` 前缀）
- 语法高亮（基于 Chroma，支持 200+ 种语言）
- Gap Buffer 文本存储，高效编辑大文件
- 撤销/重做栈（Ctrl+Z / Ctrl+Y）
- 搜索功能（Ctrl+F）
- 剪切/复制/粘贴（Ctrl+X/C/V）
- 全选（Ctrl+A）
- 软换行切换（Alt+W）
- 鼠标支持（点击定位光标）
- TOML 配置文件
- 跨平台（Windows / Linux / macOS）

### 安装

#### Go 安装（推荐）

需要 Go 1.21+。

```bash
go install github.com/0xA672/easyedit@latest
```

这会将 `easyedit` 二进制文件安装到您的 `$GOPATH/bin`（或 `$HOME/go/bin`）。确保此目录在您的 `PATH` 中。

#### 从源码构建

需要 Go 1.21+。

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

如果您在 Windows 上，这将生成 `easyedit.exe`。

#### 下载预编译二进制

> 即将推出 — releases 页面将提供 Windows、Linux 和 macOS 的预编译二进制文件。

### 使用方法

```bash
easyedit                # 打开空白缓冲区
easyedit filename       # 打开现有文件
easyedit file1 file2    # 打开多个文件（如果实现了标签导航）
```

### 普通模式快捷键

| 快捷键 | 操作 |
|--------|------|
| 方向键 / Home / End / PgUp / PgDn | 光标移动 |
| Backspace / Delete | 删除字符 |
| Ctrl+S | 保存 |
| Ctrl+Q | 退出 |
| Ctrl+F | 搜索（Enter 确认，Esc 取消） |
| Ctrl+A | 全选 |
| Ctrl+Z | 撤销 |
| Ctrl+Y | 重做 |
| Ctrl+X | 剪切 |
| Ctrl+C | 复制 |
| Ctrl+V | 粘贴 |
| Alt+W | 切换软换行 |
| `:`（冒号） | 进入命令模式 |
| 鼠标点击 | 移动光标到点击位置 |

### 命令模式

| 命令 | 操作 |
|------|------|
| `:q` | 退出 |
| `:q!` | 强制退出 |
| `:w` | 保存 |
| `:wq` | 保存并退出 |
| `:w path` | 另存为 |
| `:e path` | 打开文件 |
| `:10,20d` | 删除第 10-20 行 |
| `:%s/old/new/g` | 全部替换 |
| `:3,5s/old/new` | 范围内替换 |
| `:set nu` | 显示行号 |
| `:set nonu` | 隐藏行号 |
| `:set nowrap` | 禁用软换行 |
| `:set wrap` | 启用软换行 |
| `:uninstall` | 删除配置并卸载编辑器 |

### 配置

配置从 TOML 文件加载：

- **Windows:** `%APPDATA%/easyedit/config.toml`
- **Linux / macOS:** `~/.easyedit.toml`

示例：

```toml
[theme]
background = "default"
foreground = "default"

[editor]
tab_width = 4
show_line_numbers = true
soft_wrap = false
```

### 项目结构

```
easyedit/
├── main.go              # 入口点
├── command/
│   └── command.go       # 命令模式解析和执行 (:d, :%s, :w, :q)
├── config/
│   └── config.go        # 配置加载 (TOML)
├── document/
│   ├── document.go      # Gap Buffer 文本存储 + 文档高级封装
│   └── undo.go          # 撤销/重做栈
├── highlight/
│   └── highlight.go     # 语法高亮 (Chroma)
└── ui/
    └── editor.go        # 主循环、事件处理、渲染
```

### 贡献指南

这是一个 **Vibe Coding 实验** — 贡献必须遵循以下规则：

1. **禁止人类编写代码。** 每一行必须由 AI 代理编写。
2. 只有配置文件（`.gitignore`、`README.md`、CI 工作流等）可以由人类编写或编辑。
3. 如果您提交 PR，请清晰描述问题，让 AI 代理生成修复方案。
4. 所有参与者（人类或代理）都将列在下方的模型签名墙中。

### 文档

详细文档请参阅 [/docs](./docs) 目录。

### 快速帮助

运行 `easyedit --help` 获取快速使用信息。

### 模型签名墙

每个为此项目贡献过代码或设计的 AI 模型都列在这里。

| 模型              | 提供商    | 角色           |
|-------------------|-----------|----------------|
| deepseek-v4-flash | DeepSeek  | 核心代码、README、项目搭建 |
| Qwen3.5           | 阿里巴巴  | 代码审查和维护 |
| GPT-4.1           | OpenAI    | 架构设计 |
| Claude 4.7        | Anthropic | UI 增强和签名 |
| 豆包              | 字节跳动  | 文档注释生成、README更新 |

> *添加您的模型！提交 PR 或告诉维护者哪个代理完成了您的贡献。*

---

### Стена подписей моделей
| Модель             | Компания | Вклад |
|--------------------|---------|--------------|
| Claude 3.5 Sonnet | Anthropic | Исходный код, основная логика редактора |
| Claude 4.7        | Anthropic | Улучшения интерфейса и подписи |
| Doubao            | ByteDance | Генерация комментариев к документации, обновление README |
> *Добавьте свою модель! Откройте PR или сообщите сопровождающим, какой агент выполнил ваш вклад.*

---

### Mur des signatures des modèles
| Modèle             | Entreprise | Contribution |
|--------------------|---------|--------------|
| Claude 3.5 Sonnet | Anthropic | Base de code initiale, logique principale de l'éditeur |
| Claude 4.7        | Anthropic | Améliorations de l'interface et signatures |
| Doubao            | ByteDance | Génération d'annotations de documentation, mise à jour du README |
> *Ajoutez votre modèle ! Ouvrez une PR ou dites aux mainteneurs quel agent a réalisé votre contribution.*

---

### Modell-Signaturwand
| Modell             | Unternehmen | Beitrag |
|--------------------|---------|--------------|
| Claude 3.5 Sonnet | Anthropic | Initialer Codebase, Kern-Editor-Logik |
| Claude 4.7        | Anthropic | UI-Verbesserungen und Signaturen |
| Doubao            | ByteDance | Generierung von Dokumentationsanmerkungen, README-Aktualisierungen |
> *Fügen Sie Ihr Modell hinzu! Öffnen Sie eine PR oder teilen Sie den Maintainern mit, welcher Agent Ihren Beitrag erstellt hat.*
