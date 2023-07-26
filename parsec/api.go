package parsec

import "fmt"

type Ord = comparable

type Lex[K Ord] func(input string) ([]Token[K], error)

type Parser[K Ord, R any] interface {
	Parse([]Token[K]) Output[K, R]
}

// Parser Impl
type parser[K Ord, R any] func([]Token[K]) Output[K, R]

func (p parser[K, R]) Parse(toks []Token[K]) Output[K, R] {
	return p(toks)
}

// Output
// If Success == true, it means that the candidates field is valid, even when it is empty.
// If Success == false, error will be not null
// The Error field stores the far-est error that has even been seen, even when tokens are successfully parsed.
type Output[K Ord, R any] struct {
	Success    bool
	Candidates []Result[K, R]
	*Error
}

type Result[K Ord, R any] struct {
	Val  R
	next tokSeq[K] // rest of tokens
}

type Error struct {
	Pos
	Msg string
}

func (e *Error) Error() string {
	if vp, ok := e.Pos.(VirtualPos); ok {
		return fmt.Sprintf("%s in %s", e.Msg, vp)
	}
	idx, end, col, ln := e.Pos.Loc()
	return fmt.Sprintf("%s in pos %d-%d line %d col %d", e.Msg, idx+1, end+1, ln+1, col+1)
}