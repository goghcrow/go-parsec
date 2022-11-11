package lexer

type TokenKind int

type Token struct {
	TokenKind
	Pos
	Lexeme string
}

func (t *Token) String() string {
	return t.Lexeme
}
