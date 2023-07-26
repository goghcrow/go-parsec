package example

import "github.com/goghcrow/lexer"

//goland:noinspection GoSnakeCaseUsage
const (
	PLUS TokenKind = iota + 1 // "+"
	SUB                       // "-"
	MUL                       // "*"
	DIV                       // "/"
	MOD                       // "%"
	EXP                       // "^"

	GT // ">"
	GE // ">="
	LT // "<"
	LE // "<="
	EQ // "=="
	NE // "!="

	LOGIC_NOT // "!"
	LOGIC_AND // "&&"
	LOGIC_OR  // "||"
)

// 自定义操作符
var userDefinedOperators = []lexer.Operator[TokenKind]{
	{PLUS, "+", lexer.BP_PREFIX, lexer.PREFIX},
	{SUB, "+", lexer.BP_PREFIX, lexer.PREFIX}, // NEGATE

	{PLUS, "+", lexer.BP_TERM, lexer.INFIX_L},
	{SUB, "-", lexer.BP_TERM, lexer.INFIX_L},
	{MUL, "*", lexer.BP_FACTOR, lexer.INFIX_L},
	{DIV, "/", lexer.BP_FACTOR, lexer.INFIX_L},
	{MOD, "%", lexer.BP_FACTOR, lexer.INFIX_L},
	{EXP, "^", lexer.BP_EXP, lexer.INFIX_R},

	{LE, "<=", lexer.BP_CMP, lexer.INFIX_N},
	{LT, "<", lexer.BP_CMP, lexer.INFIX_N},
	{GE, ">=", lexer.BP_CMP, lexer.INFIX_N},
	{GT, ">", lexer.BP_CMP, lexer.INFIX_N},
	{EQ, "==", lexer.BP_EQ, lexer.INFIX_N},
	{NE, "!=", lexer.BP_EQ, lexer.INFIX_N},

	{LOGIC_OR, "||", lexer.BP_LOGIC_OR, lexer.INFIX_L},
	{LOGIC_AND, "&&", lexer.BP_LOGIC_AND, lexer.INFIX_L},
	{LOGIC_NOT, "!", lexer.BP_PREFIX, lexer.PREFIX},
}
