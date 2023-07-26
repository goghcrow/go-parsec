package lexer

import (
	"fmt"
	"strings"
	"testing"
)

type tokKind int

const (
	Number tokKind = iota
	Ident
	NumId
	Space
	Comma
)

func stroftk(k tokKind) string {
	return map[tokKind]string{
		Number: "<num>",
		Ident:  "<id>",
		NumId:  "<numid>",
		Space:  "<space>",
		Comma:  ",",
	}[k]
}

func TestLexer(t *testing.T) {
	for _, tt := range []struct {
		input  string
		expect string
		lexer  *Lexer[tokKind]
	}{
		{
			"123",
			"<num>/123",
			BuildLexer[tokKind](func(lex *Lexicon[tokKind]) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",")
			}),
		},
		{
			"123,456",
			"<num>/123🍌,/,🍌<num>/456",
			BuildLexer[tokKind](func(lex *Lexicon[tokKind]) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",")
			}),
		},
		{
			"123,456,789",
			"<num>/123🍌<num>/456🍌<num>/789",
			BuildLexer[tokKind](func(lex *Lexicon[tokKind]) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",").Skip()
			}),
		},
		{
			"123, abc, 456, def, ",
			"<num>/123🍌<id>/abc🍌<num>/456🍌<id>/def",
			BuildLexer[tokKind](func(lex *Lexicon[tokKind]) {
				lex.Regex(Number, "\\d+")
				lex.Regex(Ident, "[a-zA-Z]\\w*")
				lex.Regex(Space, "\\s+").Skip()
				lex.Str(Comma, ",").Skip()
			}),
		},
		{
			"123, abc, 456, def, ",
			"<numid>/123🍌<numid>/abc🍌<numid>/456🍌<numid>/def",
			BuildLexer[tokKind](func(lex *Lexicon[tokKind]) {
				lex.Regex(NumId, "\\d+|(?:[a-zA-Z]\\w*)")
				lex.Regex(Space, "\\s+").Skip()
				lex.Str(Comma, ",").Skip()
			}),
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			toks := tt.lexer.MustLex(tt.input)
			actual := fmtToks(toks)
			if actual != tt.expect {
				t.Errorf("expect %s actual %s", tt.expect, actual)
			}
		})
	}
}

func fmtToks(toks []*Token[tokKind]) string {
	xs := make([]string, len(toks))
	for i, t := range toks {
		xs[i] = fmt.Sprintf("%s/%s", stroftk(t.kind), t.lexeme)
	}
	return strings.Join(xs, "🍌")
}
