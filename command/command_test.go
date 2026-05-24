package command

import (
	"errors"
	"testing"
)

// ────────────────────────── Parse tests ──────────────────────────

func TestParseEmpty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if err.Error() != "empty command" {
		t.Fatalf("got %q, want %q", err.Error(), "empty command")
	}
}

func TestParseWhitespace(t *testing.T) {
	_, err := Parse("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only input")
	}
}

func TestParseQuit(t *testing.T) {
	cmd, err := Parse("q")
	if err != nil {
		t.Fatalf("Parse('q') error: %v", err)
	}
	if cmd.Kind != CmdQuit {
		t.Fatalf("kind = %d, want CmdQuit", cmd.Kind)
	}
}

func TestParseForceQuit(t *testing.T) {
	cmd, err := Parse("q!")
	if err != nil {
		t.Fatalf("Parse('q!') error: %v", err)
	}
	if cmd.Kind != CmdForceQuit {
		t.Fatalf("kind = %d, want CmdForceQuit", cmd.Kind)
	}
}

func TestParseWrite(t *testing.T) {
	cmd, err := Parse("w")
	if err != nil {
		t.Fatalf("Parse('w') error: %v", err)
	}
	if cmd.Kind != CmdWrite {
		t.Fatalf("kind = %d, want CmdWrite", cmd.Kind)
	}
}

func TestParseWriteWithPath(t *testing.T) {
	cmd, err := Parse("w /tmp/out.txt")
	if err != nil {
		t.Fatalf("Parse('w /tmp/out.txt') error: %v", err)
	}
	if cmd.Kind != CmdWrite {
		t.Fatalf("kind = %d, want CmdWrite", cmd.Kind)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "/tmp/out.txt" {
		t.Fatalf("args = %v, want ['/tmp/out.txt']", cmd.Args)
	}
}

func TestParseWriteQuit(t *testing.T) {
	cmd, err := Parse("wq")
	if err != nil {
		t.Fatalf("Parse('wq') error: %v", err)
	}
	if cmd.Kind != CmdWriteQuit {
		t.Fatalf("kind = %d, want CmdWriteQuit", cmd.Kind)
	}
}

func TestParseWriteQuitBang(t *testing.T) {
	cmd, err := Parse("wq!")
	if err != nil {
		t.Fatalf("Parse('wq!') error: %v", err)
	}
	if cmd.Kind != CmdWriteQuit {
		t.Fatalf("kind = %d, want CmdWriteQuit", cmd.Kind)
	}
}

func TestParseEdit(t *testing.T) {
	cmd, err := Parse("e main.go")
	if err != nil {
		t.Fatalf("Parse('e main.go') error: %v", err)
	}
	if cmd.Kind != CmdEdit {
		t.Fatalf("kind = %d, want CmdEdit", cmd.Kind)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "main.go" {
		t.Fatalf("args = %v, want ['main.go']", cmd.Args)
	}
}

func TestParseDelete(t *testing.T) {
	cmd, err := Parse("d")
	if err != nil {
		t.Fatalf("Parse('d') error: %v", err)
	}
	if cmd.Kind != CmdDelete {
		t.Fatalf("kind = %d, want CmdDelete", cmd.Kind)
	}
}

func TestParseDeleteWithRange(t *testing.T) {
	cmd, err := Parse("10,20d")
	if err != nil {
		t.Fatalf("Parse('10,20d') error: %v", err)
	}
	if cmd.Kind != CmdDelete {
		t.Fatalf("kind = %d, want CmdDelete", cmd.Kind)
	}
	if cmd.Range.Start != 9 {
		t.Fatalf("Range.Start = %d, want 9", cmd.Range.Start)
	}
	if cmd.Range.End != 19 {
		t.Fatalf("Range.End = %d, want 19", cmd.Range.End)
	}
}

func TestParseDeleteWithAll(t *testing.T) {
	cmd, err := Parse("%d")
	if err != nil {
		t.Fatalf("Parse('%%d') error: %v", err)
	}
	if cmd.Kind != CmdDelete {
		t.Fatalf("kind = %d, want CmdDelete", cmd.Kind)
	}
	if !cmd.Range.All {
		t.Fatal("Range.All should be true")
	}
}

func TestParseSubstitute(t *testing.T) {
	cmd, err := Parse("s/old/new/")
	if err != nil {
		t.Fatalf("Parse('s/old/new/') error: %v", err)
	}
	if cmd.Kind != CmdSubstitute {
		t.Fatalf("kind = %d, want CmdSubstitute", cmd.Kind)
	}
	if cmd.SubOld != "old" {
		t.Fatalf("SubOld = %q, want %q", cmd.SubOld, "old")
	}
	if cmd.SubNew != "new" {
		t.Fatalf("SubNew = %q, want %q", cmd.SubNew, "new")
	}
}

func TestParseSubstituteGlobal(t *testing.T) {
	cmd, err := Parse("s/old/new/g")
	if err != nil {
		t.Fatalf("Parse('s/old/new/g') error: %v", err)
	}
	if cmd.Kind != CmdSubstitute {
		t.Fatalf("kind = %d, want CmdSubstitute", cmd.Kind)
	}
	if cmd.SubFlg != "g" {
		t.Fatalf("SubFlg = %q, want %q", cmd.SubFlg, "g")
	}
}

func TestParseSubstituteIgnoreCase(t *testing.T) {
	cmd, err := Parse("s/old/new/gi")
	if err != nil {
		t.Fatalf("Parse('s/old/new/gi') error: %v", err)
	}
	if cmd.SubFlg != "gi" {
		t.Fatalf("SubFlg = %q, want %q", cmd.SubFlg, "gi")
	}
}

func TestParseSubstituteRange(t *testing.T) {
	cmd, err := Parse("3,5s/old/new/")
	if err != nil {
		t.Fatalf("Parse('3,5s/old/new/') error: %v", err)
	}
	if cmd.Kind != CmdSubstitute {
		t.Fatalf("kind = %d, want CmdSubstitute", cmd.Kind)
	}
	if cmd.Range.Start != 2 {
		t.Fatalf("Range.Start = %d, want 2", cmd.Range.Start)
	}
	if cmd.Range.End != 4 {
		t.Fatalf("Range.End = %d, want 4", cmd.Range.End)
	}
}

func TestParseSubstitutePercent(t *testing.T) {
	cmd, err := Parse("%s/old/new/")
	if err != nil {
		t.Fatalf("Parse('%%s/old/new/') error: %v", err)
	}
	if cmd.Kind != CmdSubstitute {
		t.Fatalf("kind = %d, want CmdSubstitute", cmd.Kind)
	}
	if !cmd.Range.All {
		t.Fatal("Range.All should be true")
	}
}

func TestParseSet(t *testing.T) {
	cmd, err := Parse("set nu")
	if err != nil {
		t.Fatalf("Parse('set nu') error: %v", err)
	}
	if cmd.Kind != CmdSet {
		t.Fatalf("kind = %d, want CmdSet", cmd.Kind)
	}
	if cmd.SetOpt == nil {
		t.Fatal("SetOpt should not be nil")
	}
	if cmd.SetOpt.Name != "nu" {
		t.Fatalf("SetOpt.Name = %q, want %q", cmd.SetOpt.Name, "nu")
	}
}

func TestParseSetWithValue(t *testing.T) {
	cmd, err := Parse("set tabwidth=4")
	if err != nil {
		t.Fatalf("Parse('set tabwidth=4') error: %v", err)
	}
	if cmd.Kind != CmdSet {
		t.Fatalf("kind = %d, want CmdSet", cmd.Kind)
	}
	if cmd.SetOpt.Name != "tabwidth" {
		t.Fatalf("SetOpt.Name = %q, want %q", cmd.SetOpt.Name, "tabwidth")
	}
	if cmd.SetOpt.Value != "4" {
		t.Fatalf("SetOpt.Value = %q, want %q", cmd.SetOpt.Value, "4")
	}
}

func TestParseGoto(t *testing.T) {
	cmd, err := Parse("42")
	if err != nil {
		t.Fatalf("Parse('42') error: %v", err)
	}
	if cmd.Kind != CmdGoto {
		t.Fatalf("kind = %d, want CmdGoto", cmd.Kind)
	}
	if cmd.Range.Start != 41 {
		t.Fatalf("Range.Start = %d, want 41", cmd.Range.Start)
	}
}

func TestParsePercentGoto(t *testing.T) {
	// % without a command is incomplete (need a command after range marker)
	_, err := Parse("%")
	if err == nil {
		t.Fatal("expected error for '%' without command")
	}
}

func TestParseUninstall(t *testing.T) {
	cmd, err := Parse("uninstall")
	if err != nil {
		t.Fatalf("Parse('uninstall') error: %v", err)
	}
	if cmd.Kind != CmdUninstall {
		t.Fatalf("kind = %d, want CmdUninstall", cmd.Kind)
	}
}

func TestParseUnknownCommand(t *testing.T) {
	_, err := Parse("xyz")
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestParseIncompleteAfterRange(t *testing.T) {
	// "10," is treated as goto (ambiguous range is resolved as goto)
	cmd, err := Parse("10,")
	if err != nil {
		t.Fatalf("Parse('10,') should not error, got: %v", err)
	}
	if cmd.Kind != CmdGoto {
		t.Fatalf("kind = %d, want CmdGoto", cmd.Kind)
	}
	if cmd.Range.Start != 9 {
		t.Fatalf("Range.Start = %d, want 9", cmd.Range.Start)
	}
}

// ────────────────────────── parseSetOption tests ──────────────────────────

func TestParseSetOptionBool(t *testing.T) {
	opt := parseSetOption("nu")
	if opt == nil {
		t.Fatal("expected non-nil option")
	}
	if opt.Name != "nu" {
		t.Fatalf("Name = %q, want %q", opt.Name, "nu")
	}
	if opt.Value != "" {
		t.Fatalf("Value = %q, want empty", opt.Value)
	}
}

func TestParseSetOptionKeyValue(t *testing.T) {
	opt := parseSetOption("tabwidth=8")
	if opt == nil {
		t.Fatal("expected non-nil option")
	}
	if opt.Name != "tabwidth" {
		t.Fatalf("Name = %q, want %q", opt.Name, "tabwidth")
	}
	if opt.Value != "8" {
		t.Fatalf("Value = %q, want %q", opt.Value, "8")
	}
}

func TestParseSetOptionWithSpaces(t *testing.T) {
	opt := parseSetOption("  indent width  =  4  ")
	if opt == nil {
		t.Fatal("expected non-nil option")
	}
	if opt.Name != "indent width" {
		t.Fatalf("Name = %q, want %q", opt.Name, "indent width")
	}
	if opt.Value != "4" {
		t.Fatalf("Value = %q, want %q", opt.Value, "4")
	}
}

func TestParseSetOptionEmpty(t *testing.T) {
	opt := parseSetOption("")
	if opt != nil {
		t.Fatal("expected nil for empty input")
	}
}

// ────────────────────────── ExpandRange tests ──────────────────────────

func TestExpandRangeAll(t *testing.T) {
	r := Range{All: true}
	start, end := r.ExpandRange(10)
	if start != 0 {
		t.Fatalf("start = %d, want 0", start)
	}
	if end != 9 {
		t.Fatalf("end = %d, want 9", end)
	}
}

func TestExpandRangeExplicit(t *testing.T) {
	r := Range{Start: 2, End: 5}
	start, end := r.ExpandRange(10)
	if start != 2 {
		t.Fatalf("start = %d, want 2", start)
	}
	if end != 5 {
		t.Fatalf("end = %d, want 5", end)
	}
}

func TestExpandRangeDefaultStart(t *testing.T) {
	r := Range{Start: -1, End: -1}
	start, end := r.ExpandRange(10)
	if start != 0 {
		t.Fatalf("start = %d, want 0", start)
	}
	if end != -1 {
		t.Fatalf("end = %d, want -1", end)
	}
}

func TestExpandRangeToEnd(t *testing.T) {
	r := Range{Start: 3, End: -2}
	start, end := r.ExpandRange(10)
	if start != 3 {
		t.Fatalf("start = %d, want 3", start)
	}
	if end != 9 {
		t.Fatalf("end = %d, want 9", end)
	}
}

// ────────────────────────── CompileRegex tests ──────────────────────────

func TestCompileRegexBasic(t *testing.T) {
	cmd, _ := Parse("s/hello/world/")
	re, err := cmd.CompileRegex()
	if err != nil {
		t.Fatalf("CompileRegex error: %v", err)
	}
	if !re.MatchString("hello") {
		t.Fatal("regex should match 'hello'")
	}
	if re.MatchString("Hello") {
		t.Fatal("regex should not match 'Hello' (case-sensitive)")
	}
}

func TestCompileRegexCaseInsensitive(t *testing.T) {
	cmd, _ := Parse("s/hello/world/i")
	re, err := cmd.CompileRegex()
	if err != nil {
		t.Fatalf("CompileRegex error: %v", err)
	}
	if !re.MatchString("HELLO") {
		t.Fatal("regex should match 'HELLO' (case-insensitive)")
	}
}

func TestCompileRegexInvalid(t *testing.T) {
	cmd, _ := Parse("s/[invalid/world/")
	_, err := cmd.CompileRegex()
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

// ────────────────────────── MakeGlobal tests ──────────────────────────

func TestMakeGlobalTrue(t *testing.T) {
	cmd, _ := Parse("s/a/b/g")
	if !cmd.MakeGlobal() {
		t.Fatal("expected global=true for 'g' flag")
	}
}

func TestMakeGlobalFalse(t *testing.T) {
	cmd, _ := Parse("s/a/b/")
	if cmd.MakeGlobal() {
		t.Fatal("expected global=false for no flag")
	}
}

func TestMakeGlobalCombined(t *testing.T) {
	cmd, _ := Parse("s/a/b/gi")
	if !cmd.MakeGlobal() {
		t.Fatal("expected global=true for 'gi' flags")
	}
}

// ────────────────────────── String() tests ──────────────────────────

func TestStringQuit(t *testing.T) {
	cmd, _ := Parse("q")
	if s := cmd.String(); s != "q" {
		t.Fatalf("String() = %q, want %q", s, "q")
	}
}

func TestStringForceQuit(t *testing.T) {
	cmd, _ := Parse("q!")
	if s := cmd.String(); s != "q!" {
		t.Fatalf("String() = %q, want %q", s, "q!")
	}
}

func TestStringWrite(t *testing.T) {
	cmd, _ := Parse("w")
	if s := cmd.String(); s != "w" {
		t.Fatalf("String() = %q, want %q", s, "w")
	}
}

func TestStringWriteWithPath(t *testing.T) {
	cmd, _ := Parse("w /tmp/out.txt")
	if s := cmd.String(); s != "w /tmp/out.txt" {
		t.Fatalf("String() = %q, want %q", s, "w /tmp/out.txt")
	}
}

func TestStringWriteQuit(t *testing.T) {
	cmd, _ := Parse("wq")
	if s := cmd.String(); s != "wq" {
		t.Fatalf("String() = %q, want %q", s, "wq")
	}
}

func TestStringEdit(t *testing.T) {
	cmd, _ := Parse("e /path/to/file.go")
	if s := cmd.String(); s != "e /path/to/file.go" {
		t.Fatalf("String() = %q, want %q", s, "e /path/to/file.go")
	}
}

func TestStringSubstitute(t *testing.T) {
	cmd, _ := Parse("s/old/new/g")
	if s := cmd.String(); s != "s/old/new/g" {
		t.Fatalf("String() = %q, want %q", s, "s/old/new/g")
	}
}

func TestStringDeleteRange(t *testing.T) {
	cmd, _ := Parse("3,5d")
	want := "3,5d"
	if s := cmd.String(); s != want {
		t.Fatalf("String() = %q, want %q", s, want)
	}
}

func TestStringSet(t *testing.T) {
	cmd, _ := Parse("set nu")
	if s := cmd.String(); s != "set nu" {
		t.Fatalf("String() = %q, want %q", s, "set nu")
	}
}

func TestStringSetWithValue(t *testing.T) {
	cmd, _ := Parse("set tabwidth=4")
	if s := cmd.String(); s != "set tabwidth=4" {
		t.Fatalf("String() = %q, want %q", s, "set tabwidth=4")
	}
}

func TestStringGoto(t *testing.T) {
	cmd, _ := Parse("42")
	want := "42"
	if s := cmd.String(); s != want {
		t.Fatalf("String() = %q, want %q", s, want)
	}
}

func TestStringUninstall(t *testing.T) {
	cmd, _ := Parse("uninstall")
	if s := cmd.String(); s != "uninstall" {
		t.Fatalf("String() = %q, want %q", s, "uninstall")
	}
}

// ────────────────────────── Executor tests ──────────────────────────

func TestExecutorDelete(t *testing.T) {
	ex := &Executor{
		DeleteRange: func(start, end int) (string, error) {
			return "content", nil
		},
		GetLineText: func(line int) string {
			if line >= 5 {
				return ""
			}
			return "a line"
		},
	}
	cmd, _ := Parse("d")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
}

func TestExecutorDeleteRange(t *testing.T) {
	var deletedStart, deletedEnd int
	ex := &Executor{
		DeleteRange: func(start, end int) (string, error) {
			deletedStart, deletedEnd = start, end
			return "content", nil
		},
	}
	cmd, _ := Parse("2,5d")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if deletedStart != 1 {
		t.Fatalf("deletedStart = %d, want 1", deletedStart)
	}
	if deletedEnd != 4 {
		t.Fatalf("deletedEnd = %d, want 4", deletedEnd)
	}
}

func TestExecutorDeleteNilCallback(t *testing.T) {
	ex := &Executor{} // DeleteRange is nil
	cmd, _ := Parse("d")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil DeleteRange")
	}
}

func TestExecutorWrite(t *testing.T) {
	var writtenPath string
	ex := &Executor{
		WriteFile: func(path string) error {
			writtenPath = path
			return nil
		},
	}
	cmd, _ := Parse("w /tmp/test.txt")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if writtenPath != "/tmp/test.txt" {
		t.Fatalf("writtenPath = %q, want %q", writtenPath, "/tmp/test.txt")
	}
}

func TestExecutorWriteNilCallback(t *testing.T) {
	ex := &Executor{}
	cmd, _ := Parse("w")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil WriteFile")
	}
}

func TestExecutorQuit(t *testing.T) {
	var forceQuit bool
	ex := &Executor{
		Quit: func(force bool) error {
			forceQuit = force
			return nil
		},
	}
	cmd, _ := Parse("q")
	shouldQuit, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !shouldQuit {
		t.Fatal("shouldQuit = false, want true")
	}
	if forceQuit {
		t.Fatal("force should be false for 'q'")
	}
}

func TestExecutorForceQuit(t *testing.T) {
	var forceQuit bool
	ex := &Executor{
		Quit: func(force bool) error {
			forceQuit = force
			return nil
		},
	}
	cmd, _ := Parse("q!")
	shouldQuit, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !shouldQuit {
		t.Fatal("shouldQuit = false, want true")
	}
	if !forceQuit {
		t.Fatal("force should be true for 'q!'")
	}
}

func TestExecutorQuitNilCallback(t *testing.T) {
	ex := &Executor{}
	cmd, _ := Parse("q")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil Quit")
	}
}

func TestExecutorWriteQuit(t *testing.T) {
	var written bool
	var quit bool
	ex := &Executor{
		WriteFile: func(path string) error {
			written = true
			return nil
		},
		Quit: func(force bool) error {
			quit = true
			return nil
		},
	}
	cmd, _ := Parse("wq")
	shouldQuit, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !shouldQuit {
		t.Fatal("shouldQuit = false, want true")
	}
	if !written {
		t.Fatal("WriteFile was not called")
	}
	if !quit {
		t.Fatal("Quit was not called")
	}
}

func TestExecutorEdit(t *testing.T) {
	var openedPath string
	ex := &Executor{
		OpenFile: func(path string) error {
			openedPath = path
			return nil
		},
	}
	cmd, _ := Parse("e main.go")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if openedPath != "main.go" {
		t.Fatalf("openedPath = %q, want %q", openedPath, "main.go")
	}
}

func TestExecutorEditNoPath(t *testing.T) {
	ex := &Executor{
		OpenFile: func(path string) error {
			return nil
		},
	}
	cmd, _ := Parse("e")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for 'e' without filename")
	}
}

func TestExecutorEditNilCallback(t *testing.T) {
	ex := &Executor{}
	cmd, _ := Parse("e main.go")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil OpenFile")
	}
}

func TestExecutorGoto(t *testing.T) {
	var gotoLine int
	ex := &Executor{
		GotoLine: func(line int) error {
			gotoLine = line
			return nil
		},
	}
	cmd, _ := Parse("10")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if gotoLine != 9 {
		t.Fatalf("gotoLine = %d, want 9", gotoLine)
	}
}

func TestExecutorGotoNilCallback(t *testing.T) {
	ex := &Executor{}
	cmd, _ := Parse("10")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil GotoLine")
	}
}

func TestExecutorSet(t *testing.T) {
	var setName string
	ex := &Executor{
		SetOption: func(name, value string) error {
			setName = name
			return nil
		},
	}
	cmd, _ := Parse("set nu")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if setName != "nu" {
		t.Fatalf("setName = %q, want %q", setName, "nu")
	}
}

func TestExecutorSetNilCallback(t *testing.T) {
	ex := &Executor{}
	cmd, _ := Parse("set nu")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error for nil SetOption")
	}
}

func TestExecutorUninstall(t *testing.T) {
	var uninstalled bool
	ex := &Executor{
		Uninstall: func() error {
			uninstalled = true
			return nil
		},
	}
	cmd, _ := Parse("uninstall")
	shouldQuit, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !shouldQuit {
		t.Fatal("shouldQuit = false, want true (uninstall quits)")
	}
	if !uninstalled {
		t.Fatal("Uninstall was not called")
	}
}

func TestExecutorDeleteError(t *testing.T) {
	ex := &Executor{
		DeleteRange: func(start, end int) (string, error) {
			return "", errors.New("disk full")
		},
	}
	cmd, _ := Parse("d")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error from DeleteRange")
	}
}

func TestExecutorWriteError(t *testing.T) {
	ex := &Executor{
		WriteFile: func(path string) error {
			return errors.New("permission denied")
		},
	}
	cmd, _ := Parse("w")
	_, err := ex.Execute(cmd)
	if err == nil {
		t.Fatal("expected error from WriteFile")
	}
}

func TestExecutorShowMsg(t *testing.T) {
	var msg string
	ex := &Executor{
		DeleteRange: func(start, end int) (string, error) {
			return "deleted content", nil
		},
		ShowMsg: func(m string) {
			msg = m
		},
	}
	cmd, _ := Parse("d")
	_, err := ex.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if msg == "" {
		t.Fatal("expected ShowMsg to be called")
	}
}

// ────────────────────────── Edge cases ──────────────────────────

func TestParseRangeWithSpaces(t *testing.T) {
	cmd, err := Parse("  %  d")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if cmd.Kind != CmdDelete {
		t.Fatalf("kind = %d, want CmdDelete", cmd.Kind)
	}
	if !cmd.Range.All {
		t.Fatal("Range.All should be true")
	}
}

func TestParseCommaToEnd(t *testing.T) {
	cmd, err := Parse("10,d")
	if err != nil {
		t.Fatalf("Parse('10,d') error: %v", err)
	}
	if cmd.Range.Start != 9 {
		t.Fatalf("Range.Start = %d, want 9", cmd.Range.Start)
	}
	if cmd.Range.End != -2 {
		t.Fatalf("Range.End = %d, want -2 (end marker)", cmd.Range.End)
	}
}

func TestParseSubstituteAlternateDelimiter(t *testing.T) {
	cmd, err := Parse("s|old|new|")
	if err != nil {
		t.Fatalf("Parse('s|old|new|') error: %v", err)
	}
	if cmd.Kind != CmdSubstitute {
		t.Fatalf("kind = %d, want CmdSubstitute", cmd.Kind)
	}
	if cmd.SubOld != "old" {
		t.Fatalf("SubOld = %q, want %q", cmd.SubOld, "old")
	}
	if cmd.SubNew != "new" {
		t.Fatalf("SubNew = %q, want %q", cmd.SubNew, "new")
	}
}

func TestParseSubstituteMissingDelimiter(t *testing.T) {
	_, err := Parse("s/")
	if err == nil {
		t.Fatal("expected error for incomplete substitute")
	}
}

func TestParseSetStartsWithS(t *testing.T) {
	// "set" starts with 's', should parse as CmdSet not CmdSubstitute
	cmd, err := Parse("set")
	if err != nil {
		t.Fatalf("Parse('set') error: %v", err)
	}
	if cmd.Kind != CmdSet {
		t.Fatalf("kind = %d, want CmdSet", cmd.Kind)
	}
}
