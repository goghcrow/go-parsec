package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Repetitive
// ----------------------------------------------------------------

// Rep :: p[a] -> p[list[a]]
// 重复 n 次(n>=0), 按路径从短到长返回结果
func Rep(p Parser) Parser {
	p = RepR(p)
	return newParser(func(toks []*lexer.Token) Output {
		out := p.Parse(toks)
		if out.Success {
			return successWithErr(reverse(out.Candidates), out.Error)
		}
		return out
	})
}

// RepSc :: p[a] -> p[list[a]]
// 消费尽可能多的 p, 如果零次, 则返回 p[empty_list], 不会失败
func RepSc(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		// 层序遍历, 每层更新结果(从 root 到该层节点的路径), 返回最后一层的结果(根节点到叶子节点路径)
		// Candidates 为每个节点的分叉数
		xs := []Result{{Val: anySlice(), toks: toks}} // root
		for {
			steps := xs
			xs = []Result{}
			for _, step := range steps {
				out := p.Parse(step.toks)
				err = betterError(err, out.Error)
				if out.Success {
					for _, candidate := range out.Candidates {
						// 必须消费掉 token, 重复 nil 死循环
						if !step.toks.equals(candidate.toks) {
							xs = append(xs, Result{
								Val:  append(step.Val.([]interface{}), candidate.Val),
								toks: candidate.toks,
							})
						}
					}
				}
			}
			// 无下一层, 说明 steps 是最后一层, 循环出口
			if len(xs) == 0 {
				xs = steps
				break
			}
		}
		return successWithErr(xs, err)
	})
}

// RepR :: p[a] -> p[list[a]]
// 重复 n 次(n>=0), 按路径从长到短返回结果
func RepR(p Parser) Parser {
	return newParser(func(toks []*lexer.Token) Output {
		var err *Error
		// 层序遍历, 穷举所有根节点到非根节点的路径, Candidates 为每个节点的分叉数
		xs := []Result{{Val: anySlice(), toks: toks}} // root
		for i := 0; i < len(xs); i++ {
			step := xs[i]
			out := p.Parse(step.toks)
			err = betterError(err, out.Error)
			if out.Success {
				for _, candidate := range out.Candidates {
					// 必须消费掉 token, 重复 nil 死循环
					if !step.toks.equals(candidate.toks) {
						xs = append(xs, Result{
							Val:  append(step.Val.([]interface{}), candidate.Val),
							toks: candidate.toks,
						})
					}
				}
			}
		}
		return successWithErr(xs, err)
	})
}
