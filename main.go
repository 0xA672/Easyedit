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
`)
}
