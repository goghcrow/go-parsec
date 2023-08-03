package parsec

import "fmt"

func NewRule[K Ord, R any]() *SyntaxRule[K, R] {
	return &SyntaxRule[K, R]{}
}

type SyntaxRule[K Ord, R any] struct {
	Pattern Parser[K, R]
}

func (r *SyntaxRule[K, R]) Parse(toks []Token[K]) Output[K, R] {
	if r.Pattern == nil {
		panic("Rule has not been initialized. Pattern is required before calling parse.")
	}
	return r.Pattern.Parse(toks)
}

// Parser
// SyntaxRule 已经实现了 Parser 接口, 但是类型推导不大性, 加个 Helper 函数
func (r *SyntaxRule[K, R]) Parser() Parser[K, R] {
	return r
}

func ExpectEOF[K Ord, R any](out Output[K, R]) Output[K, R] {
	if !out.Success {
		return out
	}
	if len(out.Candidates) == 0 {
		return fail[K, R](newError(EOFPos, "No result is returned."))
	}

	var xs []Result[K, R]
	err := out.Error
	for _, candidate := range out.Candidates {
		if len(candidate.next) == 0 {
			xs = append(xs, candidate)
		} else {
			pso := beginPos(candidate.next)
			msg := fmt.Sprintf("The parser cannot reach the end of file, stops %s in %s",
				candidate.next[0], Pos(candidate.next[0]))
			err = betterError(err, newError(pso, msg))
		}
	}
	return newOutput(xs, err, len(xs) != 0)
}

func ExpectSingleResult[K Ord, R any](out Output[K, R]) (R, error) {
	if !out.Success {
		return *new(R), out.Error
	}
	if len(out.Candidates) == 0 {
		return *new(R), newError(UnknownPos, "No result is returned.")
	}
	if len(out.Candidates) != 1 {
		return *new(R), newError(UnknownPos, "Multiple results are returned.")
	}
	return out.Candidates[0].Val, nil
}
