package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Ambiguity Resolving
// ----------------------------------------------------------------

// Amb :: p[a] -> p[list[a]]
// Consumes x and merge group result by consumed tokens.
func Amb(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		branches := p.Parse(toks)
		if !branches.Success {
			return branches
		}

		group := make(map[*lexer.Token][]Result)
		for _, r := range branches.Candidates {
			k := r.toks.mapKey()
			group[k] = append(group[k], r)
		}

		xs := make([]Result, 0, len(group))
		for _, rs := range group {
			xss := make([]interface{}, len(rs))
			for i, v := range rs {
				xss[i] = v.Val
			}
			xs = append(xs, Result{xss, rs[0].toks})
		}
		return successX(xs, branches.Error)
	})
}
