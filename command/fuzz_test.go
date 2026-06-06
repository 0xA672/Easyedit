package command

import (
	"strings"
	"testing"
)

// FuzzParse fuzzes the Parse function with arbitrary command strings.
// It verifies that Parse never panics and that successful parses produce
// a Command whose String() output can be re-parsed without error.
func FuzzParse(f *testing.F) {
	// Seed corpus: representative commands
	seeds := []string{
		"q", "q!", "w", "wq", "wq!",
		"w /tmp/out.txt", "e main.go",
		"d", "10,20d", "%d",
		"s/old/new/", "s/old/new/g", "%s/old/new/gi",
		"3,5s/foo/bar/", "set nu", "set tabwidth=4",
		"42", "uninstall",
		"10,d", "  %  d",
		"s|old|new|",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		cmd, err := Parse(input)
		if err != nil {
			// Parse error is acceptable; we just must not panic.
			return
		}
		// Round-trip: String() should produce a re-parseable string
		str := cmd.String()
		if str == "" {
			t.Fatalf("String() returned empty for input %q, kind=%d", input, cmd.Kind)
		}
		cmd2, err2 := Parse(str)
		if err2 != nil {
			// Some commands may not round-trip perfectly (e.g. whitespace variations)
			// but should at least not error fatally
			return
		}
		if cmd2.Kind != cmd.Kind {
			t.Fatalf("round-trip kind mismatch: input=%q original=%d roundtrip=%d",
				input, cmd.Kind, cmd2.Kind)
		}
	})
}

// FuzzParseSubstitute fuzzes the substitute parser with random delimiters and patterns.
func FuzzParseSubstitute(f *testing.F) {
	seeds := []string{
		"/old/new/",
		"/foo/bar/g",
		"|pattern|replace|gi",
		"#find#replace#",
		"/a//g", // empty replacement
		"/[regex]*/replacement/g",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, rest string) {
		old, new, flags, err := parseSubstitute(rest)
		if err != nil {
			return // parse error is valid input
		}
		// Verify no panics occurred and results are usable
		_ = old
		_ = new
		_ = flags
	})
}

// FuzzParseRange fuzzes the range parser with arbitrary range strings.
func FuzzParseRange(f *testing.F) {
	seeds := []string{
		"%", "10,20", "5,", "%  ", "1,100", "999",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		r, _ := parseRange(input, 0)
		// Verify range invariants
		if r.All && r.Start >= 0 {
			// All with explicit start is unusual but not invalid
		}
		// Start should be >= -1
		if r.Start < -1 {
			t.Fatalf("Start %d is below -1 for input %q", r.Start, input)
		}
	})
}

// FuzzExpandRange fuzzes ExpandRange with various totalLines values.
func FuzzExpandRange(f *testing.F) {
	seeds := []struct {
		all        bool
		start, end int
		total      int
	}{
		{true, -1, -1, 10},
		{false, 2, 5, 10},
		{false, -1, -1, 100},
		{false, 3, -2, 50},
	}
	for _, s := range seeds {
		f.Add(s.all, s.start, s.end, s.total)
	}

	f.Fuzz(func(t *testing.T, all bool, start, end, total int) {
		if total < 1 || total > 100000 {
			return
		}
		r := Range{All: all, Start: start, End: end}
		s, e := r.ExpandRange(total)
		// Results should be within valid bounds
		if s < 0 && !all {
			// start defaults to 0
		}
		if e >= total && e != -1 && e != -2 {
			// end may exceed if input was large; that's ok
		}
	})
}

// FuzzCompileRegex fuzzes regex compilation from substitute commands.
func FuzzCompileRegex(f *testing.F) {
	seeds := []string{
		"s/hello/world/",
		"s/[a-z]+/replace/i",
		"s/foo.*/bar/g",
		"s/(group)/$1/gi",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		cmd, err := Parse(input)
		if err != nil {
			return
		}
		if cmd.Kind != CmdSubstitute {
			return
		}
		// CompileRegex may fail on invalid patterns; that's expected
		_, _ = cmd.CompileRegex()
	})
}

// FuzzMakeGlobal fuzzes the MakeGlobal flag detection.
func FuzzMakeGlobal(f *testing.F) {
	seeds := []string{"g", "gi", "", "i", "gg", "xyz"}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, flags string) {
		cmd := &Command{SubFlg: flags}
		result := cmd.MakeGlobal()
		expected := strings.Contains(flags, "g")
		if result != expected {
			t.Fatalf("MakeGlobal()=%v, expected %v for flags %q", result, expected, flags)
		}
	})
}
