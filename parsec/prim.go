package parsec

import "fmt"

// ----------------------------------------------------------------
// Token Filter
// ----------------------------------------------------------------

// Nil
// 不消耗 token, 返回 nil
func Nil[K Ord, R any]() Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		return success([]Result[K, R]{{next: toks}})
	})
}

// Succ
// 即 Unit, Return, 不消耗 token, 返回固定值
func Succ[K Ord, R any](v R) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		return success([]Result[K, R]{{Val: v, next: toks}})
	})
}

// Fail
// 不消耗 token, 永远失败
func Fail[K Ord, R any](msg string) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		var pos Pos = EOFPos
		if len(toks) != 0 {
			pos = toks[0]
		}
		return newOutput[K, R]([]Result[K, R]{}, newError(pos, msg), false)
	})
}

// Any
// 消耗任意一个 token
func Any[K Ord]() Parser[K, Token[K]] {
	return parser[K, Token[K]](func(toks []Token[K]) Output[K, Token[K]] {
		if len(toks) == 0 {
			return fail[K, Token[K]](unableToConsumeToken(EOFToken[K](), "any token"))
		}
		return success([]Result[K, Token[K]]{{Val: toks[0], next: toks[1:]}})
	})
}

// Str
// 按 文本匹配 token
func Str[K Ord](toMatch string) Parser[K, Token[K]] {
	return parser[K, Token[K]](func(toks []Token[K]) Output[K, Token[K]] {
		if len(toks) == 0 {
			return fail[K, Token[K]](unableToConsumeToken(EOFToken[K](), toMatch))
		}
		if toks[0].Lexeme() != toMatch {
			return fail[K, Token[K]](unableToConsumeToken(toks[0], toMatch))
		}
		return success([]Result[K, Token[K]]{{Val: toks[0], next: toks[1:]}})
	})
}

// Tok
// 按 TokenKind 匹配 token
func Tok[K Ord](toMatch K) Parser[K, Token[K]] {
	return parser[K, Token[K]](func(toks []Token[K]) Output[K, Token[K]] {
		if len(toks) == 0 {
			return fail[K, Token[K]](unableToConsumeToken(EOFToken[K](), fmt.Sprintf("token<%v>", toMatch)))
		}
		if toks[0].Kind() != toMatch {
			return fail[K, Token[K]](unableToConsumeToken(toks[0], fmt.Sprintf("token<%v>", toMatch)))
		}
		return success([]Result[K, Token[K]]{{Val: toks[0], next: toks[1:]}})
	})
}
