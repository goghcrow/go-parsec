package example

import (
	"strconv"
	"testing"

	. "github.com/goghcrow/go-parsec"
	. "github.com/goghcrow/go-parsec/lexer"
)

func TestRec(t *testing.T) {
	const (
		Number TokenKind = iota + 1
		Add
		Sub
		Mul
		Div
		LParen
		RParen
		Space
	)

	lex := BuildLexer(func(lex *Lexicon) {
		lex.Regex(Number, `\d+(\.\d+)?`)
		lex.Oper(Add, "+")
		lex.Oper(Sub, "-")
		lex.Oper(Mul, "*")
		lex.Oper(Div, "/")
		lex.Str(LParen, "(")
		lex.Str(RParen, ")")
		lex.Regex(Space, `\s+`).Skip()
	})

	TERM := NewRule()
	FACTOR := NewRule()
	EXP := NewRule()

	str2num := func(s string) float64 {
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return num
	}

	applyNum := func(v interface{}) interface{} {
		return str2num(v.(*Token).Lexeme)
	}
	applyUnary := func(v interface{}) interface{} {
		xs := v.([]interface{})
		u := xs[0].(*Token)
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
		oper := b.([]interface{})[0].(*Token)
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

	// TERM
	//  	= NUMBER
	//  	= ('+' | '-') TERM
	//  	= '(' EXP ')'
	// FACTOR
	//  	= TERM
	//  	= FACTOR ('*' | '/') TERM
	// EXP
	//  	= FACTOR
	//  	= EXP ('+' | '-') FACTOR
	TERM.Pattern = Alt(
		Tok(Number).Map(applyNum),
		Seq(Alt(Str("+"), Str("-")), TERM).Map(applyUnary),
		KMid(Str("("), EXP, Str(")")),
	)
	FACTOR.Pattern = LRecSc(
		TERM,
		Seq(Alt(Str("*"), Str("/")), TERM),
		applyBinary,
	)
	EXP.Pattern = LRecSc(
		FACTOR,
		Seq(Alt(Str("+"), Str("-")), FACTOR),
		applyBinary,
	)

	eval := func(s string) float64 {
		toks := lex.MustLex(s)
		out := EXP.Parse(toks)
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
