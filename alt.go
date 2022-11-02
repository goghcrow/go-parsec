package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Alternative
// ----------------------------------------------------------------

// Alt :: p[a] -> p[b] -> p[c] -> ... -> p[a|b|c...]
// 返回所有可能结果, 当 ps 全部失败时失败
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
		return newOutput(xs, err, success)
	})
}

// Opt :: p[a] -> p[a|nil]
func Opt(p Parser) Parser {
	return Alt(p, Nil())
}

// OptSc :: p[a] -> p[a|nil]
// Opt 返回两种结果, OptSc 返回一种结果, 只有 p 失败才返回 nil
func OptSc(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		out := p.Parse(toks)
		if out.Success {
			return out
		}
		// nil with error
		return successWithErr([]Result{{toks: toks}}, out.Error)
	})
}
