package document

import (
	"strings"
	"testing"
)

// FuzzGapBufferOps fuzzes GapBuffer Insert/Delete/Slice/RuneAt with random
// positions and text. It verifies core invariants: Len >= 0, Content matches
// an independent string reconstruction.
func FuzzGapBufferOps(f *testing.F) {
	type op struct {
		kind int // 0=insert, 1=delete, 2=slice
		pos  int
		text string
		len  int
	}
	seeds := []string{
		"hello",
		"abc\ndef\n",
		"你好世界",
		"",
		"a\nb\nc\nd",
		"\t\ttabbed",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, text string) {
		runes := []rune(text)
		expected := string(runes) // rune round-trip may change invalid UTF-8
		gb := NewGapBufferFromRunes(runes)

		// Verify initial state
		if gb.Len() != len(runes) {
			t.Fatalf("initial Len()=%d, want %d", gb.Len(), len(runes))
		}
		if gb.Content() != expected {
			t.Fatalf("initial Content() mismatch")
		}

		// Test RuneAt at all valid positions
		for i := 0; i < gb.Len(); i++ {
			ch := gb.RuneAt(i)
			if ch != runes[i] {
				t.Fatalf("RuneAt(%d)=%c, want %c", i, ch, runes[i])
			}
		}

		// Test out-of-range RuneAt
		if gb.RuneAt(-1) != 0 {
			t.Fatal("RuneAt(-1) should return 0")
		}
		if gb.RuneAt(gb.Len()) != 0 {
			t.Fatal("RuneAt(Len) should return 0")
		}

		// Test Slice full
		fullSlice := gb.Slice(0, gb.Len())
		if string(fullSlice) != expected {
			t.Fatalf("Slice(0,Len) mismatch")
		}

		// Test insert at start
		gb.Insert(0, []rune("X"))
		if !strings.HasPrefix(gb.Content(), "X") {
			t.Fatal("Insert at 0 failed")
		}

		// Test delete
		gb.Delete(0, 1)
		if gb.Content() != expected {
			t.Fatal("Delete after insert failed to restore")
		}
	})
}

// FuzzPosLineColRoundtrip fuzzes PosToLineCol -> LineColToPos round-trip
// on random content. Every position should survive the round-trip.
func FuzzPosLineColRoundtrip(f *testing.F) {
	seeds := []string{
		"", "a", "ab", "a\nb", "abc\ndef\nghi",
		"你好\n世界", "\n\n\n", "no-trailing",
		"trailing\n", "multi\nline\ntext\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, text string) {
		gb := NewGapBufferFromRunes([]rune(text))
		total := gb.Len()

		// Test round-trip for all positions
		for pos := 0; pos <= total; pos++ {
			line, col := gb.PosToLineCol(pos)
			back := gb.LineColToPos(line, col)
			if back != pos {
				t.Fatalf("round-trip failed: pos=%d -> (%d,%d) -> %d, text=%q",
					pos, line, col, back, text)
			}
		}

		// Test LineCount consistency
		lc := gb.LineCount()
		if lc < 1 {
			t.Fatalf("LineCount()=%d, should be >= 1", lc)
		}

		// Verify last line is accessible
		lastLine := lc - 1
		_ = gb.Line(lastLine)
		_ = gb.LineLen(lastLine)
	})
}

// FuzzDetectEncoding fuzzes encoding detection with arbitrary byte slices.
func FuzzDetectEncoding(f *testing.F) {
	seeds := [][]byte{
		{},
		{0xEF, 0xBB, 0xBF, 'h', 'e', 'l', 'l', 'o'}, // UTF-8 BOM
		{0xFF, 0xFE, 'h', 0, 'i', 0},                  // UTF-16LE BOM
		{0xFE, 0xFF, 0, 'h', 0, 'i'},                  // UTF-16BE BOM
		[]byte("hello world"),                          // plain ASCII
		[]byte("你好世界"),                               // UTF-8 CJK
		{0x80, 0x81, 0x82},                             // invalid UTF-8
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		enc := DetectEncoding(data)
		// Encoding should always be a known value
		if enc < EncUnknown || enc > EncLatin1 {
			t.Fatalf("DetectEncoding returned unknown value %d", enc)
		}

		// DecodeToUTF8 should never panic
		str, decEnc, err := DecodeToUTF8(data)
		if err != nil {
			return // errors are acceptable for malformed data
		}
		_ = str
		_ = decEnc
	})
}

// FuzzGapBufferInsertDelete fuzzes interleaved insert and delete operations.
func FuzzGapBufferInsertDelete(f *testing.F) {
	f.Add("hello world", 5, "X", 0, 1)
	f.Add("", 0, "abc", 1, 1)
	f.Add("abcdef", 3, "123", 2, 2)

	f.Fuzz(func(t *testing.T, text string, insPos int, insText string, delPos int, delLen int) {
		runes := []rune(text)
		gb := NewGapBufferFromRunes(runes)

		// Clamp insert position
		if insPos < 0 || insPos > gb.Len() {
			insPos = gb.Len() / 2
			if insPos < 0 {
				insPos = 0
			}
		}
		insRunes := []rune(insText)
		if len(insRunes) == 0 {
			insRunes = []rune{'x'}
		}

		// Insert
		gb.Insert(insPos, insRunes)
		newLen := len(runes) + len(insRunes)
		if gb.Len() != newLen {
			t.Fatalf("after insert: Len()=%d, want %d", gb.Len(), newLen)
		}

		// Clamp delete
		if delPos < 0 || delPos >= gb.Len() {
			delPos = 0
		}
		if delLen <= 0 {
			delLen = 1
		}
		if delPos+delLen > gb.Len() {
			delLen = gb.Len() - delPos
		}
		if delLen <= 0 {
			return
		}

		// Delete
		gb.Delete(delPos, delLen)
		wantLen := newLen - delLen
		if gb.Len() != wantLen {
			t.Fatalf("after delete: Len()=%d, want %d", gb.Len(), wantLen)
		}

		// Content should not panic
		_ = gb.Content()
	})
}

// FuzzDocumentOps fuzzes Document wrapper methods with random strings.
func FuzzDocumentOps(f *testing.F) {
	seeds := []string{
		"", "hello", "a\nb\nc", "你好\n世界",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, text string) {
		d := NewDocumentFromString(text)

		// LineCount >= 1
		if d.LineCount() < 1 {
			t.Fatalf("LineCount()=%d", d.LineCount())
		}

		// Access each line
		for i := 0; i < d.LineCount(); i++ {
			line := d.Line(i)
			_ = line
			ll := d.LineLen(i)
			if ll < 0 {
				t.Fatalf("LineLen(%d)=%d", i, ll)
			}
			indent := d.GetIndent(i)
			_ = indent
		}

		// Test insert/delete
		if d.Len() > 0 {
			ch := d.DeleteRuneAt(0)
			if ch == 0 {
				t.Fatal("DeleteRuneAt(0) returned 0 on non-empty doc")
			}
			d.InsertRune(0, ch) // restore
		}

		d.InsertText(0, "X")
		if !strings.HasPrefix(d.Content(), "X") {
			t.Fatal("InsertText at 0 failed")
		}
	})
}

// FuzzStripBOM fuzzes the StripBOM function with random bytes.
func FuzzStripBOM(f *testing.F) {
	seeds := [][]byte{
		{},
		{0xEF, 0xBB, 0xBF},
		{0xEF, 0xBB, 0xBF, 'a'},
		{0xFF, 0xFE, 'a', 0},
		{0xFE, 0xFF, 0, 'a'},
		{'n', 'o', 'b', 'o', 'm'},
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		result := StripBOM(data)
		// Result should never be longer than input
		if len(result) > len(data) {
			t.Fatalf("StripBOM result longer than input")
		}
		// Result should never be nil if input is non-nil
		if data != nil && result == nil {
			t.Fatal("StripBOM returned nil for non-nil input")
		}
	})
}
