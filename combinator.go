package parsec

// List :: p[a] -> p[s] -> p[list[a]]
func List(p, s Parser) Parser { return Apply(Seq(p, Rep(Seq(s, p))), applyList) }

// ListSc :: p[a] -> p[s] -> p[list[a]]
func ListSc(p, s Parser) Parser { return Apply(Seq(p, RepSc(Seq(s, p))), applyList) }

// applyList :: (a, list[(sep, a)]) -> list[a]
func applyList(v interface{}) interface{} {
	var xs []interface{}
	a := v.([]interface{})
	xs = append(xs, a[0])
	for _, it := range a[1].([]interface{}) {
		xs = append(xs, it.([]interface{})[1])
	}
	return xs
}
