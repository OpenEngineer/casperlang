package main

import (
	"strings"
)

type Parens struct {
	TokenData
	content []Token
}

func NewParens(content []Token, ctx Context) *Parens {
	return &Parens{newTokenData(ctx), content}
}

func (t *Parens) Empty() bool {
	return len(t.content) == 0
}

func IsParens(t Token) bool {
	_, ok := t.(*Parens)
	return ok
}

func IsNonEmptyParens(t_ Token) bool {
	t, ok := t_.(*Parens)
	if ok {
		return !t.Empty()
	} else {
		return false
	}
}

func AssertParens(t_ Token) *Parens {
	t, ok := t_.(*Parens)
	if ok {
		return t
	} else {
		panic("expected *Parens")
	}
}

func (t *Parens) Content() []Token {
	return t.content
}

func (t *Parens) Dump() string {
	var b strings.Builder

	b.WriteString("(")

	for i, item := range t.content {
		b.WriteString(item.Dump())

		if i < len(t.content)-1 {
			b.WriteString(" ")
		}
	}

	b.WriteString(")")

	return b.String()
}

func GroupParens(ts []Token, ew ErrorWriter) Token {
	ctx := ts[0].Context()

	inner := GroupTokens(ts[1:len(ts)-1], ew)

	return NewParens(inner, ctx)
}
