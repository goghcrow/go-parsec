package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// Apply :: p[a] -> (a -> b) -> p[b]
// 📢 the data structural of v is topological equivalent to syntax structural of p
func Apply(p Parser, f func(v interface{}) interface{}) Parser {
	return parser(func(toks []*lexer.Token) Output {
		out := p.Parse(toks)
		if !out.Success {
			return out
		}
		xs := make([]Result, len(out.Candidates))
		for i, x := range out.Candidates {
			xs[i] = Result{f(x.Val), x.next}
		}
		return successWithErr(xs, out.Error)
	})
}

// Lazy :: (() -> p[a]) -> p[a]
func Lazy(thunk func() Parser) Parser {
	return parser(func(toks []*lexer.Token) Output {
		return thunk().Parse(toks)
	})
}
