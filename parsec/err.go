package parsec

// ----------------------------------------------------------------
// Error Recovering
// ----------------------------------------------------------------

// Err :: p[a] -> err -> p[a]
// p 如果失败, 替换错误信息, 提供更准确错误信息
// e.g. Err(Alt(Tok(Int), Tok(Float)), "expect number")
func Err[K Ord, R any](p Parser[K, R], msg string) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		if msg == "" {
			msg = branches.Msg
		}
		return fail[K, R](newError(branches.Pos, msg))
	})
}

// ErrD :: p[a] -> err -> -> a -> p[a]
// p 如果失败, 返回默认值并替换错误信息, 返回成功, 不消耗 toks, 用来进行错误回复
func ErrD[K Ord, R any](p Parser[K, R], msg string, defaultValue R) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		branches := p.Parse(toks)
		if branches.Success {
			return branches
		}
		if msg == "" {
			msg = branches.Msg
		}
		return successWithErr(
			[]Result[K, R]{{Val: defaultValue, next: toks}},
			newError(branches.Pos, msg),
		)
	})
}
