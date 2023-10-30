package lexer

import (
	"fmt"
)

type Ord = comparable

// BuildLexer
// 这里没有取最长匹配, 而是首次匹配, so, 注意规则顺序
// 具体可以参考 example/lexicon.go
func BuildLexer[K Ord](f func(lexicon *Lexicon[K])) *Lexer[K] {
	lex := NewLexicon[K]()
	f(&lex)
	return NewLexer(lex)
}

func NewLexer[K Ord](lex Lexicon[K]) *Lexer[K] {
	return &Lexer[K]{Lexicon: lex}
}

func (l *Lexer[K]) MustLex(input string) []*Token[K] {
	toks, err := l.Lex(input)
	if err != nil {
		panic(err)
	}
	return toks
}

func (l *Lexer[K]) Lex(input string) ([]*Token[K], error) {
	l.input = []rune(input)
	l.Pos = Pos{}
	var toks []*Token[K]
	for {
		t, keep, err := l.next()
		if err != nil {
			return toks, err
		}
		if t == nil {
			break
		}
		if keep {
			toks = append(toks, t)
		}
	}
	return toks, nil
}

type Lexer[K Ord] struct {
	Lexicon[K]
	Pos
	input []rune
}

func (l *Lexer[K]) next() (tok *Token[K], keep bool, err error) {
	if l.Idx >= len(l.input) {
		return nil, true, nil
	}

	pos := l.Pos
	sub := string(l.input[l.Idx:])
	for _, rl := range l.Lexicon.rules {
		offset := rl.match(sub)
		if offset >= 0 {
			matched := l.input[l.Idx : l.Idx+offset]
			for _, r := range matched {
				l.Move(r)
			}
			pos.IdxEnd = l.Pos.Idx
			return &Token[K]{kind: rl.K, lexeme: string(matched), Pos: pos}, rl.keep, nil
		}
	}
	return nil, false, fmt.Errorf("syntax error in %s: nothing token matched", l.Pos)
}

func (l *Lexer[K]) Move(r rune) {
	l.Idx++
	if r == '\n' {
		l.Line++
		l.Col = 0
	} else {
		l.Col++
	}
}

// ----------------------------------------------------------------
// Token
// ----------------------------------------------------------------

type Token[K Ord] struct {
	Pos
	kind   K
	lexeme string
}

func (t *Token[K]) Kind() K        { return t.kind }
func (t *Token[K]) Lexeme() string { return t.lexeme }
func (t *Token[K]) String() string { return t.lexeme }

// func (t *Token[K]) String() string { return fmt.Sprintf("<'%s', %v>", t.Lexeme, t.Kind) }

// ----------------------------------------------------------------
// Position or Location
// ----------------------------------------------------------------

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

// ----------------------------------------------------------------
// Error
// ----------------------------------------------------------------

type Error struct {
	Pos
	Msg string
}

func (e *Error) Error() string {
	idx, end, col, ln := e.Pos.Loc()
	return fmt.Sprintf("%s in pos %d-%d line %d col %d", e.Msg, idx+1, end+1, ln+1, col+1)
}
