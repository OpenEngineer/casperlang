package main

import "strings"

type Value interface {
	Token

	Type() Type

	CheckTypeNames(scope Scope, ew ErrorWriter)

	Eval(scope Scope, ew ErrorWriter) Value // updated type shouldn't change when calling this function

	Update(type_ Type, ctx Context) Value
	// TODO: serialization
}

type ValueData struct {
	TokenData
	type_ Type
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

func newValueData(type_ Type, ctx Context) ValueData {
	return ValueData{newTokenData(ctx), type_}
}

func NewValueData(type_ Type, ctx Context) *ValueData {
	vd := newValueData(type_, ctx)

	return &vd
}

func NewGenericValue(type_ Type, ctx Context) *ValueData {
	return NewValueData(type_, ctx)
}

func (v *ValueData) Update(type_ Type, ctx Context) Value {
	return &ValueData{newTokenData(ctx), type_}
}

func (v *ValueData) Type() Type {
	return v.type_
}

func (v *ValueData) Eval(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *ValueData) Dump() string {
	if v.type_ == nil {
		return "<value>"
	} else {
		return "<value::" + v.type_.Dump() + ">"
	}
}

func (v *ValueData) CheckTypeNames(scope Scope, ew ErrorWriter) {
	return
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
