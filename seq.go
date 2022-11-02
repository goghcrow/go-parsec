package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Sequential
// ----------------------------------------------------------------

// Seq :: p[a] -> p[b] -> p[c] -> ... -> p[(a,b,c...)]
// 对 ps 进行 foldLeft, append 收集数据
func Seq(ps ...Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		// 层序遍历, ps 代表层次(每层使用的 p), 每层更新结果(从 root 到该层节点的路径),
		// 返回根节点到所有叶子节点的路径
		xs := []Result{{Val: anySlice(), toks: toks}} // root
		for _, p := range ps {
			if len(xs) == 0 { // 快速失败
				break
			}

			steps := xs
			xs = []Result{}
			for _, step := range steps {
				out := p.Parse(step.toks)
				err = betterError(err, out.Error)
				if out.Success {
					for _, candidate := range out.Candidates {
						xs = append(xs, Result{
							Val:  append(step.Val.([]interface{}), candidate.Val),
							toks: candidate.toks,
						})
					}
				}
			}
		}
		return newOutput(xs, err, len(xs) != 0)
	})
}

// KLeft :: p[a] -> p[b] -> p[a]
func KLeft(p1, p2 Parser) Parser {
	return Apply(Seq(p1, p2), func(v interface{}) interface{} { return anyIndex(v, 0) })
}

// KRight :: p[a] -> p[b] -> p[b]
func KRight(p1, p2 Parser) Parser {
	return Apply(Seq(p1, p2), func(v interface{}) interface{} { return anyIndex(v, 1) })
}

// KMid :: p[a] -> p[b] -> p[c] -> p[b]
func KMid(p1, p2, p3 Parser) Parser {
	return Apply(Seq(p1, p2, p3), func(v interface{}) interface{} { return anyIndex(v, 1) })
}
