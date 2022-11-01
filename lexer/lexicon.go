package lexer

import (
	"regexp"
	"strings"
)

const NotMatched = -1

type Rule struct {
	keep bool
	TokenKind
	match func(string) int // 匹配返回 EndRuneCount , 失败返回 NotMatched
}

func (r *Rule) Skip() *Rule { r.keep = false; return r }

// Lexicon Lexical grammar
type Lexicon struct {
	rules []*Rule
}

func NewLexicon() Lexicon {
	return Lexicon{}
}

func (l *Lexicon) Rule(r Rule) *Rule                       { l.rules = append(l.rules, &r); return &r }
func (l *Lexicon) Str(k TokenKind) *Rule                   { return l.Rule(str(k)) }
func (l *Lexicon) Keyword(k TokenKind) *Rule               { return l.Rule(keyword(k)) }
func (l *Lexicon) Regex(k TokenKind, pattern string) *Rule { return l.Rule(reg(k, pattern)) }
func (l *Lexicon) PrimOper(k TokenKind) *Rule              { return l.Rule(primOper(k)) }
func (l *Lexicon) Oper(k TokenKind) *Rule {
	if IsIdentOp(string(k)) {
		return l.Keyword(k)
	} else {
		return l.Str(k)
	}
}

func str(k TokenKind) Rule {
	tok := string(k)
	return Rule{true, k, func(s string) int {
		if strings.HasPrefix(s, tok) {
			return runeCount(tok)
		} else {
			return NotMatched
		}
	}}
}

var keywordPostfix = regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)

func keyword(k TokenKind) Rule {
	kw := string(k)
	return Rule{true, k, func(s string) int {
		completedWord := strings.HasPrefix(s, kw) &&
			!keywordPostfix.MatchString(s[len(kw):])
		if completedWord {
			return runeCount(kw)
		} else {
			return NotMatched
		}
	}}
}

func reg(k TokenKind, pattern string) Rule {
	startWith := regexp.MustCompile("^" + pattern)
	return Rule{true, k, func(s string) int {
		found := startWith.FindString(s)
		if found == "" {
			return NotMatched
		} else {
			return runeCount(found)
		}
	}}
}

// primOper . ? 内置操作符的优先级高于自定义操作符, 且不是匹配最长, 需要特殊处理
// e.g 比如自定义操作符 .^. 不能匹配成 [`.`, `^.`]
func primOper(k TokenKind) Rule {
	op := string(k)
	return Rule{true, k, func(s string) int {
		if !strings.HasPrefix(s, op) {
			return NotMatched
		}
		completedOper := len(s) == len(op) || !HasOperPrefix(s[len(op):])
		if completedOper {
			return runeCount(op)
		} else {
			return NotMatched
		}
	}}
}
