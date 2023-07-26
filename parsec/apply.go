package parsec

// Apply :: p[a] -> (a -> b) -> p[b]
// 📢 the data structural of v is topological equivalent to syntax structural of p
func Apply[K Ord, From, To any](
	p Parser[K, From],
	f func(v From) To,
) Parser[K, To] {
	return parser[K, To](func(toks []Token[K]) Output[K, To] {
		out := p.Parse(toks)
		if !out.Success {
			return failOf[K, From, To](out)
		}
		xs := make([]Result[K, To], len(out.Candidates))
		for i, x := range out.Candidates {
			xs[i] = Result[K, To]{f(x.Val /*, tokenRange(toks, x.next)*/), x.next}
		}
		return successWithErr(xs, out.Error)
	})
}

// Bind :: p[a] -> (a->p[b]) -> p[b]
func Bind[K Ord, R1, R2 any](
	p Parser[K, R1],
	k func(R1) Parser[K, R2],
) Parser[K, R2] {
	return Combine2(p, k)
}

// Lazy :: (() -> p[a]) -> p[a]
func Lazy[K Ord, R any](thunk func() Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		return thunk().Parse(toks)
	})
}
