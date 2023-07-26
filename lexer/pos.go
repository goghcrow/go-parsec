package lexer

import "fmt"

type Pos struct {
	Idx    int // include
	IdxEnd int // exclude
	Col    int
	Line   int
}

func (p Pos) Loc() (idx /*include*/, idxEnd /*exclude*/, col, ln int) {
	return p.Idx, p.IdxEnd, p.Col, p.Line
}
func (p Pos) String() string {
	return fmt.Sprintf("pos %d-%d line %d col %d", p.Idx+1, p.IdxEnd+1, p.Line+1, p.Col+1)
}

func (p Pos) Merge(other Pos) Pos {
	if other.Idx <= p.Idx {
		panic("expect right pos")
	}
	p.IdxEnd = other.IdxEnd
	return p
}

func (p Pos) Span(runes []rune) string { return string(runes[p.Idx:p.IdxEnd]) }
