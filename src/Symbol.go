package main

import (
	"reflect"
	"strings"
)

type Symbol struct {
	TokenData
	value string
}

var validSymbols = []string{
	"=",
	"+",
	"-",
	"*",
	"&",
	"^",
	"/",
	";",
	"(",
	")",
	",",
	"|",
	"[",
	"]",
	"{",
	"}",
	"==",
	"!",
	"!=",
	">=",
	">",
	"<=",
	"<",
	"&&",
	"||",
	"'",
	".",
	"%",
	":",
	"::",
	"?",
}

func NewSymbol(value string, ctx Context) *Symbol {
	return &Symbol{newTokenData(ctx), value}
}

func (t *Symbol) Value() string {
	return t.value
}

func (t *Symbol) Dump() string {
	return "Symbol" + t.value
}

func (t *Symbol) ToWord() *Word {
	return NewWord(t.Value(), t.Context())
}

func IsValidSymbol(s string) bool {
	for _, vs := range validSymbols {
		if vs == s {
			return true
		}
	}

	return false
}

func MaybeValidSymbol(s string) bool {
	for _, vs := range validSymbols {
		if strings.HasPrefix(vs, s) {
			return true
		}
	}

	return false
}

func AssertSymbol(t_ Token) *Symbol {
	t, ok := t_.(*Symbol)
	if !ok {
		panic("expected *Symbol, got " + reflect.TypeOf(t_).String())
	}

	return t
}

func IsSymbol(t_ Token, v ...string) bool {
	t, ok := t_.(*Symbol)

	if len(v) == 0 {
		return ok
	} else if len(v) == 1 {
		return ok && v[0] == t.value
	} else {
		panic("expected 0 or 1 parameters")
	}
}

func IsGroupSymbol(t Token) bool {
	return IsGroupOpenSymbol(t) || IsGroupCloseSymbol(t)
}

func IsGroupOpenSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		return t.Value() == "\\(" || strings.ContainsAny(t.value, "{[(")
	} else {
		return false
	}
}

func IsGroupCloseSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		return strings.ContainsAny(t.value, "}])")
	} else {
		return false
	}
}

func IsOperatorSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		_, ok_ := operators[t.value]
		return ok_
	} else {
		return false
	}
}

func ContainsSymbol(ts []Token, symbol string) bool {
	for _, t := range ts {
		if IsSymbol(t, symbol) {
			return true
		}
	}

	return false
}

func ContainsSymbolBefore(ts []Token, symbol string, other string) bool {
	for _, t := range ts {
		if IsSymbol(t, other) {
			return false
		} else if IsSymbol(t, symbol) {
			return true
		}
	}

	return false
}

func (open *Symbol) IsGroupMatch(close_ *Symbol) bool {
	switch open.Value() {
	case "(":
		return close_.Value() == ")"
	case "\\(":
		return close_.Value() == ")"
	case "{":
		return close_.Value() == "}"
	case "[":
		return close_.Value() == "]"
	default:
		return false
	}
}

func MatchingGroupCloseSymbol(openSymbol string) string {
	switch openSymbol {
	case "(", "\\(":
		return ")"
	case "[":
		return "]"
	case "{":
		return "}"
	default:
		panic("expected (, \\(, [ or {")
	}
}

func IsPostfixOperatorSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		op, ok_ := operators[t.value]
		if ok_ {
			return op.Postfix
		}
	}

	return false
}

func IsPrefixOperatorSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		op, ok_ := operators[t.value]
		if ok_ {
			return op.Prefix
		}
	}

	return false
}

func IsInfixOperatorSymbol(t_ Token) bool {
	t, ok := t_.(*Symbol)

	if ok {
		op, ok_ := operators[t.value]
		if ok_ {
			return op.Infix
		}
	}

	return false
}

func SplitByCommas(ts []Token, ew ErrorWriter) [][]Token {
	// look for commas
	parts := [][]Token{}

	prev := -1
	count := 0
	for i, t := range ts {
		if IsGroupOpenSymbol(t) {
			count += 1
		} else if IsGroupCloseSymbol(t) {
			count -= 1
		} else if IsSymbol(t, ",") && count == 0 {
			if prev+1 == i {
				ew.Add(t.Context().Error("expected expression before comma"))
			} else {
				parts = append(parts, ts[prev+1:i])
				prev = i
			}
		}
	}

	if prev < len(ts)-1 {
		parts = append(parts, ts[prev+1:])
	}

	if count != 0 {
		panic("unmatched group")
	}

	return parts
}

func FindNextSymbol(ts []Token, start int, symbol string) int {
	for i := start; i < len(ts); i++ {
		if IsSymbol(ts[i], symbol) {
			return i
		}
	}

	return -1
}

func FindAllNextSymbols(ts []Token, start int, symbol string) []int {
	res := []int{}

	for i := start; i < len(ts); i++ {
		if IsSymbol(ts[i], symbol) {
			res = append(res, i)
		}
	}

	return res
}

func FindLastSymbol(ts []Token, symbol string) int {
	for i := len(ts) - 1; i >= 0; i-- {
		if IsSymbol(ts[i], symbol) {
			return i
		}
	}

	return -1
}

func FindGroupMatch(ts []Token, start int, open *Symbol, ew ErrorWriter) int {
	openSymbol := open.Value()

	stack := []*Symbol{}

	closeSymbol := MatchingGroupCloseSymbol(openSymbol)

	noErrors := true

Outer:
	for i := start + 1; i < len(ts); i++ {
		t := ts[i]

		switch {
		case IsSymbol(t, closeSymbol) && len(stack) == 0:
			if noErrors {
				return i
			} else {
				return -1
			}
		case IsGroupOpenSymbol(t):
			stack = append(stack, AssertSymbol(t))
		case IsGroupCloseSymbol(t):
			s := AssertSymbol(t)

			if len(stack) == 0 {
				ew.Add(s.Context().Error("unmatched \"" + s.Value() + "\""))
				noErrors = false
			} else {
				last := stack[len(stack)-1]
				if !last.IsGroupMatch(s) {
					ew.Add(s.Context().Error("unmatched \"" + last.Value() + "\""))
					ew.Add(s.Context().Error("unmatched \"" + s.Value() + "\""))
					noErrors = false
				} else {
					stack = stack[0 : len(stack)-1]
				}
			}
		default:
			continue Outer
		}
	}

	ew.Add(open.Context().Error("unmatched \"" + openSymbol + "\""))
	noErrors = false

	return -1
}
