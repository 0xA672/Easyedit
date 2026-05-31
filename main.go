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
//	:%s/old/new/g    Replace all (regex)
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
//
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
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/0xA672/Easyedit/ui"
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

// stringSliceFlag allows multiple -e options.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	showHelp := flag.Bool("help", false, "print this help message")
	streamMode := flag.Bool("s", false, "stream mode: read from stdin, apply edits, write to stdout (sed-like)")
	streamModeLong := flag.Bool("stream", false, "stream mode: read from stdin, apply edits, write to stdout (sed-like)")

	var editScripts stringSliceFlag
	flag.Var(&editScripts, "e", "edit script to apply in stream mode (e.g., 's/old/new/g'). May be repeated.")

	inPlace := flag.Bool("i", false, "in-place editing (used with -s)")

	flag.Usage = printUsage
	flag.Parse()

	// --version and --help always work, even with -s
	if *showVersion {
		fmt.Printf("EasyEdit %s\n", Version)
		os.Exit(0)
	}
	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	isStream := *streamMode || *streamModeLong
	if isStream {
		runStreamMode(editScripts, *inPlace)
		return
	}

	editor := ui.NewEditor()
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

// runStreamMode processes input line by line, applying sed-like rules.
// It handles both stdin → stdout and in‑place file editing atomically.
func runStreamMode(scripts []string, inPlace bool) {
	if len(scripts) == 0 {
		fmt.Fprintln(os.Stderr, "Error: at least one -e flag is required in stream mode")
		os.Exit(1)
	}

	// Parse each script into an internal Rule
	var rules []Rule
	for _, scr := range scripts {
		parsed, err := parseRule(scr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing script %q: %v\n", scr, err)
			os.Exit(1)
		}
		rules = append(rules, parsed)
	}

	if inPlace {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: filename required for in-place editing")
			os.Exit(1)
		}
		filename := args[0]

		// Process the file line by line into a temporary file
		tmpFile, err := os.CreateTemp(filepath.Dir(filename), ".easyedit-*.tmp")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
			os.Exit(1)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath) // clean up on failure

		// Open input file
		inFile, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		writer := bufio.NewWriter(tmpFile)

		// Process line by line
		for scanner.Scan() {
			line := scanner.Text()
			newLine := applyRulesToLine(line, rules)
			if _, err := writer.WriteString(newLine + "\n"); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to temp file: %v\n", err)
				os.Exit(1)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		writer.Flush()
		tmpFile.Close()

		// Atomic replace: rename temp file to original
		if err := os.Rename(tmpPath, filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error replacing file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Read from stdin, write to stdout (line by line)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			newLine := applyRulesToLine(line, rules)
			fmt.Println(newLine)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}
}

// Rule represents one sed‑like command: either a substitution or a delete.
type Rule struct {
	IsDelete bool
	// For substitution
	Pattern *regexp.Regexp
	Repl    string
	Global  bool
	// For delete
	DeletePattern *regexp.Regexp
}

// parseRule parses a single sed‑style command.
// Supported forms:
//
//	s/pattern/replacement/flags   (any delimiter, e.g. s|a|b|)
//	/pattern/d                    (delete lines matching pattern)
//
// Flags: 'g' for global replacement.
// Delimiter is the first character after 's' (or after the slash for delete).
// Escaped delimiters inside pattern/replacement are supported (e.g. s/\/foo/bar/).
func parseRule(script string) (Rule, error) {
	script = strings.TrimSpace(script)
	if len(script) < 3 {
		return Rule{}, fmt.Errorf("too short")
	}

	// Delete command: /pattern/d
	if script[0] == '/' && strings.HasSuffix(script, "/d") {
		delim := '/'
		// Find closing delimiter before the trailing "/d"
		endIdx := len(script) - 2 // position of the slash before 'd'
		pattern, err := unescapeDelimited(script[1:endIdx], delim)
		if err != nil {
			return Rule{}, fmt.Errorf("invalid delete pattern: %v", err)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return Rule{}, fmt.Errorf("invalid regex in delete: %v", err)
		}
		return Rule{IsDelete: true, DeletePattern: re}, nil
	}

	// Substitute command: s<delim>pattern<delim>replacement<delim>flags
	if script[0] != 's' {
		return Rule{}, fmt.Errorf("unsupported command, expected 's' or '/pattern/d'")
	}
	if len(script) < 4 {
		return Rule{}, fmt.Errorf("substitution command too short")
	}
	delim := rune(script[1])
	// Find the three occurrences of the delimiter, respecting escapes
	parts, err := splitDelimited(script[2:], delim)
	if err != nil || len(parts) < 2 {
		return Rule{}, fmt.Errorf("invalid substitution syntax: missing delimiters")
	}
	pattern := parts[0]
	replacement := parts[1]
	flags := ""
	if len(parts) > 2 {
		flags = parts[2]
	}
	global := strings.Contains(flags, "g")

	re, err := regexp.Compile(pattern)
	if err != nil {
		return Rule{}, fmt.Errorf("invalid regex: %v", err)
	}
	return Rule{
		IsDelete: false,
		Pattern:  re,
		Repl:     replacement,
		Global:   global,
	}, nil
}

// splitDelimited splits a string by the given delimiter, respecting backslash escapes.
// Returns a slice of substrings between delimiters (excluding the delimiters themselves).
// Example: splitDelimited("foo/bar/baz", '/') -> ["foo", "bar", "baz"]
func splitDelimited(s string, delim rune) ([]string, error) {
	var parts []string
	var cur strings.Builder
	escaped := false
	for _, ch := range s {
		if !escaped && ch == delim {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		if !escaped && ch == '\\' {
			escaped = true
			continue
		}
		if escaped {
			// add the escaped character literally (including the backslash? No, sed unescapes)
			cur.WriteRune(ch)
			escaped = false
		} else {
			cur.WriteRune(ch)
		}
	}
	if escaped {
		return nil, fmt.Errorf("unterminated escape")
	}
	// Last part after final delimiter may be empty (flags)
	parts = append(parts, cur.String())
	return parts, nil
}

// unescapeDelimited removes backslashes used to escape the delimiter inside a pattern.
// E.g., unescapeDelimited("foo\/bar", '/') -> "foo/bar"
func unescapeDelimited(s string, delim rune) (string, error) {
	var out strings.Builder
	escaped := false
	for _, ch := range s {
		if !escaped && ch == '\\' {
			escaped = true
			continue
		}
		if escaped {
			if ch == delim || ch == '\\' {
				out.WriteRune(ch)
			} else {
				out.WriteRune('\\')
				out.WriteRune(ch)
			}
			escaped = false
		} else {
			out.WriteRune(ch)
		}
	}
	if escaped {
		return "", fmt.Errorf("unterminated escape")
	}
	return out.String(), nil
}

// applyRulesToLine applies all rules to a single line and returns the result.
// For substitution rules, it performs either first-match or global replacement.
// For delete rules, it returns an empty string if the line matches, otherwise unchanged.
func applyRulesToLine(line string, rules []Rule) string {
	for _, rule := range rules {
		if rule.IsDelete {
			if rule.DeletePattern.MatchString(line) {
				return "" // line will be omitted (caller handles newline)
			}
		} else {
			if rule.Global {
				line = rule.Pattern.ReplaceAllString(line, rule.Repl)
			} else {
				// Replace only the first match
				idx := rule.Pattern.FindStringIndex(line)
				if idx != nil {
					line = line[:idx[0]] + rule.Pattern.ReplaceAllString(line[idx[0]:idx[1]], rule.Repl) + line[idx[1]:]
				}
			}
		}
	}
	return line
}

// printUsage prints localized help text.
func printUsage() {
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
  -e <script>  指定编辑脚本（可以多次使用 -e）
  -i           原地修改文件（与 -s 一起使用）
示例:
  easyedit              打开空白文件
  easyedit main.go      打开 main.go 进行编辑
  echo "hello foo" | easyedit -s -e 's/foo/bar/g'
  cat input.txt | easyedit -s -e 's/old/new/g' > output.txt
  easyedit -s -i -e 's/a/b/g' config.txt
  easyedit -s -e '/^#/d' -e 's/foo/bar/g' file.txt

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
  -e <script>    Edit script to apply (can be used multiple times)
  -i             In-place file editing (used with -s)
Examples:
  easyedit                              Open a blank file
  easyedit main.go                      Open main.go for editing
  echo "hello foo" | easyedit -s -e 's/foo/bar/g'
  cat input.txt | easyedit -s -e 's/old/new/g' > output.txt
  easyedit -s -i -e 's/a/b/g' config.txt
  easyedit -s -e '/^#/d' -e 's/foo/bar/g' file.txt

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
