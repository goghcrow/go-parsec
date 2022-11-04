package lexer

import (
	"github.com/goghcrow/go-parsec/lexer"
)

func NewBuiltinLexer(udOpers []lexer.Operator) *lexer.Lexer {
	return &lexer.Lexer{
		Lexicon: NewBuiltinLexicon(udOpers),
	}
}

// '前缀的为 psuido
//
//goland:noinspection GoSnakeCaseUsage
const (
	QUESTION lexer.TokenKind = iota + 1 // "?"
	ARROW                               // "->"

	IF   // "if"
	THEN // "then"
	ELSE // "else"

	NAME // "'name"
	NUM  // "'num"
	STR  // "'str"
	TIME // "'time"

	TRUE  // "true"
	FALSE // "false"

	COMMA // ","
	DOT   // "."
	SPACE // "'space"

	LEFT_PAREN    // "("
	RIGHT_PAREN   // ")"
	LEFT_BRACKET  // "["
	RIGHT_BRACKET // "]"
	LEFT_BRACE    // "{"
	RIGHT_BRACE   // "}"

	COLON // ":"
)

type Keyword struct {
	lexer.TokenKind
	Lexeme string
}

var keywords = []Keyword{
	{IF, "if"},
	{THEN, "then"},
	{ELSE, "else"},
}

var builtInOpers = []lexer.Operator{
	{DOT, ".", lexer.BP_MEMBER, lexer.INFIX_L},
	{ARROW, "->", lexer.BP_MEMBER, lexer.INFIX_R},
	{QUESTION, "?:", lexer.BP_COND, lexer.INFIX_R},
}

func NewBuiltinLexicon(userDefinedOpers []lexer.Operator) lexer.Lexicon {
	l := lexer.Lexicon{}

	l.Regex(SPACE, "\\s+").Skip()
	l.Str(COLON, ":")
	l.Str(COMMA, ",")

	l.Str(LEFT_PAREN, "(")
	l.Str(RIGHT_PAREN, ")")
	l.Str(LEFT_BRACKET, "[")
	l.Str(RIGHT_BRACKET, "]")
	l.Str(LEFT_BRACE, "{")
	l.Str(RIGHT_BRACE, "}")

	for _, kw := range keywords {
		l.Keyword(kw.TokenKind, kw.Lexeme)
	}

	// 内置的操作符优先级高于自定义操作符
	for _, oper := range lexer.SortOpers(builtInOpers) {
		l.PrimOper(oper.TokenKind, oper.Lexeme)
	}

	// 自定义操作符
	for _, oper := range lexer.SortOpers(userDefinedOpers) {
		l.Oper(oper.TokenKind, oper.Lexeme)
	}

	l.Str(TRUE, "true")
	l.Str(FALSE, "false")

	// 移除数字前的 [+-]?, [+-]? 被处理成一元操作符, 实际上变成没有负数字面量, 语义不变
	l.Regex(NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)+(?:[eE][-+]?[0-9]+)?") // float
	l.Regex(NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)?(?:[eE][-+]?[0-9]+)+") // float
	l.Regex(NUM, "0b(?:0|1[0-1]*)")                                     // int
	l.Regex(NUM, "0x(?:0|[1-9a-fA-F][0-9a-fA-F]*)")                     // int
	l.Regex(NUM, "0o(?:0|[1-7][0-7]*)")                                 // int
	l.Regex(NUM, "(?:0|[1-9][0-9]*)")                                   // int

	l.Regex(STR, "\"(?:[^\"\\\\]*|\\\\[\"\\\\trnbf\\/]|\\\\u[0-9a-fA-F]{4})*\"")
	l.Regex(STR, "`[^`]*`") // raw string

	l.Regex(TIME, "'[^`\"']*'")

	l.Regex(NAME, "[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*") // 支持 unicode, 不能以数字开头

	return l
}
