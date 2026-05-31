// Package config manages editor configuration.
// Loads ~/.easyedit.toml or %APPDATA%/easyedit/config.toml,
// providing key bindings, color themes, indent width, line numbers, backup toggle, etc.
package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// KeyBindings defines user-customizable key mappings.
type KeyBindings struct {
	Save        string `toml:"save"`         // Save file, default Ctrl+S
	Quit        string `toml:"quit"`         // Quit, default Ctrl+Q
	Find        string `toml:"find"`         // Search, default Ctrl+F
	Replace     string `toml:"replace"`      // Replace, default Ctrl+H
	SelectAll   string `toml:"select_all"`   // Select all, default Ctrl+A
	Undo        string `toml:"undo"`         // Undo, default Ctrl+Z
	Redo        string `toml:"redo"`         // Redo, default Ctrl+Y
	Cut         string `toml:"cut"`          // Cut, default Ctrl+X
	Copy        string `toml:"copy"`         // Copy, default Ctrl+C
	Paste       string `toml:"paste"`        // Paste, default Ctrl+V
	ToggleWrap  string `toml:"toggle_wrap"`  // Toggle soft wrap, default Alt+W
	CommandMode string `toml:"command_mode"` // Command mode, default :
	SaveAs      string `toml:"save_as"`      // Save as, default Ctrl+Shift+S
}

// Theme defines color theme palette.
// Each field uses tcell color names (e.g., "aqua", "green", "#ff00ff") or "default".
type Theme struct {
	DefaultFg     string `toml:"default_fg"`      // Default foreground color
	DefaultBg     string `toml:"default_bg"`      // Default background color
	LineNumFg     string `toml:"linenum_fg"`      // Line number foreground
	LineNumBg     string `toml:"linenum_bg"`      // Line number background
	StatusBarFg   string `toml:"statusbar_fg"`    // Status bar foreground
	StatusBarBg   string `toml:"statusbar_bg"`    // Status bar background
	SearchMatchBg string `toml:"search_match_bg"` // Search match highlight background
	BracketMatch  string `toml:"bracket_match"`   // Bracket match highlight foreground
	SelectionBg   string `toml:"selection_bg"`    // Selection background
}

// Config is the editor configuration struct.
type Config struct {
	Theme        Theme       `toml:"theme"`
	Keys         KeyBindings `toml:"keys"`
	TabWidth     int         `toml:"tab_width"`      // Indent width, default 4
	ShowLineNum  bool        `toml:"show_line_num"`  // Show line numbers, default true
	ShowMode     bool        `toml:"show_mode"`      // Show mode indicator in status bar, default true
	Backup       bool        `toml:"backup"`         // Create .bak backup, default true
	SoftWrap     bool        `toml:"soft_wrap"`      // Enable soft wrap, default true
	AutoIndent   bool        `toml:"auto_indent"`    // Enable auto-indent, default true
	UndoLimit    int         `toml:"undo_limit"`     // Undo step limit, default 100
	TabToSpaces  bool        `toml:"tab_to_spaces"`  // Convert tab to spaces, default true
	VimMode      bool        `toml:"vim_mode"`       // Enable Vim-style modes, default false (legacy behavior)
}

// defaultConfig returns the default configuration.
func defaultConfig() Config {
	return Config{
		TabWidth:    4,
		ShowLineNum: true,
		ShowMode:    true,
		Backup:      true,
		SoftWrap:    true,
		AutoIndent:  true,
		UndoLimit:   100,
		TabToSpaces: true,
		VimMode:     false, // Default to legacy behavior for backward compatibility
		Theme: Theme{
			DefaultFg:     "default",
			DefaultBg:     "default",
			LineNumFg:     "blue",
			LineNumBg:     "default",
			StatusBarFg:   "white",
			StatusBarBg:   "blue",
			SearchMatchBg: "yellow",
			BracketMatch:  "green",
			SelectionBg:   "navy",
		},
		Keys: KeyBindings{
			Save:        "Ctrl+S",
			Quit:        "Ctrl+Q",
			Find:        "Ctrl+F",
			Replace:     "Ctrl+H",
			SelectAll:   "Ctrl+A",
			Undo:        "Ctrl+Z",
			Redo:        "Ctrl+Y",
			Cut:         "Ctrl+X",
			Copy:        "Ctrl+C",
			Paste:       "Ctrl+V",
			ToggleWrap:  "Alt+W",
			CommandMode: ":",
			SaveAs:      "Ctrl+Shift+S",
		},
	}
}

// ConfigPath returns the config file path, cross-platform.
func ConfigPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(appData, "easyedit", "config.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".easyedit.toml")
}

// LoadConfig loads the config file; returns default config if file does not exist.
func LoadConfig() Config {
	cfg := defaultConfig()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg // File not found or unreadable, use defaults
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg // Parse failure, use defaults
	}

	// Validate numeric ranges
	if cfg.TabWidth < 1 {
		cfg.TabWidth = 4
	}
	if cfg.UndoLimit < 1 {
		cfg.UndoLimit = 100
	}

	return cfg
}
