package main

import (
	"strconv"
)

type Int struct {
	ValueData
	i int64
}

func NewInt(i int64, ctx Context) *Int {
	return &Int{newValueData(ctx), i}
}

func (t *Int) Dump() string {
	return strconv.FormatInt(t.i, 10)
}

func (t *Int) Value() int64 {
	return t.i
}

func (v *Int) TypeName() string {
	return "Int"
}

func IsInt(t Token) bool {
	_, ok := t.(*Int)
	return ok
}

func AssertInt(t_ Token) *Int {
	t, ok := t_.(*Int)
	if ok {
		return t
	} else {
		panic("expected *Int")
	}
}

func (v *Int) SetConstructors(cs []Call) Value {
	return &Int{ValueData{newTokenData(v.Context()), cs}, v.i}
}

func (v *Int) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *Int) SubVars(stack *Stack) Value {
	return v
}
