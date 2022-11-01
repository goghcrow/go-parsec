package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Token Filters
// ----------------------------------------------------------------

func Nil() Parser {
	return newParser(func(toks []*lexer.Token) Output {
		return success(Result{toks: toks})
	})
}

func Str(toMatch string) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		if len(toks) == 0 {
			return fail(unableToConsumeToken(eof))
		}
		if toks[0].Lexeme != toMatch {
			return fail(unableToConsumeToken(toks[0]))
		}
		return success(Result{toks[0], toks[1:]})
	})
}

func Tok(toMatch lexer.TokenKind) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		if len(toks) == 0 {
			return fail(unableToConsumeToken(eof))
		}
		if toks[0].TokenKind != toMatch {
			return fail(unableToConsumeToken(toks[0]))
		}
		return success(Result{toks[0], toks[1:]})
	})
}
