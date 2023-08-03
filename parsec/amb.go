package parsec

// ----------------------------------------------------------------
// Ambiguity Resolving, 歧义合并器
// ----------------------------------------------------------------

// Amb :: p[a] -> p[list[a]]
// Consumes x and merge group result by consumed tokens.
func Amb[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		branches := p.Parse(toks)
		if !branches.Success {
			return failOf[K, R, []R](branches)
		}

		group := make(map[Token[K]][]Result[K, R])
		for _, r := range branches.Candidates {
			k := beginTok(r.next)
			group[k] = append(group[k], r)
		}

		xs := make([]Result[K, []R], 0, len(group))
		for _, vals := range group {
			merged := sliceMap(vals, func(v Result[K, R]) R { return v.Val })
			xs = append(xs, Result[K, []R]{merged, vals[0].next})
		}
		return successWithErr(xs, branches.Error)
	})
}
