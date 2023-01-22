package parsec

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/goghcrow/lexer"
)

const (
	Number lexer.TokenKind = iota + 1
	Add
	Space
	Ident
	Comma
)

func stroftk(k lexer.TokenKind) string {
	return map[lexer.TokenKind]string{
		Number: "<num>",
		Add:    "+",
		Space:  "<space>",
		Ident:  "<id>",
		Comma:  ",",
	}[k]
}

var lex = lexer.BuildLexer(func(lex *lexer.Lexicon) {
	lex.Regex(Number, "\\d+")
	lex.Regex(Ident, "[a-zA-Z]\\w*")
	lex.Regex(Space, "\\s+").Skip()
	lex.Str(Comma, ",").Skip()
	lex.Str(Add, "+")
})

func TestParser(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		p       Parser
		success bool
		result  string
		error   string
	}{
		{
			name:    "Parser: str",
			input:   "123,456",
			p:       Str("123"),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: str",
			input:   "123,456",
			p:       Str("456"),
			success: false,
			error:   "Unable to consume token `123` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: tok",
			input:   "123,456",
			p:       Tok(Number),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: alt",
			input:   "123,456",
			p:       Alt(Tok(Number), Tok(Ident)),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       Seq(Tok(Number), Tok(Ident)),
			success: false,
			error:   "Unable to consume token `456` in pos 5-8 line 1 col 5",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       Seq(Tok(Number), Tok(Number)),
			success: true,
			result:  "{v=[123 456], next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KLeft(Tok(Number), Seq(Tok(Number), Tok(Number))),
			success: true,
			result:  "{v=123, next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KMid(Tok(Number), Tok(Number), Tok(Number)),
			success: true,
			result:  "{v=456, next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KRight(Seq(Tok(Number), Tok(Number)), Tok(Number)),
			success: true,
			result:  "{v=789, next=}",
		},
		{
			name:    "Parser: opt",
			input:   "123,456",
			p:       Opt(Tok(Number)),
			success: true,
			result:  "{v=123, next=<num>/456}🍊{v=<nil>, next=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       OptSc(Tok(Number)),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       OptSc(Tok(Ident)),
			success: true,
			result:  "{v=<nil>, next=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       RepSc(Tok(Number)),
			success: true,
			result:  "{v=[123 456], next=}",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       RepSc(Tok(Ident)),
			success: true,
			result:  "{v=[], next=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: repr",
			input:   "123,456",
			p:       RepR(Tok(Number)),
			success: true,
			result:  "{v=[], next=<num>/123🍌<num>/456}🍊{v=[123], next=<num>/456}🍊{v=[123 456], next=}",
		},
		{
			name:    "Parser: rep",
			input:   "123,456",
			p:       Rep(Tok(Number)),
			success: true,
			result:  "{v=[123 456], next=}🍊{v=[123], next=<num>/456}🍊{v=[], next=<num>/123🍌<num>/456}",
		},
		{
			name:  "Parser: apply",
			input: "123,456",
			p: RepR(Tok(Number)).Map(func(toks interface{}) interface{} {
				var xs []string
				for _, v := range toks.([]interface{}) {
					xs = append(xs, v.(*lexer.Token).Lexeme)
				}
				return strings.Join(xs, ";")
			}),
			success: true,
			result:  "{v=, next=<num>/123🍌<num>/456}🍊{v=123, next=<num>/456}🍊{v=123;456, next=}",
		},
		{
			name:  "Parser: errd",
			input: "a",
			p: ErrDef(Tok(Number).Map(func(v interface{}) interface{} {
				num, _ := strconv.Atoi(v.(*lexer.Token).Lexeme)
				return num
			}), "This is not a number!", 42),
			success: true,
			result:  "{v=42, next=<id>/a}",
			error:   "This is not a number! in pos 1-2 line 1 col 1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := lex.MustLex(tt.input)
			out := tt.p.Parse(toks)
			if tt.success {
				xs := succeed(out)
				actual := fmtResults(xs)
				if actual != tt.result {
					t.Errorf("expect %s actual %s", tt.result, actual)
				}
				if tt.error != "" {
					actual = out.Error.Error()
					if actual != tt.error {
						t.Errorf("expect %s actual %s", tt.error, actual)
					}
				}
			} else {
				if out.Success {
					t.Errorf("expect fail actual success")
				}
				actual := out.Error.Error()
				if actual != tt.error {
					t.Errorf("expect %s actual %s", tt.error, actual)
				}
			}
		})
	}
}

func TestAmbParser(t *testing.T) {
	TERM := NewRule()
	EXPR := NewRule()

	// TERM
	//		= NUMBER
	//		= + EXPR
	// EXPR
	//		= TERM
	//		= EXPR | (+ EXPR)
	TERM.Pattern = Alt(
		Tok(Number).Map(func(v interface{}) interface{} { return v.(*lexer.Token).Lexeme }),
		KRight(Str("+"), EXPR).Map(func(v interface{}) interface{} { return fmt.Sprintf("(+ %s)", v) }),
	)
	EXPR.Pattern = Amb(
		LRecSc(
			TERM,
			Alt(EXPR, Seq(Str("+"), EXPR)),
			func(a, b interface{}) interface{} {
				s := a.(string)
				t, ok := b.(string)
				if ok {
					return fmt.Sprintf(`(%s . %s)`, s, t)
				} else {
					t := b.([]interface{}) // [token, string]
					return fmt.Sprintf(`(%s + %s)`, s, t[1])
				}
			},
		)).
		Map(func(v interface{}) interface{} {
			var ss []string
			for _, v := range v.([]interface{}) {
				ss = append(ss, v.(string))
			}
			if len(ss) == 1 {
				return ss[0]
			} else {
				return "[" + strings.Join(ss, ", ") + "]"
			}
		})

	for _, tt := range []struct {
		name    string
		input   string
		p       Parser
		success bool
		result  string
		error   string
	}{
		{
			name:    "Parser: amb, 1",
			input:   "1",
			p:       EXPR,
			success: true,
			result:  "{v=1, next=}",
		},
		{
			name:    "Parser: amb, +1",
			input:   "+1",
			p:       EXPR,
			success: true,
			result:  "{v=(+ 1), next=}",
		},
		{
			name:    "Parser: amb, 1+2",
			input:   "1+2",
			p:       EXPR,
			success: true,
			result:  "{v=[(1 . (+ 2)), (1 + 2)], next=}",
		},
		{
			name:    "Parser: amb, 1+2+3",
			input:   "1+2+3",
			p:       EXPR,
			success: true,
			result:  "{v=[(1 . (+ [(2 . (+ 3)), (2 + 3)])), (1 + [(2 . (+ 3)), (2 + 3)])], next=}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := lex.MustLex(tt.input)
			out := tt.p.Parse(toks)
			if tt.success {
				xs := succeed(out)
				actual := fmtResults(xs)
				if actual != tt.result {
					t.Errorf("expect %s actual %s", tt.result, actual)
				}
				if tt.error != "" {
					actual = out.Error.Error()
					if actual != tt.error {
						t.Errorf("expect %s actual %s", tt.error, actual)
					}
				}
			} else {
				if out.Success {
					t.Errorf("expect fail actual success")
				}
				actual := out.Error.Error()
				if actual != tt.error {
					t.Errorf("expect %s actual %s", tt.error, actual)
				}
			}
		})
	}
}

func TestFailure(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		p       Parser
		success bool
		result  string
		error   string
	}{
		{
			name:  "Failure: alt",
			input: "123,456",
			p: Alt(
				Tok(Comma),
				Tok(Space),
			),
			success: false,
			error:   "Unable to consume token `123` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Failure: seq",
			input: "123,456",
			p: Seq(
				Tok(Ident),
				Tok(Number),
			),
			success: false,
			error:   "Unable to consume token `123` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Failure: seq",
			input: "123,456",
			p: Seq(
				Tok(Number),
				Tok(Ident),
			),
			success: false,
			error:   "Unable to consume token `456` in pos 5-8 line 1 col 5",
		},
		{
			name:  "Failure: apply",
			input: "123,456",
			p: Apply(Tok(Comma), func(v interface{}) interface{} {
				return nil
			}),
			success: false,
			error:   "Unable to consume token `123` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Failure: rep_sc + seq",
			input:   "1a 2b 3c d e",
			p:       RepSc(Seq(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=[[1 a] [2 b] [3 c]], next=<id>/d🍌<id>/e}",
			error:   "Unable to consume token `d` in pos 10-11 line 1 col 10",
		},
		{
			name:    "Failure: rep_sc + seq",
			input:   "1a 2b 3c d e",
			p:       Rep(Seq(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=[[1 a] [2 b] [3 c]], next=<id>/d🍌<id>/e}🍊{v=[[1 a] [2 b]], next=<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}🍊{v=[[1 a]], next=<num>/2🍌<id>/b🍌<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}🍊{v=[], next=<num>/1🍌<id>/a🍌<num>/2🍌<id>/b🍌<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}",
			// 返回最远的错误
			error: "Unable to consume token `d` in pos 10-11 line 1 col 10",
		},
		{
			name:    "Failure: rep_sc + alt",
			input:   "1 a b 2 c 3",
			p:       RepSc(Alt(Tok(Number), Seq(Tok(Ident), Tok(Ident)))),
			success: true,
			result:  "{v=[1 [a b] 2], next=<id>/c🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` in pos 11-12 line 1 col 11",
		},
		{
			name:    "Failure: rep_sc + alt",
			input:   "1 a b 2 c 3",
			p:       Rep(Alt(Tok(Number), Seq(Tok(Ident), Tok(Ident)))),
			success: true,
			result:  "{v=[1 [a b] 2], next=<id>/c🍌<num>/3}🍊{v=[1 [a b]], next=<num>/2🍌<id>/c🍌<num>/3}🍊{v=[1], next=<id>/a🍌<id>/b🍌<num>/2🍌<id>/c🍌<num>/3}🍊{v=[], next=<num>/1🍌<id>/a🍌<id>/b🍌<num>/2🍌<id>/c🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` in pos 11-12 line 1 col 11",
		},
		{
			name:    "Failure: rep_sc + opt",
			input:   "a b c d e f g 3",
			p:       RepSc(OptSc(Seq(Tok(Ident), Tok(Ident)))),
			success: true,
			result:  "{v=[[a b] [c d] [e f]], next=<id>/g🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` in pos 15-16 line 1 col 15",
		},
		{
			name:    "Failure: rep_sc + opt",
			input:   "a b c d e f g 3",
			p:       Rep(OptSc(Seq(Tok(Ident), Tok(Ident)))),
			success: true,
			result:  "{v=[[a b] [c d] [e f]], next=<id>/g🍌<num>/3}🍊{v=[[a b] [c d]], next=<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}🍊{v=[[a b]], next=<id>/c🍌<id>/d🍌<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}🍊{v=[], next=<id>/a🍌<id>/b🍌<id>/c🍌<id>/d🍌<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` in pos 15-16 line 1 col 15",
		},
		{
			name:    "Failure: err",
			input:   "a",
			p:       Err(Tok(Number), "This is not a number!"),
			success: false,
			result:  "",
			error:   "This is not a number! in pos 1-2 line 1 col 1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := lex.MustLex(tt.input)
			out := tt.p.Parse(toks)
			if tt.success {
				xs := succeed(out)
				actual := fmtResults(xs)
				if actual != tt.result {
					t.Errorf("expect %s actual %s", tt.result, actual)
				}
				if tt.error != "" {
					actual = out.Error.Error()
					if actual != tt.error {
						t.Errorf("expect %s actual %s", tt.error, actual)
					}
				}
			} else {
				if out.Success {
					t.Errorf("expect fail actual success")
				}
				actual := out.Error.Error()
				if actual != tt.error {
					t.Errorf("expect %s actual %s", tt.error, actual)
				}
			}
		})
	}
}

func succeed(out Output) []Result {
	if out.Success {
		return out.Candidates
	}
	panic(out)
}

func fmtResults(results []Result) string {
	xs := make([]string, len(results))
	for i, r := range results {
		xs[i] = fmt.Sprintf("{v=%v, next=%s}", r.Val, fmtToks(r.next))
	}
	return strings.Join(xs, "🍊")
}

func fmtToks(toks []*lexer.Token) string {
	xs := make([]string, len(toks))
	for i, t := range toks {
		xs[i] = fmt.Sprintf("%s/%s", stroftk(t.TokenKind), t.Lexeme)
	}
	return strings.Join(xs, "🍌")
}
