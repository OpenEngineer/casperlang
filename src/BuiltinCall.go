package main

import (
	"strings"
)

type BuiltinCall struct {
	ValueData
	name  string
	args  []Value
	links map[string][]Func
	eval  func(self *BuiltinCall, ew ErrorWriter) Value
}

func NewBuiltinCall(name string, args []Value, links map[string][]Func, eval func(self *BuiltinCall, ew ErrorWriter) Value, ctx Context) *BuiltinCall {
	return &BuiltinCall{newValueData(ctx), name, args, links, eval}
}

func (f *BuiltinCall) Name() string {
	return f.name
}

func (f *BuiltinCall) Dump() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(f.name)

	for _, arg := range f.args {
		b.WriteString(" ")
		b.WriteString(arg.Dump())
	}

	b.WriteString(")")

	return b.String()
}

func (f *BuiltinCall) TypeName() string {
	if isConstructorName(f.name) {
		return f.name
	} else {
		return ""
	}
}

func (f *BuiltinCall) NumArgs() int {
	return len(f.args)
}

func (f *BuiltinCall) Args() []Value {
	res := make([]Value, len(f.args))

	for i, arg := range f.args {
		res[i] = arg
	}

	return res
}

func (f *BuiltinCall) SetConstructors(cs []Call) Value {
	return &BuiltinCall{ValueData{newTokenData(f.Context()), cs}, f.name, f.args, f.links, f.eval}
}

func (f *BuiltinCall) Link(scope Scope, ew ErrorWriter) Value {
	return f
}

func (f *BuiltinCall) Eval(ew ErrorWriter) Value {
	v := f.eval(f, ew)

	cs := f.Constructors()
	if isConstructorName(f.name) {
		cs = append(make([]Call, 0), cs...)
		cs = append(cs, f)
	}

	v = v.SetConstructors(cs)
	return v
}
