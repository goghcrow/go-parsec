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
