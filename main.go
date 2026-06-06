// EasyEdit — Hybrid terminal text editor.
//
// Usage:
//
// easyedit                Open a blank file
// easyedit <filename>     Open specified file
//
// Shortcuts:
//
// Ctrl+S    Save
// Ctrl+Q    Quit
// Ctrl+F    Search
// Ctrl+H    Replace
// Ctrl+A    Select All
// Ctrl+Z    Undo
// Ctrl+Y    Redo
// Ctrl+X    Cut
// Ctrl+C    Copy
// Ctrl+V    Paste
// Alt+W     Toggle soft wrap
// :         Enter command mode
//
// Command mode (press : to enter):
//
// :q         Quit
// :q!        Force quit (no save)
// :w         Save
// :wq        Save and quit
// :w <path>  Save as
// :e <path>  Open file
// :10,20d    Delete lines 10-20
// :%s/old/new/g    Replace all
// :3,5s/old/new    Replace in range
// :set nu    Show line numbers
// :set nonu  Hide line numbers
// :42        Jump to line 42
// :uninstall Remove editor binary, config, and backup files
//
// Config file (auto-loaded):
//
// Linux/macOS: ~/.easyedit.toml
// Windows:     %APPDATA%/easyedit/config.toml
//
// Dependencies:
//   - github.com/gdamore/tcell/v2      (terminal handling)
//   - github.com/alecthomas/chroma/v2   (syntax highlighting)
//   - github.com/atotto/clipboard       (system clipboard)
//   - github.com/BurntSushi/toml        (config parsing)
//   - github.com/mattn/go-runewidth     (character width)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/0xA672/Easyedit/ui"
)

// Version holds the editor version. Priority:
//  1. runtime/debug.ReadBuildInfo() — picks up module tag from 'go install'
//  2. -ldflags "-X main.Version=x.y.z" — injected at link time
//  3. "0.2.0-dev" — hardcoded fallback
var Version = "0.2.0-dev"

// init 初始化编辑器版本信息，优先从构建信息中获取版本号，其次使用链接时注入的版本，最后使用硬编码的默认版本。
func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v != "" && v != "(devel)" {
			Version = v
		}
	}
}

// main 程序的入口函数，处理命令行参数，初始化编辑器实例，启动编辑器主循环。
func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	showHelp := flag.Bool("help", false, "print this help message")
	streamMode := flag.Bool("s", false, "stream mode: read from stdin, apply edits, write to stdout (sed-like)")
	streamModeLong := flag.Bool("stream", false, "stream mode: read from stdin, apply edits, write to stdout (sed-like)")
	editScript := flag.String("e", "", "edit script to apply in stream mode (e.g., 's/old/new/g')")
	inPlace := flag.Bool("i", false, "in-place editing (used with -s)")
	flag.Usage = printUsage
	flag.Parse()
	// Check if stream mode is enabled
	isStream := *streamMode || *streamModeLong
	if isStream {
		runStreamMode(*editScript, *inPlace)
		return
	}
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

// runStreamMode 运行流模式，提供类似sed的批量编辑能力，从输入读取内容，应用编辑脚本后输出结果，支持原地修改文件。
func runStreamMode(script string, inPlace bool) {
	if script == "" {
		fmt.Fprintln(os.Stderr, "Error: -e flag is required in stream mode")
		os.Exit(1)
	}
	// Parse the script (simple sed-like syntax for now)
	// Supports: s/old/new/g, s/old/new/, d (delete lines matching pattern)
	rules := parseScript(script)
	if inPlace {
		// Read from file argument
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: filename required for in-place editing")
			os.Exit(1)
		}
		filename := args[0]
		// Read file content
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		// Apply rules
		result := applyRules(string(content), rules)
		// Write back to file
		if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Read from stdin, write to stdout
		scanner := bufio.NewScanner(os.Stdin)
		var input strings.Builder
		for scanner.Scan() {
			input.WriteString(scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		// Apply rules
		result := applyRules(input.String(), rules)
		// Write to stdout (trim trailing newline added by loop)
		fmt.Print(strings.TrimSuffix(result, "\n"))
	}
}

// parseScript 解析编辑脚本，将输入的脚本按分隔符分割为多个独立的编辑规则，支持分号或换行分隔多个命令。
func parseScript(script string) []string {
	// Simple parser: split by semicolon or handle multiple -e flags
	// For now, treat the whole script as one rule or split by common delimiters
	var rules []string
	// Handle multiple commands separated by semicolons or newlines
	parts := strings.FieldsFunc(script, func(r rune) bool {
		return r == ';' || r == '\n'
	})
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			rules = append(rules, part)
		}
	}
	return rules
}

// applyRules 将解析后的编辑规则应用到输入内容上，支持字符串替换、匹配行删除等sed风格的编辑操作。
func applyRules(content string, rules []string) string {
	result := content
	for _, rule := range rules {
		// Parse sed-like syntax: s/pattern/replacement/flags
		if strings.HasPrefix(rule, "s/") {
			parts := strings.SplitN(rule[2:], "/", 3)
			if len(parts) >= 3 {
				pattern := parts[0]
				replacement := parts[1]
				flags := ""
				if len(parts) == 3 {
					flags = parts[2]
				}
				// Simple string replacement (not full regex for safety)
				if strings.Contains(flags, "g") {
					result = strings.ReplaceAll(result, pattern, replacement)
				} else {
					result = strings.Replace(result, pattern, replacement, 1)
				}
			}
		} else if strings.HasPrefix(rule, "/") && strings.HasSuffix(rule, "/d") && len(rule) > 2 {
			// Delete lines matching pattern: /pattern/d
			pattern := rule[1 : len(rule)-2]
			lines := strings.Split(result, "\n")
			var filtered []string
			for _, line := range lines {
				if !strings.Contains(line, pattern) {
					filtered = append(filtered, line)
				}
			}
			result = strings.Join(filtered, "\n")
		}
	}
	return result
}

// printUsage 打印程序的使用帮助信息，会根据系统语言自动选择中文或英文的帮助内容，适配不同语言的用户。
func printUsage() {
	// Detect system language for localized help
	lang := os.Getenv("LANG")
	isZh := len(lang) >= 2 && lang[:2] == "zh"
	if isZh {
		os.Stdout.WriteString(`EasyEdit — 混合终端文本编辑器
用法:
  easyedit [flags] [file]
  echo "text" | easyedit -s -e 's/old/new/g'
  easyedit -s -i -e 's/foo/bar/g' file.txt
标志:
  --help       打印此帮助信息
  --version    打印版本并退出
  -s, --stream 流模式：从 stdin 读取，应用编辑规则后输出到 stdout（sed 平替）
  -e <script>  指定编辑脚本（与 -s 一起使用）
  -i           原地修改文件（与 -s 一起使用）
示例:
  easyedit              打开空白文件
  easyedit main.go      打开 main.go 进行编辑
  echo "hello foo" | easyedit -s -e 's/foo/bar/g'
  cat input.txt | easyedit -s -e 's/old/new/g' > output.txt
  easyedit -s -i -e 's/a/b/g' config.txt
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
  echo "text" | easyedit -s -e 's/old/new/g'
  easyedit -s -i -e 's/foo/bar/g' file.txt
Flags:
  --help         Print this help message
  --version      Print version and exit
  -s, --stream   Stream mode: read from stdin, apply edits, write to stdout (sed alternative)
  -e <script>    Edit script to apply (used with -s)
  -i             In-place file editing (used with -s)
  Examples:
  easyedit                              Open a blank file
  easyedit main.go                      Open main.go for editing
  echo "hello foo" | easyedit -s -e 's/foo/bar/g'
  cat input.txt | easyedit -s -e 's/old/new/g' > output.txt
  easyedit -s -i -e 's/a/b/g' config.txt
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
