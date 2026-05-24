package ui

import (
	"testing"
)

// ────────────────────────── autoPair tests ──────────────────────────

func TestAutoPairOpenParen(t *testing.T) {
	close, ok := autoPair('(')
	if !ok {
		t.Fatal("autoPair('(') should return ok=true")
	}
	if close != ')' {
		t.Fatalf("autoPair('(') = %c, want ')'", close)
	}
}

func TestAutoPairOpenSquare(t *testing.T) {
	close, ok := autoPair('[')
	if !ok {
		t.Fatal("autoPair('[') should return ok=true")
	}
	if close != ']' {
		t.Fatalf("autoPair('[') = %c, want ']'", close)
	}
}

func TestAutoPairOpenCurly(t *testing.T) {
	close, ok := autoPair('{')
	if !ok {
		t.Fatal("autoPair('{') should return ok=true")
	}
	if close != '}' {
		t.Fatalf("autoPair('{') = %c, want '}'", close)
	}
}

func TestAutoPairSingleQuote(t *testing.T) {
	close, ok := autoPair('\'')
	if !ok {
		t.Fatal("autoPair('\\'') should return ok=true")
	}
	if close != '\'' {
		t.Fatalf("autoPair('\\'') = %c, want '\\''", close)
	}
}

func TestAutoPairDoubleQuote(t *testing.T) {
	close, ok := autoPair('"')
	if !ok {
		t.Fatal("autoPair('\"') should return ok=true")
	}
	if close != '"' {
		t.Fatalf("autoPair('\"') = %c, want '\"'", close)
	}
}

func TestAutoPairBacktick(t *testing.T) {
	close, ok := autoPair('`')
	if !ok {
		t.Fatal("autoPair('`') should return ok=true")
	}
	if close != '`' {
		t.Fatalf("autoPair('`') = %c, want '`'", close)
	}
}

func TestAutoPairNonPair(t *testing.T) {
	pairs := []rune{'a', '1', ' ', '\n', '>', '<', '|', '/'}
	for _, ch := range pairs {
		_, ok := autoPair(ch)
		if ok {
			t.Fatalf("autoPair('%c') should return ok=false", ch)
		}
	}
}

func TestAutoPairReturnsZeroWhenFalse(t *testing.T) {
	close, ok := autoPair('x')
	if ok {
		t.Fatal("expected ok=false")
	}
	if close != 0 {
		t.Fatalf("expected 0, got %c", close)
	}
}
