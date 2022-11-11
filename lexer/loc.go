package lexer

import (
	"fmt"
)

type Positionable interface {
	GetPos() Pos
}

func (p Pos) GetPos() Pos { return p }

type Pos struct {
	Idx    int // include
	IdxEnd int // exclude
	Col    int
	Line   int
}

var UnknownPos = Pos{-1, -1, -1, -1}

// Move Cursor
func (p *Pos) Move(r rune) {
	p.Idx++
	if r == '\n' {
		p.Line++
		p.Col = 0
	} else {
		p.Col++
	}
}

func (p Pos) Merge(other Pos) Pos {
	if other.Idx <= p.Idx {
		panic("expect right pos")
	}
	p.IdxEnd = other.IdxEnd
	return p
}

func (p Pos) String() string {
	return fmt.Sprintf("pos %d-%d line %d col %d", p.Idx+1, p.IdxEnd+1, p.Line+1, p.Col+1)
}

func (p Pos) Span(runes []rune) string { return string(runes[p.Idx:p.IdxEnd]) }
