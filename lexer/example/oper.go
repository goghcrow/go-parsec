package example

import . "github.com/goghcrow/go-parsec/lexer"

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
var userDefinedOperators = []Operator{
	{PLUS, "+", BP_PREFIX, PREFIX},
	{SUB, "+", BP_PREFIX, PREFIX}, // NEGATE

	{PLUS, "+", BP_TERM, INFIX_L},
	{SUB, "-", BP_TERM, INFIX_L},
	{MUL, "*", BP_FACTOR, INFIX_L},
	{DIV, "/", BP_FACTOR, INFIX_L},
	{MOD, "%", BP_FACTOR, INFIX_L},
	{EXP, "^", BP_EXP, INFIX_R},

	{LE, "<=", BP_CMP, INFIX_N},
	{LT, "<", BP_CMP, INFIX_N},
	{GE, ">=", BP_CMP, INFIX_N},
	{GT, ">", BP_CMP, INFIX_N},
	{EQ, "==", BP_EQ, INFIX_N},
	{NE, "!=", BP_EQ, INFIX_N},

	{LOGIC_OR, "||", BP_LOGIC_OR, INFIX_L},
	{LOGIC_AND, "&&", BP_LOGIC_AND, INFIX_L},
	{LOGIC_NOT, "!", BP_PREFIX, PREFIX},
}
