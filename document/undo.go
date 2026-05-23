// Package document's undo.go implements the undo/redo system.
//
// Design:
// Uses the Command Pattern — each edit operation is recorded as an UndoStep.
// The stack capacity is configurable (default 100 steps); oldest records are dropped when exceeded.
// Merge strategy: consecutive insert operations are merged into one step to reduce redundancy.
package document

// UndoKind represents the operation type.
type UndoKind int

const (
	UndoInsert UndoKind = iota // Insert operation
	UndoDelete                 // Delete operation
	UndoReplace                // Replace operation
)

// Position represents a position in the document.
type Position struct {
	Line int // Line number (0-indexed)
	Col  int // Column number (0-indexed)
}

// UndoStep records a single undoable operation step.
type UndoStep struct {
	Kind    UndoKind // Operation type
	Pos     int      // Global position where the operation occurred
	Text    []rune   // Inserted text (Insert) or replaced text (Replace)
	Deleted []rune   // Deleted text (Delete) or pre-replacement text (Replace)
}

// UndoStack manages the undo/redo stacks.
type UndoStack struct {
	steps    []UndoStep // Undo stack
	redo     []UndoStep // Redo stack
	limit    int        // Max steps
	lastKind UndoKind   // Last operation type, used for merge detection
	lastSize int        // Last operation text length
}

// NewUndoStack creates an undo stack with the given step limit.
func NewUndoStack(limit int) *UndoStack {
	if limit < 1 {
		limit = 100
	}
	return &UndoStack{
		steps: make([]UndoStep, 0, limit),
		redo:  make([]UndoStep, 0, limit),
		limit: limit,
	}
}

// Push records an operation. Consecutive same-type inserts are merged to reduce steps.
func (us *UndoStack) Push(kind UndoKind, pos int, text, deleted []rune) {
	// Merge strategy: same type as last, consecutive position, insert-only
	canMerge := kind == UndoInsert &&
		us.lastKind == UndoInsert &&
		len(us.steps) > 0 &&
		pos == us.steps[len(us.steps)-1].Pos+us.lastSize &&
		kind == us.lastKind

	if canMerge {
		last := &us.steps[len(us.steps)-1]
		last.Text = append(last.Text, text...)
		us.lastSize += len(text)
		// Clear redo stack on merge
		us.redo = us.redo[:0]
		return
	}

	step := UndoStep{
		Kind:    kind,
		Pos:     pos,
		Text:    append([]rune{}, text...),     // Deep copy
		Deleted: append([]rune{}, deleted...), // Deep copy
	}

	us.steps = append(us.steps, step)
	us.lastKind = kind
	us.lastSize = len(text)

	// Exceed limit, drop oldest step
	if len(us.steps) > us.limit {
		us.steps = us.steps[1:]
	}

	// New operation clears redo stack
	us.redo = us.redo[:0]
}

// CanUndo checks whether undo is available.
func (us *UndoStack) CanUndo() bool {
	return len(us.steps) > 0
}

// CanRedo checks whether redo is available.
func (us *UndoStack) CanRedo() bool {
	return len(us.redo) > 0
}

// Undo pops and returns the last operation step for undoing.
func (us *UndoStack) Undo() (UndoStep, bool) {
	if len(us.steps) == 0 {
		return UndoStep{}, false
	}
	step := us.steps[len(us.steps)-1]
	us.steps = us.steps[:len(us.steps)-1]
	us.redo = append(us.redo, step)
	us.lastKind = -1 // Reset merge state
	return step, true
}

// Redo pops and returns the top of the redo stack.
func (us *UndoStack) Redo() (UndoStep, bool) {
	if len(us.redo) == 0 {
		return UndoStep{}, false
	}
	step := us.redo[len(us.redo)-1]
	us.redo = us.redo[:len(us.redo)-1]
	us.steps = append(us.steps, step)
	us.lastKind = -1
	return step, true
}

// Clear clears all steps.
func (us *UndoStack) Clear() {
	us.steps = us.steps[:0]
	us.redo = us.redo[:0]
	us.lastKind = -1
}

// UndoCount returns the number of steps in the undo stack.
func (us *UndoStack) UndoCount() int {
	return len(us.steps)
}

// RedoCount returns the number of steps in the redo stack.
func (us *UndoStack) RedoCount() int {
	return len(us.redo)
}
