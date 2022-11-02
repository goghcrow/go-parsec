package example

import (
	"github.com/goghcrow/go-parsec/lexer"
	"testing"
)

func TestSimple(t *testing.T) {
	const (
		Number lexer.TokenKind = iota + 1
		Space
		Comma
	)
	var lex = lexer.BuildLexer(func(lex *lexer.Lexicon) {
		lex.Regex(Number, "\\d+")
		lex.Regex(Space, "\\s+").Skip()
		lex.Str(Comma, ",").Skip()
	})
	toks := lex.MustLex("1, 2,3, 4,5")
	for _, tok := range toks {
		t.Logf("%s in %s", tok, tok.Loc)
	}
}

func TestExample(t *testing.T) {
	lex := NewBuiltinLexer(userDefinedOperators)
	toks := lex.MustLex("(1 + 2) * 3 < 5 == false")
	for _, tok := range toks {
		t.Logf("%s in %s", tok, tok.Loc)
	}
}
