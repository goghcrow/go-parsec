package lexer

import (
	"regexp"
	"sort"
	"strings"
)

// lexer auxiliary functions

// ----------------------------------------------------------------
// Common Regex
// ----------------------------------------------------------------

const (
	RegFloat = "(?:[-+]?(?:0|[1-9][0-9]*)(?:[.][0-9]+)+(?:[eE][-+]?[0-9]+)?)" +
		"|" +
		"(?:[-+]?(?:0|[1-9][0-9]*)(?:[.][0-9]+)?(?:[eE][-+]?[0-9]+)+)"
	RegInt = "(?:[-+]?0b(?:0|1[0-1]*))" +
		"|" +
		"(?:[-+]?0x(?:0|[1-9a-fA-F][0-9a-fA-F]*))" +
		"|" +
		"(?:[-+]?0o(?:0|[1-7][0-7]*))" +
		"|" +
		"(?:[-+]?(?:0|[1-9][0-9]*))"
	RegStr = "(?:\"(?:[^\"\\\\]*|\\\\[\"\\\\trnbf\\/]|\\\\u[0-9a-fA-F]{4})*\")" +
		"|" +
		"(?:`[^`]*`)" // need to strconv.Quote
	RegIdent = "[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*"
	RegOper  = "[:!#\\$%\\^&\\*\\+\\./<=>\\?@\\\\ˆ\\|~-]+" // ref operator.go
)

// ----------------------------------------------------------------
// BP BindingPower, Precedence & Fixity
// ----------------------------------------------------------------

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
	BP_MEMBER // . -> []
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

// ----------------------------------------------------------------
// Operator & Ident
// ----------------------------------------------------------------

type Operator[TokenKind comparable] struct {
	TokenKind TokenKind
	Lexeme    string
	BP
	Fixity
}

const (
	// 允许自定义操作符字符列表
	operators = ":!#$%^&*+./<=>?@\\ˆ|~-"
	// operators = "!$%&*+\\-./:<=>?@^|~" // 取消 # 与 ˆ
)

var (
	idReg = regexp.MustCompile("^" + RegIdent + "$")
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
func SortOpers[TokenKind comparable](ops []Operator[TokenKind]) []Operator[TokenKind] {
	sort.SliceStable(ops, func(i, j int) bool {
		x := ops[i].Lexeme
		y := ops[j].Lexeme
		if x == y || len(x) == len(y) {
			return false
		}
		return len(x) > len(y)
	})
	return ops
}
