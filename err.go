package parsec

import "github.com/goghcrow/lexer"

// ----------------------------------------------------------------
// Error Recovering
// ----------------------------------------------------------------

// Err :: p[a] -> err -> p[a]
// p 如果失败, 替换错误信息
func Err(p Parser, msg string) Parser {
	return parser(func(toks []*lexer.Token) Output {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		return fail(newError(branches.Pos, msg))
	})
}

// ErrDef :: p[a] -> err -> -> a -> p[a]
// p 如果失败, 返回默认值并替换错误信息
// 不会失败
func ErrDef(p Parser, msg string, def interface{}) Parser {
	return parser(func(toks []*lexer.Token) Output {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		return successWithErr([]Result{{def, toks}}, newError(branches.Pos, msg))
	})
}
