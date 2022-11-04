package parsec

// List :: p[a] -> p[s] -> p[list[a]]
func List(p, s Parser) Parser { return Apply(Seq(p, Rep(Seq(s, p))), applyList) }

// ListSc :: p[a] -> p[s] -> p[list[a]]
func ListSc(p, s Parser) Parser { return Apply(Seq(p, RepSc(Seq(s, p))), applyList) }

// applyList :: (a, list[(sep, a)]) -> list[a]
func applyList(v interface{}) interface{} {
	a := v.([]interface{})
	fst := a[0]
	rest := a[1].([]interface{})
	xs := make([]interface{}, 1+len(rest))
	xs[0] = fst
	for i, it := range rest {
		xs[i+1] = it.([]interface{})[1]
	}
	return xs
}
