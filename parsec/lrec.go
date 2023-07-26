package parsec

// ----------------------------------------------------------------
// Left Recursive
// ----------------------------------------------------------------

// LRec :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 返回多种可能的结果
// P  → Pq|p
// 消除直接左递归
// P  → pP’
// P’ → qP’|ε
// pq, pqq, pqqq, ...
func LRec[K Ord /*Fst extends R,*/, Sec, R any](
	p Parser[K, R],
	q Parser[K, Sec],
	f func(R, Sec) R,
) Parser[K, R] {
	return Apply(Seq2(p, Rep(q)), applyLRec(f))
}

// LRecSc :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Equivalent to seq(a, rep_sc(b)))
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 只返回一个消费尽可能多 token 的结果
func LRecSc[K Ord /*Fst extends R,*/, Sec, R any](
	p Parser[K, R],
	q Parser[K, Sec],
	f func(R, Sec) R,
) Parser[K, R] {
	return Apply(Seq2(p, RepSc(q)), applyLRec(f))
}

// applyLRec :: (a -> b -> c) -> ((a, list[b]) -> a)
func applyLRec[Sec, R any](f func(R, Sec) R) func(Cons[R, []Sec]) R {
	return func(v Cons[R, []Sec]) R {
		x := v.Car
		for _, tail := range v.Cdr {
			x = f(x, tail)
		}
		return x
	}
}
