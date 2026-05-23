// Package highlight implements syntax highlighting.
//
// Uses the chroma lexer to automatically select the appropriate lexer
// based on file extension, tokenizing text and mapping tokens to terminal colors.
//
// Design:
// The chroma library provides lexing and styling. We create our own style (easyedit style),
// mapped to tcell colors. Tokenized results are cached and only re-parsed when the document changes.
//
// Supported file types:
// - .py  Python
// - .js / .jsx / .mjs  JavaScript
// - .ts / .tsx  TypeScript
// - .json  JSON
// - .sh / .bash  Shell
// - .html / .htm  HTML
// - .css  CSS
// - .md / .markdown  Markdown
// - .go  Go
// - .rs  Rust
// - .yaml / .yml  YAML
// - .java  Java
// - .c / .h  C
// - .cpp / .hpp / .cc  C++
// - .rb  Ruby
package highlight

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gdamore/tcell/v2"
)

// TokenStyle stores the foreground and background colors for a token.
type TokenStyle struct {
	Fg tcell.Color
	Bg tcell.Color
}

// Highlighter manages syntax highlighting.
type Highlighter struct {
	lexer    chroma.Lexer   // chroma lexer
	style    *chroma.Style  // chroma style
	tokens   []chroma.Token // cached tokenization result
	lastText string         // text from last tokenization, for cache comparison
}

// NewHighlighter creates a highlighter based on file path.
func NewHighlighter(filePath string) *Highlighter {
	h := &Highlighter{}
	h.initLexer(filePath)
	h.style = styles.Get("monokai")
	if h.style == nil {
		// Fallback to built-in style
		h.style = styles.Fallback
	}
	return h
}

// initLexer initializes the lexer based on file extension.
func (h *Highlighter) initLexer(filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))
	// Look up lexer by extension
	h.lexer = lexers.Get(ext)
	if h.lexer == nil {
		// Try looking up by filename
		h.lexer = lexers.Get(filepath.Base(filePath))
	}
}

// SetFile updates the lexer when switching files.
func (h *Highlighter) SetFile(filePath string) {
	h.initLexer(filePath)
	h.tokens = nil
	h.lastText = ""
}

// Tokenize tokenizes the text and caches the result.
// Returns cached tokens if the text hasn't changed.
func (h *Highlighter) Tokenize(text string) []chroma.Token {
	if text == h.lastText && h.tokens != nil {
		return h.tokens
	}

	h.lastText = text

	if h.lexer == nil {
		// No matching lexer, return nil (plain text)
		h.tokens = nil
		return nil
	}

	// Perform lexing with chroma
	iterator, err := h.lexer.Tokenise(nil, text)
	if err != nil {
		h.tokens = nil
		return nil
	}

	h.tokens = iterator.Tokens()
	return h.tokens
}

// GetTokenStyle maps a chroma token type to terminal colors.
// Uses the monokai theme style mapping.
func (h *Highlighter) GetTokenStyle(t chroma.TokenType) TokenStyle {
	if h.style == nil {
		return TokenStyle{Fg: tcell.ColorDefault, Bg: tcell.ColorDefault}
	}

	entry := h.style.Get(t)
	// In chroma v2, StyleEntry is a value type; check if it's the zero value
	if entry.Colour.String() == "" && entry.Background.String() == "" {
		return TokenStyle{Fg: tcell.ColorDefault, Bg: tcell.ColorDefault}
	}

	fg := tcell.ColorDefault
	bg := tcell.ColorDefault

	if entry.Colour.String() != "" {
		fg = hexToTcell(entry.Colour.String())
	}
	if entry.Background.String() != "" {
		bg = hexToTcell(entry.Background.String())
	}

	return TokenStyle{Fg: fg, Bg: bg}
}

// hexToTcell converts a chroma color string (e.g. "#ff8800") to tcell.Color.
func hexToTcell(hex string) tcell.Color {
	// chroma format like "#ff8800"
	if strings.HasPrefix(hex, "#") {
		return tcell.GetColor(hex)
	}
	// Direct color name
	return tcell.GetColor(hex)
}

// IsBracket checks if a rune is a bracket character.
func IsBracket(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == '{' || ch == '}'
}

// MatchingBracket returns the matching bracket for a given bracket rune.
func MatchingBracket(ch rune) rune {
	switch ch {
	case '(':
		return ')'
	case ')':
		return '('
	case '[':
		return ']'
	case ']':
		return '['
	case '{':
		return '}'
	case '}':
		return '{'
	}
	return 0
}
