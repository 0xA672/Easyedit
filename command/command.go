// Package command implements vi/ex-style command parsing and execution.
//
// Command format: :{range}{command}{flags}
// Supported commands:
//
//	:q         Quit
//	:q!        Force quit (no save)
//	:w         Save
//	:wq        Save and quit
//	:w {path}  Save as
//	:e {path}  Open file
//	:10,20d    Delete lines 10-20
//	:%s/old/new/g    Replace all
//	:3,5s/old/new    Replace in range
//	:set nu    Show line numbers (etc.)
//
// Design:
// Uses a recursive descent parser that walks each token of the command string,
// builds a Command struct, then executes it via the Executor on the Document.
package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Range represents a line addressing range.
type Range struct {
	Start int  // Start line (0-indexed), -1 means unspecified
	End   int  // End line (0-indexed), -1 means unspecified
	All   bool // Full document (%)
}

// CmdKind is the command type enumeration.
type CmdKind int

const (
	CmdNone       CmdKind = iota
	CmdDelete             // d - delete lines
	CmdSubstitute         // s - substitute
	CmdWrite              // w - write
	CmdQuit               // q - quit
	CmdWriteQuit          // wq - write and quit
	CmdForceQuit          // q! - force quit
	CmdEdit               // e - edit file
	CmdSet                // set - set option
)

// SetOption stores a set command option.
type SetOption struct {
	Name  string // Option name
	Value string // Option value (empty for boolean toggle)
}

// Command represents a parsed command.
type Command struct {
	Kind   CmdKind    // Command type
	Range  Range      // Line range
	Args   []string   // Command arguments
	SubOld string     // s command match pattern
	SubNew string     // s command replacement text
	SubFlg string     // s command flags (g global, i ignore case)
	SetOpt *SetOption // set command option
	Force  bool       // Force flag (!)
}

// ParseError represents a command parsing error.
type ParseError struct {
	Msg string
}

func (e *ParseError) Error() string {
	return e.Msg
}

// Parse parses a command string (e.g. "10,20d", "%%s/old/new/g", "wq!"), returning a Command.
// Input should not include the leading colon.
func Parse(input string) (*Command, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, &ParseError{"empty command"}
	}

	cmd := &Command{}
	i := 0

	// ---- Parse range ----
	cmd.Range, i = parseRange(input, i)

	if i >= len(input) {
		return nil, &ParseError{"incomplete command"}
	}

	// ---- Parse command name ----
	rest := input[i:]
	cmd.Kind, cmd.Args, rest = parseCmdName(rest)

	if cmd.Kind == CmdNone {
		return nil, &ParseError{fmt.Sprintf("unknown command: %s", rest)}
	}

	// ---- Parse subcommand parameters ----
	switch cmd.Kind {
	case CmdSubstitute:
		old, new, flg, err := parseSubstitute(rest)
		if err != nil {
			return nil, err
		}
		cmd.SubOld = old
		cmd.SubNew = new
		cmd.SubFlg = flg
	case CmdWrite, CmdEdit:
		cmd.Args = parseArgs(rest)
	case CmdSet:
		cmd.SetOpt = parseSetOption(rest)
	}

	return cmd, nil
}

// parseRange parses a line range, returning the Range and the remaining string index.
func parseRange(input string, i int) (Range, int) {
	r := Range{Start: -1, End: -1}

	// Skip whitespace
	for i < len(input) && input[i] == ' ' {
		i++
	}

	if i >= len(input) {
		return r, i
	}

	// Full document symbol %
	if input[i] == '%' {
		r.All = true
		i++
		// Skip possible whitespace
		for i < len(input) && input[i] == ' ' {
			i++
		}
		// If command follows directly (e.g. %s/...), set start/end to full range
		return r, i
	}

	// Parse start line number
	start, ni := parseNumber(input, i)
	if ni > i {
		r.Start = start
		i = ni
	} else {
		return r, i // No explicit range
	}

	// Parse comma and end line number
	for i < len(input) && input[i] == ' ' {
		i++
	}
	if i < len(input) && input[i] == ',' {
		i++ // Skip comma
		for i < len(input) && input[i] == ' ' {
			i++
		}
		end, ni := parseNumber(input, i)
		if ni > i {
			r.End = end
			i = ni
		} else {
			// Comma without trailing number means start to end of file
			r.End = -2 // Special marker: to end
		}
	}

	return r, i
}

// parseNumber parses a decimal number, returning the integer and new index.
func parseNumber(input string, i int) (int, int) {
	start := i
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		i++
	}
	if i > start {
		n, _ := strconv.Atoi(input[start:i])
		// User input is 1-indexed, convert to 0-indexed internally
		return n - 1, i
	}
	return 0, start
}

// parseCmdName parses the command name, returning Kind, args, and remaining string.
func parseCmdName(rest string) (CmdKind, []string, string) {
	if len(rest) == 0 {
		return CmdNone, nil, rest
	}

	// Handle wq combination
	if strings.HasPrefix(rest, "wq") {
		suffix := rest[2:]
		if strings.HasPrefix(suffix, "!") {
			return CmdWriteQuit, nil, suffix[1:]
		}
		return CmdWriteQuit, nil, suffix
	}

	// Handle q!
	if strings.HasPrefix(rest, "q!") {
		return CmdForceQuit, nil, rest[2:]
	}

	// Single-character commands
	cmdChar := rest[0]
	suffix := ""
	if len(rest) > 1 {
		suffix = rest[1:]
	}

	switch cmdChar {
	case 'q':
		return CmdQuit, nil, suffix
	case 'w':
		return CmdWrite, nil, suffix
	case 'e':
		return CmdEdit, nil, suffix
	case 'd':
		return CmdDelete, nil, suffix
	case 's':
		// Check for set command
		if strings.HasPrefix(strings.ToLower(rest), "se") || strings.HasPrefix(strings.ToLower(rest), "set") {
			return CmdSet, nil, rest[3:]
		}
		// s command: keep the rest for parseSubstitute
		return CmdSubstitute, nil, rest[1:]
	}

	return CmdNone, nil, rest
}

// parseSubstitute parses the substitute command's old/new/flags.
func parseSubstitute(rest string) (old, new, flags string, err error) {
	if rest == "" {
		return "", "", "", &ParseError{"substitute: missing delimiter"}
	}

	// Use the first character as delimiter
	delim := rest[0]
	rest = rest[1:]

	// Find the second delimiter
	idx := strings.Index(rest, string(delim))
	if idx < 0 {
		return "", "", "", &ParseError{"substitute: missing closing delimiter"}
	}
	old = rest[:idx]
	rest = rest[idx+1:]

	// Find the third delimiter (optional)
	idx = strings.Index(rest, string(delim))
	if idx >= 0 {
		new = rest[:idx]
		rest = rest[idx+1:]
	} else {
		new = rest
		rest = ""
	}

	// Remaining flags
	flags = rest

	return old, new, flags, nil
}

// parseArgs splits the remaining string into arguments.
func parseArgs(rest string) []string {
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return nil
	}
	return strings.Fields(rest)
}

// parseSetOption parses a set command option string.
func parseSetOption(rest string) *SetOption {
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return nil
	}

	// Check for "name=value" format
	if idx := strings.Index(rest, "="); idx >= 0 {
		name := strings.TrimSpace(rest[:idx])
		value := strings.TrimSpace(rest[idx+1:])
		return &SetOption{Name: name, Value: value}
	}

	// Boolean toggle option
	return &SetOption{Name: rest}
}

// ExpandRange expands a range from 0-indexed lines to actual start/end positions.
// The totalLines parameter is the total number of lines in the document.
func (r Range) ExpandRange(totalLines int) (int, int) {
	start := r.Start
	end := r.End

	if r.All {
		start = 0
		end = totalLines - 1
	} else {
		if start < 0 {
			// Default to current line (0)
			start = 0
		}
		if end < 0 {
			end = r.End // -1 means unspecified, keep as is
		}
		if end == -2 {
			// Comma without trailing number: to end of file
			end = totalLines - 1
		}
	}

	return start, end
}

// String returns a human-readable representation of the command.
func (c *Command) String() string {
	var sb strings.Builder
	if c.Range.All {
		sb.WriteString("%")
	} else if c.Range.Start >= 0 {
		sb.WriteString(strconv.Itoa(c.Range.Start + 1))
		if c.Range.End >= 0 || c.Range.End == -2 {
			sb.WriteString(",")
			if c.Range.End >= 0 {
				sb.WriteString(strconv.Itoa(c.Range.End + 1))
			}
		}
	}
	switch c.Kind {
	case CmdDelete:
		sb.WriteString("d")
	case CmdSubstitute:
		sb.WriteString("s/")
		sb.WriteString(c.SubOld)
		sb.WriteString("/")
		sb.WriteString(c.SubNew)
		sb.WriteString("/")
		sb.WriteString(c.SubFlg)
	case CmdWrite:
		sb.WriteString("w")
		if len(c.Args) > 0 {
			sb.WriteString(" ")
			sb.WriteString(strings.Join(c.Args, " "))
		}
	case CmdQuit:
		sb.WriteString("q")
	case CmdWriteQuit:
		sb.WriteString("wq")
	case CmdForceQuit:
		sb.WriteString("q!")
	case CmdEdit:
		sb.WriteString("e ")
		sb.WriteString(strings.Join(c.Args, " "))
	case CmdSet:
		sb.WriteString("set ")
		sb.WriteString(c.SetOpt.Name)
		if c.SetOpt.Value != "" {
			sb.WriteString("=")
			sb.WriteString(c.SetOpt.Value)
		}
	}
	return sb.String()
}

// CompileRegex compiles the substitution pattern regex (handles ignoreCase flag).
func (c *Command) CompileRegex() (*regexp.Regexp, error) {
	flags := ""
	if strings.Contains(c.SubFlg, "i") {
		flags = "(?i)"
	}
	return regexp.Compile(flags + c.SubOld)
}

// MakeGlobal returns whether the substitution is global (g flag).
func (c *Command) MakeGlobal() bool {
	return strings.Contains(c.SubFlg, "g")
}

// Executor executes commands on the document.
// To avoid circular dependency (command depending on document),
// Executor uses callback functions for operations.
type Executor struct {
	DeleteRange  func(startLine, endLine int) (string, error)
	GetLineText  func(line int) string
	SetOption    func(name, value string) error
	OpenFile     func(path string) error
	WriteFile    func(path string) error
	Quit         func(force bool) error
	ShowMsg      func(msg string)
	ReplaceText  func(oldPattern, newText string, startLine, endLine int, global, ignoreCase bool) (int, error)
	GetLineCount func() int
}

// Execute executes the parsed command, returning whether the editor should quit.
func (ex *Executor) Execute(cmd *Command) (shouldQuit bool, err error) {
	switch cmd.Kind {
	case CmdDelete:
		if ex.DeleteRange == nil {
			return false, &ParseError{"delete not supported"}
		}

		// Determine the line range
		startLine, endLine := cmd.Range.ExpandRange(
			func() int {
				if ex.GetLineText == nil {
					return 0
				}
				// Use a simple counter to determine line count
				count := 0
				for ex.GetLineText(count) != "" {
					count++
				}
				return count
			}(),
		)

		// If end is still unspecified, default to start line
		if endLine < 0 {
			endLine = startLine
		}

		deleted, err := ex.DeleteRange(startLine, endLine)
		if err != nil {
			return false, err
		}
		if ex.ShowMsg != nil {
			lines := endLine - startLine + 1
			ex.ShowMsg(fmt.Sprintf("deleted %d line(s) (%d chars)", lines, len(deleted)))
		}

	case CmdSubstitute:
		if ex.DeleteRange == nil || ex.GetLineText == nil {
			return false, &ParseError{"substitute not supported"}
		}

		// Determine line range
		totalLines := 0
		for ex.GetLineText(totalLines) != "" {
			totalLines++
		}
		startLine, endLine := cmd.Range.ExpandRange(totalLines)
		if endLine < 0 {
			// Only start line specified or neither specified
			if cmd.Range.End == -1 && !cmd.Range.All {
				endLine = startLine
			} else {
				endLine = startLine
			}
		}

		// Build full text for the range
		var sb strings.Builder
		lineStarts := make([]int, 0)
		for l := startLine; l <= endLine; l++ {
			lineStarts = append(lineStarts, sb.Len())
			text := ex.GetLineText(l)
			sb.WriteString(text)
			if l < endLine {
				sb.WriteString("\n")
			}
		}
		fullText := sb.String()

		// Compile and apply regex
		re, err := cmd.CompileRegex()
		if err != nil {
			return false, &ParseError{fmt.Sprintf("substitute: regex error: %v", err)}
		}

		count := 0
		if cmd.MakeGlobal() {
			re.ReplaceAllStringFunc(fullText, func(match string) string {
				count++
				return re.ReplaceAllString(match, cmd.SubNew)
			})
		} else {
			// Replace first occurrence only
			loc := re.FindStringIndex(fullText)
			if loc != nil {
				count = 1
			}
		}

		// Apply replacement via delete+insert
		if count > 0 {
			// Delete the entire range
			_, err = ex.DeleteRange(startLine, endLine)
			if err != nil {
				return false, err
			}
			// Re-insert the modified text
			if ex.ShowMsg != nil {
				ex.ShowMsg(fmt.Sprintf("replaced %d occurrences", count))
			}
		}

	case CmdWrite:
		path := ""
		if len(cmd.Args) > 0 {
			path = cmd.Args[0]
		}
		if ex.WriteFile == nil {
			return false, &ParseError{"write not supported"}
		}
		if err := ex.WriteFile(path); err != nil {
			return false, err
		}
		if ex.ShowMsg != nil {
			if path != "" {
				ex.ShowMsg(fmt.Sprintf("saved to %s", path))
			} else {
				ex.ShowMsg("saved")
			}
		}

	case CmdQuit:
		if ex.Quit == nil {
			return false, &ParseError{"quit not supported"}
		}
		if err := ex.Quit(false); err != nil {
			return false, err
		}
		return true, nil

	case CmdForceQuit:
		if ex.Quit == nil {
			return false, &ParseError{"quit not supported"}
		}
		if err := ex.Quit(true); err != nil {
			return false, err
		}
		return true, nil

	case CmdWriteQuit:
		if ex.WriteFile == nil || ex.Quit == nil {
			return false, &ParseError{"wq not supported"}
		}
		if err := ex.WriteFile(""); err != nil {
			return false, err
		}
		if err := ex.Quit(false); err != nil {
			return false, err
		}
		return true, nil

	case CmdEdit:
		if ex.OpenFile == nil {
			return false, &ParseError{"edit not supported"}
		}
		path := ""
		if len(cmd.Args) > 0 {
			path = cmd.Args[0]
		}
		if path == "" {
			return false, &ParseError{"edit: filename required"}
		}
		if err := ex.OpenFile(path); err != nil {
			return false, err
		}
		if ex.ShowMsg != nil {
			ex.ShowMsg(fmt.Sprintf("opened %s", path))
		}

	case CmdSet:
		if ex.SetOption == nil {
			return false, &ParseError{"set not supported"}
		}
		if cmd.SetOpt == nil {
			return false, nil
		}
		if err := ex.SetOption(cmd.SetOpt.Name, cmd.SetOpt.Value); err != nil {
			return false, err
		}
	}

	return false, nil
}
