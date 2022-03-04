package main

import (
	"reflect"
)

type RunFunc func(ioc IOContext) Value

type IO struct {
	ValueData
	Run RunFunc
}

func NewIO(Run RunFunc, ctx Context) *IO {
	return &IO{newValueData(ctx), Run}
}

func IsIO(v_ Value) bool {
	_, ok := v_.(*IO)
	return ok
}

func AssertIO(v_ Value) *IO {
	v, ok := v_.(*IO)
	if ok {
		return v
	} else {
		panic("expected *IO, got " + reflect.TypeOf(v_).String())
	}
}

func (v *IO) TypeName() string {
	return "IO"
}

func (v *IO) Dump() string {
	return "IO"
}

func (v *IO) SetConstructors(cs []Call) Value {
	return &IO{ValueData{newTokenData(v.Context()), cs}, v.Run}
}

func (v *IO) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *IO) SubVars(stack *Stack) Value {
	return v
}

func ResolveIO(res Value, ctx Context, ew ErrorWriter) Value {
	concrete, _ := EvalUntil(res, func(tn string) bool {
		return tn == "IO"
	}, ew)

	if concrete == nil || !ew.Empty() {
		if ew.Empty() {
			ew.Add(ctx.Error("expected IO return value"))
		}
		return nil
	} else {
		return concrete
	}
}
