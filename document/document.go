// Package document implements Gap Buffer based text storage and editing.
//
// Design:
// Gap Buffer is one of the classic in-memory structures for text editors.
// It splits the buffer into three parts:
//   [left half] [gap] [right half]
// The cursor position corresponds to the gap start. Inserting at the cursor
// writes characters into the gap and advances the gap pointer without moving
// large amounts of data. Deleting simply expands the gap.
//
// Complexity:
// - Insert/Delete: amortized O(1)
// - Line-based access: requires scanning, O(n), but can be optimized with line start caching
package document

import (
	"unicode/utf8"
)

// GapBuffer is the gap buffer data structure.
type GapBuffer struct {
	buf      []rune // storage for all runes
	gapStart int    // gap start position (inclusive)
	gapEnd   int    // gap end position (exclusive)
}

// NewGapBuffer creates a gap buffer with the specified initial capacity.
func NewGapBuffer(capacity int) *GapBuffer {
	if capacity < 16 {
		capacity = 16
	}
	buf := make([]rune, capacity)
	return &GapBuffer{
		buf:      buf,
		gapStart: 0,
		gapEnd:   capacity,
	}
}

// NewGapBufferFromRunes creates a GapBuffer from a rune slice.
func NewGapBufferFromRunes(runes []rune) *GapBuffer {
	n := len(runes)
	cap := n + 16
	buf := make([]rune, cap)
	copy(buf, runes)
	return &GapBuffer{
		buf:      buf,
		gapStart: n,
		gapEnd:   cap,
	}
}

// Len returns the number of valid characters (excluding the gap).
func (gb *GapBuffer) Len() int {
	return len(gb.buf) - (gb.gapEnd - gb.gapStart)
}

// Cap returns the total buffer capacity.
func (gb *GapBuffer) Cap() int {
	return len(gb.buf)
}

// grow increases capacity to at least the needed size.
func (gb *GapBuffer) grow(needed int) {
	oldCap := len(gb.buf)
	newCap := oldCap * 2
	if newCap < needed {
		newCap = needed + 16
	}
	newBuf := make([]rune, newCap)
	// Copy left half
	copy(newBuf, gb.buf[:gb.gapStart])
	// Shift right half
	rightLen := oldCap - gb.gapEnd
	newGapEnd := newCap - rightLen
	copy(newBuf[newGapEnd:], gb.buf[gb.gapEnd:oldCap])
	gb.buf = newBuf
	gb.gapEnd = newGapEnd
}

// moveGap moves the gap to position pos (0 <= pos <= Len()).
func (gb *GapBuffer) moveGap(pos int) {
	if pos < 0 {
		pos = 0
	}
	if pos > gb.Len() {
		pos = gb.Len()
	}

	currPos := gb.gapStart

	if pos < currPos {
		// Move left: shift data from the left end to the gap end
		moveLen := currPos - pos
		srcStart := pos
		dstStart := gb.gapEnd - moveLen
		copy(gb.buf[dstStart:dstStart+moveLen], gb.buf[srcStart:srcStart+moveLen])
		gb.gapStart = pos
		gb.gapEnd = dstStart
	} else if pos > currPos {
		// Move right: shift data from the right start to the gap start
		moveLen := pos - currPos
		srcStart := gb.gapEnd
		dstStart := currPos
		copy(gb.buf[dstStart:dstStart+moveLen], gb.buf[srcStart:srcStart+moveLen])
		gb.gapStart = pos
		gb.gapEnd = srcStart + moveLen
	}
	// pos == currPos: no move needed
}

// Insert inserts a rune slice at position pos.
func (gb *GapBuffer) Insert(pos int, text []rune) {
	if len(text) == 0 {
		return
	}
	gb.moveGap(pos)
	// Grow if the gap is too small
	needed := len(gb.buf) + len(text) - (gb.gapEnd - gb.gapStart)
	cap := len(gb.buf)
	if needed > cap {
		gb.grow(needed)
	}
	// Write into the gap
	copy(gb.buf[gb.gapStart:], text)
	gb.gapStart += len(text)
}

// Delete deletes 'length' characters starting at position pos.
func (gb *GapBuffer) Delete(pos, length int) {
	if length <= 0 || pos < 0 || pos+length > gb.Len() {
		return
	}
	gb.moveGap(pos)
	gb.gapEnd += length
}

// RuneAt returns the rune at position pos.
func (gb *GapBuffer) RuneAt(pos int) rune {
	if pos < 0 || pos >= gb.Len() {
		return 0
	}
	if pos < gb.gapStart {
		return gb.buf[pos]
	}
	return gb.buf[pos+(gb.gapEnd-gb.gapStart)]

}

// Slice returns a slice of runes from start to end (exclusive).
// Clamps out-of-range indices.
func (gb *GapBuffer) Slice(start, end int) []rune {
	l := gb.Len()
	if start < 0 {
		start = 0
	}
	if end > l {
		end = l
	}
	if start >= end {
		return []rune{}
	}

	result := make([]rune, end-start)
	// Calculate positions considering the gap
	// The logical-to-physical index shift is gapSize = gb.gapEnd - gb.gapStart
	for i := 0; i < end-start; i++ {
		logicalPos := start + i
		if logicalPos < gb.gapStart {
			result[i] = gb.buf[logicalPos]
		} else {
			result[i] = gb.buf[logicalPos+(gb.gapEnd-gb.gapStart)]
		}
	}
	return result
}

// Content returns the complete text content as a string.
func (gb *GapBuffer) Content() string {
	return string(gb.Slice(0, gb.Len()))
}

// LineStart returns the position of the first character on the given line.
// Returns -1 if the line is out of range.
func (gb *GapBuffer) LineStart(line int) int {
	if line < 0 {
		return -1
	}

	l := gb.Len()
	if l == 0 {
		if line == 0 {
			return 0
		}
		return -1
	}

	// Count newlines to find the line
	count := 0
	if line == 0 {
		return 0
	}

	for i := 0; i < l; i++ {
		if gb.RuneAt(i) == '\n' {
			count++
			if count == line {
				return i + 1
			}
		}
	}

	return -1 // Line out of range
}

// LineEnd returns the position of the first character after the given line,
// excluding the newline character. Returns the buffer length if the line has
// no trailing newline.
func (gb *GapBuffer) LineEnd(line int) int {
	l := gb.Len()
	start := gb.LineStart(line)
	if start < 0 {
		return -1
	}

	// Scan for next newline from the start position
	for i := start; i < l; i++ {
		if gb.RuneAt(i) == '\n' {
			return i
		}
	}

	return l // Last line, no trailing newline
}

// LineCount returns the total number of lines.
// An empty buffer counts as 1 line. Content ending with a newline counts
// the trailing newline as starting an additional empty line.
func (gb *GapBuffer) LineCount() int {
	l := gb.Len()
	if l == 0 {
		return 1
	}

	count := 1
	for i := 0; i < l; i++ {
		if gb.RuneAt(i) == '\n' {
			count++
		}
	}
	return count
}

// Line returns the content of the specified line as a string.
func (gb *GapBuffer) Line(line int) string {
	start := gb.LineStart(line)
	if start < 0 {
		return ""
	}
	end := gb.LineEnd(line)
	if end < 0 {
		return ""
	}
	return string(gb.Slice(start, end))
}

// LineLen returns the number of characters on the given line.
func (gb *GapBuffer) LineLen(line int) int {
	return len([]rune(gb.Line(line)))
}

// PosToLineCol converts a global position to (line, column).
func (gb *GapBuffer) PosToLineCol(pos int) (int, int) {
	l := gb.Len()
	if pos < 0 {
		pos = 0
	}
	if pos > l {
		pos = l
	}

	line := 0
	col := 0
	for i := 0; i < pos; i++ {
		if gb.RuneAt(i) == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return line, col
}

// LineColToPos converts (line, column) to a global position.
// Clamps out-of-range values.
func (gb *GapBuffer) LineColToPos(line, col int) int {
	l := gb.Len()
	totalLines := gb.LineCount()

	if line >= totalLines {
		line = totalLines - 1
	}
	if line < 0 {
		line = 0
	}

	start := gb.LineStart(line)
	if start < 0 {
		return l
	}

	lineContent := gb.Line(line)
	lineLen := len([]rune(lineContent))

	if col > lineLen {
		col = lineLen
	}
	if col < 0 {
		col = 0
	}

	return start + col
}

// Indentation returns the whitespace indentation string for the given line.
func (gb *GapBuffer) Indentation(line int) string {
	l := gb.Line(line)
	runes := []rune(l)
	var indent []rune
	for _, ch := range runes {
		if ch == ' ' || ch == '\t' {
			indent = append(indent, ch)
		} else {
			break
		}
	}
	return string(indent)
}

// String returns the content as a string.
func (gb *GapBuffer) String() string {
	return gb.Content()
}

// BytePos converts a rune position to a byte position.
func (gb *GapBuffer) BytePos(runePos int) int {
	s := gb.Content()
	if runePos <= 0 {
		return 0
	}
	runes := []rune(s)
	if runePos >= len(runes) {
		return len(s)
	}
	// Accumulate byte offset from rune positions
	bytePos := 0
	for i := 0; i < runePos; i++ {
		bytePos += utf8.RuneLen(runes[i])
	}
	return bytePos
}

// Document is the high-level document structure wrapping GapBuffer.
type Document struct {
	Buffer   *GapBuffer // underlying gap buffer
	FilePath string     // file path
	Modified bool       // modified flag
}

// NewDocument creates a new empty document.
func NewDocument() *Document {
	return &Document{
		Buffer: NewGapBuffer(16),
	}
}

// NewDocumentFromString creates a document from a string.
func NewDocumentFromString(s string) *Document {
	return &Document{
		Buffer: NewGapBufferFromRunes([]rune(s)),
	}
}

// InsertRune inserts a single rune at position pos.
func (d *Document) InsertRune(pos int, ch rune) {
	d.Buffer.Insert(pos, []rune{ch})
	d.Modified = true
}

// InsertText inserts a string at position pos.
func (d *Document) InsertText(pos int, text string) {
	if len(text) == 0 {
		return
	}
	d.Buffer.Insert(pos, []rune(text))
	d.Modified = true
}

// DeleteRange deletes 'length' characters starting at pos, returning the deleted text.
func (d *Document) DeleteRange(pos, length int) []rune {
	deleted := d.Buffer.Slice(pos, pos+length)
	d.Buffer.Delete(pos, length)
	d.Modified = true
	return deleted
}

// DeleteRuneAt deletes the character at position pos, returning the deleted rune.
func (d *Document) DeleteRuneAt(pos int) rune {
	if pos < 0 || pos >= d.Buffer.Len() {
		return 0
	}
	ch := d.Buffer.RuneAt(pos)
	d.Buffer.Delete(pos, 1)
	d.Modified = true
	return ch
}

// Content returns the complete text content.
func (d *Document) Content() string {
	return d.Buffer.Content()
}

// Line returns the content of the specified line.
func (d *Document) Line(n int) string {
	return d.Buffer.Line(n)
}

// LineCount returns the total number of lines.
func (d *Document) LineCount() int {
	return d.Buffer.LineCount()
}

// LineLen returns the number of characters on the given line.
func (d *Document) LineLen(line int) int {
	return d.Buffer.LineLen(line)
}

// Len returns the total number of characters.
func (d *Document) Len() int {
	return d.Buffer.Len()
}

// GetIndent returns the indentation string for the given line.
func (d *Document) GetIndent(line int) string {
	return d.Buffer.Indentation(line)
}

// PosToLineCol converts a global position to (line, column).
func (d *Document) PosToLineCol(pos int) (int, int) {
	return d.Buffer.PosToLineCol(pos)
}

// LineColToPos converts (line, column) to a global position.
func (d *Document) LineColToPos(line, col int) int {
	return d.Buffer.LineColToPos(line, col)
}
