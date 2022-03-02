package main

import (
	"fmt"
)

type Float struct {
	ValueData
	f float64
}

func NewFloat(f float64, ctx Context) *Float {
	return &Float{newValueData(ctx), f}
}

func (t *Float) Dump() string {
	return fmt.Sprintf("%g", t.f)
}

func (t *Float) Value() float64 {
	return t.f
}

func (v *Float) TypeName() string {
	return "Float"
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

func (v *Float) SetConstructors(cs []Call) Value {
	return &Float{ValueData{newTokenData(v.Context()), cs}, v.f}
}

func (v *Float) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *Float) Eval(stack *Stack, ew ErrorWriter) Value {
	return v
}
