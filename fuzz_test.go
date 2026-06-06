package main

import (
	"testing"
)

// FuzzApplyRules fuzzes applyRules with random content and rule strings.
func FuzzApplyRules(f *testing.F) {
	type seed struct {
		content string
		rule    string
	}
	seeds := []seed{
		{"hello world", "s/hello/goodbye/g"},
		{"foo bar foo", "s/foo/baz/"},
		{"aaa", "s/a/b/g"},
		{"hello\nworld\nfoo", "/world/d"},
		{"test", "s/t/T/"},
		{"nothing", "x/invalid"},
		{"", "s/a/b/g"},
		{"abc", "s/a//g"},
		{"hello", "s/hello/world/gi"},
		{"line1\nline2\nline3", "/line2/d"},
	}
	for _, s := range seeds {
		f.Add(s.content, s.rule)
	}

	f.Fuzz(func(t *testing.T, content, rule string) {
		rules := []string{rule}
		// applyRules should never panic
		result := applyRules(content, rules)
		_ = result
	})
}

// FuzzApplyRulesMultiple fuzzes applyRules with multiple random rules.
func FuzzApplyRulesMultiple(f *testing.F) {
	seeds := []string{
		"s/a/b/g;s/c/d/g",
		"/pattern/d;s/old/new/",
		"s/x/y/g",
		"",
		";;;",
		"s/a/b/g\n/pattern/d",
	}
	for _, s := range seeds {
		f.Add(s, "test content here")
	}

	f.Fuzz(func(t *testing.T, script, content string) {
		rules := parseScript(script)
		// applyRules should never panic
		result := applyRules(content, rules)
		_ = result
	})
}

// FuzzParseScript fuzzes parseScript with arbitrary script strings.
func FuzzParseScript(f *testing.F) {
	seeds := []string{
		"s/old/new/g",
		"s/a/b/g;s/c/d/",
		"/pattern/d",
		"",
		"   ",
		";;;",
		"\n\n\n",
		"s/a/b/g\n/pattern/d\ns/x/y/",
		"s/a/b/g ; s/c/d/",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, script string) {
		rules := parseScript(script)
		// Rules should be a valid slice (may be empty)
		if rules == nil && script != "" {
			// parseScript returns nil for scripts with only delimiters
			return
		}
		for _, r := range rules {
			if r == "" {
				t.Fatalf("parseScript returned empty rule for input %q", script)
			}
		}
	})
}
