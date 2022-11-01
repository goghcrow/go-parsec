package lexer

import (
	"regexp"
	"sort"
	"strings"
)

// BP BindingPower, Precedence
// 这里使用 float 是因为可以更精细定义自定义操作符的优先级
// e.g. 如果需要区分前后缀操作符优先级, 可以自己调整
type BP float32

//goland:noinspection GoSnakeCaseUsage
const (
	BP_NONE       BP = iota
	BP_LEFT_BRACE    // {
	BP_COND          // ?:
	BP_LOGIC_OR      // ||
	BP_LOGIC_AND     // &&
	BP_EQ            // == !=
	BP_CMP           // < > <= >=
	BP_TERM          // + -
	BP_FACTOR        // * / %
	BP_EXP           // ^
	BP_PREFIX        // - !
	BP_POSTFIX
	BP_CALL   // ()
	BP_MEMBER // . []
)

// Fixity Associativity
type Fixity int

//goland:noinspection GoSnakeCaseUsage
const (
	NA Fixity = iota
	PREFIX
	INFIX_N
	INFIX_L
	INFIX_R
	POSTFIX
)

type Operator struct {
	TokenKind
	BP
	Fixity
}

const (
	// 允许自定义操作符字符列表
	operators = ":!#$%^&*+./<=>?@\\ˆ|~-"
)

var (
	idReg = regexp.MustCompile("^[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*$")
	opReg = regexp.MustCompile("^[" + regexp.QuoteMeta(operators) + "]+$")
)

func HasOperPrefix(s string) bool {
	for _, r := range []rune(operators) {
		if strings.HasPrefix(s, string(r)) {
			return true
		}
	}
	return false
}

func IsIdentOp(name string) bool { return idReg.MatchString(name) }

func IsOp(s string) bool { return opReg.MatchString(s) }

// SortOpers 📢 因为 lexer 是按顺序匹配, 对于多字符的符号操作符需要注意顺序, 多字符放在单字符之前, ident 操作符不需要
// e.g. ! 需要放在 != 之后, > 需要放在 >= 之后
// e.g. 如果定义 & 需要放在  && 之后
// 使用 ops 之前, 需要先排下序
func SortOpers(ops []Operator) []Operator {
	sort.SliceStable(ops, func(i, j int) bool {
		x := ops[i].TokenKind
		y := ops[j].TokenKind
		if x == y || len(x) == len(y) {
			return false
		}
		return len(x) > len(y)
	})
	return ops
}
