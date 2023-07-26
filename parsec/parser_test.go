package parsec

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/goghcrow/lexer"
)

type tokKind = int
type token = Token[tokKind]

const (
	Number tokKind = iota + 1
	Add
	Space
	Ident
	Comma
)

func stroftk(k tokKind) string {
	return map[tokKind]string{
		Number: "<num>",
		Add:    "+",
		Space:  "<space>",
		Ident:  "<id>",
		Comma:  ",",
	}[k]
}

var lex = lexer.BuildLexer(func(lex *lexer.Lexicon[tokKind]) {
	lex.Regex(Number, "\\d+")
	lex.Regex(Ident, "[a-zA-Z]\\w*")
	lex.Regex(Space, "\\s+").Skip()
	lex.Str(Comma, ",").Skip()
	lex.Str(Add, "+")
})

var lexForCombinator = lexer.BuildLexer(func(lex *lexer.Lexicon[tokKind]) {
	lex.Regex(Number, "\\d+")
	lex.Regex(Ident, "[a-zA-Z]\\w*")
	lex.Regex(Space, "\\s+").Skip()
	lex.Str(Comma, ",")
})

func mustLex(s string) []Token[tokKind] {
	toks := lex.MustLex(s)
	xs := make([]Token[tokKind], len(toks))
	for i, t := range toks {
		xs[i] = t
	}
	return xs
}
func mustLexForCombinator(s string) []Token[tokKind] {
	toks := lexForCombinator.MustLex(s)
	xs := make([]Token[tokKind], len(toks))
	for i, t := range toks {
		xs[i] = t
	}
	return xs
}

func TestParser(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		p       func(toks []token) (bool, string, string)
		success bool
		result  string
		error   string
	}{
		{
			name:    "Parser: str",
			input:   "123,456",
			p:       wrap(Str[tokKind]("123")),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: str",
			input:   "123,456",
			p:       wrap(Str[tokKind]("456")),
			success: false,
			error:   "Unable to consume token `123` expect `456` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: tok",
			input:   "123,456",
			p:       wrap(Tok(Number)),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: alt",
			input:   "123,456",
			p:       wrap(Alt(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=123, next=<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt",
			input:   "abc,def",
			p:       wrap(Alt(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=abc, next=<id>/def}",
			error:   "Unable to consume token `abc` expect `token<1>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt",
			input:   "123,456",
			p:       wrap(Alt(Alt(Tok(Number), Tok(Ident)), Alt(Tok(Ident), Tok(Number)))),
			success: true,
			result:  "{v=123, next=<num>/456}🍊{v=123, next=<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt",
			input:   "abc,def",
			p:       wrap(Alt(Alt(Tok(Number), Tok(Ident)), Alt(Tok(Ident), Tok(Number)))),
			success: true,
			result:  "{v=abc, next=<id>/def}🍊{v=abc, next=<id>/def}",
			error:   "Unable to consume token `abc` expect `token<1>` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Parser: alt",
			input: "123,456",
			p: wrap(Apply[tokKind, Either[token, token], string](Alt2(Tok(Number), Tok(Ident)), func(either Either[token, token]) string {
				return either.Left.Lexeme()
			})),
			success: true,
			result:  "{v=123, next=<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt_sc",
			input:   "123,456",
			p:       wrap(AltSc(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=123, next=<num>/456}",
			error:   "",
		},
		{
			name:    "Parser: alt_sc",
			input:   "abc,def",
			p:       wrap(AltSc(Tok(Number), Tok(Ident))),
			success: true,
			result:  "{v=abc, next=<id>/def}",
			error:   "Unable to consume token `abc` expect `token<1>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt_sc",
			input:   "123,456",
			p:       wrap(AltSc(Alt(Tok(Number), Tok(Ident)), Alt(Tok(Ident), Tok(Number)))),
			success: true,
			result:  "{v=123, next=<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Parser: alt_sc",
			input: "123,456",
			p: wrap(AltSc(Apply(Tok(Ident), func(v token) string {
				return "alt1: " + v.Lexeme()
			}), Apply(Alt(Tok(Ident), Tok(Number)), func(v token) string {
				return "alt2: " + v.Lexeme()
			}))),
			success: true,
			result:  "{v=alt2: 123, next=<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: alt_sc",
			input:   "abc,def",
			p:       wrap(AltSc(Alt(Tok(Number), Tok(Ident)), Alt(Tok(Ident), Tok(Number)))),
			success: true,
			result:  "{v=abc, next=<id>/def}",
			error:   "Unable to consume token `abc` expect `token<1>` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Parser: alt_sc",
			input: "abc,def",
			p: wrap(AltSc(Apply(Tok(Number), func(v token) string {
				return "alt1: " + v.Lexeme()
			}), Apply(Alt(Tok(Ident), Tok(Number)), func(v token) string {
				return "alt2: " + v.Lexeme()
			}))),
			success: true,
			result:  "{v=alt2: abc, next=<id>/def}",
			error:   "Unable to consume token `abc` expect `token<1>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       wrap(Seq(Tok(Number), Tok(Ident))),
			success: false,
			error:   "Unable to consume token `456` expect `token<4>` in pos 5-8 line 1 col 5",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       wrap(Seq2(Tok(Number), Tok(Ident))),
			success: false,
			error:   "Unable to consume token `456` expect `token<4>` in pos 5-8 line 1 col 5",
		},
		{
			name:    "Parser: seq",
			input:   "123,456",
			p:       wrap(Seq(Tok(Number), Tok(Number))),
			success: true,
			result:  "{v=[123 456], next=}",
		},
		{
			name:  "Parser: seq",
			input: "123,456",
			p: wrap(Apply[tokKind, Cons[token, token], []string](Seq2(Tok(Number), Tok(Number)), func(t2 Cons[token, token]) []string {
				return []string{t2.Car.Lexeme(), t2.Cdr.Lexeme()}
			})),
			success: true,
			result:  "{v=[123 456], next=}",
		},
		{
			name:    "Parser: seq",
			input:   "123,456,a",
			p:       wrap(Seq(Tok(Number), Tok(Number))),
			success: true,
			result:  "{v=[123 456], next=<id>/a}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       wrap(KLeft(Tok(Number), Seq(Tok(Number), Tok(Number)))),
			success: true,
			result:  "{v=123, next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       wrap(KMid(Tok(Number), Tok(Number), Tok(Number))),
			success: true,
			result:  "{v=456, next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       wrap(KRight(Seq(Tok(Number), Tok(Number)), Tok(Number))),
			success: true,
			result:  "{v=789, next=}",
		},
		{
			name:    "Parser: kleft, kmid, kright",
			input:   "123,456,789",
			p:       wrap(KRight(Tok(Number), Seq(Tok(Number), Tok(Number)))),
			success: true,
			result:  "{v=[456 789], next=}",
		},
		{
			name:    "Parser: opt",
			input:   "123,456",
			p:       wrap(Opt(Tok(Number))),
			success: true,
			result:  "{v=123, next=<num>/456}🍊{v=<nil>, next=<num>/123🍌<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       wrap(OptSc(Tok(Number))),
			success: true,
			result:  "{v=123, next=<num>/456}",
		},
		{
			name:    "Parser: opt_sc",
			input:   "123,456",
			p:       wrap(OptSc(Tok(Ident))),
			success: true,
			result:  "{v=<nil>, next=<num>/123🍌<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       wrap(RepSc(Tok(Number))),
			success: true,
			result:  "{v=[123 456], next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: rep_sc",
			input:   "123,456",
			p:       wrap(RepSc(Tok(Ident))),
			success: true,
			result:  "{v=[], next=<num>/123🍌<num>/456}",
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Parser: repr",
			input:   "123,456",
			p:       wrap(RepR(Tok(Number))),
			success: true,
			result:  "{v=[], next=<num>/123🍌<num>/456}🍊{v=[123], next=<num>/456}🍊{v=[123 456], next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: rep",
			input:   "123,456",
			p:       wrap(Rep(Tok(Number))),
			success: true,
			result:  "{v=[123 456], next=}🍊{v=[123], next=<num>/456}🍊{v=[], next=<num>/123🍌<num>/456}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: rep_n",
			input:   "123,456,789",
			p:       wrap(RepN(Tok(Number), 0)),
			success: true,
			result:  "{v=[], next=<num>/123🍌<num>/456🍌<num>/789}",
			error:   "",
		},
		{
			name:    "Parser: rep_n",
			input:   "123,456,789",
			p:       wrap(RepN(Tok(Number), 1)),
			success: true,
			result:  "{v=[123], next=<num>/456🍌<num>/789}",
			error:   "",
		},
		{
			name:    "Parser: rep_n",
			input:   "123,456,789",
			p:       wrap(RepN(Tok(Number), 2)),
			success: true,
			result:  "{v=[123 456], next=<num>/789}",
			error:   "",
		},
		{
			name:    "Parser: rep_n",
			input:   "123,456,789",
			p:       wrap(RepN(Tok(Number), 3)),
			success: true,
			result:  "{v=[123 456 789], next=}",
			error:   "",
		},
		{
			name:    "Parser: rep_n",
			input:   "123,456,789",
			p:       wrap(RepN(Tok(Number), 4)),
			success: false,
			result:  "",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: list",
			input:   "123 + 456 + 789",
			p:       wrap(List(Tok(Number), Tok(Add))),
			success: true,
			result:  "{v=[123 456 789], next=}🍊{v=[123 456], next=+/+🍌<num>/789}🍊{v=[123], next=+/+🍌<num>/456🍌+/+🍌<num>/789}",
			error:   "Nothing to consume` expect `token<2>` in <EOF>",
		},
		{
			name:    "Parser: list",
			input:   "123 + 456 + 789",
			p:       wrap(ListSc(Tok(Number), Tok(Add))),
			success: true,
			result:  "{v=[123 456 789], next=}",
			error:   "Nothing to consume` expect `token<2>` in <EOF>",
		},
		{
			name:    "Parser: list",
			input:   "123 + 456 + 789",
			p:       wrap(ListN(Tok(Number), Tok(Add), 2)),
			success: true,
			result:  "{v=[123 456], next=+/+🍌<num>/789}",
			error:   "",
		},
		{
			name:    "Parser: list",
			input:   "123 + 456 + 789",
			p:       wrap(ListN(Tok(Number), Tok(Add), 3)),
			success: true,
			result:  "{v=[123 456 789], next=}",
			error:   "",
		},
		{
			name:  "Parser: apply",
			input: "123,456",
			p: wrap(Apply[tokKind, []token, string](RepR(Tok(Number)), func(toks []token) string {
				var xs []string
				for _, tok := range toks {
					xs = append(xs, tok.Lexeme())
				}
				return strings.Join(xs, ";")
			})),
			success: true,
			result:  "{v=, next=<num>/123🍌<num>/456}🍊{v=123, next=<num>/456}🍊{v=123;456, next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:  "Parser: errd",
			input: "a",
			p: wrap(ErrD(Apply(Tok(Number), func(tok token) int {
				num, _ := strconv.Atoi(tok.Lexeme())
				return num
			}), "This is not a number!", 42)),
			success: true,
			result:  "{v=42, next=<id>/a}",
			error:   "This is not a number! in pos 1-2 line 1 col 1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := mustLex(tt.input)
			succ, out, err := tt.p(toks)
			if tt.success != succ {
				t.Errorf("[succ]expect %v actual %v", tt.success, succ)
			}
			if out != tt.result {
				t.Errorf("[out]expect %s actual %s", tt.result, out)
			}
			if err != tt.error {
				t.Errorf("[err]expect %s actual %s", tt.error, err)
			}
		})
	}
}

func TestCombine(t *testing.T) {
	COUNT := NewRule[tokKind, int]()
	NAME := NewRule[tokKind, string]()
	NAME_LIST := NewRule[tokKind, []string]()

	COUNT.Pattern = Apply(Tok[tokKind](Number), func(v Token[tokKind]) int {
		num, _ := strconv.Atoi(v.Lexeme())
		return num
	})
	NAME.Pattern = Apply(Tok[tokKind](Ident), func(v Token[tokKind]) string {
		return v.Lexeme()
	})
	NAME_LIST.Pattern = Combine2(Parser[tokKind, int](COUNT), func(cnt int) Parser[tokKind, []string] {
		if cnt < 1 {
			return Fail[tokKind, []string](fmt.Sprintf("illegal number of names: %d, it should >= 1", cnt))
		} else {
			return ListN(Parser[tokKind, string](NAME), Str[tokKind](","), cnt)
		}
	})

	for _, tt := range []struct {
		name    string
		input   string
		p       func(toks []token) (bool, string, string)
		success bool
		result  string
		error   string
	}{
		{
			name:  "Parser: combinator 0",
			input: "a aa aaaa",
			p: wrap(Combine(Str[tokKind]("a"), func(t Token[tokKind]) Parser[tokKind, Token[tokKind]] {
				return Str[tokKind](t.Lexeme() + t.Lexeme())
			}, func(t Token[tokKind]) Parser[tokKind, Token[tokKind]] {
				return Str[tokKind](t.Lexeme() + t.Lexeme())
			})),
			success: true,
			result:  "{v=aaaa, next=}",
			error:   "",
		},
		{
			name:  "Parser: combinator 0",
			input: "a aa",
			p: wrap(Combine(Str[tokKind]("a"), func(t Token[tokKind]) Parser[tokKind, Token[tokKind]] {
				return Str[tokKind](t.Lexeme() + t.Lexeme())
			}, func(t Token[tokKind]) Parser[tokKind, Token[tokKind]] {
				return Str[tokKind](t.Lexeme() + t.Lexeme())
			})),
			success: false,
			result:  "",
			error:   "Nothing to consume` expect `aaaa` in <EOF>",
		},
		{
			name:    "Parser: combinator 0",
			input:   "0",
			p:       wrap(NAME_LIST.Parser()),
			success: false,
			result:  "",
			error:   "illegal number of names: 0, it should >= 1 in <EOF>",
		},
		{
			name:    "Parser: combinator 1 foo",
			input:   "1 foo",
			p:       wrap(NAME_LIST.Parser()),
			success: true,
			result:  "{v=[foo], next=}",
			error:   "",
		},
		{
			name:    "Parser: combinator 2 foo,bar",
			input:   "2 foo,bar",
			p:       wrap(NAME_LIST.Parser()),
			success: true,
			result:  "{v=[foo bar], next=}",
			error:   "",
		},
		{
			name:    "Parser: combinator 2 foo,bar,baz",
			input:   "3 foo,bar,baz",
			p:       wrap(NAME_LIST.Parser()),
			success: true,
			result:  "{v=[foo bar baz], next=}",
			error:   "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := mustLexForCombinator(tt.input)
			succ, out, err := tt.p(toks)
			if tt.success != succ {
				t.Errorf("[succ]expect %v actual %v", tt.success, succ)
			}
			if out != tt.result {
				t.Errorf("[out]expect %s actual %s", tt.result, out)
			}
			if err != tt.error {
				t.Errorf("[err]expect %s actual %s", tt.error, err)
			}
		})
	}
}

func TestAmbParser(t *testing.T) {
	var TERM = NewRule[tokKind, string]()
	EXPR := NewRule[tokKind, string]()
	expr := Parser[tokKind, string](EXPR)

	// TERM
	//		= NUMBER
	//		= + EXPR
	// EXPR
	//		= TERM
	//		= EXPR | (+ EXPR)
	tok2str := func(t token) string { return t.Lexeme() }
	fmtPlus := func(v string) string {
		return fmt.Sprintf("(+ %s)", v)
	}
	TERM.Pattern = Alt[tokKind, string](
		Apply(Tok(Number), tok2str),
		Apply(KRight[tokKind, token, string](Str[tokKind]("+"), EXPR), fmtPlus),
	)
	lrec := Amb(LRecSc[tokKind, Either[string, Cons[token, string]], string](
		TERM,
		Alt2(expr, Seq2(Str[tokKind]("+"), expr)),
		func(a string, b Either[string, Cons[token, string]]) string {
			if b.IsLeft() {
				return fmt.Sprintf(`(%s . %s)`, a, b.Left)
			} else {
				return fmt.Sprintf("(%s + %s)", a, b.Right.Cdr)
			}
		},
	))
	EXPR.Pattern = Apply(lrec, func(xs []string) string {
		var ss []string
		for _, x := range xs {
			ss = append(ss, x)
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
		p       func(toks []token) (bool, string, string)
		success bool
		result  string
		error   string
	}{
		{
			name:    "Parser: amb, 1",
			input:   "1",
			p:       wrap(EXPR.Parser()),
			success: true,
			result:  "{v=1, next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: amb, +1",
			input:   "+1",
			p:       wrap(EXPR.Parser()),
			success: true,
			result:  "{v=(+ 1), next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: amb, 1+2",
			input:   "1+2",
			p:       wrap(EXPR.Parser()),
			success: true,
			result:  "{v=[(1 . (+ 2)), (1 + 2)], next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
		{
			name:    "Parser: amb, 1+2+3",
			input:   "1+2+3",
			p:       wrap(EXPR.Parser()),
			success: true,
			result:  "{v=[(1 . (+ [(2 . (+ 3)), (2 + 3)])), (1 + [(2 . (+ 3)), (2 + 3)])], next=}",
			error:   "Nothing to consume` expect `token<1>` in <EOF>",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := mustLex(tt.input)
			succ, out, err := tt.p(toks)
			if tt.success != succ {
				t.Errorf("[succ]expect %v actual %v", tt.success, succ)
			}
			if out != tt.result {
				t.Errorf("[out]expect %s actual %s", tt.result, out)
			}
			if err != tt.error {
				t.Errorf("[err]expect %s actual %s", tt.error, err)
			}
		})
	}
}

func TestFailure(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		p       func(toks []token) (bool, string, string)
		success bool
		result  string
		error   string
	}{
		{
			name:  "Failure: alt",
			input: "123,456",
			p: wrap(Alt(
				Tok(Comma),
				Tok(Space),
			)),
			success: false,
			error:   "Unable to consume token `123` expect `token<5>` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Failure: seq",
			input: "123,456",
			p: wrap(Seq(
				Tok(Ident),
				Tok(Number),
			)),
			success: false,
			error:   "Unable to consume token `123` expect `token<4>` in pos 1-4 line 1 col 1",
		},
		{
			name:  "Failure: seq",
			input: "123,456",
			p: wrap(Seq(
				Tok(Number),
				Tok(Ident),
			)),
			success: false,
			error:   "Unable to consume token `456` expect `token<4>` in pos 5-8 line 1 col 5",
		},
		{
			name:  "Failure: apply",
			input: "123,456",
			p: wrap(Apply(Tok(Comma), func(v token) any {
				return nil
			})),
			success: false,
			error:   "Unable to consume token `123` expect `token<5>` in pos 1-4 line 1 col 1",
		},
		{
			name:    "Failure: rep_sc + seq",
			input:   "1a 2b 3c d e",
			p:       wrap(RepSc(Seq(Tok(Number), Tok(Ident)))),
			success: true,
			result:  "{v=[[1 a] [2 b] [3 c]], next=<id>/d🍌<id>/e}",
			error:   "Unable to consume token `d` expect `token<1>` in pos 10-11 line 1 col 10",
		},
		{
			name:    "Failure: rep_sc + seq",
			input:   "1a 2b 3c d e",
			p:       wrap(Rep(Seq(Tok(Number), Tok(Ident)))),
			success: true,
			result:  "{v=[[1 a] [2 b] [3 c]], next=<id>/d🍌<id>/e}🍊{v=[[1 a] [2 b]], next=<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}🍊{v=[[1 a]], next=<num>/2🍌<id>/b🍌<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}🍊{v=[], next=<num>/1🍌<id>/a🍌<num>/2🍌<id>/b🍌<num>/3🍌<id>/c🍌<id>/d🍌<id>/e}",
			// 返回最远的错误
			error: "Unable to consume token `d` expect `token<1>` in pos 10-11 line 1 col 10",
		},
		{
			name:  "Failure: rep_sc + alt",
			input: "1 a b 2 c 3",
			p: wrap(RepSc(Apply(Alt2(Tok(Number), Seq2(Tok(Ident), Tok(Ident))), func(ei Either[token, Cons[token, token]]) string {
				if ei.IsLeft() {
					return ei.Left.Lexeme()
				} else {
					return fmt.Sprintf("%s", []string{ei.Right.Car.Lexeme(), ei.Right.Cdr.Lexeme()})
				}
			}))),
			success: true,
			result:  "{v=[1 [a b] 2], next=<id>/c🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` expect `token<4>` in pos 11-12 line 1 col 11",
		},
		{
			name:  "Failure: rep_sc + alt",
			input: "1 a b 2 c 3",
			p: wrap(Rep(Apply(Alt2(Tok(Number), Seq2(Tok(Ident), Tok(Ident))), func(ei Either[token, Cons[token, token]]) string {
				if ei.IsLeft() {
					return ei.Left.Lexeme()
				} else {
					return fmt.Sprintf("%s", []string{ei.Right.Car.Lexeme(), ei.Right.Cdr.Lexeme()})
				}
			}))),
			success: true,
			result:  "{v=[1 [a b] 2], next=<id>/c🍌<num>/3}🍊{v=[1 [a b]], next=<num>/2🍌<id>/c🍌<num>/3}🍊{v=[1], next=<id>/a🍌<id>/b🍌<num>/2🍌<id>/c🍌<num>/3}🍊{v=[], next=<num>/1🍌<id>/a🍌<id>/b🍌<num>/2🍌<id>/c🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` expect `token<4>` in pos 11-12 line 1 col 11",
		},
		{
			name:    "Failure: rep_sc + opt",
			input:   "a b c d e f g 3",
			p:       wrap(RepSc(OptSc(Seq(Tok(Ident), Tok(Ident))))),
			success: true,
			result:  "{v=[[a b] [c d] [e f]], next=<id>/g🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` expect `token<4>` in pos 15-16 line 1 col 15",
		},
		{
			name:    "Failure: rep_sc + opt",
			input:   "a b c d e f g 3",
			p:       wrap(Rep(OptSc(Seq(Tok(Ident), Tok(Ident))))),
			success: true,
			result:  "{v=[[a b] [c d] [e f]], next=<id>/g🍌<num>/3}🍊{v=[[a b] [c d]], next=<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}🍊{v=[[a b]], next=<id>/c🍌<id>/d🍌<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}🍊{v=[], next=<id>/a🍌<id>/b🍌<id>/c🍌<id>/d🍌<id>/e🍌<id>/f🍌<id>/g🍌<num>/3}",
			// Seq(Tok(Ident), Tok(Ident)) 解析到 3 失败
			error: "Unable to consume token `3` expect `token<4>` in pos 15-16 line 1 col 15",
		},
		{
			name:    "Failure: err",
			input:   "a",
			p:       wrap(Err(Tok(Number), "This is not a number!")),
			success: false,
			result:  "",
			error:   "This is not a number! in pos 1-2 line 1 col 1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			toks := mustLex(tt.input)
			succ, out, err := tt.p(toks)
			if tt.success != succ {
				t.Errorf("[succ]expect %v actual %v", tt.success, succ)
			}
			if out != tt.result {
				t.Errorf("[out]expect %s actual %s", tt.result, out)
			}
			if err != tt.error {
				t.Errorf("[err]expect %s actual %s", tt.error, err)
			}
		})
	}
}

func wrap[R any](p Parser[tokKind, R]) func(toks []token) (bool, string, string) {
	return func(toks []token) (bool, string, string) {
		return outOf(p.Parse(toks))
	}
}

func outOf[R any](out Output[tokKind, R]) (bool, string, string) {
	if out.Success {
		if out.Error == nil {
			return true, fmtResults(out.Candidates), ""
		}
		return true, fmtResults(out.Candidates), out.Error.Error()
	} else {
		return false, "", out.Error.Error()
	}
}

func fmtResults[R any](results []Result[tokKind, R]) string {
	xs := make([]string, len(results))
	for i, r := range results {
		if tok, ok := any(r.Val).(token); ok && tok != nil {
			xs[i] = fmt.Sprintf("{v=%s, next=%s}", tok.Lexeme(), fmtToks(r.next))
		} else {
			xs[i] = fmt.Sprintf("{v=%v, next=%s}", r.Val, fmtToks(r.next))
		}
	}
	return strings.Join(xs, "🍊")
}

func fmtToks(toks []token) string {
	xs := make([]string, len(toks))
	for i, t := range toks {
		xs[i] = fmt.Sprintf("%s/%s", stroftk(t.Kind()), t.Lexeme())
	}
	return strings.Join(xs, "🍌")
}
