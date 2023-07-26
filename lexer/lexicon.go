package lexer

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const NotMatched = -1

type Rule[K Ord] struct {
	keep  bool
	K     K
	match func(string) int // 匹配返回 EndRuneCount , 失败返回 NotMatched
}

func (r *Rule[K]) Skip() *Rule[K] { r.keep = false; return r }

// Lexicon Lexical grammar
type Lexicon[K Ord] struct {
	rules []*Rule[K]
}

func NewLexicon[K Ord]() Lexicon[K] {
	return Lexicon[K]{}
}

func (l *Lexicon[K]) Rule(r Rule[K]) *Rule[K] {
	l.rules = append(l.rules, &r)
	return &r
}
func (l *Lexicon[K]) Str(k K, s string) *Rule[K] { return l.Rule(str(k, s)) }
func (l *Lexicon[K]) Keyword(k K, s string) *Rule[K] {
	return l.Rule(keyword(k, s))
}
func (l *Lexicon[K]) Regex(k K, pattern string) *Rule[K] {
	return l.Rule(regex(k, pattern))
}
func (l *Lexicon[TokenKind]) PrimOper(k TokenKind, oper string) *Rule[TokenKind] {
	return l.Rule(primOper(k, oper))
}
func (l *Lexicon[TokenKind]) Oper(k TokenKind, oper string) *Rule[TokenKind] {
	if IsIdentOp(oper) {
		return l.Keyword(k, oper)
	} else {
		return l.Str(k, oper)
	}
}

func str[K Ord](k K, str string) Rule[K] {
	return Rule[K]{true, k, func(s string) int {
		if strings.HasPrefix(s, str) {
			return runeCount(str)
		} else {
			return NotMatched
		}
	}}
}

var keywordPostfix = regexp.MustCompile(`^[a-zA-Z\d\p{L}_]+`)

func keyword[TokenKind comparable](k TokenKind, kw string) Rule[TokenKind] {
	return Rule[TokenKind]{true, k, func(s string) int {
		// golang regexp 不支持 lookahead
		completedWord := strings.HasPrefix(s, kw) &&
			!keywordPostfix.MatchString(s[len(kw):])
		if completedWord {
			return runeCount(kw)
		} else {
			return NotMatched
		}
	}}
}

func regex[K Ord](k K, pattern string) Rule[K] {
	startWith := regexp.MustCompile("^(?:" + pattern + ")")
	return Rule[K]{true, k, func(s string) int {
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
func primOper[TokenKind comparable](k TokenKind, oper string) Rule[TokenKind] {
	return Rule[TokenKind]{true, k, func(s string) int {
		if !strings.HasPrefix(s, oper) {
			return NotMatched
		}
		completedOper := len(s) == len(oper) || !HasOperPrefix(s[len(oper):])
		if completedOper {
			return runeCount(oper)
		} else {
			return NotMatched
		}
	}}
}

func runeCount(s string) int { return utf8.RuneCountInString(s) }
