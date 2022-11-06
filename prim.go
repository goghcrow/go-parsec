package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Token Filters
// ----------------------------------------------------------------

func Nil() Parser {
	return parser(func(toks []*lexer.Token) Output {
		return success([]Result{{next: toks}})
	})
}

func Str(toMatch string) Parser {
	return parser(func(toks []*lexer.Token) Output {
		if len(toks) == 0 {
			return fail(unableToConsumeToken(eof))
		}
		if toks[0].Lexeme != toMatch {
			return fail(unableToConsumeToken(toks[0]))
		}
		// 消费 toks[0], toks[1:] 为剩余 token 序列
		return success([]Result{{toks[0], toks[1:]}})
	})
}

func Tok(toMatch lexer.TokenKind) Parser {
	return parser(func(toks []*lexer.Token) Output {
		if len(toks) == 0 {
			return fail(unableToConsumeToken(eof))
		}
		if toks[0].TokenKind != toMatch {
			return fail(unableToConsumeToken(toks[0]))
		}
		// 消费 toks[0], toks[1:] 为剩余 token 序列
		return success([]Result{{toks[0], toks[1:]}})
	})
}
