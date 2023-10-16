package parsec

// ----------------------------------------------------------------
// Repetitive
// ----------------------------------------------------------------

// Rep :: p[a] -> p[list[a]]
// 重复 n 次(n>=0), 按路径从长到短返回结果
func Rep[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	repR := RepR[K, R](p)
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		out := repR.Parse(toks)
		if out.Success {
			return successWithErr(reverse(out.Candidates), out.Error)
		}
		return out
	})
}

// RepSc :: p[a] -> p[list[a]]
// 消费尽可能多的 p, 如果零次, 则返回 p[empty_list], 不会失败
// Rep|RepR 返回所有层的结果, RepSc 返回最深一层结果
func RepSc[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		var err *Error
		// 层序遍历, 每层更新结果(从 root 到该层节点的路径), 返回最后一层的结果(根节点到叶子节点路径)
		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
		for {
			var nxs []Result[K, []R]
			for _, x := range xs {
				out := p.Parse(x.next)
				err = betterError(err, out.Error)
				if out.Success {
					for _, candidate := range out.Candidates {
						// 必须消费掉 token, 重复 nil 死循环
						if !toksEqual(x.next, candidate.next) {
							nxs = append(nxs, Result[K, []R]{
								Val:  concat(x.Val, candidate.Val),
								next: candidate.next,
							})
						}
					}
				}
			}
			if len(nxs) == 0 {
				break
			}
			xs = nxs
		}
		return successWithErr(xs, err)
	})
}

// RepR :: p[a] -> p[list[a]]
// 重复 n 次(n>=0), 按路径从短(empty)到长返回结果
func RepR[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		var err *Error
		// 层序遍历, 穷举所有根节点到非根节点的路径, Candidates 为每个节点的分叉数
		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
		for i := 0; i < len(xs); i++ {
			step := xs[i]
			out := p.Parse(step.next)
			err = betterError(err, out.Error)
			if out.Success {
				for _, candidate := range out.Candidates {
					// 必须消费掉 token, 重复 nil 死循环
					if !toksEqual(step.next, candidate.next) {
						xs = append(xs, Result[K, []R]{
							Val:  concat(step.Val, candidate.Val),
							next: candidate.next,
						})
					}
				}
			}
		}
		return successWithErr(xs, err)
	})
}

// RepN :: p[a] -> int -> p[list[a]]
// 即 Count, 重复 n 次
func RepN[K TK, R any](p Parser[K, R], cnt int) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		var err *Error
		// 层序遍历, 每层更新结果(从 root 到该层节点的路径), 返回最后一层的结果(根节点到叶子节点路径)
		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
		for i := 0; i < cnt; i++ {
			var nxs []Result[K, []R]
			for _, x := range xs {
				out := p.Parse(x.next)
				err = betterError(err, out.Error)
				if out.Success {
					// if !x.next.equals(candidate.next) {}
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

		return successWithErr(xs, err)
	})
}

// ----------------------------------------------------------------
// List
// ----------------------------------------------------------------

// List :: p[a] -> p[s] -> p[list[a]]
// 返回所有可能的序列
func List[K TK, R, Sep any](p Parser[K, R], s Parser[K, Sep]) Parser[K, []R] {
	return Apply(Seq2(p, Rep(KRight(s, p))), applyList[R, Sep])
}

// ListSc :: p[a] -> p[s] -> p[list[a]]
// 返回最长序列
func ListSc[K TK, R, Sep any](p Parser[K, R], s Parser[K, Sep]) Parser[K, []R] {
	return Apply(Seq2(p, RepSc(KRight(s, p))), applyList[R, Sep])
}

// ListN :: p[a] -> p[s] -> int -> p[list[a]]
// 返回固定数量序列
func ListN[K TK, R, Sep any](p Parser[K, R], s Parser[K, Sep], cnt int) Parser[K, []R] {
	if cnt < 1 {
		return Succ[K, []R]([]R{})
	} else if cnt == 1 {
		return Apply(p, func(v R) []R { return []R{v} })
	} else {
		return Apply(Seq2(p, RepN(KRight(s, p), cnt-1)), applyList[R, Sep])
	}
}

// applyList :: (a, Cons[(sep, a)]) -> list[a]
func applyList[R, Sep any](v Cons[R, []R]) []R {
	return concat([]R{v.Car}, v.Cdr...)
}
