package parsec

// ----------------------------------------------------------------
// Left Recursive
// ----------------------------------------------------------------

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

// LRec :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Equivalent to seq(a, rep(b))
// f(f(f(a, b1), b2), b3) ...
// returns multiple possible results
func LRec(p, q Parser, f func(a, b interface{}) interface{}) Parser {
	return Apply(Seq(p, Rep(q)), applyLRec(f))
}

// LRecSc :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// returns a single result which consumes tokens as much as possible,
func LRecSc(p, q Parser, f func(a, b interface{}) interface{}) Parser {
	return Apply(Seq(p, RepSc(q)), applyLRec(f))
}
