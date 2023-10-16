package parsec

// ----------------------------------------------------------------
// Sequential
// ----------------------------------------------------------------

// Seq :: p[a] -> p[b] -> p[c] -> ... -> p[(a,b,c...)]
// 顺次匹配, 对 ps 进行 foldLeft, append 收集数据
func Seq[K TK, R any](ps ...Parser[K, R]) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		var err *Error
		// 层序遍历, ps 代表层次(每层使用的 p), 每层更新结果(从 root 到该层节点的路径),
		// 返回根节点到所有叶子节点的路径
		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
		for _, p := range ps {
			var nxs []Result[K, []R]
			for _, x := range xs {
				out := p.Parse(x.next)
				err = betterError(err, out.Error)
				if out.Success {
					for _, candidate := range out.Candidates {
						nxs = append(nxs, Result[K, []R]{
							Val:  concat(x.Val, candidate.Val),
							next: candidate.next,
						})
					}
				}
			}
			if len(nxs) == 0 {
				return fail[K, []R](err)
			}
			xs = nxs
		}
		return newOutput(xs, err, len(xs) != 0)
	})
}

func Seq2[K TK, R1, R2 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
) Parser[K, Cons[R1, R2]] {
	return parser[K, Cons[R1, R2]](func(toks []Token[K]) Output[K, Cons[R1, R2]] {
		out1 := p1.Parse(toks)
		if !out1.Success {
			return failOf[K, R1, Cons[R1, R2]](out1)
		}
		var xs []Result[K, Cons[R1, R2]]
		err := out1.Error
		for _, step := range out1.Candidates {
			out2 := p2.Parse(step.next)
			err = betterError(err, out2.Error)
			if out2.Success {
				for _, candidate := range out2.Candidates {
					xs = append(xs, Result[K, Cons[R1, R2]]{
						Val:  Cons[R1, R2]{Car: step.Val, Cdr: candidate.Val},
						next: candidate.next,
					})
				}
			}
		}
		return newOutput(xs, err, len(xs) != 0)
	})
}
func Seq3[K TK, R1, R2, R3 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
	p3 Parser[K, R3],
) Parser[K, Cons[R1, Cons[R2, R3]]] {
	return Seq2(p1, Seq2(p2, p3))
}
func Seq4[K TK, R1, R2, R3, R4 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
	p3 Parser[K, R3],
	p4 Parser[K, R4],
) Parser[K, Cons[R1, Cons[R2, Cons[R3, R4]]]] {
	var p = Seq2(p1, Seq2(p2, Seq2(p3, p4)))
	return p
}
func Seq5[K TK, R1, R2, R3, R4, R5 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
	p3 Parser[K, R3],
	p4 Parser[K, R4],
	p5 Parser[K, R5],
) Parser[K, Cons[R1, Cons[R2, Cons[R3, Cons[R4, R5]]]]] {
	// tmp vars making goland happy !
	tl := Seq2(p2, Seq2(p3, Seq2(p4, p5)))
	var p = Seq2(p1, tl)
	return p
}

// ----------------------------------------------------------------
// Combine
// ----------------------------------------------------------------

// Combine :: p[a] -> (a->p[b]) -> (b->p[c]) -> p[c]
// 这里实际上是 >>= Bind 函数, CPS写法, 用来表达上下文文法
// 这里牵扯到 Monad, monadic parsing,
// 因为 Monad 是 context sensitive 所以可以通过 combine 可以处理 CSG,
// 其他函数只能处理 CFG, 这个写的挺好的
// https://stackoverflow.com/questions/7861903/what-are-the-benefits-of-applicative-parsing-over-monadic-parsing
func Combine[K TK, R any](
	p Parser[K, R],
	ks ...func(R) Parser[K, R], // continuations
) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		out1 := p.Parse(toks)
		if !out1.Success {
			return out1
		}

		xs := out1.Candidates
		err := out1.Error

		for _, k := range ks {
			var nxs []Result[K, R]
			for _, x := range xs {
				out := k(x.Val).Parse(x.next)
				err = betterError(err, out.Error)
				if out.Success {
					// 如果需要 concat 用 seq
					nxs = append(nxs, out.Candidates...)
				}
			}
			if len(nxs) == 0 {
				return fail[K, R](err)
			}
			xs = nxs
		}
		return newOutput(xs, err, len(xs) != 0)
	})
}

// Combine2 p[a] -> (a->p[b]) -> p[b]
func Combine2[K TK, R1, R2 any](
	p Parser[K, R1],
	k func(R1) Parser[K, R2],
) Parser[K, R2] {
	return parser[K, R2](func(toks []Token[K]) Output[K, R2] {
		out1 := p.Parse(toks)
		if !out1.Success {
			return failOf[K, R1, R2](out1)
		}

		var xs []Result[K, R2]
		err := out1.Error
		for _, step := range out1.Candidates {
			out := k(step.Val).Parse(step.next)
			err = betterError(err, out.Error)
			if out.Success {
				xs = append(xs, out.Candidates...)
			}
		}

		return newOutput(xs, err, len(xs) != 0)
	})
}
func Combine3[K TK, R1, R2, R3 any](
	p Parser[K, R1],
	k1 func(R1) Parser[K, R2],
	k2 func(R2) Parser[K, R3],
) Parser[K, R3] {
	return parser[K, R3](func(toks []Token[K]) Output[K, R3] {
		return Combine2(Combine2(p, k1), k2).Parse(toks)
	})
}
func Combine4[K TK, R1, R2, R3, R4 any](
	p Parser[K, R1],
	k1 func(R1) Parser[K, R2],
	k2 func(R2) Parser[K, R3],
	k3 func(R3) Parser[K, R4],
) Parser[K, R4] {
	return parser[K, R4](func(toks []Token[K]) Output[K, R4] {
		return Combine2(Combine3(p, k1, k2), k3).Parse(toks)
	})
}
func Combine5[K TK, R1, R2, R3, R4, R5 any](
	p Parser[K, R1],
	k1 func(R1) Parser[K, R2],
	k2 func(R2) Parser[K, R3],
	k3 func(R3) Parser[K, R4],
	k4 func(R4) Parser[K, R5],
) Parser[K, R5] {
	return parser[K, R5](func(toks []Token[K]) Output[K, R5] {
		return Combine2(Combine4(p, k1, k2, k3), k4).Parse(toks)
	})
}

// ----------------------------------------------------------------
// KLeft / KRight / KMid
// ----------------------------------------------------------------

// KLeft :: p[a] -> p[b] -> p[a]
func KLeft[K TK, A, B any](
	p1 Parser[K, A],
	p2 Parser[K, B],
) Parser[K, A] {
	return Apply(Seq2(p1, p2), C2ar[A, B])
}

// KRight :: p[a] -> p[b] -> p[b]
func KRight[K TK, A, B any](
	p1 Parser[K, A],
	p2 Parser[K, B],
) Parser[K, B] {
	return Apply(Seq2(p1, p2), C2dr[A, B])
}

// KMid :: p[a] -> p[b] -> p[c] -> p[b]
func KMid[K TK, A, B, C any](
	p1 Parser[K, A],
	p2 Parser[K, B],
	p3 Parser[K, C],
) Parser[K, B] {
	return Apply(Seq3(p1, p2, p3), C3adr[A, B, C])
}
