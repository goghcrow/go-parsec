package parsec

import (
	"github.com/goghcrow/go-parsec/lexer"
)

// ----------------------------------------------------------------
// Parser Impl
// ----------------------------------------------------------------

type parser struct {
	parse func([]*lexer.Token) Output
}

func newParser(p func([]*lexer.Token) Output) Parser           { return &parser{p} }
func (p *parser) Parse(toks []*lexer.Token) Output             { return p.parse(toks) }
func (p *parser) Map(f func(v interface{}) interface{}) Parser { return Apply(p, f) }

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

func (t tokSeq) loc() lexer.Loc {
	if len(t) == 0 {
		return lexer.UnknownLoc
	} else {
		return t[0].Loc
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
	if e1.Loc == lexer.UnknownLoc { // eof
		return e1
	}
	if e2.Loc == lexer.UnknownLoc { // eof
		return e2
	}
	if e1.Loc.Pos < e2.Loc.Pos {
		return e2
	} else if e1.Loc.Pos > e2.Loc.Pos {
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

func newError(loc lexer.Loc, msg string) *Error { return &Error{Loc: loc, Msg: msg} }
func unableToConsumeToken(tok *lexer.Token) *Error {
	return &Error{
		Loc: tok.Loc,
		Msg: "Unable to consume token `" + tok.String() + "`",
	}
}

// ----------------------------------------------------------------
// Other
// ----------------------------------------------------------------

var eof = &lexer.Token{
	TokenKind: -1,
	Loc:       lexer.UnknownLoc,
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
