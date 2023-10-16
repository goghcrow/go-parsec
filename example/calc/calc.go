package calc

import (
	"strconv"

	"github.com/goghcrow/go-parsec/lexer"
	. "github.com/goghcrow/go-parsec/parsec"
)

type Val = float64

func Calc(s string) Val {
	xs := lex.MustLex(s)
	toks := make([]Token[TokenKind], len(xs))
	for i, t := range xs {
		toks[i] = t
	}
	out := parser.Parse(toks)
	result, err := ExpectSingleResult(ExpectEOF(out))
	if err != nil {
		panic(err)
	}
	return result
}

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

var lex *lexer.Lexer[TokenKind]
var parser Parser[TokenKind, Val]

func init() {
	lex = lexer.BuildLexer(func(lex *lexer.Lexicon[TokenKind]) {
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
		Apply(Seq2(Alt(strOf("+"), strOf("-")), term), applyUnary),
		KMid(strOf("("), exp, strOf(")")),
	)
	FACTOR.Pattern = LRecSc(
		term,
		Seq2(Alt(strOf("*"), strOf("/")), term),
		applyBinary,
	)
	EXP.Pattern = LRecSc(
		factor,
		Seq2(Alt(strOf("+"), strOf("-")), factor),
		applyBinary,
	)

	parser = EXP
}
