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
	"fmt"
	"os"

	"easyedit/ui"
)

func main() {
	editor := ui.NewEditor()

	// If a command-line argument is provided, open that file
	if len(os.Args) > 1 {
		filePath := os.Args[1]
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
