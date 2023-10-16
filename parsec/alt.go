package parsec

// ----------------------------------------------------------------
// Alternative, Choice
// ----------------------------------------------------------------

// Alt :: p[a] -> p[b] -> p[c] -> ... -> p[a|b|c...]
// 返回所有可能结果, 当 ps 全部失败时失败
// foldr (<|>) mzero ps
func Alt[K TK, R any](ps ...Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		var xs []Result[K, R]
		var err *Error
		var succ bool
		for _, p := range ps {
			out := p.Parse(toks)
			err = betterError(err, out.Error)
			if out.Success {
				xs = append(xs, out.Candidates...)
				succ = true
			}
		}
		return newOutput(xs, err, succ)
	})
}

func Alt2[K TK, R1, R2 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
) Parser[K, Either[R1, R2]] {
	mkLeft := resultOf[K, R1, Either[R1, R2]](Left[R1, R2])
	mkRight := resultOf[K, R2, Either[R1, R2]](Right[R1, R2])
	return parser[K, Either[R1, R2]](func(toks []Token[K]) Output[K, Either[R1, R2]] {
		var xs []Result[K, Either[R1, R2]]

		out1 := p1.Parse(toks)
		if out1.Success {
			xs = append(xs, sliceMap(out1.Candidates, mkLeft)...)
		}

		out2 := p2.Parse(toks)
		if out2.Success {
			xs = append(xs, sliceMap(out2.Candidates, mkRight)...)
		}

		succ := out1.Success || out2.Success
		err := betterError(out1.Error, out2.Error)
		return newOutput(xs, err, succ)
	})
}
func Alt3[K TK, T1, T2, T3 any](
	p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
) Parser[K, Either[T1, Either[T2, T3]]] {
	return Alt2(p1, Alt2(p2, p3))
}
func Alt4[K TK, T1, T2, T3, T4 any](
	p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
	p4 Parser[K, T4],
) Parser[K, Either[T1, Either[T2, Either[T3, T4]]]] {
	return Alt2(p1, Alt3(p2, p3, p4))
}
func Alt5[K TK, T1, T2, T3, T4, T5 any](p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
	p4 Parser[K, T4],
	p5 Parser[K, T5],
) Parser[K, Either[T1, Either[T2, Either[T3, Either[T4, T5]]]]] {
	return Alt2(p1, Alt4(p2, p3, p4, p5))
}

// AltSc :: p[a] -> p[b] -> p[c] -> ... -> p[a|b|c...]
// 返回第一个结果, 当 ps 全部失败时失败
func AltSc[K TK, R any](ps ...Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		var err *Error
		for _, p := range ps {
			out := p.Parse(toks)
			err = betterError(err, out.Error)
			if out.Success {
				return successWithErr[K, R](out.Candidates, err)
			}
		}
		return fail[K, R](err)
	})
}

func AltSc2[K TK, R1, R2 any](
	p1 Parser[K, R1],
	p2 Parser[K, R2],
) Parser[K, Either[R1, R2]] {
	mkLeft := resultOf[K, R1, Either[R1, R2]](Left[R1, R2])
	mkRight := resultOf[K, R2, Either[R1, R2]](Right[R1, R2])
	return parser[K, Either[R1, R2]](func(toks []Token[K]) Output[K, Either[R1, R2]] {
		var err *Error

		out1 := p1.Parse(toks)
		err = betterError(err, out1.Error)
		if out1.Success {
			return successWithErr(sliceMap(out1.Candidates, mkLeft), err)
		}

		out2 := p2.Parse(toks)
		err = betterError(err, out2.Error)
		if out2.Success {
			return successWithErr(sliceMap(out2.Candidates, mkRight), err)
		}

		return fail[K, Either[R1, R2]](err)
	})
}
func AltSc3[K TK, T1, T2, T3 any](
	p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
) Parser[K, Either[T1, Either[T2, T3]]] {
	return AltSc2(p1, AltSc2(p2, p3))
}
func AltSc4[K TK, T1, T2, T3, T4 any](
	p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
	p4 Parser[K, T4],
) Parser[K, Either[T1, Either[T2, Either[T3, T4]]]] {
	return AltSc2(p1, AltSc3(p2, p3, p4))
}
func AltSc5[K TK, T1, T2, T3, T4, T5 any](p1 Parser[K, T1],
	p2 Parser[K, T2],
	p3 Parser[K, T3],
	p4 Parser[K, T4],
	p5 Parser[K, T5],
) Parser[K, Either[T1, Either[T2, Either[T3, Either[T4, T5]]]]] {
	return AltSc2(p1, AltSc4(p2, p3, p4, p5))
}

// ----------------------------------------------------------------
// Optional
// ----------------------------------------------------------------

// Opt :: p[a] -> p[a|nil]
// Alt 返回失败, Opt & OptSc 不返回失败, p 错误不消耗 token
func Opt[K TK, R any](p Parser[K, R]) Parser[K, R /*Option[R]*/] {
	return Alt(p, Nil[K, R]())
}

// OptSc :: p[a] -> p[a|nil]
// Opt 返回两种结果, OptSc 返回一种结果, 只有 p 失败才返回 nil
func OptSc[K TK, R any](p Parser[K, R]) Parser[K, R /*Option[R]*/] {
	return AltSc(p, Nil[K, R]())
}
