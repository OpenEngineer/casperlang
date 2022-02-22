package main

import (
	"fmt"
)

type Float struct {
	ValueData
	f float64
}

func NewFloat(f float64, ctx Context) *Float {
	return &Float{newValueData(NewFloatType(), ctx), f}
}

func (v *Float) Update(type_ Type, ctx Context) Value {
	return &Float{newValueData(type_, ctx), v.f}
}

func (t *Float) Dump() string {
	return fmt.Sprintf("%g", t.f)
}

func (t *Float) Value() float64 {
	return t.f
}

func IsFloat(t Token) bool {
	_, ok := t.(*Float)
	return ok
}

func AssertFloat(t_ Token) *Float {
	t, ok := t_.(*Float)

	if ok {
		return t
	} else {
		panic("expected *Float")
	}
}

func (v *Float) Eval(scope Scope, ew ErrorWriter) Value {
	return v
}
