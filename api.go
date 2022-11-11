package parsec

import (
	"fmt"
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Interface
// ----------------------------------------------------------------

type Parser interface {
	Parse([]*lexer.Token) Output
	Map(mapper func(v interface{}) interface{}) Parser
}

func NewRule() *SyntaxRule { return &SyntaxRule{} }

type SyntaxRule struct {
	Pattern Parser
}

func (r *SyntaxRule) Parse(toks []*lexer.Token) Output             { return r.Pattern.Parse(toks) }
func (r *SyntaxRule) Map(f func(v interface{}) interface{}) Parser { return Apply(r, f) }

// Parser Impl
type parser func([]*lexer.Token) Output

func (p parser) Parse(toks []*lexer.Token) Output             { return p(toks) }
func (p parser) Map(f func(v interface{}) interface{}) Parser { return Apply(p, f) }

// Output
// If successful===true, it means that the candidates field is valid, even when it is empty.
// If successful===false, error will be not null
// The error field stores the farest error that has even been seen, even when tokens are successfully parsed.
type Output struct {
	Success    bool
	Candidates []Result
	*Error
}

type Result struct {
	Val  interface{}
	next tokSeq // rest of tokens
}

type Error struct {
	lexer.Pos
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("%s in %s", e.Msg, e.Pos) }

func ExpectEOF(out Output) Output {
	if !out.Success {
		return out
	}
	if len(out.Candidates) == 0 {
		return fail(newError(lexer.UnknownPos, "No result is returned."))
	}

	var xs []Result
	err := out.Error
	for _, candidate := range out.Candidates {
		if len(candidate.next) == 0 {
			xs = append(xs, candidate)
		} else {
			err = betterError(err, newError(candidate.next.pos(),
				fmt.Sprintf("The parser cannot reach the end of file, stops %s in %s",
					candidate.next[0], candidate.next[0].Pos)))
		}
	}
	return newOutput(xs, err, len(xs) != 0)
}

func ExpectSingleResult(out Output) (interface{}, error) {
	if !out.Success {
		// return Result{}, newError(lexer.UnknownPos, out.Error.Error())
		return Result{}, out.Error
	}
	if len(out.Candidates) == 0 {
		return Result{}, newError(lexer.UnknownPos, "No result is returned.")
	}
	if len(out.Candidates) != 1 {
		return Result{}, newError(lexer.UnknownPos, "Multiple results are returned.")
	}
	return out.Candidates[0].Val, nil
}
