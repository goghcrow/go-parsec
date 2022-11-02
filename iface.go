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

func NewRule() *Rule { return &Rule{} }

type Rule struct {
	Pattern Parser
}

func (i *Rule) Parse(toks []*lexer.Token) Output             { return i.Pattern.Parse(toks) }
func (i *Rule) Map(f func(v interface{}) interface{}) Parser { return Apply(i, f) }

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
	toks tokSeq // rest of tokens
}

type Error struct {
	Loc lexer.Loc
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("%s in %s", e.Msg, e.Loc) }

func ExpectEOF(out Output) Output {
	if !out.Success {
		return out
	}
	if len(out.Candidates) == 0 {
		return fail(newError(lexer.UnknownLoc, "No result is returned."))
	}

	var xs []Result
	err := out.Error
	for _, candidate := range out.Candidates {
		if len(candidate.toks) == 0 {
			xs = append(xs, candidate)
		} else {
			err = betterError(err, newError(candidate.toks.loc(),
				fmt.Sprintf("The parser cannot reach the end of file, stops %s in %s",
					candidate.toks[0], candidate.toks[0].Loc)))
		}
	}
	return newOutput(xs, err, len(xs) != 0)
}

func ExpectSingleResult(out Output) (interface{}, error) {
	if !out.Success {
		// return Result{}, newError(lexer.UnknownLoc, out.Error.Error())
		return Result{}, out.Error
	}
	if len(out.Candidates) == 0 {
		return Result{}, newError(lexer.UnknownLoc, "No result is returned.")
	}
	if len(out.Candidates) != 1 {
		return Result{}, newError(lexer.UnknownLoc, "Multiple results are returned.")
	}
	return out.Candidates[0].Val, nil
}
