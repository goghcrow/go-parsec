package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Sequential
// ----------------------------------------------------------------

// Seq :: p[a] -> p[b] -> p[c] -> ... -> p[(a,b,c...)]
func Seq(ps ...Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		xs := []Result{{Val: emptySlice(), toks: toks}}
		for _, p := range ps {
			if len(xs) == 0 {
				break
			}

			steps := xs
			xs = []Result{}
			for _, step := range steps {
				out := p.Parse(step.toks)
				err = betterError(err, out.Error)
				if out.Success {
					for _, candidate := range out.Candidates {
						xs = append(xs, Result{
							Val:  append(step.Val.([]interface{}), candidate.Val),
							toks: candidate.toks,
						})
					}
				}
			}
		}
		return resultOrError(xs, err, len(xs) != 0)
	})
}

// KLeft :: p[a] -> p[b] -> p[a]
func KLeft(p1, p2 Parser) Parser {
	return Apply(Seq(p1, p2), func(v interface{}) interface{} { return v.([]interface{})[0] })
}

// KRight :: p[a] -> p[b] -> p[b]
func KRight(p1, p2 Parser) Parser {
	return Apply(Seq(p1, p2), func(v interface{}) interface{} { return v.([]interface{})[1] })
}

// KMid :: p[a] -> p[b] -> p[c] -> p[b]
func KMid(p1, p2, p3 Parser) Parser {
	return Apply(Seq(p1, p2, p3), func(v interface{}) interface{} { return v.([]interface{})[1] })
}
