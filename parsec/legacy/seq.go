package legacy

// // Seq :: p[a] -> p[b] -> p[c] -> ... -> p[(a,b,c...)]
// // 对 ps 进行 foldLeft, append 收集数据
// func Seq[K Ord, R any](ps ...Parser[K, R]) Parser[K, []R] {
// 	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
// 		var err *Error
// 		// 层序遍历, ps 代表层次(每层使用的 p), 每层更新结果(从 root 到该层节点的路径),
// 		// 返回根节点到所有叶子节点的路径
// 		xs := []Result[K, []R]{{Val: []R{}, next: toks}}
// 		for _, p := range ps {
// 			if len(xs) == 0 { // 快速失败
// 				break
// 			}
//
// 			steps := xs
// 			xs = []Result[K, []R]{}
// 			for _, step := range steps {
// 				out := p.Parse(step.next)
// 				err = betterError(err, out.Error)
// 				if out.Success {
// 					for _, candidate := range out.Candidates {
// 						xs = append(xs, Result[K, []R]{
// 							Val:  concat(step.Val, candidate.Val),
// 							next: candidate.next,
// 						})
// 					}
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }

// // Combine :: p[a] -> (x->p[b]) -> (x->p[c]) -> p[c]
// // CPS 风格
// func Combine[K Ord, R any](
// 	p Parser[K, R],
// 	ks ...func(R) Parser[K, R], // continuations
// ) Parser[K, R] {
// 	return parser[K, R](func(toks []Token[K]) Output[K, R] {
// 		out1 := p.Parse(toks)
// 		if !out1.Success {
// 			return out1
// 		}
//
// 		xs := out1.Candidates
// 		err := out1.Error
// 		for _, k := range ks {
// 			if len(xs) == 0 {
// 				break
// 			}
//
// 			steps := xs
// 			xs = []Result[K, R]{}
// 			for _, step := range steps {
// 				out := k(step.Val).Parse(step.next)
// 				err = betterError(err, out.Error)
// 				if out.Success {
// 					xs = append(xs, out.Candidates...)
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
