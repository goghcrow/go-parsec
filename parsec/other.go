package parsec

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------
// Alias & Other Combinators
// ----------------------------------------------------------------

func Unit[K Ord, R any](v R) Parser[K, R]   { return Succ[K, R](v) }
func Return[K Ord, R any](v R) Parser[K, R] { return Succ[K, R](v) }

func Label[K Ord, R any](p Parser[K, R], msg string) Parser[K, R] { return Err(p, msg) }

// Try 支持 lookaheadN
// 错误发生时不消耗 state, 其他跟 p 一样, TrySc 是传统的 parsec 的 Try 语义
func Try[K Ord, R any](p Parser[K, R]) Parser[K, R]   { return Opt(p) }
func TrySc[K Ord, R any](p Parser[K, R]) Parser[K, R] { return OptSc(p) }

func Map[K Ord, F, T any](p Parser[K, F], f func(v F) T) Parser[K, T] { return Apply(p, f) }

// Bind :: p[a] -> (a->p[b]) -> p[b]
func Bind[K Ord, R1, R2 any](p Parser[K, R1], k func(R1) Parser[K, R2]) Parser[K, R2] {
	return Combine2(p, k)
}
func FlatMap[K Ord, R1, R2 any](p Parser[K, R1], k func(R1) Parser[K, R2]) Parser[K, R2] {
	return Combine2(p, k)
}

// Lazy :: (() -> p[a]) -> p[a]
func Lazy[K Ord, R any](thunk func() Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		return thunk().Parse(toks)
	})
}

// Between
// do{ _ <- open; x <- p; _ <- close; return x }
// e.g. braces  = between (symbol "{") (symbol "}")
func Between[K Ord, A, B, C any](open Parser[K, A], p Parser[K, B], close Parser[K, C]) Parser[K, B] {
	return KMid(open, p, close)
}

func Count[K Ord, R any](p Parser[K, R], cnt int) Parser[K, []R] { return RepN(p, cnt) }

func Many[K Ord, R any](p Parser[K, R]) Parser[K, []R]   { return Rep(p) }
func ManyR[K Ord, R any](p Parser[K, R]) Parser[K, []R]  { return RepR(p) }
func ManySc[K Ord, R any](p Parser[K, R]) Parser[K, []R] { return RepSc(p) }

func Many1[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, Many[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}
func Many1R[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, ManyR[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}
func Many1Sc[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, ManySc[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}

func Skip[K Ord, R any](p Parser[K, R]) Parser[K, R] {
	return Apply[K, R, R](Opt(p), func(v R) R { return *new(R) })
}
func SkipSc[K Ord, R any](p Parser[K, R]) Parser[K, R] {
	return Apply[K, R, R](OptSc(p), func(v R) R { return *new(R) })
}

// SkipMany 应用 p >= 0 次, 跳过结果
// do{ _ <- many p; return ()} <|> return ()
func SkipMany[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return Skip(Many1(p)) // skip 已经包换 nil, 所以这里用 Many1, many 会重复[]
}
func SkipManyR[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	// return Skip(Many1R(p)) // 这样顺序不对, 需要把 nil 提到前面
	alt := Alt(Nil[K, []R](), Many1R(p))
	return Apply(alt, func(v []R) []R { return []R{} })
}
func SkipManySc[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return SkipSc(ManySc(p))
}

// SkipMany1 应用 p >= 1 次, 跳过结果
// 注意 Skip(Many1(p)) != SkipMany1(p)
// do{ _ <- p; skipMany p }
func SkipMany1[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipMany(p))
}
func SkipMany1R[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipManyR(p))
}
func SkipMany1Sc[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipManySc(p))
}

func TrimSc[K Ord, R any](p Parser[K, R], cut Parser[K, R]) Parser[K, R] {
	return KMid(ManySc(cut), p, ManySc(cut))
}

// LookAhead
// peek p 的值, 如果失败会消费 token, 如果不期望消费可以 LookAhead(Try(p))
func LookAhead[K Ord, R any](p Parser[K, R]) Parser[K, []R] {
	return parser[K, []R](func(toks []Token[K]) Output[K, []R] {
		out := p.Parse(toks)
		if !out.Success {
			return failOf[K, R, []R](out)
		}
		xs := make([]R, len(out.Candidates))
		for i, candidate := range out.Candidates {
			xs[i] = candidate.Val
		}
		res := []Result[K, []R]{{Val: xs, next: toks}}
		return successWithErr(res, out.Error)
	})
}

// NotFollowedBy 只有在 p 匹配失败时才成功, 不消耗 token, 可以用来实现最长匹配
// 在传统 parsec 中可以用来在识别 keywords,
// e.g. 识别 let 需要确保关键词后面没有合法的标识符(e.g. lets)
// 可以写成 let := Left(Str("let"), NotFollowedBy(Regex(`[\d\w]+`)))
// try (do{ c <- try p; unexpected (show c) } <|> return () )
// e.g. KLeft(Tok(Number), NotFollowedBy(Tok(Add)))
func NotFollowedBy[K Ord, R any](p Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		out := p.Parse(toks)
		if !out.Success {
			return success([]Result[K, R]{{next: toks}})
		}
		stringify := func(c Result[K, R]) string { return fmt.Sprintf("`%v`", c.Val) }
		xs := sliceMap(out.Candidates, stringify)
		return fail[K, R](newError(beginPos(toks), "unexpect "+strings.Join(xs, " or ")))
	})
}
