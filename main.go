// EasyEdit — Hybrid terminal text editor.
//
// Usage:
//
//	easyedit                Open a blank file
//	easyedit <filename>     Open specified file
//
// Shortcuts:
//
//	Ctrl+S    Save
//	Ctrl+Q    Quit
//	Ctrl+F    Search
//	Ctrl+H    Replace
//	Ctrl+A    Select All
//	Ctrl+Z    Undo
//	Ctrl+Y    Redo
//	Ctrl+X    Cut
//	Ctrl+C    Copy
//	Ctrl+V    Paste
//	Alt+W     Toggle soft wrap
//	:         Enter command mode
//
// Command mode (press : to enter):
//
//	:q         Quit
//	:q!        Force quit (no save)
//	:w         Save
//	:wq        Save and quit
//	:w <path>  Save as
//	:e <path>  Open file
//	:10,20d    Delete lines 10-20
//	:%s/old/new/g    Replace all
//	:3,5s/old/new    Replace in range
//	:set nu    Show line numbers
//	:set nonu  Hide line numbers
//	:42        Jump to line 42
//	:uninstall Remove editor binary, config, and backup files
//
// Config file (auto-loaded):
//
//	Linux/macOS: ~/.easyedit.toml
//	Windows:     %APPDATA%/easyedit/config.toml
//
// Dependencies:
//   - github.com/gdamore/tcell/v2      (terminal handling)
//   - github.com/alecthomas/chroma/v2   (syntax highlighting)
//   - github.com/atotto/clipboard       (system clipboard)
//   - github.com/BurntSushi/toml        (config parsing)
//   - github.com/mattn/go-runewidth     (character width)
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"easyedit/ui"
)

// Version holds the editor version. Priority:
//  1. runtime/debug.ReadBuildInfo() — picks up module tag from 'go install'
//  2. -ldflags "-X main.Version=x.y.z" — injected at link time
//  3. "0.2.0-dev" — hardcoded fallback
var Version = "0.2.0-dev"

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v != "" && v != "(devel)" {
			Version = v
		}
	}
}

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	showHelp := flag.Bool("help", false, "print this help message")
	flag.Usage = printUsage
	flag.Parse()

	if *showVersion {
		fmt.Printf("EasyEdit %s\n", Version)
		os.Exit(0)
	}

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	editor := ui.NewEditor()

	// If a positional argument is provided, open that file
	if args := flag.Args(); len(args) > 0 {
		filePath := args[0]
		if err := editor.OpenFile(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := editor.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	// Detect system language for localized help
	lang := os.Getenv("LANG")
	isZh := len(lang) >= 2 && lang[:2] == "zh"

	if isZh {
		os.Stdout.WriteString(`EasyEdit — 混合终端文本编辑器

用法:
  easyedit [flags] [file]

标志:
  --help       打印此帮助信息
  --version    打印版本并退出

示例:
  easyedit              打开空白文件
  easyedit main.go      打开 main.go 进行编辑

键位绑定（在编辑器内）:
  Ctrl+S    保存        Ctrl+Q    退出
  Ctrl+F    搜索        Ctrl+H    替换
  Ctrl+Z    撤销        Ctrl+Y    重做
  Ctrl+X    剪切        Ctrl+C    复制        Ctrl+V    粘贴
  Ctrl+A    全选        Alt+W     切换软换行
  :         进入命令模式

命令 (: 模式):
  :q        退出           :q!       强制退出
  :w        保存           :w <path> 另存为
  :wq       保存并退出     :e <path> 打开文件
  :42       跳转到第 42 行   :set nu   显示行号
  :%s/a/b/g 全部替换       :10,20d   删除第 10-20 行

更多信息请访问：https://github.com/0xA672/Easyedit
文档：https://github.com/0xA672/Easyedit/tree/main/docs
`)
	} else {
		os.Stdout.WriteString(`EasyEdit — Hybrid terminal text editor

Usage:
  easyedit [flags] [file]

Flags:
  --help       Print this help message
  --version    Print version and exit

Examples:
  easyedit              Open a blank file
  easyedit main.go      Open main.go for editing

Keybindings (inside the editor):
  Ctrl+S    Save        Ctrl+Q    Quit
  Ctrl+F    Search      Ctrl+H    Replace
  Ctrl+Z    Undo        Ctrl+Y    Redo
  Ctrl+X    Cut         Ctrl+C    Copy        Ctrl+V    Paste
  Ctrl+A    Select All  Alt+W     Toggle soft wrap
  :         Enter command mode

Commands (: mode):
  :q        Quit           :q!       Force quit
  :w        Save           :w <path> Save as
  :wq       Save & quit    :e <path> Open file
  :42       Go to line 42  :set nu   Show line numbers
  :%s/a/b/g Replace all    :10,20d   Delete lines 10-20

More info: https://github.com/0xA672/Easyedit
Docs: https://github.com/0xA672/Easyedit/tree/main/docs
`)
	}
}
