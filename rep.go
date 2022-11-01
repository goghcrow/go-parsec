package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Repetitive
// ----------------------------------------------------------------

// Rep :: p[a] -> p[list[a]]
func Rep(p Parser) Parser {
	p = RepR(p)
	return newParser(func(toks []*lexer.Token) Output {
		out := p.Parse(toks)
		if out.Success {
			return successX(reverse(out.Candidates), out.Error)
		}
		return out
	})
}

// RepSc :: p[a] -> p[list[a]]
func RepSc(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		xs := []Result{{Val: emptySlice(), toks: toks}}

		for {
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

			if len(xs) == 0 {
				xs = steps
				break
			}
		}
		return resultOrError(xs, err, true)
	})
}

// RepR :: p[a] -> p[list[a]]
func RepR(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		xs := []Result{{Val: emptySlice(), toks: toks}}

		for i := 0; i < len(xs); i++ {
			step := xs[i]
			out := p.Parse(step.toks)
			err = betterError(err, out.Error)
			if out.Success {
				for _, candidate := range out.Candidates {
					if !step.toks.equals(candidate.toks) {
						xs = append(xs, Result{
							Val:  append(step.Val.([]interface{}), candidate.Val),
							toks: candidate.toks,
						})
					}
				}
			}
		}
		return resultOrError(xs, err, true)
	})
}

// applyList :: (a, list[(sep, a)]) -> list[a]
func applyList(v interface{}) interface{} {
	var xs []interface{}
	a := v.([]interface{})
	xs = append(xs, a[0])
	for _, it := range a[1].([]interface{}) {
		xs = append(xs, it.([]interface{})[1])
	}
	return xs
}

// List :: p[a] -> p[s] -> p[list[a]]
func List(p, s Parser) Parser { return Apply(Seq(p, Rep(Seq(s, p))), applyList) }

// ListSc :: p[a] -> p[s] -> p[list[a]]
func ListSc(p, s Parser) Parser { return Apply(Seq(p, RepSc(Seq(s, p))), applyList) }
