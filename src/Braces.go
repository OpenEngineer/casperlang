package main

import (
	"strings"
)

type Braces struct {
	TokenData
	keys []*String
	vals [][]Token
}

func NewBraces(keys []*String, vals [][]Token, ctx Context) *Braces {
	return &Braces{newTokenData(ctx), keys, vals}
}

func IsBraces(t Token) bool {
	_, ok := t.(*Braces)
	return ok
}

func (t *Braces) Empty() bool {
	return len(t.keys) == 0
}

func IsEmptyBraces(t_ Token) bool {
	t, ok := t_.(*Braces)
	if ok {
		return t.Empty()
	} else {
		return false
	}
}

func IsNonEmptyBraces(t_ Token) bool {
	t, ok := t_.(*Braces)
	if ok {
		return !t.Empty()
	} else {
		return false
	}
}

func AssertBraces(t_ Token) *Braces {
	t, ok := t_.(*Braces)
	if ok {
		return t
	} else {
		panic("expected *Braces")
	}
}

func (t *Braces) Dump() string {
	var b strings.Builder

	b.WriteString("{")

	for i, key := range t.keys {
		b.WriteString(key.Dump())
		b.WriteString(":")
		b.WriteString(DumpTokens(t.vals[i]))
		if i < len(t.keys)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String()
}

func GroupBraces(ts []Token, ew ErrorWriter) Token {
	ctx := ts[0].Context()

	fields := SplitByCommas(ts[1:len(ts)-1], ew)

	keys := []*String{}
	vals := [][]Token{}

	for _, field := range fields {
		if len(field) < 3 {
			ew.Add(field[0].Context().Error("invalid dict syntax"))
		} else {
			keyRaw := field[0]
			colon := field[1]

			switch {
			case IsWord(keyRaw):
				keys = append(keys, AssertWord(keyRaw).ToString())
			case IsString(keyRaw):
				keys = append(keys, AssertString(keyRaw))
			default:
				ew.Add(keyRaw.Context().Error("invalid dict key"))
				continue
			}

			if !IsSymbol(colon, ":") {
				ew.Add(colon.Context().Error("invalid dict syntax"))
				continue
			}

			val := GroupTokens(field[2:], ew)
			vals = append(vals, val)
		}
	}

	return NewBraces(keys, vals, ctx)
}
