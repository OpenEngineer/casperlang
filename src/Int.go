package main

import (
	"strconv"
)

type Int struct {
	ValueData
	i int64
}

func NewInt(i int64, ctx Context) *Int {
	return &Int{newValueData(NewIntType(), ctx), i}
}

func (v *Int) Update(type_ Type, ctx Context) Value {
	return &Int{newValueData(type_, ctx), v.i}
}

func (t *Int) Dump() string {
	return strconv.FormatInt(t.i, 10)
}

func (t *Int) Value() int64 {
	return t.i
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

// creates an immutable variable
func (v *Int) Eval(scope Scope, ew ErrorWriter) Value {
	return v
}
