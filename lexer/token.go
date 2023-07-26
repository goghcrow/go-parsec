package lexer

type Token[K Ord] struct {
	Pos
	kind   K
	lexeme string
}

func (t *Token[K]) Kind() K        { return t.kind }
func (t *Token[K]) Lexeme() string { return t.lexeme }
func (t *Token[K]) String() string { return t.lexeme }

// func (t *Token[K]) String() string { return fmt.Sprintf("<'%s', %v>", t.Lexeme, t.Kind) }
