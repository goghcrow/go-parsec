package legacy

// // OptSc :: p[a] -> p[a|nil]
// // Opt 返回两种结果, OptSc 返回一种结果, 只有 p 失败才返回 nil
// func OptSc[K Ord, R any](p Parser[K, R]) Parser[K, R] {
// 	return parser[K, R](func(toks []Token[K]) Output[K, R] {
// 		out := p.Parse(toks)
// 		if out.Success {
// 			return out
// 		}
// 		// nil with error
// 		return successWithErr([]Result[K, R]{{next: toks}}, out.Error)
// 	})
// }

// 返回 option 版本, 语义更清晰

// // Opt :: p[a] -> p[a|nil]
// // Alt 返回失败, Opt 不返回失败
// func Opt[K Ord, R any](p Parser[K, R]) Parser[K, Option[R]] {
// 	return Alt(
// 		Apply(p, func(v R) Option[R] { return Some[R](v) }),
// 		Apply(Nil[K, R](), func(v R) Option[R] { return None[R]() }),
// 	)
// }
//
// // Opt :: p[a] -> p[a|nil]
// // Alt 返回失败, Opt 不返回失败
// func OptSc[K Ord, R any](p Parser[K, R]) Parser[K, Option[R]] {
// 	return AltSc(
// 		Apply(p, func(v R) Option[R] { return Some[R](v) }),
// 		Apply(Nil[K, R](), func(v R) Option[R] { return None[R]() }),
// 	)
// }
