package legacy

// // RepSc :: p[a] -> p[list[a]]
// // 消费尽可能多的 p, 如果零次, 则返回 p[empty_list], 不会失败
// func RepSc[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
// 	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
// 		var err *Error
// 		// 层序遍历, 每层更新结果(从 root 到该层节点的路径), 返回最后一层的结果(根节点到叶子节点路径)
// 		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
// 		for {
// 			steps := xs
// 			xs = []Result[K, []R]{}
// 			for _, step := range steps {
// 				out := p.Parse(step.next)
// 				err = betterError(err, out.Error)
// 				if out.Success {
// 					for _, candidate := range out.Candidates {
// 						// 必须消费掉 token, 重复 nil 死循环
// 						if !step.next.equals(candidate.next) {
// 							xs = append(xs, Result[K, []R]{
// 								Val:  concat(step.Val, candidate.Val),
// 								next: candidate.next,
// 							})
// 						}
// 					}
// 				}
// 			}
// 			// 无下一层, 说明 steps 是最后一层, 循环出口
// 			if len(xs) == 0 {
// 				xs = steps
// 				break
// 			}
// 		}
// 		return successWithErr(xs, err)
// 	})
// }

// // List :: p[a] -> p[s] -> p[list[a]]
// func List[K Ord, R, Sep any](p Parser[K, R], s Parser[K, Sep]) Parser[K, []R] {
// 	seq2 := Seq2(p, Rep(Seq2(s, p)))
// 	return Apply[K, Cons[R, []Cons[Sep, R]], []R](seq2, applyList[R, Sep])
// }
//
// // ListSc :: p[a] -> p[s] -> p[list[a]]
// func ListSc[K Ord, R, Sep any](p Parser[K, R], s Parser[K, Sep]) Parser[K, []R] {
// 	seq2 := RepSc(Seq2(s, p))
// 	return Apply[K, Cons[R, []Cons[Sep, R]], []R](Seq2(p, seq2), applyList[R, Sep])
// }
//
// // applyList :: (a, list[(sep, a)]) -> list[a]
// func applyList[R, Sep any](t2 Cons[R, []Cons[Sep, R]]) []R {
// 	fst := t2.Car
// 	rest := t2.Cdr
// 	xs := make([]R, 1+len(rest))
// 	xs[0] = fst
// 	for i, it := range rest {
// 		xs[i+1] = it.Cdr
// 	}
// 	return xs
// }
