package lexer

// TokenKind 因为要支持动态添加操作符, 所以 TokenKind 没有定义成 int 枚举
// 这里需要自己保证 type 值不重复
type TokenKind string

type Token struct {
	TokenKind
	Loc
	Lexeme string
}

func (t *Token) String() string {
	// return strconv.Quote(t.Lexeme)
	return t.Lexeme
}
