package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Error Recovering
// ----------------------------------------------------------------

// Err :: p[a] -> err -> p[a]
// Consumes x, When x fails, the error message is replaced by e.
func Err(p Parser, msg string) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		return fail(newError(branches.Loc, msg))
	})
}

// ErrDef :: p[a] -> err -> -> a -> p[a]
func ErrDef(p Parser, msg string, def interface{}) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		return successX([]Result{{def, toks}}, newError(branches.Loc, msg))
	})
}
