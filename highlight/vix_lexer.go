package highlight

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func init() {
	lexers.Register(vixLexer)
}

var vixLexer = chroma.MustNewLexer(&chroma.Config{
	Name:      "Vix",
	Aliases:   []string{"vix"},
	Filenames: []string{"*.vix", "*.vixl"},
	MimeTypes: []string{"text/x-vix"},
}, func() chroma.Rules {
	return chroma.Rules{
		"root": {
			{Pattern: `//[^\n]*`, Type: chroma.Comment},
			{Pattern: `/\*[\s\S]*?\*/`, Type: chroma.CommentMultiline},
			{Pattern: `"(?:\\.|[^"\\])*"`, Type: chroma.String},
			{Pattern: `\b0x[0-9a-fA-F]+\b`, Type: chroma.NumberHex},
			{Pattern: `\b\d+(\.\d+)?\b`, Type: chroma.Number},
			{Pattern: `\b(fn|func|var|let|const|if|else|for|in|return|while|loop)\b`, Type: chroma.Keyword},
			{Pattern: `\b(true|false|null)\b`, Type: chroma.KeywordConstant},
			{Pattern: `\b(import|export|extern|struct|interface|type|as|is|match|case|break|continue)\b`, Type: chroma.Keyword},
			{Pattern: `\b(i8|i16|i32|i64|u8|u16|u32|u64|f32|f64|str|bool|any|ptr|array|map|option|result)\b`, Type: chroma.KeywordType},
			{Pattern: `\b(print|println|len|range|panic|assert|sizeof|typeof)\b`, Type: chroma.NameBuiltin},
			{Pattern: `\+\+|--|->|=>|\.\.|\?\.`, Type: chroma.Operator},
			{Pattern: `[+\-*/%=<>!&|^~]+`, Type: chroma.Operator},
			{Pattern: `\b([a-zA-Z_]\w*)\s*\(`, Type: chroma.NameFunction},
			{Pattern: `\b[a-zA-Z_]\w*\b`, Type: chroma.Name},
			{Pattern: `[{}()\[\],;:.]`, Type: chroma.Punctuation},
			{Pattern: `\s+`, Type: chroma.Text},
		},
	}
})
