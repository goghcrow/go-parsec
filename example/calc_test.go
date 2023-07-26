package example

import (
	"strconv"
	"testing"

	"github.com/goghcrow/lexer"
	. "github.com/goghcrow/parsec"
)

func StrOf(toMatch string) Parser[TokenKind, Token[TokenKind]] { return Str[TokenKind](toMatch) }

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

	lex := lexer.BuildLexer(func(lex *lexer.Lexicon[TokenKind]) {
		lex.Regex(Number, `\d+(\.\d+)?`)
		lex.Oper(Add, "+")
		lex.Oper(Sub, "-")
		lex.Oper(Mul, "*")
		lex.Oper(Div, "/")
		lex.Str(LParen, "(")
		lex.Str(RParen, ")")
		lex.Regex(Space, `\s+`).Skip()
	})

	type Val = float64

	str2num := func(s string) Val {
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return num
	}

	applyNum := func(v Token[TokenKind]) Val {
		return str2num(v.Lexeme())
	}

	applyUnary := func(v Cons[Token[TokenKind], Val]) Val {
		switch v.Car.Lexeme() {
		case "+":
			return v.Cdr
		case "-":
			return -v.Cdr
		default:
			panic("unreached")
		}
	}
	applyBinary := func(a Val, b Cons[Token[TokenKind], Val]) Val {
		switch b.Car.Lexeme() {
		case "+":
			return a + b.Cdr
		case "-":
			return a - b.Cdr
		case "*":
			return a * b.Cdr
		case "/":
			return a / b.Cdr
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

	TERM := NewRule[TokenKind, Val]()
	FACTOR := NewRule[TokenKind, Val]()
	EXP := NewRule[TokenKind, Val]()

	term := TERM.Parser()
	factor := FACTOR.Parser()
	exp := EXP.Parser()

	TERM.Pattern = Alt(
		Apply(Tok(Number), applyNum),
		Apply(Seq2(Alt(StrOf("+"), StrOf("-")), term), applyUnary),
		KMid(StrOf("("), exp, StrOf(")")),
	)
	FACTOR.Pattern = LRecSc(
		term,
		Seq2(Alt(StrOf("*"), StrOf("/")), term),
		applyBinary,
	)
	EXP.Pattern = LRecSc(
		factor,
		Seq2(Alt(StrOf("+"), StrOf("-")), factor),
		applyBinary,
	)

	eval := func(s string) Val {
		xs := lex.MustLex(s)
		toks := make([]Token[TokenKind], len(xs))
		for i, t := range xs {
			toks[i] = t
		}
		out := EXP.Parse(toks)
		result, err := ExpectSingleResult(ExpectEOF(out))
		if err != nil {
			panic(err)
		}
		return result
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
