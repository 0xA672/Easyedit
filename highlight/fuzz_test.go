package highlight

import (
	"testing"

	"github.com/alecthomas/chroma/v2"
)

// FuzzBracketMatch fuzzes IsBracket and MatchingBracket with arbitrary runes.
func FuzzBracketMatch(f *testing.F) {
	seeds := []rune{
		'(', ')', '[', ']', '{', '}',
		'<', '>', 'a', '1', ' ', '\n', '\t',
		0, 0xFFFF, '你',
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, ch rune) {
		isBr := IsBracket(ch)
		match := MatchingBracket(ch)

		if isBr {
			// Matching bracket of a bracket should be a different bracket
			if match == 0 {
				t.Fatalf("IsBracket(%c)=true but MatchingBracket returns 0", ch)
			}
			if match == ch {
				t.Fatalf("MatchingBracket(%c)=%c, should differ", ch, match)
			}
			// Double matching should return original
			back := MatchingBracket(match)
			if back != ch {
				t.Fatalf("MatchingBracket(MatchingBracket(%c))=%c, want %c", ch, back, ch)
			}
			// Match should also be a bracket
			if !IsBracket(match) {
				t.Fatalf("MatchingBracket(%c)=%c is not a bracket", ch, match)
			}
		} else {
			// Non-brackets should return 0
			if match != 0 {
				t.Fatalf("non-bracket %c returned match %c", ch, match)
			}
		}
	})
}

// FuzzTokenize fuzzes the Tokenize method with random text and file extensions.
func FuzzTokenize(f *testing.F) {
	type seed struct {
		path string
		text string
	}
	seeds := []seed{
		{"main.go", "package main\nfunc main() {}"},
		{"test.py", "def foo():\n    pass"},
		{"app.js", "const x = () => {}"},
		{"style.css", "body { color: red; }"},
		{"data.json", `{"key": "value"}`},
		{"page.html", "<html><body></body></html>"},
		{"script.sh", "#!/bin/bash\necho hello"},
		{"unknown.xyz", "some random text"},
		{"", "no extension"},
		{"main.go", ""}, // empty text
		{"main.go", "你好世界"},
	}
	for _, s := range seeds {
		f.Add(s.path, s.text)
	}

	f.Fuzz(func(t *testing.T, path, text string) {
		h := NewHighlighter(path)
		// Tokenize should never panic
		tokens := h.Tokenize(text)
		_ = tokens

		// Tokenize again (cache hit path)
		tokens2 := h.Tokenize(text)
		_ = tokens2

		// SetFile with random path
		h.SetFile(path)
		tokens3 := h.Tokenize(text)
		_ = tokens3
	})
}

// FuzzGetTokenStyle fuzzes GetTokenStyle with random token types.
func FuzzGetTokenStyle(f *testing.F) {
	seeds := []int{0, 1, 5, 10, 100, 255, 1000}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, typeVal int) {
		h := NewHighlighter("test.go")
		// Should not panic for any token type
		style := h.GetTokenStyle(chromaTypeFromInt(typeVal))
		_ = style
	})
}

// chromaTypeFromInt maps an integer to a plausible chroma token type.
// We use the actual chroma type constants indirectly.
func chromaTypeFromInt(n int) chroma.TokenType {
	// Map to a range of valid chroma token types
	types := []chroma.TokenType{
		chroma.Text,
		chroma.Error,
		chroma.Comment,
		chroma.CommentSingle,
		chroma.CommentMultiline,
		chroma.Keyword,
		chroma.KeywordDeclaration,
		chroma.Name,
		chroma.NameFunction,
		chroma.NameClass,
		chroma.LiteralString,
		chroma.LiteralNumber,
		chroma.Operator,
		chroma.Punctuation,
	}
	if n < 0 {
		n = -n
	}
	return types[n%len(types)]
}

