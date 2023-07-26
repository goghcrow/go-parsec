package example

import (
	"testing"

	"github.com/goghcrow/lexer"
)

func TestSimple(t *testing.T) {
	const (
		Number TokenKind = iota + 1
		Space
		Comma
	)
	var lex = lexer.BuildLexer[TokenKind](func(lex *lexer.Lexicon[TokenKind]) {
		lex.Regex(Number, "\\d+")
		lex.Regex(Space, "\\s+").Skip()
		lex.Str(Comma, ",").Skip()
	})
	toks := lex.MustLex("1, 2,3, 4,5")
	for _, tok := range toks {
		t.Logf("%s in %s", tok, tok.Pos)
	}
}

func TestUserDefinedOperators(t *testing.T) {
	lex := NewBuiltinLexer(userDefinedOperators)
	toks := lex.MustLex("(1 + 2) * 3 < 5 == false")
	for _, tok := range toks {
		t.Logf("%s in %s", tok, tok.Pos)
	}
}
