package lexer

import (
	"fmt"
	"unicode/utf8"
)

func BuildLexer(f func(lexicon *Lexicon)) *Lexer {
	lex := NewLexicon()
	f(&lex)
	return NewLexer(lex)
}

func NewLexer(lex Lexicon) *Lexer { return &Lexer{Lexicon: lex} }

func (l *Lexer) Lex(input string) []*Token {
	l.input = []rune(input)
	l.Loc = Loc{}
	var toks []*Token
	for {
		t, keep := l.next()
		if t == nil {
			break
		}
		if keep {
			toks = append(toks, t)
		}
	}
	return toks
}

type Lexer struct {
	Lexicon
	Loc
	input []rune
}

func (l *Lexer) next() (*Token, bool) {
	if l.Pos >= len(l.input) {
		return nil, true
	}

	pos := l.Loc
	sub := string(l.input[l.Pos:])
	for _, rl := range l.Lexicon.rules {
		offset := rl.match(sub)
		if offset >= 0 {
			matched := l.input[l.Pos : l.Pos+offset]
			for _, r := range matched {
				l.Move(r)
			}
			pos.PosEnd = l.Loc.Pos
			return &Token{TokenKind: rl.TokenKind, Lexeme: string(matched), Loc: pos}, rl.keep
		}
	}
	panic(fmt.Errorf("syntax error in %s: nothing token matched", l.Loc))
}

func runeCount(s string) int { return utf8.RuneCountInString(s) }
