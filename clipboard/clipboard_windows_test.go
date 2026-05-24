//go:build windows

package clipboard

import (
	"testing"
	"unicode/utf16"
)

func TestUTF16FromString_ASCII(t *testing.T) {
	input := "hello"
	got := UTF16FromString(input)
	expected := utf16.Encode([]rune(input))
	if len(got) != len(expected) {
		t.Fatalf("len = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], expected[i])
		}
	}
}

func TestUTF16FromString_Chinese(t *testing.T) {
	input := "你好世界"
	got := UTF16FromString(input)
	expected := utf16.Encode([]rune(input))
	if len(got) != len(expected) {
		t.Fatalf("len = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], expected[i])
		}
	}
}

func TestUTF16FromString_Empty(t *testing.T) {
	got := UTF16FromString("")
	if got != nil {
		t.Fatalf("expected nil for empty string, got %v", got)
	}
}

func TestUTF16FromString_ContainsNull(t *testing.T) {
	// Input with embedded null — function should truncate at null
	input := "hello\x00world"
	got := UTF16FromString(input)
	expected := utf16.Encode([]rune("hello"))
	if len(got) != len(expected) {
		t.Fatalf("len = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], expected[i])
		}
	}
}

func TestUTF16FromString_Emoji(t *testing.T) {
	input := "👋🌍" // surrogate pair
	got := UTF16FromString(input)
	expected := utf16.Encode([]rune(input))
	if len(got) != len(expected) {
		t.Fatalf("len = %d, want %d", len(got), len(expected))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], expected[i])
		}
	}
}

func TestMultiByteToWideChar_ASCII(t *testing.T) {
	input := "hello"
	n := MultiByteToWideChar(CP_UTF8, 0, input, len(input), nil, 0)
	if n <= 0 {
		t.Fatalf("expected >0, got %d", n)
	}
	buf := make([]uint16, n)
	n2 := MultiByteToWideChar(CP_UTF8, 0, input, len(input), &buf[0], n)
	if n2 != n {
		t.Fatalf("second call returned %d, want %d", n2, n)
	}
	if string(utf16.Decode(buf)) != input {
		t.Fatalf("roundtrip = %q, want %q", string(utf16.Decode(buf)), input)
	}
}

func TestMultiByteToWideChar_ZeroLength(t *testing.T) {
	// Passing cbMultiByte=0 should produce zero-length result
	n := MultiByteToWideChar(CP_UTF8, 0, "hello", 0, nil, 0)
	if n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestWideCharToMultiByte_ASCII(t *testing.T) {
	runes := []rune("hello")
	u16 := utf16.Encode(runes)
	n := WideCharToMultiByte(CP_UTF8, 0, &u16[0], len(u16), nil, 0, nil, nil)
	if n <= 0 {
		t.Fatalf("expected >0, got %d", n)
	}
	buf := make([]byte, n)
	n2 := WideCharToMultiByte(CP_UTF8, 0, &u16[0], len(u16), &buf[0], n, nil, nil)
	if n2 != n {
		t.Fatalf("second call returned %d, want %d", n2, n)
	}
	if string(buf) != "hello" {
		t.Fatalf("roundtrip = %q, want %q", string(buf), "hello")
	}
}
