package parsec

// Apply :: p[a] -> (a -> b) -> p[b]
// 即 Map, the data structural of v is topological equivalent to syntax structural of p
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
