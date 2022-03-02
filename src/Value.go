package main

import (
	"fmt"
	"strings"
)

type Value interface {
	Token

	TypeName() string

	Link(scope Scope, ew ErrorWriter) Value // returns a value that no longer depends on scopes

	Constructors() []Call
	SetConstructors(cs []Call) Value
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

func (v *ValueData) SetConstructors(cs []Call) Value {
	panic("SetConstructors() not (yet) defined")
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

func EvalUntil(arg Value, cond func(string) bool, ew ErrorWriter) (Value, Value) {
	if arg == nil {
		return nil, nil
	}

	for _, c := range arg.Constructors() {
		if cond(c.TypeName()) {
			return arg, c
		}
	}

	if arg.TypeName() == "All" {
		return arg, arg
	}

	for {
		if cond(arg.TypeName()) {
			return arg, arg
		} else {
			call, ok := arg.(Call)
			if !ok {
				return nil, nil
			}

			fmt.Println("calling", call.Dump())
			arg = call.Eval(ew)
			if arg == nil || !ew.Empty() {
				return nil, nil
			}
		}
	}
}
