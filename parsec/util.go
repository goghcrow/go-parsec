package parsec

// ----------------------------------------------------------------
// Model
// ----------------------------------------------------------------

type Cons[T1, T2 any] struct {
	Car T1
	Cdr T2
}

func C2ar[T1, T2 any](t2 Cons[T1, T2]) T1                                           { return t2.Car }
func C2dr[T1, T2 any](t2 Cons[T1, T2]) T2                                           { return t2.Cdr }
func C3ar[T1, T2, T3 any](t3 Cons[T1, Cons[T2, T3]]) T1                             { return t3.Car }
func C3adr[T1, T2, T3 any](t3 Cons[T1, Cons[T2, T3]]) T2                            { return t3.Cdr.Car }
func C3ddr[T1, T2, T3 any](t3 Cons[T1, Cons[T2, T3]]) T3                            { return t3.Cdr.Cdr }
func C4ar[T1, T2, T3, T4 any](t4 Cons[T1, Cons[T2, Cons[T3, T4]]]) T1               { return t4.Car }
func C4adr[T1, T2, T3, T4 any](t4 Cons[T1, Cons[T2, Cons[T3, T4]]]) T2              { return t4.Cdr.Car }
func C4addr[T1, T2, T3, T4 any](t4 Cons[T1, Cons[T2, Cons[T3, T4]]]) T3             { return t4.Cdr.Cdr.Car }
func C4dddr[T1, T2, T3, T4 any](t4 Cons[T1, Cons[T2, Cons[T3, T4]]]) T4             { return t4.Cdr.Cdr.Cdr }
func C5ar[T1, T2, T3, T4, T5 any](t5 Cons[T1, Cons[T2, Cons[T3, Cons[T4, T5]]]]) T1 { return t5.Car }
func C5adr[T1, T2, T3, T4, T5 any](t5 Cons[T1, Cons[T2, Cons[T3, Cons[T4, T5]]]]) T2 {
	return t5.Cdr.Car
}
func C5addr[T1, T2, T3, T4, T5 any](t5 Cons[T1, Cons[T2, Cons[T3, Cons[T4, T5]]]]) T3 {
	return t5.Cdr.Cdr.Car
}
func C5adddr[T1, T2, T3, T4, T5 any](t5 Cons[T1, Cons[T2, Cons[T3, Cons[T4, T5]]]]) T4 {
	return t5.Cdr.Cdr.Cdr.Car
}
func C5addddr[T1, T2, T3, T4, T5 any](t5 Cons[T1, Cons[T2, Cons[T3, Cons[T4, T5]]]]) T5 {
	return t5.Cdr.Cdr.Cdr.Cdr
}

type Either[L, R any] struct {
	isLeft bool
	Left   L
	Right  R
}

func (e Either[L, R]) IsLeft() bool    { return e.isLeft }
func Left[L, R any](v L) Either[L, R]  { return Either[L, R]{isLeft: true, Left: v} }
func Right[L, R any](v R) Either[L, R] { return Either[L, R]{Right: v} }

type Option[T any] struct {
	ok bool
	V  T
}

func Some[T any](v T) Option[T] { return Option[T]{ok: true, V: v} }
func None[T any]() Option[T]    { return Option[T]{} }

// ----------------------------------------------------------------
// Output
// ----------------------------------------------------------------

func fail[K TK, R any](err *Error) Output[K, R] {
	return Output[K, R]{Success: false, Error: err}
}
func success[K TK, R any](xs []Result[K, R]) Output[K, R] {
	return Output[K, R]{Success: true, Candidates: xs}
}
func successWithErr[K TK, R any](xs []Result[K, R], err *Error) Output[K, R] {
	return Output[K, R]{true, xs, err}
}
func newOutput[K TK, R any](xs []Result[K, R], err *Error, success bool) Output[K, R] {
	if success {
		return successWithErr(xs, err)
	} else {
		return fail[K, R](err)
	}
}

func resultOf[K TK, TFrom, TTo any](f func(TFrom) TTo) func(from Result[K, TFrom]) Result[K, TTo] {
	return func(from Result[K, TFrom]) Result[K, TTo] {
		return Result[K, TTo]{
			Val:  f(from.Val),
			next: from.next,
		}
	}
}
func failOf[K TK, TFrom, TTo any](from Output[K, TFrom]) Output[K, TTo] {
	return Output[K, TTo]{
		Success: false,
		Error:   from.Error,
	}
}

// ----------------------------------------------------------------
// Error
// ----------------------------------------------------------------

func newError(pos Pos, msg string) *Error {
	return &Error{Pos: pos, Msg: msg}
}
func unableToConsumeToken[K TK](tok Token[K], expect string) *Error {
	var pos Pos = tok
	if vt, ok := tok.(virtualToken[K]); ok {
		pos = vt.VirtualPos
	}
	if pos == EOFPos {
		return newError(pos, "Nothing to consume expect `"+expect+"`")
	} else {
		return newError(pos, "Unable to consume token `"+tok.String()+"` expect `"+expect+"`")
	}
}

// 返回最远的错误
func betterError(e1, e2 *Error) *Error {
	if e1 == nil {
		return e2
	}
	if e2 == nil {
		return e1
	}
	if e1.Pos == EOFPos {
		return e1
	}
	if e2.Pos == EOFPos {
		return e2
	}
	idx1, _, _, _ := e1.Loc()
	idx2, _, _, _ := e2.Loc()
	if idx1 < idx2 {
		return e2
	}
	return e1
}

// ----------------------------------------------------------------
// Tokens
// ----------------------------------------------------------------

func beginTok[K TK](t []Token[K]) Token[K] {
	if len(t) == 0 {
		return nil
	} else {
		return t[0]
	}
}

func beginPos[K TK](t []Token[K]) Pos {
	if len(t) == 0 {
		return UnknownPos
	} else {
		return t[0]
	}
}

func toksEqual[K TK](t []Token[K], other []Token[K]) bool {
	if len(t) == 0 && len(other) == 0 {
		return true
	}
	if len(t) == 0 || len(other) == 0 || t[0] != other[0] {
		return false
	}
	return true
}

func tokenRange[K TK](seq []Token[K], nxt Token[K]) []Token[K] {
	if len(seq) == 0 {
		return nil
	}
	if nxt == nil {
		return seq
	}
	for i, tok := range seq {
		if tok == nxt {
			return seq[:i]
		}
	}
	panic("unreached")
}

// ----------------------------------------------------------------
// Other
// ----------------------------------------------------------------

func reverse[K TK, R any](s []Result[K, R]) []Result[K, R] {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func concat[T any](x []T, y ...T) []T {
	xs := make([]T, len(x)+len(y))
	copy(xs, x)
	copy(xs[len(x):], y)
	return xs
}

func sliceMap[TFrom, TTo any](s []TFrom, f func(TFrom) TTo) []TTo {
	t := make([]TTo, len(s))
	for i, v := range s {
		t[i] = f(v)
	}
	return t
}

func foldLeft[A, B any](xs []A, z B, op func(B, A) B) B {
	for _, x := range xs {
		z = op(z, x)
	}
	return z
}

func foldRight[A, B any](xs []A, z B, op func(A, B) B) B {
	for i := len(xs) - 1; i >= 0; i-- {
		z = op(xs[i], z)
	}
	return z
}
