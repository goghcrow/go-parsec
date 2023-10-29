package calc

import (
	"fmt"
	"strconv"

	"github.com/goghcrow/go-parsec/lexer"
	. "github.com/goghcrow/go-parsec/parsec"
)

func Calc(s string) float64 { return calculator(s) }

func Show(s string) string { return printer(s) }

var calculator = BuildParser[float64](
	func(s string) float64 {
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return num
	},
	func(op Op, a float64) float64 {
		switch op {
		case "+":
			return a
		case "-":
			return -a
		default:
			panic("unreached")
		}
	},
	func(op Op, l float64, r float64) float64 {
		switch op {
		case "+":
			return l + r
		case "-":
			return l - r
		case "*":
			return l * r
		case "/":
			return l / r
		default:
			panic("unreached")
		}
	},
)
var printer = BuildParser[string](
	func(s string) string { return s },
	func(op Op, a string) string { return fmt.Sprintf("(%s %s)", op, a) },
	func(op Op, l string, r string) string { return fmt.Sprintf("(%s %s %s)", op, l, r) },
)

type TokenKind int

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

func (k TokenKind) String() string {
	return map[TokenKind]string{
		Number: "number",
		Add:    "+",
		Sub:    "-",
		Mul:    "*",
		Div:    "/",
		LParen: "(",
		RParen: ")",
		Space:  "<space>",
	}[k]
}

type Op = string

func BuildParser[Val any](
	val func(string) Val,
	unary func(Op, Val) Val,
	binary func(Op, Val, Val) Val,
) func(s string) Val {
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

	strOf := func(toMatch string) Parser[TokenKind, Token[TokenKind]] {
		return Str[TokenKind](toMatch)
	}

	applyNum := func(v Token[TokenKind]) Val { return val(v.Lexeme()) }
	applyUnary := func(v Cons[Token[TokenKind], Val]) Val {
		return unary(v.Car.Lexeme(), v.Cdr)
	}
	applyBinary := func(a Val, b Cons[Token[TokenKind], Val]) Val {
		return binary(b.Car.Lexeme(), a, b.Cdr)
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

	sign := Seq2(AltSc(strOf("+"), strOf("-")), term)
	TERM.Pattern = AltSc(
		Apply(Tok(Number), applyNum),
		Apply(sign, applyUnary),
		KMid(strOf("("), exp, strOf(")")),
	)
	FACTOR.Pattern = LRecSc(
		term,
		Seq2(AltSc(strOf("*"), strOf("/")), term),
		applyBinary,
	)
	EXP.Pattern = LRecSc(
		factor,
		Seq2(AltSc(strOf("+"), strOf("-")), factor),
		applyBinary,
	)

	return func(s string) Val {
		xs := lex.MustLex(s)
		toks := make([]Token[TokenKind], len(xs))
		for i, t := range xs {
			toks[i] = t
		}
		out := EXP.Parse(toks)
		v, err := ExpectSingleResult(ExpectEOF(out))
		if err != nil {
			panic(err)
		}
		return v
	}
}
