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
			{`//[^\n]*`, chroma.Comment, nil},
			{`/\*[\s\S]*?\*/`, chroma.CommentMultiline, nil},
			{`"(?:\\.|[^"\\])*"`, chroma.String, nil},
			{`\b0x[0-9a-fA-F]+\b`, chroma.NumberHex, nil},
			{`\b\d+(\.\d+)?\b`, chroma.Number, nil},
			{`\b(fn|func|var|let|const|if|else|for|in|return|while|loop)\b`, chroma.Keyword, nil},
			{`\b(true|false|null)\b`, chroma.KeywordConstant, nil},
			{`\b(import|export|extern|struct|interface|type|as|is|match|case|break|continue)\b`, chroma.Keyword, nil},
			{`\b(i8|i16|i32|i64|u8|u16|u32|u64|f32|f64|str|bool|any|ptr|array|map|option|result)\b`, chroma.KeywordType, nil},
			{`\b(print|println|len|range|panic|assert|sizeof|typeof)\b`, chroma.NameBuiltin, nil},
			{`\+\+|--|->|=>|\.\.|\?\.`, chroma.Operator, nil},
			{`[+\-*/%=<>!&|^~]+`, chroma.Operator, nil},
			{`\b([a-zA-Z_]\w*)\s*\(`, chroma.NameFunction, nil},
			{`\b[a-zA-Z_]\w*\b`, chroma.Name, nil},
			{`[{}()\[\],;:.]`, chroma.Punctuation, nil},
			{`\s+`, chroma.Text, nil},
		},
	}
})
