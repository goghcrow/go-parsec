package parsec

import (
	"fmt"
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

func betterError(e1, e2 *Error) *Error {
	if e1 == nil {
		return e2
	}
	if e2 == nil {
		return e1
	}
	if e1.Loc == lexer.UnknownLoc {
		return e1
	}
	if e2.Loc == lexer.UnknownLoc {
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

func fail(err *Error) Output                  { return Output{Success: false, Error: err} }
func success(x Result) Output                 { return Output{Success: true, Candidates: []Result{x}} }
func successX(xs []Result, err *Error) Output { return Output{true, xs, err} }
func resultOrError(result []Result, err *Error, success bool) Output {
	if success {
		return successX(result, err)
	} else {
		return fail(err)
	}
}

// ----------------------------------------------------------------
// Error
// ----------------------------------------------------------------

func (e *Error) Error() string { return fmt.Sprintf("%s in %s", e.Msg, e.Loc) }

func newError(loc lexer.Loc, msg string) *Error { return &Error{Loc: loc, Msg: msg} }
func unableToConsumeToken(tok *lexer.Token) *Error {
	return &Error{
		Loc: tok.Loc,
		Msg: "Unable to consume token " + tok.String(),
	}
}

// ----------------------------------------------------------------
// Other
// ----------------------------------------------------------------

var eof = &lexer.Token{
	TokenKind: "<END-OF-FILE>",
	Loc:       lexer.UnknownLoc,
	Lexeme:    "<END-OF-FILE>",
}

func reverse(s []Result) []Result {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func emptySlice() interface{} { return []interface{}{} }
