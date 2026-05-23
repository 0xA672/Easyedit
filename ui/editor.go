// Package ui implements the terminal text editor user interface.
//
// Uses the tcell library for terminal event handling and rendering.
// The editor operates in three modes:
//   - Insert mode (default): direct text input, arrow key navigation
//   - Command mode: press : to enter, command bar at bottom
//   - Search mode: press Ctrl+F to enter, search bar at bottom
//
// Design:
// The Editor struct holds all state; the main loop in Run() handles events and rendering.
// Rendering is split into three sections: body (line numbers + text), status bar, command/search bar.
package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	"easyedit/command"
	"easyedit/config"
	"easyedit/document"
	"easyedit/highlight"
)

// Mode represents the editor mode.
type Mode int

const (
	ModeInsert Mode = iota // Insert mode
	ModeCommand            // Command mode
	ModeSearch             // Search mode
	ModeReplace            // Replace mode
)

// Editor is the core editor struct holding all state.
type Editor struct {
	config config.Config // Configuration

	doc      *document.Document // Current document
	undo     *document.UndoStack
	screen   tcell.Screen
	hl       *highlight.Highlighter // Syntax highlighter
	clipText string                 // Internal clipboard text

	// Cursor position
	cursor    int // Global position
	cursorRow int // Screen row
	cursorCol int // Screen column

	// Viewport offset
	offsetRow int // Vertical offset (first visible row)
	offsetCol int // Horizontal offset (first visible column)

	// Mode
	mode Mode

	// Selection
	selStart int  // Selection start position (-1 means no selection)
	selEnd   int  // Selection end position
	selecting bool // Whether currently selecting

	// Search
	searchQuery   string
	searchRegex   bool    // Enable regex
	searchCase    bool    // Enable case sensitivity
	searchResults []int   // Match positions
	searchIdx     int     // Current match index
	searchDir     int     // Search direction: 1 forward, -1 backward

	// Command bar
	cmdBuffer   string // Command/search input buffer
	message     string // Status bar message
	messageTime time.Time

	// Terminal dimensions
	termW int // Terminal width
	termH int // Terminal height

	running bool

	// Replace mode
	replaceOld  string
	replaceNew  string
	replaceFlag string

	// Syntax highlight cache (avoid re-tokenization every frame)
	lastContent string           // Text from last render
	hlTokens    []chroma.Token   // Cached tokenization result
	runeColors  []tcell.Color    // Foreground color per rune position (0=default)
}

// Token type alias for convenience.
type chromaToken = chroma.Token

// NewEditor creates a new Editor instance.
func NewEditor() *Editor {
	e := &Editor{
		config:    config.LoadConfig(),
		doc:       document.NewDocument(),
		undo:      document.NewUndoStack(100),
		running:   true,
		selStart:  -1,
		selEnd:    -1,
		searchDir: 1,
		mode:      ModeInsert,
	}
	e.hl = highlight.NewHighlighter("")
	return e
}

// Run is the editor main loop.
func (e *Editor) Run() error {
	var err error
	e.screen, err = tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("cannot create screen: %w", err)
	}
	if err := e.screen.Init(); err != nil {
		return fmt.Errorf("cannot init screen: %w", err)
	}
	defer e.screen.Fini()

	e.termW, e.termH = e.screen.Size()
	e.screen.EnableMouse()

	e.showMsg("EasyEdit - Press Ctrl+Q to quit, Ctrl+S to save, : for command")

	for e.running {
		e.render()
		e.screen.Show()
		ev := e.screen.PollEvent()
		e.handleEvent(ev)
	}
	return nil
}

// Quit requests the editor to quit.
func (e *Editor) Quit(force bool) error {
	if !force && e.doc.Modified {
		e.showMsg("File has unsaved changes! Use :q! to force quit.")
		return fmt.Errorf("unsaved changes")
	}
	e.running = false
	return nil
}

// OpenFile opens a file.
func (e *Editor) OpenFile(path string) error {
	if path == "" {
		return fmt.Errorf("no filename")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", path, err)
	}
	e.doc = document.NewDocumentFromString(string(data))
	e.doc.FilePath = path
	e.undo = document.NewUndoStack(e.config.UndoLimit)
	e.cursor = 0
	e.offsetRow = 0
	e.offsetCol = 0
	e.hl.SetFile(path)
	e.showMsg(fmt.Sprintf("opened %s (%d lines)", filepath.Base(path), e.doc.LineCount()))
	return nil
}

// SaveFile saves the file; if path is empty, uses the original path.
func (e *Editor) SaveFile(path string) error {
	if path == "" {
		path = e.doc.FilePath
	}
	if path == "" {
		return fmt.Errorf("no filename (use :w <path> or Ctrl+S)")
	}

	// Create backup
	if e.config.Backup {
		backupPath := path + ".bak"
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err == nil {
				os.WriteFile(backupPath, data, 0644)
			}
		}
	}

	// Write file
	content := e.doc.Content()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("cannot write %s: %w", path, err)
	}

	e.doc.FilePath = path
	e.doc.Modified = false
	e.hl.SetFile(path)
	return nil
}

// showMsg displays a message in the status bar.
func (e *Editor) showMsg(msg string) {
	e.message = msg
	e.messageTime = time.Now()
}

// handleEvent dispatches events to the appropriate mode handler.
func (e *Editor) handleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		e.termW, e.termH = ev.Size()
		return
	case *tcell.EventMouse:
		e.handleMouse(ev)
		return
	case *tcell.EventKey:
		switch e.mode {
		case ModeInsert:
			e.handleInsertMode(ev)
		case ModeCommand:
			e.handleCommandMode(ev)
		case ModeSearch:
			e.handleSearchMode(ev)
		case ModeReplace:
			e.handleReplaceMode(ev)
		}
	}
}

// ---- Mouse Events ----

func (e *Editor) handleMouse(ev *tcell.EventMouse) {
	x, y := ev.Position()
	// Clicking the status bar enters command mode
	statusLine := e.termH - 1
	if y == statusLine {
		// Click command bar area
		e.mode = ModeCommand
		e.cmdBuffer = ""
		return
	}
	// Calculate the line and column corresponding to the click
	lineNumWidth := e.lineNumWidth()
	clickCol := x - lineNumWidth
	if clickCol < 0 {
		clickCol = 0
	}
	clickRow := y + e.offsetRow
	if clickRow >= e.doc.LineCount() {
		clickRow = e.doc.LineCount() - 1
	}
	if clickRow < 0 {
		clickRow = 0
	}
	e.cursor = e.doc.LineColToPos(clickRow, clickCol)
	e.clampCursor()
	// If primary button is pressed, start selecting
	if ev.Buttons()&tcell.ButtonPrimary != 0 {
		e.selecting = true
		e.selStart = e.cursor
		e.selEnd = e.cursor
	}
}

// ---- Insert Mode Event Handling ----

func (e *Editor) handleInsertMode(ev *tcell.EventKey) {
	// Special shortcuts
	switch {
	case ev.Rune() == ':':
		// Colon enters command mode (must check before KeyRune)
		e.mode = ModeCommand
		e.cmdBuffer = ""
	case ev.Modifiers() == tcell.ModAlt && ev.Rune() == 'w':
		// Alt+W toggles soft wrap
		e.config.SoftWrap = !e.config.SoftWrap
		if e.config.SoftWrap {
			e.showMsg("Soft wrap: ON")
		} else {
			e.showMsg("Soft wrap: OFF")
		}
	case ev.Key() == tcell.KeyRune:
		e.doInsert(ev.Rune())
	case ev.Key() == tcell.KeyEnter:
		e.doNewLine()

	case ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2:
		e.doBackspace()

	case ev.Key() == tcell.KeyDelete:
		e.doDelete()

	case ev.Key() == tcell.KeyLeft:
		e.moveCursor(-1)
	case ev.Key() == tcell.KeyRight:
		e.moveCursor(1)
	case ev.Key() == tcell.KeyUp:
		e.moveUp()
	case ev.Key() == tcell.KeyDown:
		e.moveDown()
	case ev.Key() == tcell.KeyHome:
		e.moveLineStart()
	case ev.Key() == tcell.KeyEnd:
		e.moveLineEnd()
	case ev.Key() == tcell.KeyPgUp:
		e.movePageUp()
	case ev.Key() == tcell.KeyPgDn:
		e.movePageDown()

	case ev.Key() == tcell.KeyCtrlS:
		if err := e.SaveFile(""); err != nil {
			e.showMsg(err.Error())
		} else {
			e.showMsg("saved")
		}
	case ev.Key() == tcell.KeyCtrlQ:
		e.Quit(false)
	case ev.Key() == tcell.KeyCtrlF:
		e.mode = ModeSearch
		e.cmdBuffer = ""
		e.searchQuery = ""
	case ev.Key() == tcell.KeyCtrlA:
		e.selectAll()
	case ev.Key() == tcell.KeyCtrlZ:
		e.doUndo()
	case ev.Key() == tcell.KeyCtrlY:
		e.doRedo()
	case ev.Key() == tcell.KeyCtrlX:
		e.doCut()
	case ev.Key() == tcell.KeyCtrlC:
		// In terminal, Ctrl+C is copy, should not quit
		e.doCopy()
	case ev.Key() == tcell.KeyCtrlV:
		e.doPaste()

	case ev.Key() == tcell.KeyEscape:
		// Esc enters command mode
		e.mode = ModeCommand
		e.cmdBuffer = ""

	case ev.Key() == tcell.KeyTab:
		e.doInsert('\t')

	default:
		// Ignore other keys
	}
}

// ---- Command Mode ----

func (e *Editor) handleCommandMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.mode = ModeInsert
		e.cmdBuffer = ""
		e.showMsg("")
	case tcell.KeyEnter:
		e.executeCommand()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(e.cmdBuffer) > 0 {
			e.cmdBuffer = e.cmdBuffer[:len(e.cmdBuffer)-1]
		}
	case tcell.KeyRune:
		e.cmdBuffer += string(ev.Rune())
	}
}

// ---- Search Mode ----

func (e *Editor) handleSearchMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.mode = ModeInsert
		e.cmdBuffer = ""
		e.unhighlightSearch()
	case tcell.KeyEnter:
		e.searchQuery = e.cmdBuffer
		e.doSearch()
		e.mode = ModeInsert
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(e.cmdBuffer) > 0 {
			e.cmdBuffer = e.cmdBuffer[:len(e.cmdBuffer)-1]
		}
	case tcell.KeyRune:
		e.cmdBuffer += string(ev.Rune())
		// Live search preview
		e.searchQuery = e.cmdBuffer
		e.doSearch()
		e.render()
		e.screen.Show()
	case tcell.KeyCtrlR:
		e.searchRegex = !e.searchRegex
		e.searchQuery = e.cmdBuffer
		e.doSearch()
	case tcell.KeyCtrlC:
		e.searchCase = !e.searchCase
		e.searchQuery = e.cmdBuffer
		e.doSearch()
	}
}

// ---- Replace Mode ----

func (e *Editor) handleReplaceMode(ev *tcell.EventKey) {
	// Replace mode has two input steps: search pattern / replacement text
	switch ev.Key() {
	case tcell.KeyEscape:
		e.mode = ModeInsert
		e.cmdBuffer = ""
		e.replaceOld = ""
		e.replaceNew = ""
	case tcell.KeyEnter:
		if e.replaceOld == "" {
			// First line input: search pattern
			e.replaceOld = e.cmdBuffer
			e.cmdBuffer = ""
			e.showMsg("Replace with:")
		} else if e.replaceNew == "" {
			// Second line input: replacement text
			e.replaceNew = e.cmdBuffer
			e.cmdBuffer = ""
			e.doReplace()
			e.mode = ModeInsert
			e.replaceOld = ""
			e.replaceNew = ""
		}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(e.cmdBuffer) > 0 {
			e.cmdBuffer = e.cmdBuffer[:len(e.cmdBuffer)-1]
		}
	case tcell.KeyRune:
		e.cmdBuffer += string(ev.Rune())
	}
}

// ---- Edit Operations ----

func (e *Editor) doInsert(ch rune) {
	e.undo.Push(document.UndoInsert, e.cursor, []rune{ch}, nil)
	e.doc.InsertRune(e.cursor, ch)
	e.cursor++
	e.clampCursor()
}

func (e *Editor) doNewLine() {
	e.undo.Push(document.UndoInsert, e.cursor, []rune{'\n'}, nil)
	e.doc.InsertRune(e.cursor, '\n')
	e.cursor++

	// Auto-indent
	if e.config.AutoIndent {
		line, _ := e.doc.PosToLineCol(e.cursor)
		prevLine := line - 1
		if prevLine > 0 {
			indent := e.doc.GetIndent(prevLine - 1)
			// If the previous line ends with {, increase indentation level
			prevLineContent := e.doc.Line(prevLine - 1)
			lastRune := rune(0)
			if len(prevLineContent) > 0 {
				runes := []rune(prevLineContent)
				lastRune = runes[len(runes)-1]
			}
			if indent == "" && prevLineContent != "" && lastRune == '{' {
				for i := 0; i < e.config.TabWidth; i++ {
					indent += " "
				}
			}
			if indent != "" {
				e.doc.InsertText(e.cursor, indent)
				e.cursor += len([]rune(indent))
			}
		}
	}

	e.clampCursor()
}

func (e *Editor) doBackspace() {
	if e.cursor <= 0 {
		return
	}
	e.cursor--
	ch := e.doc.DeleteRuneAt(e.cursor)
	e.undo.Push(document.UndoDelete, e.cursor, nil, []rune{ch})
	e.clampCursor()
}

func (e *Editor) doDelete() {
	if e.cursor >= e.doc.Len() {
		return
	}
	ch := e.doc.DeleteRuneAt(e.cursor)
	e.undo.Push(document.UndoDelete, e.cursor, nil, []rune{ch})
	e.clampCursor()
}

// ---- Undo / Redo ----

func (e *Editor) doUndo() {
	step, ok := e.undo.Undo()
	if !ok {
		return
	}
	switch step.Kind {
	case document.UndoInsert:
		e.cursor = step.Pos
		e.doc.DeleteRange(step.Pos, len(step.Text))
	case document.UndoDelete:
		e.doc.InsertText(step.Pos, string(step.Deleted))
		e.cursor = step.Pos
	case document.UndoReplace:
		e.doc.DeleteRange(step.Pos, len(step.Text))
		e.doc.InsertText(step.Pos, string(step.Deleted))
		e.cursor = step.Pos + len(step.Deleted)
	}
	e.clampCursor()
	e.showMsg(fmt.Sprintf("undo (%d steps remaining)", e.undo.UndoCount()))
}

func (e *Editor) doRedo() {
	step, ok := e.undo.Redo()
	if !ok {
		return
	}
	switch step.Kind {
	case document.UndoInsert:
		e.doc.InsertText(step.Pos, string(step.Text))
		e.cursor = step.Pos + len(step.Text)
	case document.UndoDelete:
		e.cursor = step.Pos
		e.doc.DeleteRange(step.Pos, len(step.Deleted))
	case document.UndoReplace:
		e.doc.DeleteRange(step.Pos, len(step.Deleted))
		e.doc.InsertText(step.Pos, string(step.Text))
		e.cursor = step.Pos + len(step.Text)
	}
	e.clampCursor()
	e.showMsg(fmt.Sprintf("redo (%d steps remaining)", e.undo.RedoCount()))
}

// ---- Cut / Copy / Paste ----

func (e *Editor) doCut() {
	if e.selStart < 0 {
		return
	}
	start, end := e.getSelectionRange()
	text := e.doc.DeleteRange(start, end-start)
	e.undo.Push(document.UndoDelete, start, nil, text)
	e.cursor = start

	// Write to system clipboard
	str := string(text)
	e.clipText = str
	_ = clipboard.WriteAll(str)

	e.selStart = -1
	e.selEnd = -1
	e.clampCursor()
}

func (e *Editor) doCopy() {
	if e.selStart < 0 {
		return
	}
	start, end := e.getSelectionRange()
	text := e.doc.Buffer.Slice(start, end)
	str := string(text)
	e.clipText = str
	_ = clipboard.WriteAll(str)
	e.showMsg(fmt.Sprintf("copied %d characters", len(text)))
}

func (e *Editor) doPaste() {
	text, err := clipboard.ReadAll()
	if err != nil || text == "" {
		// Fallback to internal clipboard
		text = e.clipText
	}
	if text == "" {
		return
	}

	// Delete selection first
	if e.selStart >= 0 {
		start, end := e.getSelectionRange()
		_ = e.doc.DeleteRange(start, end-start)
		e.cursor = start
		e.selStart = -1
		e.selEnd = -1
	}

	e.undo.Push(document.UndoInsert, e.cursor, []rune(text), nil)
	e.doc.InsertText(e.cursor, text)
	e.cursor += len([]rune(text))
	e.clampCursor()
}

// ---- Selection ----

func (e *Editor) selectAll() {
	e.selStart = 0
	e.selEnd = e.doc.Len()
	e.cursor = e.selEnd
	e.clampCursor()
}

func (e *Editor) getSelectionRange() (int, int) {
	if e.selStart < 0 {
		return e.cursor, e.cursor
	}
	start := e.selStart
	end := e.selEnd
	if end < start {
		start, end = end, start
	}
	return start, end
}

// ---- Search ----

func (e *Editor) doSearch() {
	e.searchResults = e.searchResults[:0]
	e.searchIdx = -1

	if e.searchQuery == "" {
		return
	}

	content := e.doc.Content()
	if !e.searchCase {
		content = strings.ToLower(content)
		query := strings.ToLower(e.searchQuery)
		if e.searchRegex {
			re, err := regexp.Compile(query)
			if err != nil {
				e.showMsg(fmt.Sprintf("regex error: %v", err))
				return
			}
			matches := re.FindAllStringIndex(content, -1)
			for _, m := range matches {
				e.searchResults = append(e.searchResults, m[0])
			}
		} else {
			idx := 0
			for {
				pos := strings.Index(content[idx:], query)
				if pos < 0 {
					break
				}
				e.searchResults = append(e.searchResults, idx+pos)
				idx += pos + 1
			}
		}
	} else {
		if e.searchRegex {
			re, err := regexp.Compile(e.searchQuery)
			if err != nil {
				e.showMsg(fmt.Sprintf("regex error: %v", err))
				return
			}
			matches := re.FindAllStringIndex(content, -1)
			for _, m := range matches {
				e.searchResults = append(e.searchResults, m[0])
			}
		} else {
			idx := 0
			for {
				pos := strings.Index(content[idx:], e.searchQuery)
				if pos < 0 {
					break
				}
				e.searchResults = append(e.searchResults, idx+pos)
				idx += pos + 1
			}
		}
	}

	if len(e.searchResults) > 0 {
		// Find the first match after cursor
		for i, pos := range e.searchResults {
			if pos >= e.cursor {
				e.searchIdx = i
				e.cursor = pos
				break
			}
		}
		if e.searchIdx < 0 {
			e.searchIdx = 0
			e.cursor = e.searchResults[0]
		}
		e.clampCursor()
		e.showMsg(fmt.Sprintf("found %d matches", len(e.searchResults)))
	} else {
		e.showMsg("no matches found")
	}
}

func (e *Editor) unhighlightSearch() {
	e.searchResults = nil
	e.searchIdx = -1
}

// ---- Replace ----

func (e *Editor) doReplace() {
	if e.replaceOld == "" {
		return
	}

	content := e.doc.Content()
	var re *regexp.Regexp
	var err error

	if e.searchRegex {
		pattern := e.replaceOld
		if !e.searchCase {
			pattern = "(?i)" + pattern
		}
		re, err = regexp.Compile(pattern)
	} else {
		pattern := regexp.QuoteMeta(e.replaceOld)
		if !e.searchCase {
			pattern = "(?i)" + pattern
		}
		re, err = regexp.Compile(pattern)
	}

	if err != nil {
		e.showMsg(fmt.Sprintf("regex error: %v", err))
		return
	}

	// If selection exists, replace only within selection
	start := 0
	end := e.doc.Len()
	if e.selStart >= 0 {
		start, end = e.getSelectionRange()
	}

	subContent := content[start:end]
	// Count replacements
	matches := re.FindAllStringIndex(subContent, -1)
	if len(matches) == 0 {
		e.showMsg("no matches found for replacement")
		return
	}

	// Replace all
	newContent := re.ReplaceAllString(subContent, e.replaceNew)
	fullNew := content[:start] + newContent + content[end:]

	// Record undo step
	e.undo.Push(document.UndoReplace, 0, []rune(fullNew), []rune(content))
	e.doc = document.NewDocumentFromString(fullNew)
	e.cursor = 0
	e.clampCursor()
	e.showMsg(fmt.Sprintf("replaced %d occurrences", len(matches)))
}

// ---- Command Execution ----

func (e *Editor) executeCommand() {
	cmdStr := strings.TrimSpace(e.cmdBuffer)
	e.cmdBuffer = ""
	e.mode = ModeInsert

	if cmdStr == "" {
		return
	}

	cmd, err := command.Parse(cmdStr)
	if err != nil {
		e.showMsg(fmt.Sprintf("command error: %v", err))
		return
	}

	exec := &command.Executor{
		DeleteRange: func(startLine, endLine int) (string, error) {
			startPos := e.doc.LineColToPos(startLine, 0)
			endLineLen := e.doc.LineLen(endLine)
			endPos := e.doc.LineColToPos(endLine, endLineLen)
			if endPos < e.doc.Len() {
				endPos++ // include newline in range
			}
			if endPos > e.doc.Len() {
				endPos = e.doc.Len()
			}
			deleted := e.doc.DeleteRange(startPos, endPos-startPos)
			e.undo.Push(document.UndoDelete, startPos, nil, deleted)
			e.cursor = startPos
			e.clampCursor()
			return string(deleted), nil
		},
		ReplaceText: func(oldPattern, newText string, startLine, endLine int, global, ignoreCase bool) (int, error) {
			content := e.doc.Content()
			pattern := oldPattern
			if ignoreCase {
				pattern = "(?i)" + pattern
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				return 0, fmt.Errorf("regex: %w", err)
			}

			// Calculate positions for the line range
			startPos := e.doc.LineColToPos(startLine, 0)
			endLineLen := e.doc.LineLen(endLine)
			endPos := e.doc.LineColToPos(endLine, endLineLen)
			if endPos < e.doc.Len() {
				endPos++
			}

			subContent := content[startPos:endPos]
			matches := re.FindAllStringIndex(subContent, -1)
			count := len(matches)
			if count == 0 {
				return 0, nil
			}

			var newSub string
			if global || count == 1 {
				newSub = re.ReplaceAllString(subContent, newText)
			} else {
				newSub = re.ReplaceAllString(subContent, newText)
			}

			fullNew := content[:startPos] + newSub + content[endPos:]
			e.undo.Push(document.UndoReplace, 0, []rune(fullNew), []rune(content))
			e.doc = document.NewDocumentFromString(fullNew)
			e.cursor = startPos
			e.clampCursor()
			return count, nil
		},
		WriteFile: func(path string) error {
			return e.SaveFile(path)
		},
		Quit: func(force bool) error {
			return e.Quit(force)
		},
		OpenFile: func(path string) error {
			return e.OpenFile(path)
		},
		GetLineCount: func() int {
			return e.doc.LineCount()
		},
		SetOption: func(name, value string) error {
			switch name {
			case "nu", "number":
				e.config.ShowLineNum = value != "0" && value != "false"
			case "nowrap":
				e.config.SoftWrap = false
			case "wrap":
				e.config.SoftWrap = true
			case "tabwidth", "ts":
				fmt.Sscanf(value, "%d", &e.config.TabWidth)
			default:
				return fmt.Errorf("unknown option: %s", name)
			}
			return nil
		},
		ShowMsg: func(msg string) {
			e.showMsg(msg)
		},
	}

	shouldQuit, err := exec.Execute(cmd)
	if err != nil {
		e.showMsg(fmt.Sprintf("error: %v", err))
	}
	if shouldQuit {
		e.running = false
	}
}

// ---- Cursor Movement ----

func (e *Editor) moveCursor(delta int) {
	e.cursor += delta
	if e.selecting {
		e.selEnd = e.cursor
	}
	e.clampCursor()
}

func (e *Editor) moveUp() {
	line, col := e.doc.PosToLineCol(e.cursor)
	if line <= 0 {
		return
	}
	line--
	e.cursor = e.doc.LineColToPos(line, col)
	e.clampCursor()
}

func (e *Editor) moveDown() {
	line, col := e.doc.PosToLineCol(e.cursor)
	if line >= e.doc.LineCount()-1 {
		return
	}
	line++
	e.cursor = e.doc.LineColToPos(line, col)
	e.clampCursor()
}

func (e *Editor) moveLineStart() {
	line, _ := e.doc.PosToLineCol(e.cursor)
	e.cursor = e.doc.LineColToPos(line, 0)
	e.clampCursor()
}

func (e *Editor) moveLineEnd() {
	line, _ := e.doc.PosToLineCol(e.cursor)
	e.cursor = e.doc.LineColToPos(line, e.doc.LineLen(line))
	e.clampCursor()
}

func (e *Editor) movePageUp() {
	line, col := e.doc.PosToLineCol(e.cursor)
	line -= e.termH / 2
	if line < 0 {
		line = 0
	}
	e.cursor = e.doc.LineColToPos(line, col)
	e.offsetRow -= e.termH / 2
	if e.offsetRow < 0 {
		e.offsetRow = 0
	}
	e.clampCursor()
}

func (e *Editor) movePageDown() {
	line, col := e.doc.PosToLineCol(e.cursor)
	line += e.termH / 2
	if line >= e.doc.LineCount() {
		line = e.doc.LineCount() - 1
	}
	e.cursor = e.doc.LineColToPos(line, col)
	e.offsetRow += e.termH / 2
	e.clampCursor()
}

// clampCursor ensures cursor and viewport are within valid range.
func (e *Editor) clampCursor() {
	if e.cursor < 0 {
		e.cursor = 0
	}
	if e.cursor > e.doc.Len() {
		e.cursor = e.doc.Len()
	}

	// Update cursor row/col
	line, col := e.doc.PosToLineCol(e.cursor)
	e.cursorRow = line
	// Convert rune index to visual column (handles wide chars like CJK)
	lineText := e.doc.Line(line)
	runes := []rune(lineText)
	visCol := 0
	for i := 0; i < col && i < len(runes); i++ {
		if runes[i] == '\t' {
			visCol += e.tabVisualWidth(visCol)
		} else {
			w := runewidth.RuneWidth(runes[i])
			if w <= 0 {
				w = 1
			}
			visCol += w
		}
	}
	e.cursorCol = visCol
}

// lineNumWidth returns the width of the line number area.
func (e *Editor) lineNumWidth() int {
	if !e.config.ShowLineNum {
		return 0
	}
	lines := e.doc.LineCount()
	digits := 1
	for lines >= 10 {
		digits++
		lines /= 10
	}
	return digits + 2 // numbers + "│ "
}

// ---- Rendering ----

// rebuildRuneColors builds a foreground color map for each rune position from syntax tokens.
func (e *Editor) rebuildRuneColors(content string) {
	contentRunes := []rune(content)
	e.runeColors = make([]tcell.Color, len(contentRunes))
	for i := range e.runeColors {
		e.runeColors[i] = tcell.ColorDefault
	}
	if e.hlTokens == nil {
		return
	}

	// Tokens are sequential; accumulate character positions to map colors
	pos := 0
	for _, token := range e.hlTokens {
		tokenRunes := len([]rune(token.Value))
		ts := e.hl.GetTokenStyle(token.Type)
		if ts.Fg != tcell.ColorDefault {
			for i := 0; i < tokenRunes && pos+i < len(e.runeColors); i++ {
				e.runeColors[pos+i] = ts.Fg
			}
		}
		pos += tokenRunes
	}
}

func (e *Editor) render() {
	e.screen.Clear()

	// Syntax highlight tokenize (only recompute when content changes)
	content := e.doc.Content()
	if content != e.lastContent {
		e.lastContent = content
		e.hlTokens = e.hl.Tokenize(content)
		e.rebuildRuneColors(content)
	}

	lineNumWidth := e.lineNumWidth()
	textWidth := e.termW - lineNumWidth
	if textWidth < 1 {
		textWidth = 1
	}

	for y := 0; y < e.termH-1; y++ {
		docLine := y + e.offsetRow
		if docLine >= e.doc.LineCount() {
			break
		}

		lineText := e.doc.Line(docLine)
		runes := []rune(lineText)

		// Line numbers
		if e.config.ShowLineNum {
			numStr := fmt.Sprintf("%*d ", e.lineNumWidth()-2, docLine+1)
			for i, ch := range numStr {
				e.screen.SetContent(i, y, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorDefault))
			}
		}

		// Syntax highlight (use pre-cached tokens)
		if e.hl != nil && e.hlTokens != nil {
			e.renderHighlightedLine(runes, docLine, lineNumWidth, y, textWidth)
		} else {
			e.renderPlainLine(runes, docLine, lineNumWidth, y, textWidth)
		}
	}

	// Status bar
	e.renderStatusBar()

	// Cursor position
	cursorX := lineNumWidth + e.cursorCol - e.offsetCol
	cursorY := e.cursorRow - e.offsetRow

	if cursorX < 0 {
		cursorX = 0
	}
	if cursorX >= e.termW {
		cursorX = e.termW - 1
	}
	if cursorY >= 0 && cursorY < e.termH-1 {
		e.screen.ShowCursor(cursorX, cursorY)
	}
}

// renderPlainLine renders a plain text line.
func (e *Editor) renderPlainLine(runes []rune, docLine, lineNumWidth, y, textWidth int) {
	style := tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault)

	idx := e.offsetCol
	for screenX := 0; screenX < textWidth && idx < len(runes); {
		ch := runes[idx]
		var w int
		if ch == '\t' {
			w = e.tabVisualWidth(screenX)
		} else {
			w = runewidth.RuneWidth(ch)
		}
		if w <= 0 {
			w = 1
		}
		if screenX+w > textWidth {
			break
		}

		// Check selection
		globalPos := e.doc.LineColToPos(docLine, idx)
		if e.selStart >= 0 && globalPos >= e.selStart && globalPos <= e.selEnd {
			style = style.Background(tcell.ColorNavy)
		} else {
			style = style.Background(tcell.ColorDefault)
		}

		e.screen.SetContent(lineNumWidth+screenX, y, ch, nil, style)
		screenX += w
		idx++
	}
}

// renderHighlightedLine renders a syntax highlighted line.
// Uses pre-cached runeColors from render() for coloring.
func (e *Editor) renderHighlightedLine(runes []rune, docLine, lineNumWidth, y, textWidth int) {
	lineStart := e.doc.LineColToPos(docLine, 0)

	idx := e.offsetCol
	for screenX := 0; screenX < textWidth && idx < len(runes); {
		ch := runes[idx]
		var w int
		if ch == '\t' {
			w = e.tabVisualWidth(screenX)
		} else {
			w = runewidth.RuneWidth(ch)
		}
		if w <= 0 {
			w = 1
		}
		if screenX+w > textWidth {
			break
		}

		globalPos := lineStart + idx

		// Syntax color
		fg := tcell.ColorDefault
		if globalPos >= 0 && globalPos < len(e.runeColors) {
			fg = e.runeColors[globalPos]
		}
		style := tcell.StyleDefault.Foreground(fg).Background(tcell.ColorDefault)

		// Check selection
		sel := false
		if e.selStart >= 0 && globalPos >= e.selStart && globalPos <= e.selEnd {
			sel = true
			style = style.Background(tcell.ColorNavy)
		}

		// Bracket match highlight (check only cursor bracket)
		if !sel && globalPos == e.cursor && highlight.IsBracket(ch) && e.isMatchingBracket(globalPos) {
			style = style.Foreground(tcell.ColorGreen).Background(tcell.ColorDarkCyan)
		}

		// Highlight matching bracket when cursor is on a bracket
		if !sel && globalPos != e.cursor && highlight.IsBracket(ch) {
			// Check if cursor is on a bracket, and current char pairs with cursor bracket
			cursorCh := e.doc.Buffer.RuneAt(e.cursor)
			if cursorCh != 0 && highlight.IsBracket(cursorCh) {
				match := highlight.MatchingBracket(cursorCh)
				if ch == match {
					// Check if it's actually a matching pair
					if e.isMatchingBracket(e.cursor) {
						// Verify current char pairs with cursor bracket
						match2 := highlight.MatchingBracket(ch)
						if match2 == cursorCh {
							style = style.Foreground(tcell.ColorGreen).Background(tcell.ColorDarkCyan)
						}
					}
				}
			}
		}

		// Search match highlight
		if e.searchResults != nil && !sel {
			for _, pos := range e.searchResults {
				if globalPos >= pos && globalPos < pos+len([]rune(e.searchQuery)) {
					style = style.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack)
					break
				}
			}
		}

		e.screen.SetContent(lineNumWidth+screenX, y, ch, nil, style)
		screenX += w
		idx++
	}
}

// tabVisualWidth calculates the visual width of a tab at position screenX.
func (e *Editor) tabVisualWidth(screenX int) int {
	tw := e.config.TabWidth
	if tw <= 0 {
		tw = 4
	}
	return tw - (screenX % tw)
}

// isMatchingBracket checks if a matching bracket exists at position pos.
func (e *Editor) isMatchingBracket(pos int) bool {
	ch := e.doc.Buffer.RuneAt(pos)
	if ch == 0 {
		return false
	}
	match := highlight.MatchingBracket(ch)
	if match == 0 {
		return false
	}

	// Scan forward or backward
	direction := 1
	if ch == ')' || ch == ']' || ch == '}' {
		direction = -1
	}

	depth := 0
	p := pos + direction
	for p >= 0 && p < e.doc.Len() {
		c := e.doc.Buffer.RuneAt(p)
		if c == ch {
			depth++
		} else if c == match {
			if depth == 0 {
				return true
			}
			depth--
		}
		p += direction
	}
	return false
}

// renderStatusBar renders the status bar.
func (e *Editor) renderStatusBar() {
	y := e.termH - 1

	var modeStr string
	switch e.mode {
	case ModeInsert:
		modeStr = "INSERT"
	case ModeCommand:
		modeStr = "COMMAND"
	case ModeSearch:
		modeStr = "SEARCH"
	case ModeReplace:
		modeStr = "REPLACE"
	}

	line, col := e.doc.PosToLineCol(e.cursor)
	fileName := filepath.Base(e.doc.FilePath)
	if fileName == "" || fileName == "." {
		fileName = "[No Name]"
	}
	modified := ""
	if e.doc.Modified {
		modified = " ●"
	}
	undoCount := e.undo.UndoCount()
	redoCount := e.undo.RedoCount()

	// Build status bar text
	var statusText string

	switch e.mode {
	case ModeCommand:
		statusText = ":" + e.cmdBuffer
		// Render command bar
		for x := 0; x < e.termW && x < len([]rune(statusText)); x++ {
			ch := []rune(statusText)[x]
			e.screen.SetContent(x, y, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
		}
		// Fill remaining with spaces
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
		for x := len([]rune(statusText)); x < e.termW; x++ {
			e.screen.SetContent(x, y, ' ', nil, style)
		}
		return

	case ModeSearch:
		prefix := "Search: "
		if e.searchRegex {
			prefix = "Regex: "
		}
		statusText = prefix + e.cmdBuffer
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
		for x := 0; x < e.termW && x < len([]rune(statusText)); x++ {
			ch := []rune(statusText)[x]
			e.screen.SetContent(x, y, ch, nil, style)
		}
		for x := len([]rune(statusText)); x < e.termW; x++ {
			e.screen.SetContent(x, y, ' ', nil, style)
		}
		return

	case ModeReplace:
		statusText = "Replace: " + e.cmdBuffer
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
		for x := 0; x < e.termW && x < len([]rune(statusText)); x++ {
			ch := []rune(statusText)[x]
			e.screen.SetContent(x, y, ch, nil, style)
		}
		for x := len([]rune(statusText)); x < e.termW; x++ {
			e.screen.SetContent(x, y, ' ', nil, style)
		}
		return
	}

	// Message mode: display message
	if e.message != "" && time.Since(e.messageTime) < 5*time.Second {
		statusText = e.message
	} else {
		statusText = fmt.Sprintf("%s%s  %s  %d:%d  U:%d R:%d",
			modeStr, modified, fileName, line+1, col+1, undoCount, redoCount)
	}

	// Render status bar
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
	runes := []rune(statusText)
	for x := 0; x < e.termW; x++ {
		ch := rune(' ')
		if x < len(runes) {
			ch = runes[x]
		}
		e.screen.SetContent(x, y, ch, nil, style)
	}
}
