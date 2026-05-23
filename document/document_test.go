package document

import (
	"testing"
)

// ==================== GapBuffer Basics ====================

func TestNewGapBuffer(t *testing.T) {
	gb := NewGapBuffer(100)
	if gb.Len() != 0 {
		t.Fatalf("new buffer length should be 0, got %d", gb.Len())
	}
	if gb.Cap() < 100 {
		t.Fatalf("capacity should be >= 100, got %d", gb.Cap())
	}
	if gb.Content() != "" {
		t.Fatalf("content should be empty, got %q", gb.Content())
	}
}

func TestNewGapBufferMinCapacity(t *testing.T) {
	gb := NewGapBuffer(1)
	if gb.Cap() < 16 {
		t.Fatalf("minimum capacity should be 16, got %d", gb.Cap())
	}
}

func TestNewGapBufferFromRunes(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello"))
	if gb.Len() != 5 {
		t.Fatalf("length should be 5, got %d", gb.Len())
	}
	if gb.Content() != "hello" {
		t.Fatalf("content should be %q, got %q", "hello", gb.Content())
	}
}

func TestInsert(t *testing.T) {
	gb := NewGapBuffer(16)
	gb.Insert(0, []rune("abc"))
	if gb.Content() != "abc" {
		t.Fatalf("expected %q, got %q", "abc", gb.Content())
	}
	if gb.Len() != 3 {
		t.Fatalf("length should be 3, got %d", gb.Len())
	}
}

func TestInsertAtPosition(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("ac"))
	gb.Insert(1, []rune{'b'})
	if gb.Content() != "abc" {
		t.Fatalf("expected %q, got %q", "abc", gb.Content())
	}
}

func TestInsertAtEnd(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("ab"))
	gb.Insert(2, []rune{'c'})
	if gb.Content() != "abc" {
		t.Fatalf("expected %q, got %q", "abc", gb.Content())
	}
}

func TestInsertAtStart(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("bc"))
	gb.Insert(0, []rune{'a'})
	if gb.Content() != "abc" {
		t.Fatalf("expected %q, got %q", "abc", gb.Content())
	}
}

func TestInsertEmpty(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	gb.Insert(1, []rune{})
	if gb.Content() != "abc" {
		t.Fatalf("empty insert should not change content, got %q", gb.Content())
	}
}

func TestInsertGrow(t *testing.T) {
	// Insert more than initial capacity, triggering grow
	gb := NewGapBufferFromRunes([]rune("hello"))
	gb.Insert(5, []rune(" world!!!!!!!!!!!!!!!")) // exceeds 16
	want := "hello world!!!!!!!!!!!!!!!"
	if gb.Content() != want {
		t.Fatalf("expected %q, got %q", want, gb.Content())
	}
}

func TestDelete(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abcdef"))
	gb.Delete(2, 3) // delete "cde"
	if gb.Content() != "abf" {
		t.Fatalf("expected %q, got %q", "abf", gb.Content())
	}
	if gb.Len() != 3 {
		t.Fatalf("length should be 3, got %d", gb.Len())
	}
}

func TestDeleteAtStart(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	gb.Delete(0, 1)
	if gb.Content() != "bc" {
		t.Fatalf("expected %q, got %q", "bc", gb.Content())
	}
}

func TestDeleteAtEnd(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	gb.Delete(2, 1)
	if gb.Content() != "ab" {
		t.Fatalf("expected %q, got %q", "ab", gb.Content())
	}
}

func TestDeleteAll(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	gb.Delete(0, 3)
	if gb.Content() != "" {
		t.Fatalf("expected empty string, got %q", gb.Content())
	}
	if gb.Len() != 0 {
		t.Fatalf("length should be 0, got %d", gb.Len())
	}
}

func TestInsertAfterDelete(t *testing.T) {
	// Gap position changes after delete, then insert
	gb := NewGapBufferFromRunes([]rune("hello"))
	gb.Delete(1, 3) // "ho" (delete e,l,l)
	gb.Insert(1, []rune("i")) // "hio"
	gb.Insert(3, []rune("p")) // insert at end → "hiop"
	if gb.Content() != "hiop" {
		t.Fatalf("expected %q, got %q", "hiop", gb.Content())
	}
}

func TestRuneAt(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	if gb.RuneAt(0) != 'a' {
		t.Fatalf("RuneAt(0) should be 'a', got %c", gb.RuneAt(0))
	}
	if gb.RuneAt(1) != 'b' {
		t.Fatalf("RuneAt(1) should be 'b', got %c", gb.RuneAt(1))
	}
	if gb.RuneAt(2) != 'c' {
		t.Fatalf("RuneAt(2) should be 'c', got %c", gb.RuneAt(2))
	}
}

func TestRuneAtOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	if gb.RuneAt(-1) != 0 {
		t.Fatalf("out of range should return 0")
	}
	if gb.RuneAt(100) != 0 {
		t.Fatalf("out of range should return 0")
	}
}

func TestRuneAtAfterInsert(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("ac"))
	gb.Insert(1, []rune{'b'})
	if gb.RuneAt(1) != 'b' {
		t.Fatalf("RuneAt(1) after insert should be 'b', got %c", gb.RuneAt(1))
	}
}

func TestSlice(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abcdef"))
	got := string(gb.Slice(1, 4))
	if got != "bcd" {
		t.Fatalf("expected %q, got %q", "bcd", got)
	}
}

func TestSliceFull(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello"))
	got := string(gb.Slice(0, gb.Len()))
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestSliceEmpty(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	got := string(gb.Slice(1, 1))
	if got != "" {
		t.Fatalf("empty slice should be empty string, got %q", got)
	}
}

func TestSliceClamp(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	got := string(gb.Slice(-1, 10)) // should clamp
	if got != "abc" {
		t.Fatalf("expected %q, got %q", "abc", got)
	}
}

// ==================== Line Operations ====================

func TestLineCountEmpty(t *testing.T) {
	gb := NewGapBuffer(16)
	if gb.LineCount() != 1 {
		t.Fatalf("empty document line count should be 1, got %d", gb.LineCount())
	}
}

func TestLineCountSingleLine(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello"))
	if gb.LineCount() != 1 {
		t.Fatalf("single line document should have 1 line, got %d", gb.LineCount())
	}
}

func TestLineCountMultiLine(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello\nworld\n!"))
	if gb.LineCount() != 3 {
		t.Fatalf("3-line document should have 3 lines, got %d", gb.LineCount())
	}
}

func TestLineCountTrailingNewline(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello\n"))
	if gb.LineCount() != 2 {
		t.Fatalf("trailing newline should give 2 lines, got %d", gb.LineCount())
	}
}

func TestLineStart(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		line int
		want int
	}{
		{0, 0},
		{1, 4},
		{2, 8},
	}
	for _, tt := range tests {
		got := gb.LineStart(tt.line)
		if got != tt.want {
			t.Errorf("LineStart(%d) = %d, expected %d", tt.line, got, tt.want)
		}
	}
}

func TestLineStartOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	if gb.LineStart(5) != -1 {
		t.Fatalf("out of range should return -1")
	}
	if gb.LineStart(-1) != -1 {
		t.Fatalf("negative line should return -1")
	}
}

func TestLineEnd(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		line int
		want int
	}{
		{0, 3},
		{1, 7},
		{2, 10}, // last line has no newline, returns total length
	}
	for _, tt := range tests {
		got := gb.LineEnd(tt.line)
		if got != tt.want {
			t.Errorf("LineEnd(%d) = %d, expected %d", tt.line, got, tt.want)
		}
	}
}

func TestLine(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		line int
		want string
	}{
		{0, "abc"},
		{1, "def"},
		{2, "gh"},
	}
	for _, tt := range tests {
		got := gb.Line(tt.line)
		if got != tt.want {
			t.Errorf("Line(%d) = %q, expected %q", tt.line, got, tt.want)
		}
	}
}

func TestLineOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	if gb.Line(5) != "" {
		t.Fatalf("out of range should return empty string")
	}
}

func TestLineLen(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		line int
		want int
	}{
		{0, 3},
		{1, 3},
		{2, 2},
	}
	for _, tt := range tests {
		got := gb.LineLen(tt.line)
		if got != tt.want {
			t.Errorf("LineLen(%d) = %d, expected %d", tt.line, got, tt.want)
		}
	}
}

// ==================== Line/Column Conversion ====================

func TestPosToLineCol(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		pos     int
		wantL   int
		wantCol int
	}{
		{0, 0, 0},
		{2, 0, 2},
		{4, 1, 0},
		{7, 1, 3},
		{8, 2, 0},
		{9, 2, 1},
	}
	for _, tt := range tests {
		l, col := gb.PosToLineCol(tt.pos)
		if l != tt.wantL || col != tt.wantCol {
			t.Errorf("PosToLineCol(%d) = (%d,%d), expected (%d,%d)",
				tt.pos, l, col, tt.wantL, tt.wantCol)
		}
	}
}

func TestPosToLineColOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	l, col := gb.PosToLineCol(100) // should clamp to end
	if l != 0 || col != 3 {
		t.Errorf("out of range = (%d,%d), expected (0,3)", l, col)
	}
}

func TestPosToLineColNegative(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	l, col := gb.PosToLineCol(-1) // should clamp to 0
	if l != 0 || col != 0 {
		t.Errorf("negative = (%d,%d), expected (0,0)", l, col)
	}
}

func TestLineColToPos(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc\ndef\ngh"))
	tests := []struct {
		line int
		col  int
		want int
	}{
		{0, 0, 0},
		{0, 2, 2},
		{1, 0, 4},
		{1, 3, 7},
		{2, 0, 8},
		{2, 1, 9},
	}
	for _, tt := range tests {
		got := gb.LineColToPos(tt.line, tt.col)
		if got != tt.want {
			t.Errorf("LineColToPos(%d,%d) = %d, expected %d",
				tt.line, tt.col, got, tt.want)
		}
	}
}

func TestLineColToPosColOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	got := gb.LineColToPos(0, 100) // should clamp to end of line
	if got != 3 {
		t.Fatalf("column out of range = %d, expected 3", got)
	}
}

func TestLineColToPosLineOutOfRange(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	got := gb.LineColToPos(5, 0) // out-of-range line should clamp to last line
	if got != 0 {
		t.Fatalf("line out of range = %d, expected 0", got)
	}
}

func TestPosToLineColRoundtrip(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello\nworld\nfoo\nbar"))
	total := gb.Len()
	for pos := 0; pos <= total; pos++ {
		l, col := gb.PosToLineCol(pos)
		back := gb.LineColToPos(l, col)
		if back != pos {
			t.Errorf("roundtrip failed: input pos=%d → (%d,%d) → %d", pos, l, col, back)
		}
	}
}

// ==================== Indentation ====================

func TestIndentation(t *testing.T) {
	tests := []struct {
		text string
		line int
		want string
	}{
		{"    hello", 0, "    "},
		{"\t\thello", 0, "\t\t"},
		{"  \thello", 0, "  \t"},
		{"noindent", 0, ""},
		{"line1\n  line2", 1, "  "},
		{"  a\n\t\tb", 1, "\t\t"},
	}
	for _, tt := range tests {
		gb := NewGapBufferFromRunes([]rune(tt.text))
		got := gb.Indentation(tt.line)
		if got != tt.want {
			t.Errorf("Indentation(%d) text=%q → %q, expected %q",
				tt.line, tt.text, got, tt.want)
		}
	}
}

// ==================== Chinese / Unicode Support ====================

func TestInsertChinese(t *testing.T) {
	gb := NewGapBuffer(32)
	gb.Insert(0, []rune("你好世界"))
	if gb.Content() != "你好世界" {
		t.Fatalf("expected %q, got %q", "你好世界", gb.Content())
	}
}

func TestInsertMixed(t *testing.T) {
	gb := NewGapBuffer(32)
	gb.Insert(0, []rune("Hello你好World"))
	if gb.Content() != "Hello你好World" {
		t.Fatalf("expected %q, got %q", "Hello你好World", gb.Content())
	}
}

func TestDeleteChinese(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("你好世界"))
	gb.Delete(0, 2) // delete "你好"
	if gb.Content() != "世界" {
		t.Fatalf("expected %q, got %q", "世界", gb.Content())
	}
}

func TestLineCountChinese(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("你好\n世界"))
	if gb.LineCount() != 2 {
		t.Fatalf("line count should be 2, got %d", gb.LineCount())
	}
	if gb.Line(0) != "你好" {
		t.Fatalf("Line(0) = %q, expected %q", gb.Line(0), "你好")
	}
	if gb.Line(1) != "世界" {
		t.Fatalf("Line(1) = %q, expected %q", gb.Line(1), "世界")
	}
}

func TestPosToLineColChinese(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("你好世界"))
	l, col := gb.PosToLineCol(2)
	if l != 0 || col != 2 {
		t.Fatalf("PosToLineCol(2) = (%d,%d), expected (0,2)", l, col)
	}
}

// ==================== Document Wrapper ====================

func TestNewDocumentFromString(t *testing.T) {
	d := NewDocumentFromString("hello")
	if d.Content() != "hello" {
		t.Fatalf("content should be %q, got %q", "hello", d.Content())
	}
	if d.Modified {
		t.Fatalf("new document should not be marked as modified")
	}
}

func TestDocumentInsertRune(t *testing.T) {
	d := NewDocumentFromString("ab")
	d.InsertRune(1, 'x')
	if d.Content() != "axb" {
		t.Fatalf("after insert should be %q, got %q", "axb", d.Content())
	}
	if !d.Modified {
		t.Fatalf("document should be marked modified after insert")
	}
}

func TestDocumentInsertNewline(t *testing.T) {
	d := NewDocumentFromString("abc")
	d.InsertRune(1, '\n')
	if d.Content() != "a\nbc" {
		t.Fatalf("after newline insert should be %q, got %q", "a\nbc", d.Content())
	}
}

func TestDocumentInsertText(t *testing.T) {
	d := NewDocumentFromString("ac")
	d.InsertText(1, "xyz")
	if d.Content() != "axyzc" {
		t.Fatalf("after text insert should be %q, got %q", "axyzc", d.Content())
	}
}

func TestDocumentDeleteRange(t *testing.T) {
	d := NewDocumentFromString("abcdef")
	deleted := d.DeleteRange(1, 3)
	if string(deleted) != "bcd" {
		t.Fatalf("deleted text should be %q, got %q", "bcd", string(deleted))
	}
	if d.Content() != "aef" {
		t.Fatalf("after delete should be %q, got %q", "aef", d.Content())
	}
}

func TestDocumentDeleteRuneAt(t *testing.T) {
	d := NewDocumentFromString("abc")
	ch := d.DeleteRuneAt(1)
	if ch != 'b' {
		t.Fatalf("deleted rune should be 'b', got %c", ch)
	}
	if d.Content() != "ac" {
		t.Fatalf("after delete should be %q, got %q", "ac", d.Content())
	}
}

func TestDocumentDeleteRuneAtOutOfRange(t *testing.T) {
	d := NewDocumentFromString("abc")
	ch := d.DeleteRuneAt(100)
	if ch != 0 {
		t.Fatalf("out of range delete should return 0")
	}
}

// ==================== Undo/Redo ====================

func TestNewUndoStack(t *testing.T) {
	us := NewUndoStack(50)
	if us.UndoCount() != 0 {
		t.Fatalf("new stack UndoCount should be 0")
	}
	if us.RedoCount() != 0 {
		t.Fatalf("new stack RedoCount should be 0")
	}
	if us.CanUndo() {
		t.Fatalf("new stack should not be undoable")
	}
	if us.CanRedo() {
		t.Fatalf("new stack should not be redoable")
	}
}

func TestPushAndUndo(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("hello"), nil)

	if !us.CanUndo() {
		t.Fatalf("should be undoable after Push")
	}
	if us.CanRedo() {
		t.Fatalf("should not be redoable after Push")
	}

	step, ok := us.Undo()
	if !ok {
		t.Fatalf("Undo should return true")
	}
	if step.Kind != UndoInsert {
		t.Fatalf("Kind should be UndoInsert")
	}
	if step.Pos != 0 {
		t.Fatalf("Pos should be 0")
	}
	if string(step.Text) != "hello" {
		t.Fatalf("Text should be %q", "hello")
	}

	if us.UndoCount() != 0 {
		t.Fatalf("undo stack should be empty after Undo")
	}
}

func TestUndoRedoCycle(t *testing.T) {
	us := NewUndoStack(50)

	// Use different position/type operations to avoid merging
	us.Push(UndoDelete, 0, nil, []rune("a"))
	us.Push(UndoDelete, 0, nil, []rune("b"))

	us.Undo() // undo "b"
	us.Undo() // undo "a"

	if us.CanUndo() {
		t.Fatalf("should not be undoable after two Undos")
	}
	if !us.CanRedo() {
		t.Fatalf("should be redoable after Undo")
	}

	step, ok := us.Redo()
	if !ok {
		t.Fatalf("Redo should return true")
	}
	if step.Kind != UndoDelete || string(step.Deleted) != "a" {
		t.Fatalf("first Redo should be Delete 'a', got Kind=%d Deleted=%q", step.Kind, string(step.Deleted))
	}

	step, ok = us.Redo()
	if !ok {
		t.Fatalf("second Redo should return true")
	}
	if step.Kind != UndoDelete || string(step.Deleted) != "b" {
		t.Fatalf("second Redo should be Delete 'b', got Kind=%d Deleted=%q", step.Kind, string(step.Deleted))
	}

	if us.CanRedo() {
		t.Fatalf("should not be redoable after two Redos")
	}
}

func TestNewPushClearsRedo(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Undo()
	if !us.CanRedo() {
		t.Fatalf("should be redoable after Undo")
	}
	us.Push(UndoInsert, 0, []rune("b"), nil) // new operation should clear redo stack
	if us.CanRedo() {
		t.Fatalf("new Push should clear redo stack")
	}
}

func TestMergeConsecutiveInserts(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoInsert, 1, []rune("b"), nil) // consecutive insert should merge
	us.Push(UndoInsert, 2, []rune("c"), nil) // consecutive insert should merge

	if us.UndoCount() != 1 {
		t.Fatalf("three consecutive inserts should merge to 1 step, got %d steps", us.UndoCount())
	}

	step, _ := us.Undo()
	if string(step.Text) != "abc" {
		t.Fatalf("merged text should be 'abc', got %q", string(step.Text))
	}
}

func TestNoMergeNonConsecutiveInserts(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoInsert, 5, []rune("b"), nil) // non-consecutive position, should not merge

	if us.UndoCount() != 2 {
		t.Fatalf("non-consecutive inserts should produce 2 steps, got %d steps", us.UndoCount())
	}
}

func TestNoMergeDifferentTypes(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoDelete, 0, nil, []rune("a")) // different type, should not merge

	if us.UndoCount() != 2 {
		t.Fatalf("different types should produce 2 steps, got %d steps", us.UndoCount())
	}
}

func TestUndoAfterMerge(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoInsert, 1, []rune("b"), nil)

	step, ok := us.Undo()
	if !ok {
		t.Fatalf("Undo should succeed")
	}
	if string(step.Text) != "ab" {
		t.Fatalf("merged undo text should be 'ab', got %q", string(step.Text))
	}

	if us.CanUndo() {
		t.Fatalf("stack should be empty after Undo")
	}
}

func TestClearUndoStack(t *testing.T) {
	us := NewUndoStack(50)
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoInsert, 1, []rune("b"), nil)
	us.Undo() // has redo history
	us.Clear()

	if us.UndoCount() != 0 {
		t.Fatalf("undo stack should be empty after Clear, got %d", us.UndoCount())
	}
	if us.RedoCount() != 0 {
		t.Fatalf("redo stack should be empty after Clear, got %d", us.RedoCount())
	}
	if us.CanUndo() {
		t.Fatalf("should not be undoable after Clear")
	}
	if us.CanRedo() {
		t.Fatalf("should not be redoable after Clear")
	}
}

func TestUndoStepDeepCopy(t *testing.T) {
	// Verify that Push deep-copies text/deleted slices
	us := NewUndoStack(50)
	text := []rune("hello")
	us.Push(UndoInsert, 0, text, nil)
	text[0] = 'X' // modify original slice

	step, _ := us.Undo()
	if string(step.Text) != "hello" {
		t.Fatalf("deep copy failed: expected 'hello', got %q", string(step.Text))
	}
}

func TestUndoStackLimit(t *testing.T) {
	us := NewUndoStack(3) // max 3 steps
	us.Push(UndoInsert, 0, []rune("a"), nil)
	us.Push(UndoInsert, 1, []rune("b"), nil)
	us.Push(UndoInsert, 2, []rune("c"), nil)
	us.Push(UndoInsert, 3, []rune("d"), nil) // exceeds limit, drop oldest

	// Note: consecutive inserts merge
	step, _ := us.Undo()
	// Since "a","b","c","d" are all consecutive, they merge to 1 step "abcd"
	if string(step.Text) != "abcd" {
		t.Fatalf("merged should be 'abcd', got %q", string(step.Text))
	}
}

// ==================== GapBuffer Edge Cases ====================

func TestInsertAndDeleteMultiple(t *testing.T) {
	gb := NewGapBuffer(16)
	for i := 0; i < 10; i++ {
		gb.Insert(i, []rune{'a' + rune(i)})
	}
	want := "abcdefghij"
	if gb.Content() != want {
		t.Fatalf("expected %q, got %q", want, gb.Content())
	}

	// Delete every other position
	gb.Delete(1, 1) // delete 'b'
	gb.Delete(3, 1) // delete 'e'
	gb.Delete(5, 1) // delete 'h'
	if gb.Content() != "acdfgij" {
		t.Fatalf("after multiple deletes should be %q, got %q", "acdfgij", gb.Content())
	}
}

func TestGapBufferMoveGapLeft(t *testing.T) {
	// Ensure gap move left is correct
	gb := NewGapBufferFromRunes([]rune("abcdef"))
	gb.Insert(3, []rune("XYZ")) // insert before d, gap at 3
	// now content is abcXYZdef
	gb.Insert(0, []rune("!")) // move to start, gap moves left
	if gb.Content() != "!abcXYZdef" {
		t.Fatalf("after gap move left should be %q, got %q", "!abcXYZdef", gb.Content())
	}
}

func TestGapBufferMoveGapRight(t *testing.T) {
	// Ensure gap move right is correct
	gb := NewGapBufferFromRunes([]rune("abcdef"))
	gb.Insert(3, []rune("XYZ")) // insert before d, gap at 3
	gb.Insert(9, []rune("!"))   // move to end, gap moves right
	if gb.Content() != "abcXYZdef!" {
		t.Fatalf("after gap move right should be %q, got %q", "abcXYZdef!", gb.Content())
	}
}

func TestLargeInsert(t *testing.T) {
	// Large data insert, testing grow
	gb := NewGapBuffer(32)
	text := make([]rune, 1000)
	for i := range text {
		text[i] = 'a' + rune(i%26)
	}
	gb.Insert(0, text)
	got := gb.Content()
	if len([]rune(got)) != 1000 {
		t.Fatalf("length should be 1000, got %d", len([]rune(got)))
	}
	if string(got) != string(text) {
		t.Fatalf("large insert content mismatch")
	}
}

func TestContentAfterGapOperations(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello world"))
	gb.Delete(5, 6)       // delete " world"
	gb.Insert(5, []rune("!!")) // insert "!!"
	if gb.Content() != "hello!!" {
		t.Fatalf("expected %q, got %q", "hello!!", gb.Content())
	}
}

// ==================== Content/BytePos ====================

func TestContent(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune(""))
	if gb.Content() != "" {
		t.Fatalf("empty buffer Content() should be empty string")
	}
}

func TestBytePos(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("hello"))
	got := gb.BytePos(0)
	if got != 0 {
		t.Fatalf("BytePos(0) = %d, expected 0", got)
	}
	// Only test ASCII
	got = gb.BytePos(4)
	if got != 4 {
		t.Fatalf("BytePos(4) = %d, expected 4", got)
	}
}

// ==================== Document File Metadata ====================

func TestDocumentFilePath(t *testing.T) {
	d := NewDocument()
	if d.FilePath != "" {
		t.Fatalf("new document FilePath should be empty")
	}
}

func TestDocumentLineCount(t *testing.T) {
	d := NewDocumentFromString("a\nb\nc")
	if d.LineCount() != 3 {
		t.Fatalf("LineCount() = %d, expected 3", d.LineCount())
	}
}

func TestDocumentLen(t *testing.T) {
	d := NewDocumentFromString("hello")
	if d.Len() != 5 {
		t.Fatalf("Len() = %d, expected 5", d.Len())
	}
}

// ==================== Slice Copy / Comparison Tests ====================

func TestSliceReturnsCopy(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abc"))
	slice := gb.Slice(1, 3)
	slice[0] = 'X' // modify the copy

	if gb.Content() != "abc" {
		t.Fatalf("Slice should return a copy; modifying it should not change buffer content")
	}
}

// ==================== Empty Buffer Edge Cases ====================

func TestEmptyBufferLine(t *testing.T) {
	gb := NewGapBuffer(16)
	if gb.Line(0) != "" {
		t.Fatalf("empty buffer Line(0) should be empty string")
	}
	if gb.LineLen(0) != 0 {
		t.Fatalf("empty buffer LineLen(0) should be 0")
	}
	l, col := gb.PosToLineCol(0)
	if l != 0 || col != 0 {
		t.Fatalf("empty buffer PosToLineCol(0) = (%d,%d)", l, col)
	}
}

func TestEmptyBufferLineColToPos(t *testing.T) {
	gb := NewGapBuffer(16)
	got := gb.LineColToPos(0, 0)
	if got != 0 {
		t.Fatalf("empty buffer LineColToPos(0,0) = %d, expected 0", got)
	}
}

// ==================== Scenario-Based Integration Tests ===================

func TestTypingSimulation(t *testing.T) {
	// Simulate typing: insert characters sequentially
	gb := NewGapBuffer(16)
	text := "Hello, 世界!"
	runes := []rune(text)
	for i, ch := range runes {
		gb.Insert(i, []rune{ch})
	}
	if gb.Content() != text {
		t.Fatalf("after typing simulation should be %q, got %q", text, gb.Content())
	}
}

func TestBackspaceSimulation(t *testing.T) {
	// Simulate backspace
	gb := NewGapBufferFromRunes([]rune("hello"))
	for i := 0; i < 3; i++ {
		gb.Delete(gb.Len()-1, 1)
	}
	if gb.Content() != "he" {
		t.Fatalf("after backspace should be 'he', got %q", gb.Content())
	}
}

func TestInsertAndDeleteSamePosition(t *testing.T) {
	// Repeated insert/delete at the same position, verifying correct gap movement
	gb := NewGapBufferFromRunes([]rune("ab"))
	gb.Insert(1, []rune{'1'}) // a1b
	gb.Delete(1, 1)           // ab
	gb.Insert(1, []rune{'2'}) // a2b
	gb.Delete(1, 1)           // ab
	gb.Insert(1, []rune{'3'}) // a3b
	if gb.Content() != "a3b" {
		t.Fatalf("after repeated insert/delete should be 'a3b', got %q", gb.Content())
	}
}

func TestMultipleInsertsDifferentPositions(t *testing.T) {
	gb := NewGapBuffer(16)
	gb.Insert(0, []rune("world"))
	gb.Insert(0, []rune("hello "))
	gb.Insert(11, []rune("!!!"))
	if gb.Content() != "hello world!!!" {
		t.Fatalf("expected %q, got %q", "hello world!!!", gb.Content())
	}
}

func TestRepeatedDeleteAtStart(t *testing.T) {
	gb := NewGapBufferFromRunes([]rune("abcdef"))
	for i := 0; i < 3; i++ {
		gb.Delete(0, 1)
	}
	if gb.Content() != "def" {
		t.Fatalf("after three deletes at start should be 'def', got %q", gb.Content())
	}
}
