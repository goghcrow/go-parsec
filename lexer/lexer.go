package lexer

import (
	"fmt"
	"unicode/utf8"
)

// BuildLexer
// 这里没有取最长匹配, 而是首次匹配, so, 注意规则顺序
// 具体可以参考 example/lexicon.go
func BuildLexer(f func(lexicon *Lexicon)) *Lexer {
	lex := NewLexicon()
	f(&lex)
	return NewLexer(lex)
}

func NewLexer(lex Lexicon) *Lexer { return &Lexer{Lexicon: lex} }

func (l *Lexer) MustLex(input string) []*Token {
	toks, err := l.Lex(input)
	if err != nil {
		panic(err)
	}
	return toks
}

func (l *Lexer) Lex(input string) ([]*Token, error) {
	l.input = []rune(input)
	l.Loc = Loc{}
	var toks []*Token
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

type Lexer struct {
	Lexicon
	Loc
	input []rune
}

func (l *Lexer) next() (tok *Token, keep bool, err error) {
	if l.Pos >= len(l.input) {
		return nil, true, nil
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
			return &Token{TokenKind: rl.TokenKind, Lexeme: string(matched), Loc: pos}, rl.keep, nil
		}
	}
	return nil, false, fmt.Errorf("syntax error in %s: nothing token matched", l.Loc)
}

func runeCount(s string) int { return utf8.RuneCountInString(s) }
