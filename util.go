package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// tokSeq
// ----------------------------------------------------------------

type tokSeq []*lexer.Token

func (t tokSeq) mapKey() *lexer.Token {
	if len(t) == 0 {
		return nil
	} else {
		return t[0]
	}
}

func (t tokSeq) pos() lexer.Pos {
	if len(t) == 0 {
		return lexer.UnknownPos
	} else {
		return t[0].Pos
	}
}

func (t tokSeq) equals(other tokSeq) bool {
	if len(t) == 0 && len(other) == 0 {
		return true
	}
	if len(t) == 0 || len(other) == 0 || t[0] != other[0] {
		return false
	}
	return true
}

// 返回最远的错误
func betterError(e1, e2 *Error) *Error {
	if e1 == nil {
		return e2
	}
	if e2 == nil {
		return e1
	}
	if e1.Pos == lexer.UnknownPos { // eof
		return e1
	}
	if e2.Pos == lexer.UnknownPos { // eof
		return e2
	}
	if e1.Pos.Idx < e2.Pos.Idx {
		return e2
	} else if e1.Pos.Idx > e2.Pos.Idx {
		return e1
	}
	return e1
}

// ----------------------------------------------------------------
// Output
// ----------------------------------------------------------------

func fail(err *Error) Output                        { return Output{Success: false, Error: err} }
func success(xs []Result) Output                    { return Output{Success: true, Candidates: xs} }
func successWithErr(xs []Result, err *Error) Output { return Output{true, xs, err} }
func newOutput(xs []Result, err *Error, success bool) Output {
	if success {
		return successWithErr(xs, err)
	} else {
		return fail(err)
	}
}

// ----------------------------------------------------------------
// Error
// ----------------------------------------------------------------

func newError(pos lexer.Pos, msg string) *Error { return &Error{Pos: pos, Msg: msg} }
func unableToConsumeToken(tok *lexer.Token) *Error {
	return &Error{
		Pos: tok.Pos,
		Msg: "Unable to consume token `" + tok.String() + "`",
	}
}

// ----------------------------------------------------------------
// Other
// ----------------------------------------------------------------

var eof = &lexer.Token{
	TokenKind: -1,
	Pos:       lexer.UnknownPos,
	Lexeme:    "<END-OF-FILE>",
}

func reverse(s []Result) []Result {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func anySlice() interface{} { return []interface{}{} }

func anyIndex(v interface{}, i int) interface{} { return v.([]interface{})[i] }
