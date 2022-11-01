package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Alternative
// ----------------------------------------------------------------

// Alt :: p[a] -> p[b] -> p[c] -> ... -> p[a|b|c...]
func Alt(ps ...Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var xs []Result
		var err *Error
		var success bool
		for _, p := range ps {
			out := p.Parse(toks)
			err = betterError(err, out.Error)
			if out.Success {
				xs = append(xs, out.Candidates...)
				success = true
			}
		}
		return resultOrError(xs, err, success)
	})
}

// Opt :: p[a] -> p[a|nil]
func Opt(p Parser) Parser {
	return Alt(p, Nil())
}

// OptSc :: p[a] -> p[a|nil]
func OptSc(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		out := p.Parse(toks)
		if out.Success {
			return out
		}
		return successX([]Result{{toks: toks}}, out.Error)
	})
}
