package highlight

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// ────────────────────────── IsBracket tests ──────────────────────────

func TestIsBracketOpenParen(t *testing.T) {
	if !IsBracket('(') {
		t.Fatal("'(' should be a bracket")
	}
}

func TestIsBracketCloseParen(t *testing.T) {
	if !IsBracket(')') {
		t.Fatal("')' should be a bracket")
	}
}

func TestIsBracketOpenSquare(t *testing.T) {
	if !IsBracket('[') {
		t.Fatal("'[' should be a bracket")
	}
}

func TestIsBracketCloseSquare(t *testing.T) {
	if !IsBracket(']') {
		t.Fatal("']' should be a bracket")
	}
}

func TestIsBracketOpenCurly(t *testing.T) {
	if !IsBracket('{') {
		t.Fatal("'{' should be a bracket")
	}
}

func TestIsBracketCloseCurly(t *testing.T) {
	if !IsBracket('}') {
		t.Fatal("'}' should be a bracket")
	}
}

func TestIsBracketNonBracket(t *testing.T) {
	if IsBracket('a') {
		t.Fatal("'a' should not be a bracket")
	}
	if IsBracket('1') {
		t.Fatal("'1' should not be a bracket")
	}
	if IsBracket('<') {
		t.Fatal("'<' should not be a bracket (not in our bracket set)")
	}
	if IsBracket(' ') {
		t.Fatal("space should not be a bracket")
	}
}

// ────────────────────────── MatchingBracket tests ──────────────────────────

func TestMatchingBracketParens(t *testing.T) {
	if m := MatchingBracket('('); m != ')' {
		t.Fatalf("MatchingBracket('(') = %c, want ')'", m)
	}
	if m := MatchingBracket(')'); m != '(' {
		t.Fatalf("MatchingBracket(')') = %c, want '('", m)
	}
}

func TestMatchingBracketSquare(t *testing.T) {
	if m := MatchingBracket('['); m != ']' {
		t.Fatalf("MatchingBracket('[') = %c, want ']'", m)
	}
	if m := MatchingBracket(']'); m != '[' {
		t.Fatalf("MatchingBracket(']') = %c, want '['", m)
	}
}

func TestMatchingBracketCurly(t *testing.T) {
	if m := MatchingBracket('{'); m != '}' {
		t.Fatalf("MatchingBracket('{') = %c, want '}'", m)
	}
	if m := MatchingBracket('}'); m != '{' {
		t.Fatalf("MatchingBracket('}') = %c, want '{'", m)
	}
}

func TestMatchingBracketNonBracket(t *testing.T) {
	if m := MatchingBracket('x'); m != 0 {
		t.Fatalf("MatchingBracket('x') = %c, want 0", m)
	}
	if m := MatchingBracket('\n'); m != 0 {
		t.Fatalf("MatchingBracket('\\n') = %d, want 0", m)
	}
}

// ────────────────────────── hexToTcell tests ──────────────────────────

func TestHexToTcellValid(t *testing.T) {
	color := hexToTcell("#ff8800")
	if color == tcell.ColorDefault {
		t.Fatal("expected a valid color, got ColorDefault")
	}
	// Check the color string representation
	_, _, _ = color.RGB() // ensure it doesn't panic
}

func TestHexToTcellColorName(t *testing.T) {
	color := hexToTcell("red")
	if color == tcell.ColorDefault {
		t.Fatal("expected a valid color, got ColorDefault")
	}
}

func TestHexToTcellEmpty(t *testing.T) {
	color := hexToTcell("")
	_ = color // should not panic
}

// ────────────────────────── Highlighter tests ──────────────────────────

func TestNewHighlighterGo(t *testing.T) {
	h := NewHighlighter("main.go")
	if h == nil {
		t.Fatal("NewHighlighter returned nil")
	}
	if h.lexer == nil {
		t.Fatal("lexer should be non-nil for .go file")
	}
	if h.style == nil {
		t.Fatal("style should be non-nil")
	}
}

func TestNewHighlighterUnknownExt(t *testing.T) {
	h := NewHighlighter("file.xyzabc")
	if h == nil {
		t.Fatal("NewHighlighter returned nil")
	}
	// Unknown extension should get nil lexer (plain text)
	_ = h
}

func TestHighlighterTokenizeGo(t *testing.T) {
	h := NewHighlighter("main.go")
	tokens := h.Tokenize(`package main

import "fmt"

func main() {
	fmt.Println("hello")
}`)
	if len(tokens) == 0 {
		t.Fatal("expected non-empty tokens for Go code")
	}
	// First token should be keyword or something meaningful
	foundPkg := false
	for _, tok := range tokens {
		if tok.Value == "package" {
			foundPkg = true
			break
		}
	}
	if !foundPkg {
		t.Fatal("expected 'package' keyword in tokenized Go code")
	}
}

func TestHighlighterTokenizePlainText(t *testing.T) {
	h := NewHighlighter("file.xyzabc")
	tokens := h.Tokenize("just plain text")
	// Unknown lexer should return nil tokens
	if tokens != nil {
		t.Fatalf("expected nil tokens for plain text, got %d tokens", len(tokens))
	}
}

func TestHighlighterTokenizeCache(t *testing.T) {
	h := NewHighlighter("main.go")
	text := `package main`
	tokens1 := h.Tokenize(text)
	tokens2 := h.Tokenize(text) // Should use cache
	if len(tokens1) != len(tokens2) {
		t.Fatalf("cache returned different token count: %d vs %d", len(tokens1), len(tokens2))
	}
}

func TestHighlighterTokenizeInvalidSyntax(t *testing.T) {
	h := NewHighlighter("test.py")
	tokens := h.Tokenize("   \n\n\n") // minimal valid-ish content
	_ = tokens                        // should not panic
}

func TestHighlighterSetFile(t *testing.T) {
	h := NewHighlighter("main.go")
	h.Tokenize("package main")
	h.SetFile("test.py")
	if h.lexer == nil {
		t.Fatal("lexer should be non-nil for .py file")
	}
	// Cache should be cleared
	if h.tokens != nil {
		t.Fatal("tokens should be nil after SetFile")
	}
}

func TestHighlighterSetFileUnknownExt(t *testing.T) {
	h := NewHighlighter("main.go")
	h.SetFile("noext")
	if h.lexer != nil {
		t.Fatal("lexer should be nil for unknown extension")
	}
}

func TestGetTokenStyleNoStyle(t *testing.T) {
	h := &Highlighter{} // no style set
	ts := h.GetTokenStyle(0)
	if ts.Fg != tcell.ColorDefault {
		t.Fatalf("Fg = %v, want ColorDefault", ts.Fg)
	}
	if ts.Bg != tcell.ColorDefault {
		t.Fatalf("Bg = %v, want ColorDefault", ts.Bg)
	}
}

func TestGetTokenStyleMonokai(t *testing.T) {
	h := NewHighlighter("main.go")
	ts := h.GetTokenStyle(0) // Keyword type tokens
	// Should have some style from monokai theme
	_ = ts
}

// ────────────────────────── Highlighter nil lexer tokenize ──────────────────────────

func TestHighlighterNilLexerTokenize(t *testing.T) {
	h := &Highlighter{lexer: nil}
	tokens := h.Tokenize("some text")
	if tokens != nil {
		t.Fatal("expected nil tokens when lexer is nil")
	}
}
