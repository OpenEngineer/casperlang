package main

import (
	"strings"
)

type Brackets struct {
	TokenData
	values [][]Token
}

func NewBrackets(values [][]Token, ctx Context) *Brackets {
	return &Brackets{newTokenData(ctx), values}
}

func (t *Brackets) Empty() bool {
	return len(t.values) == 0
}

func IsBrackets(t Token) bool {
	_, ok := t.(*Brackets)
	return ok
}

func IsEmptyBrackets(t_ Token) bool {
	t, ok := t_.(*Brackets)
	if ok {
		return t.Empty()
	} else {
		return false
	}
}

func IsNonEmptyBrackets(t_ Token) bool {
	t, ok := t_.(*Brackets)
	if ok {
		return !t.Empty()
	} else {
		return false
	}
}

func AssertBrackets(t_ Token) *Brackets {
	t, ok := t_.(*Brackets)
	if ok {
		return t
	} else {
		panic("expected *Brackets")
	}
}

func (t *Brackets) Dump() string {
	var b strings.Builder

	b.WriteString("[")

	for i, v := range t.values {
		b.WriteString(DumpTokens(v))
		if i < len(t.values)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String()
}

func GroupBrackets(ts []Token, ew ErrorWriter) Token {
	ctx := ts[0].Context()

	fields := SplitByCommas(ts[1:len(ts)-1], ew)

	resFields := [][]Token{}

	for _, field := range fields {
		resField := GroupTokens(field, ew)
		resFields = append(resFields, resField)
	}

	return NewBrackets(resFields, ctx)
}
