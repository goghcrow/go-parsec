package lexer

type TokenKind int

type Token struct {
	TokenKind
	Loc
	Lexeme string
}

func (t *Token) String() string {
	return t.Lexeme
}
