package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// setConfigDir redirects the config path to dir for the duration of the test.
// Returns a cleanup function that restores the original environment.
func setConfigDir(t *testing.T, dir string) func() {
	t.Helper()
	if runtime.GOOS == "windows" {
		orig := os.Getenv("APPDATA")
		os.Setenv("APPDATA", dir)
		return func() {
			if orig == "" {
				os.Unsetenv("APPDATA")
			} else {
				os.Setenv("APPDATA", orig)
			}
		}
	}
	orig := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	return func() {
		if orig == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", orig)
		}
	}
}

func TestDefaultConfigTabWidth(t *testing.T) {
	cfg := defaultConfig()
	if cfg.TabWidth != 4 {
		t.Fatalf("TabWidth = %d, want 4", cfg.TabWidth)
	}
}

func TestDefaultConfigShowLineNum(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.ShowLineNum {
		t.Fatal("ShowLineNum should default to true")
	}
}

func TestDefaultConfigBackup(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.Backup {
		t.Fatal("Backup should default to true")
	}
}

func TestDefaultConfigSoftWrap(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.SoftWrap {
		t.Fatal("SoftWrap should default to true")
	}
}

func TestDefaultConfigAutoIndent(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.AutoIndent {
		t.Fatal("AutoIndent should default to true")
	}
}

func TestDefaultConfigUndoLimit(t *testing.T) {
	cfg := defaultConfig()
	if cfg.UndoLimit != 100 {
		t.Fatalf("UndoLimit = %d, want 100", cfg.UndoLimit)
	}
}

func TestDefaultConfigTabToSpaces(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.TabToSpaces {
		t.Fatal("TabToSpaces should default to true")
	}
}

func TestDefaultConfigThemeColors(t *testing.T) {
	cfg := defaultConfig()
	if cfg.Theme.DefaultFg != "default" {
		t.Fatalf("DefaultFg = %q, want %q", cfg.Theme.DefaultFg, "default")
	}
	if cfg.Theme.DefaultBg != "default" {
		t.Fatalf("DefaultBg = %q, want %q", cfg.Theme.DefaultBg, "default")
	}
	if cfg.Theme.LineNumFg != "blue" {
		t.Fatalf("LineNumFg = %q, want %q", cfg.Theme.LineNumFg, "blue")
	}
	if cfg.Theme.StatusBarBg != "blue" {
		t.Fatalf("StatusBarBg = %q, want %q", cfg.Theme.StatusBarBg, "blue")
	}
	if cfg.Theme.SearchMatchBg != "yellow" {
		t.Fatalf("SearchMatchBg = %q, want %q", cfg.Theme.SearchMatchBg, "yellow")
	}
	if cfg.Theme.BracketMatch != "green" {
		t.Fatalf("BracketMatch = %q, want %q", cfg.Theme.BracketMatch, "green")
	}
	if cfg.Theme.SelectionBg != "navy" {
		t.Fatalf("SelectionBg = %q, want %q", cfg.Theme.SelectionBg, "navy")
	}
}

func TestDefaultConfigKeys(t *testing.T) {
	cfg := defaultConfig()
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Save", cfg.Keys.Save, "Ctrl+S"},
		{"Quit", cfg.Keys.Quit, "Ctrl+Q"},
		{"Find", cfg.Keys.Find, "Ctrl+F"},
		{"Replace", cfg.Keys.Replace, "Ctrl+H"},
		{"SelectAll", cfg.Keys.SelectAll, "Ctrl+A"},
		{"Undo", cfg.Keys.Undo, "Ctrl+Z"},
		{"Redo", cfg.Keys.Redo, "Ctrl+Y"},
		{"Cut", cfg.Keys.Cut, "Ctrl+X"},
		{"Copy", cfg.Keys.Copy, "Ctrl+C"},
		{"Paste", cfg.Keys.Paste, "Ctrl+V"},
		{"ToggleWrap", cfg.Keys.ToggleWrap, "Alt+W"},
		{"CommandMode", cfg.Keys.CommandMode, ":"},
		{"SaveAs", cfg.Keys.SaveAs, "Ctrl+Shift+S"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Should return defaults when config file doesn't exist
	cfg := LoadConfig()
	if cfg.TabWidth != 4 {
		t.Fatalf("TabWidth = %d, want default 4", cfg.TabWidth)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	// On Windows: APPDATA/dir/easyedit/config.toml
	// On Unix: HOME/.easyedit.toml
	var configPath string
	if runtime.GOOS == "windows" {
		configDir := filepath.Join(dir, "easyedit")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		configPath = filepath.Join(configDir, "config.toml")
	} else {
		configPath = filepath.Join(dir, ".easyedit.toml")
	}

	content := []byte(`
tab_width = 8
show_line_num = false
backup = false
soft_wrap = false
auto_indent = false
undo_limit = 50
tab_to_spaces = false

[theme]
default_fg = "red"
default_bg = "black"

[keys]
save = "Ctrl+Shift+S"
`)
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cleanup := setConfigDir(t, dir)
	defer cleanup()

	cfg := LoadConfig()
	if cfg.TabWidth != 8 {
		t.Fatalf("TabWidth = %d, want 8", cfg.TabWidth)
	}
	if cfg.ShowLineNum {
		t.Fatal("ShowLineNum should be false")
	}
	if cfg.Backup {
		t.Fatal("Backup should be false")
	}
	if cfg.SoftWrap {
		t.Fatal("SoftWrap should be false")
	}
	if cfg.AutoIndent {
		t.Fatal("AutoIndent should be false")
	}
	if cfg.UndoLimit != 50 {
		t.Fatalf("UndoLimit = %d, want 50", cfg.UndoLimit)
	}
	if cfg.TabToSpaces {
		t.Fatal("TabToSpaces should be false")
	}
	if cfg.Theme.DefaultFg != "red" {
		t.Fatalf("Theme.DefaultFg = %q, want %q", cfg.Theme.DefaultFg, "red")
	}
	if cfg.Theme.DefaultBg != "black" {
		t.Fatalf("Theme.DefaultBg = %q, want %q", cfg.Theme.DefaultBg, "black")
	}
	if cfg.Keys.Save != "Ctrl+Shift+S" {
		t.Fatalf("Keys.Save = %q, want %q", cfg.Keys.Save, "Ctrl+Shift+S")
	}
}

func TestLoadConfigInvalidTOML(t *testing.T) {
	dir := t.TempDir()
	var configPath string
	if runtime.GOOS == "windows" {
		configDir := filepath.Join(dir, "easyedit")
		os.MkdirAll(configDir, 0755)
		configPath = filepath.Join(configDir, "config.toml")
	} else {
		configPath = filepath.Join(dir, ".easyedit.toml")
	}
	content := []byte(`invalid toml [[[`)
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cleanup := setConfigDir(t, dir)
	defer cleanup()

	// Should fall back to defaults
	cfg := LoadConfig()
	if cfg.TabWidth != 4 {
		t.Fatalf("TabWidth = %d, want default 4", cfg.TabWidth)
	}
}

func TestLoadConfigTabWidthClamp(t *testing.T) {
	dir := t.TempDir()
	var configPath string
	if runtime.GOOS == "windows" {
		configDir := filepath.Join(dir, "easyedit")
		os.MkdirAll(configDir, 0755)
		configPath = filepath.Join(configDir, "config.toml")
	} else {
		configPath = filepath.Join(dir, ".easyedit.toml")
	}
	content := []byte(`tab_width = 0`)
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cleanup := setConfigDir(t, dir)
	defer cleanup()

	cfg := LoadConfig()
	if cfg.TabWidth != 4 {
		t.Fatalf("TabWidth = %d, want 4 (clamped)", cfg.TabWidth)
	}
}

func TestLoadConfigUndoLimitClamp(t *testing.T) {
	dir := t.TempDir()
	var configPath string
	if runtime.GOOS == "windows" {
		configDir := filepath.Join(dir, "easyedit")
		os.MkdirAll(configDir, 0755)
		configPath = filepath.Join(configDir, "config.toml")
	} else {
		configPath = filepath.Join(dir, ".easyedit.toml")
	}
	content := []byte(`undo_limit = 0`)
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cleanup := setConfigDir(t, dir)
	defer cleanup()

	cfg := LoadConfig()
	if cfg.UndoLimit != 100 {
		t.Fatalf("UndoLimit = %d, want 100 (clamped)", cfg.UndoLimit)
	}
}

func TestConfigPathWindows(t *testing.T) {
	if dir := ConfigPath(); dir == "" {
		t.Fatal("ConfigPath() should not be empty")
	}
}
