package main

import (
	"strings"
)

type Token interface {
	Dump() string
	Context() Context
}

type TokenData struct {
	ctx Context
}

func newTokenData(ctx Context) TokenData {
	return TokenData{ctx}
}

func (t *TokenData) Context() Context {
	return t.ctx
}

func DumpTokens(ts []Token) string {
	var b strings.Builder

	for i, t := range ts {
		b.WriteString(t.Dump())
		if i < len(ts)-1 {
			b.WriteString(" ")
		}
	}

	return b.String()
}

func IsLiteral(t Token) bool {
	return IsInt(t) || IsFloat(t) || IsString(t)
}

func IsContainer(t Token) bool {
	return IsDict(t) || IsList(t)
}

func unwrapParens(s string) string {
	n := len(s)
	if n <= 2 {
		return s
	} else if s[0] == '(' && s[n-1] == ')' {
		return s[1 : n-1]
	} else {
		return s
	}
}
