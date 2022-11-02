package parsec

// ----------------------------------------------------------------
// Left Recursive
// ----------------------------------------------------------------

// LRec :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 返回多种可能的结果
// P → Pq|p 消除直接左递归
// P  → pP’
// P’ → qP’|ε
// pq, pqq, pqqq, ...
func LRec(p, q Parser, f func(a, b interface{}) interface{}) Parser {
	return Apply(Seq(p, Rep(q)), applyLRec(f))
}

// LRecSc :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Equivalent to seq(a, rep_sc(b)))
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 只返回一个消费尽可能多 token 的结果
func LRecSc(p, q Parser, f func(a, b interface{}) interface{}) Parser {
	return Apply(Seq(p, RepSc(q)), applyLRec(f))
}

// applyLRec :: (a -> b -> a) -> ((a, list[b]) -> a)
func applyLRec(f func(a, b interface{}) interface{}) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		a := v.([]interface{})
		x := a[0]
		for _, tail := range a[1].([]interface{}) {
			x = f(x, tail)
		}
		return x
	}
}
