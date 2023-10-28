package parsec

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------
// Alias & Other Combinators
// ----------------------------------------------------------------

// 在 ts 版本基础上, 扩展了部分新的 combinator 和 别名

func Unit[K TK, R any](v R) Parser[K, R]   { return Succ[K, R](v) }
func Return[K TK, R any](v R) Parser[K, R] { return Succ[K, R](v) }

func Label[K TK, R any](p Parser[K, R], msg string) Parser[K, R] { return Err(p, msg) }

// Try 支持 lookaheadN
// 错误发生时不消耗 state, 其他跟 p 一样, TrySc 是传统的 parsec 的 Try 语义
func Try[K TK, R any](p Parser[K, R]) Parser[K, R]   { return Opt(p) }
func TrySc[K TK, R any](p Parser[K, R]) Parser[K, R] { return OptSc(p) }

func Map[K TK, F, T any](p Parser[K, F], f func(v F) T) Parser[K, T] { return Apply(p, f) }

// Bind :: p[a] -> (a->p[b]) -> p[b]
func Bind[K TK, R1, R2 any](p Parser[K, R1], k func(R1) Parser[K, R2]) Parser[K, R2] {
	return Combine2(p, k)
}
func FlatMap[K TK, R1, R2 any](p Parser[K, R1], k func(R1) Parser[K, R2]) Parser[K, R2] {
	return Combine2(p, k)
}

// Lazy :: (() -> p[a]) -> p[a]
func Lazy[K TK, R any](thunk func() Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		return thunk().Parse(toks)
	})
}

// Between
// do{ _ <- open; x <- p; _ <- close; return x }
// e.g. braces  = between (symbol "{") (symbol "}")
func Between[K TK, A, B, C any](open Parser[K, A], p Parser[K, B], close Parser[K, C]) Parser[K, B] {
	return KMid(open, p, close)
}

func Count[K TK, R any](p Parser[K, R], cnt int) Parser[K, []R] { return RepN(p, cnt) }

func Many[K TK, R any](p Parser[K, R]) Parser[K, []R]   { return Rep(p) }
func ManyR[K TK, R any](p Parser[K, R]) Parser[K, []R]  { return RepR(p) }
func ManySc[K TK, R any](p Parser[K, R]) Parser[K, []R] { return RepSc(p) }

func Many1[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, Many[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}
func Many1R[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, ManyR[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}
func Many1Sc[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return Apply(Seq2[K, R](p, ManySc[K, R](p)), func(c Cons[R, []R]) []R {
		return concat([]R{c.Car}, c.Cdr...)
	})
}

func Skip[K TK, R any](p Parser[K, R]) Parser[K, R] {
	return Apply[K, R, R](Opt(p), func(v R) R { return *new(R) })
}
func SkipSc[K TK, R any](p Parser[K, R]) Parser[K, R] {
	return Apply[K, R, R](OptSc(p), func(v R) R { return *new(R) })
}

// SkipMany 应用 p >= 0 次, 跳过结果
// do{ _ <- many p; return ()} <|> return ()
func SkipMany[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return Skip(Many1(p)) // skip 已经包换 nil, 所以这里用 Many1, many 会重复[]
}
func SkipManyR[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	// return Skip(Many1R(p)) // 这样顺序不对, 需要把 nil 提到前面
	alt := Alt(Nil[K, []R](), Many1R(p))
	return Apply(alt, func(v []R) []R { return []R{} })
}
func SkipManySc[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return SkipSc(ManySc(p))
}

// SkipMany1 应用 p >= 1 次, 跳过结果
// 注意 Skip(Many1(p)) != SkipMany1(p)
// do{ _ <- p; skipMany p }
func SkipMany1[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipMany(p))
}
func SkipMany1R[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipManyR(p))
}
func SkipMany1Sc[K TK, R any](p Parser[K, R]) Parser[K, []R] {
	return KRight(p, SkipManySc(p))
}

func Trim[K TK, C, R any](p Parser[K, R], cut Parser[K, C]) Parser[K, R] {
	return KMid(Many(cut), p, Many(cut))
}
func TrimSc[K TK, C, R any](p Parser[K, R], cut Parser[K, C]) Parser[K, R] {
	return KMid(ManySc(cut), p, ManySc(cut))
}

// SepBy p 被 sep 分隔的 >=0 个 p, 不以 seq 结尾
// sepBy1 p sep <|> return []
func SepBy[K TK, S, R any](p Parser[K, R], sep Parser[K, S]) Parser[K, []R] {
	return Opt(SepBy1(p, sep))
}

// SepBy1 p 被 sep 分隔的 >=1 个 p, 不以 seq 结尾
// do{ x <- p; xs <- many (sep >> p); return (x:xs) }
func SepBy1[K TK, S, R any](p Parser[K, R], sep Parser[K, S]) Parser[K, []R] {
	return List(p, sep)
}

// SepBySc p 被 sep 分隔的 >=0 个 p, 不以 seq 结尾
// sepBy1 p sep <|> return []
func SepBySc[K TK, S, R any](p Parser[K, R], sep Parser[K, S]) Parser[K, []R] {
	return OptSc(SepBy1Sc(p, sep))
}

// SepBy1Sc p 被 sep 分隔的 >=1 个 p, 不以 seq 结尾
// do{ x <- p; xs <- many (sep >> p); return (x:xs) }
func SepBy1Sc[K TK, S, R any](p Parser[K, R], sep Parser[K, S]) Parser[K, []R] {
	return ListSc(p, sep)
}

// LookAhead
// peek p 的值, 如果失败会消费 token, 如果不期望消费可以 LookAhead(Try(p))
func LookAhead[K TK, R any](p Parser[K, R]) Parser[K, []R] {
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
func NotFollowedBy[K TK, R any](p Parser[K, R]) Parser[K, R] {
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
