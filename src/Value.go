package main

import (
	"strings"
)

type Value interface {
	Token

	TypeName() string
	Constructors() []Call
	SetConstructors(cs []Call) Value

	Link(scope Scope, ew ErrorWriter) Value // returns a value that no longer depends on scopes
	SubVars(stack *Stack) Value             // substitute variables
}

type ValueData struct {
	TokenData
	constructors []Call
}

func IsValue(t Token) bool {
	_, ok := t.(Value)
	return ok
}

func AssertValue(t_ Token) Value {
	t, ok := t_.(Value)
	if ok {
		return t
	} else {
		panic("expected Value")
	}
}

func newValueData(ctx Context) ValueData {
	return ValueData{newTokenData(ctx), make([]Call, 0)}
}

func (v *ValueData) Constructors() []Call {
	return v.constructors
}

func DumpValues(ts []Value) string {
	var b strings.Builder

	for i, t := range ts {
		b.WriteString(t.Dump())
		if i < len(ts)-1 {
			b.WriteString(" ")
		}
	}

	return b.String()
}
