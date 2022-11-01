package parsec

import (
	"fmt"
	"github.com/goghcrow/go-parsec/lexer"
	"strconv"
	"strings"
	"testing"
)

func TestRec(t *testing.T) {
	const (
		Number lexer.TokenKind = "'num"
		Add                    = "+"
		Sub                    = "-"
		Mul                    = "*"
		Div                    = "/"
		LParen                 = "("
		RParen                 = ")"
		Space                  = "'space"
	)

	lex := lexer.BuildLexer(func(lex *lexer.Lexicon) {
		lex.Regex(Number, `\d+(\.\d+)?`)
		lex.Oper(Add)
		lex.Oper(Sub)
		lex.Oper(Mul)
		lex.Oper(Div)
		lex.Str(LParen)
		lex.Str(RParen)
		lex.Regex(Space, `\s+`).Skip()
	})

	term := NewRule()
	factor := NewRule()
	exp := NewRule()

	/*
		TERM
		  = NUMBER
		  = ('+' | '-') TERM
		  = '(' EXP ')'
	*/
	str2num := func(s string) float64 {
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return num
	}

	applyNum := func(v interface{}) interface{} {
		return str2num(v.(*lexer.Token).Lexeme)
	}
	applyUnary := func(v interface{}) interface{} {
		xs := v.([]interface{})
		u := xs[0].(*lexer.Token)
		rhs := xs[1].(float64)
		switch u.Lexeme {
		case "+":
			return rhs
		case "-":
			return -rhs
		default:
			panic("unreached")
		}
	}
	applyBinary := func(a, b interface{}) interface{} {
		lhs := a.(interface{}).(float64)
		oper := b.([]interface{})[0].(*lexer.Token)
		rhs := b.([]interface{})[1].(float64)
		switch oper.Lexeme {
		case "+":
			return lhs + rhs
		case "-":
			return lhs - rhs
		case "*":
			return lhs * rhs
		case "/":
			return lhs / rhs
		default:
			panic("unreached")
		}
	}
	term.Pattern = Alt(
		Tok(Number).Map(applyNum),
		Seq(Alt(Str("+"), Str("-")), term).Map(applyUnary),
		KMid(Str("("), exp, Str(")")),
	)
	/*
		FACTOR
		  = TERM
		  = FACTOR ('*' | '/') TERMx
	*/
	factor.Pattern = LRecSc(
		term,
		Seq(Alt(Str("*"), Str("/")), term),
		applyBinary,
	)
	/*
		EXP
		  = FACTOR
		  = EXP ('+' | '-') FACTOR
	*/
	exp.Pattern = LRecSc(
		factor,
		Seq(Alt(Str("+"), Str("-")), factor),
		applyBinary,
	)

	eval := func(s string) float64 {
		toks := lex.Lex(s)
		out := exp.Parse(toks)
		result, err := ExpectSingleResult(ExpectEOF(out))
		if err != nil {
			panic(err)
		}
		return result.(float64)
	}

	for _, tt := range []struct {
		input  string
		expect float64
	}{
		{"1", 1},
		{"+1.5", 1.5},
		{"-0.5", -0.5},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"1 / 2", 0.5},
		{"1 + 2 * 3 + 4", 11},
		{"(1 + 2) * (3 + 4)", 21},
		{"1.2--3.4", 4.6},
	} {
		t.Run(tt.input, func(t *testing.T) {
			v := eval(tt.input)
			if tt.expect != v {
				t.Errorf("expect %f actual %f", tt.expect, v)
			}
		})
	}
}

func TestParser(t *testing.T) {
	const (
		Number lexer.TokenKind = "<num>"
		Ident                  = "<id>"
		Space                  = "<space>"
		Add                    = "+"
		Comma                  = ","
	)

	var lex = lexer.BuildLexer(func(lex *lexer.Lexicon) {
		lex.Regex(Number, "\\d+")
		lex.Regex(Ident, "[a-zA-Z]\\w*")
		lex.Regex(Space, "\\s+").Skip()
		lex.Str(Comma).Skip()
		lex.Str(Add)
	})

	term := NewRule()
	expr := NewRule()

	term.Pattern = Alt(
		Tok(Number).Map(func(v interface{}) interface{} { return v.(*lexer.Token).Lexeme }),
		KRight(Str("+"), expr).Map(func(v interface{}) interface{} { return fmt.Sprintf("(+ %s)", v) }),
	)
	expr.Pattern = Amb(
		LRecSc(
			term,
			Alt(expr, Seq(Str("+"), expr)),
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
		)).Map(func(v interface{}) interface{} {
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
			name:    "Parser: str",
			input:   "123,456",
			p:       Str("123"),
			success: true,
			result:  "{v=123, toks=<num>/456}",
		},
		{
			name:    "Parser: str",
			input:   "123,456",
			p:       Str("456"),
			success: false,
			error:   "Unable to consume token 123 in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: tok",
			input:   "123,456",
			p:       Tok(Number),
			success: true,
			result:  "{v=123, toks=<num>/456}",
		},
		{
			name:    "Parser: alt",
			input:   "123,456",
			p:       Alt(Tok(Number), Tok(Ident)),
			success: true,
			result:  "{v=123, toks=<num>/456}",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       Seq(Tok(Number), Tok(Ident)),
			success: false,
			error:   "Unable to consume token 456 in pos 5-8 line 1 col 5",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       Seq(Tok(Number), Tok(Number)),
			success: true,
			result:  "{v=[123 456], toks=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KLeft(Tok(Number), Seq(Tok(Number), Tok(Number))),
			success: true,
			result:  "{v=123, toks=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KMid(Tok(Number), Tok(Number), Tok(Number)),
			success: true,
			result:  "{v=456, toks=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       KRight(Seq(Tok(Number), Tok(Number)), Tok(Number)),
			success: true,
			result:  "{v=789, toks=}",
		},
		{
			name:    "Parser: opt",
			input:   "123,456",
			p:       Opt(Tok(Number)),
			success: true,
			result:  "{v=123, toks=<num>/456}🍊{v=<nil>, toks=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       OptSc(Tok(Number)),
			success: true,
			result:  "{v=123, toks=<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       OptSc(Tok(Ident)),
			success: true,
			result:  "{v=<nil>, toks=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       RepSc(Tok(Number)),
			success: true,
			result:  "{v=[123 456], toks=}",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       RepSc(Tok(Ident)),
			success: true,
			result:  "{v=[], toks=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: repr",
			input:   "123,456",
			p:       RepR(Tok(Number)),
			success: true,
			result:  "{v=[], toks=<num>/123🍌<num>/456}🍊{v=[123], toks=<num>/456}🍊{v=[123 456], toks=}",
		},
		{
			name:    "Parser: rep",
			input:   "123,456",
			p:       Rep(Tok(Number)),
			success: true,
			result:  "{v=[123 456], toks=}🍊{v=[123], toks=<num>/456}🍊{v=[], toks=<num>/123🍌<num>/456}",
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
			result:  "{v=, toks=<num>/123🍌<num>/456}🍊{v=123, toks=<num>/456}🍊{v=123;456, toks=}",
		},
		{
			name:  "Parser: errd",
			input: "a",
			p: ErrDef(Tok(Number).Map(func(v interface{}) interface{} {
				num, _ := strconv.Atoi(v.(*lexer.Token).Lexeme)
				return num
			}), "This is not a number!", 42),
			success: true,
			result:  "{v=42, toks=<id>/a}",
			error:   "This is not a number! in pos 1-2 line 1 col 1",
		},
		{
			name:    "Parser: amb, 1",
			input:   "1",
			p:       expr,
			success: true,
			result:  "{v=1, toks=}",
		},
		{
			name:    "Parser: amb, +1",
			input:   "+1",
			p:       expr,
			success: true,
			result:  "{v=(+ 1), toks=}",
		},
		{
			name:    "Parser: amb, 1+2",
			input:   "1+2",
			p:       expr,
			success: true,
			result:  "{v=[(1 . (+ 2)), (1 + 2)], toks=}",
		},
		{
			name:    "Parser: amb, 1+2+3",
			input:   "1+2+3",
			p:       expr,
			success: true,
			result:  "{v=[(1 . (+ [(2 . (+ 3)), (2 + 3)])), (1 + [(2 . (+ 3)), (2 + 3)])], toks=}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := lex.Lex(tt.input)
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
						t.Errorf("expect %s actual %s", tt.result, actual)
					}
				}
			} else {
				if out.Success {
					t.Errorf("expect fail actual success")
				}
				actual := out.Error.Error()
				if actual != tt.error {
					t.Errorf("expect %s actual %s", tt.result, actual)
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
		xs[i] = fmt.Sprintf("{v=%v, toks=%s}", r.Val, fmtToks(r.toks))
	}
	return strings.Join(xs, "🍊")
}

func fmtToks(toks []*lexer.Token) string {
	xs := make([]string, len(toks))
	for i, t := range toks {
		xs[i] = fmt.Sprintf("%s/%s", t.TokenKind, t.Lexeme)
	}
	return strings.Join(xs, "🍌")
}
