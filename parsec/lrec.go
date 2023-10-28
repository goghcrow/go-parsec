package parsec

// ----------------------------------------------------------------
// Left Recursive
// ----------------------------------------------------------------

// LRec :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 返回多种可能的结果
// P  → Pq|p
// 消除直接左递归
// P  → pP’
// P’ → qP’|ε
// pq, pqq, pqqq, ...
func LRec[K TK /*Fst extends R,*/, Sec, R any](
	p Parser[K, R],
	q Parser[K, Sec],
	f func(R, Sec) R,
) Parser[K, R] {
	return Apply(Seq2(p, Rep(q)), applyLRec(f))
}

// LRecSc :: p[a] -> p[b] -> ((a b) -> c) -> p[c]
// Equivalent to seq(a, rep_sc(b)))
// Returns the result of f(f(f(a, b1), b2), b3) .... If no b succeeds, it returns a
// 只返回一个消费尽可能多 token 的结果
func LRecSc[K TK /*Fst extends R,*/, Sec, R any](
	p Parser[K, R],
	q Parser[K, Sec],
	f func(R, Sec) R,
) Parser[K, R] {
	return Apply(Seq2(p, RepSc(q)), applyLRec(f))
}

// applyLRec :: (a -> b -> c) -> ((a, list[b]) -> a)
func applyLRec[Sec, R any](f func(R, Sec) R) func(Cons[R, []Sec]) R {
	return func(v Cons[R, []Sec]) R {
		x := v.Car
		for _, tail := range v.Cdr {
			x = f(x, tail)
		}
		return x
	}
}

// Chainl 构造左结合双目运算符解析, 可以用来处理左递归文法
// parse >=0 次被 op 分隔的 p, 返回左结合调用 f 得到的值, 如果 0 次, 返回默认值 x
// chainr1 p op <|> return x
func Chainl[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
	x R,
) Parser[K, R] {
	return Alt(Chainl1(p, op), Succ[K, R](x))
}

// Chainl1 构造左结合双目运算符解析, 可以用来处理左递归文法
// parse >=1 次被 op 分隔的 p, 返回左结合调用 f 得到的值
// do { x <- p; rest x } where rest x = do{ f <- op ; y <- p ; rest (f x y) } <|> return x
func Chainl1[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
) Parser[K, R] {
	var chain1Rest func(lval R) Parser[K, R]
	chain1Rest = func(lval R) Parser[K, R] {
		opv := Combine2(op, func(f func(R, R) R) Parser[K, R] {
			// 左结合: 优先匹配 p(即 term), 然后递归的匹配 term op
			return Combine2(p, func(rval R) Parser[K, R] {
				return chain1Rest(f(lval, rval))
			})
		})
		return Alt(opv, Succ[K, R](lval))
	}
	return Combine2(p, chain1Rest)
}

// Chainr 构造右结合双目运算符解析
// parse >=0 次被 op 分隔的 p, 返回右结合调用 f 得到的值, 如果 0 次, 返回默认值 x
// chainr1 p op <|> return x
func Chainr[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
	x R,
) Parser[K, R] {
	return Alt(Chainr1(p, op), Succ[K, R](x))
}

// Chainr1 构造右结合双目运算符解析
// parse >=1 次被 op 分隔的 p, 返回右结合调用 f 得到的值
// do{ x <- p; rest x } where rest x = do{ f <- op ; y <- scan ; return (f x y)  } <|> return x
func Chainr1[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
) Parser[K, R] {
	return Combine2(p, func(lval R) Parser[K, R] {
		seq := Combine2(op, func(f func(R, R) R) Parser[K, R] {
			// 右结合就是自然地递归下降
			return Combine2(Chainr1(p, op), func(rval R) Parser[K, R] {
				return Succ[K, R](f(lval, rval))
			})
		})
		return Alt(seq, Succ[K, R](lval))
	})
}

// ChainlSc 构造左结合双目运算符解析, 可以用来处理左递归文法
// parse >=0 次被 op 分隔的 p, 返回左结合调用 f 得到的值, 如果 0 次, 返回默认值 x
// chainr1 p op <|> return x
func ChainlSc[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
	x R,
) Parser[K, R] {
	return AltSc(Chainl1Sc(p, op), Succ[K, R](x))
}

// Chainl1Sc 构造左结合双目运算符解析, 可以用来处理左递归文法
// parse >=1 次被 op 分隔的 p, 返回左结合调用 f 得到的值
// do { x <- p; rest x } where rest x = do{ f <- op ; y <- p ; rest (f x y) } <|> return x
func Chainl1Sc[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
) Parser[K, R] {
	var chain1Rest func(lval R) Parser[K, R]
	chain1Rest = func(lval R) Parser[K, R] {
		opv := Combine2(op, func(f func(R, R) R) Parser[K, R] {
			// 左结合: 优先匹配 p(即 term), 然后递归的匹配 term op
			return Combine2(p, func(rval R) Parser[K, R] {
				return chain1Rest(f(lval, rval))
			})
		})
		return AltSc(opv, Succ[K, R](lval))
	}
	return Combine2(p, chain1Rest)
}

// ChainrSc 构造右结合双目运算符解析
// parse >=0 次被 op 分隔的 p, 返回右结合调用 f 得到的值, 如果 0 次, 返回默认值 x
// chainr1 p op <|> return x
func ChainrSc[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
	x R,
) Parser[K, R] {
	return AltSc(Chainr1Sc(p, op), Succ[K, R](x))
}

// Chainr1Sc 构造右结合双目运算符解析
// parse >=1 次被 op 分隔的 p, 返回右结合调用 f 得到的值
// do{ x <- p; rest x } where rest x = do{ f <- op ; y <- scan ; return (f x y)  } <|> return x
func Chainr1Sc[K TK, R any](
	p Parser[K, R],
	op Parser[K, func(R, R) R],
) Parser[K, R] {
	return Combine2(p, func(lval R) Parser[K, R] {
		seq := Combine2(op, func(f func(R, R) R) Parser[K, R] {
			// 右结合就是自然地递归下降
			return Combine2(Chainr1Sc(p, op), func(rval R) Parser[K, R] {
				return Succ[K, R](f(lval, rval))
			})
		})
		return AltSc(seq, Succ[K, R](lval))
	})
}
