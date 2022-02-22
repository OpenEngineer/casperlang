package main

import "strings"

type BuiltinCall struct {
	ctx   Context
	name  string
	args  []Value
	type_ Type
	eval  func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value
}

func NewBuiltinCall(name string, args []Value, type_ Type, eval func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value, ctx Context) *BuiltinCall {
	return &BuiltinCall{ctx, name, args, type_, eval}
}

func (f *BuiltinCall) Context() Context {
	return f.ctx
}

func (f *BuiltinCall) Type() Type {
	return f.type_
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

func (f *BuiltinCall) Eval(scope Scope, ew ErrorWriter) Value {
	return f.eval(f, scope, ew)
}

func (f *BuiltinCall) Update(type_ Type, ctx Context) Value {
	return f.UpdateArgs(f.args, type_, ctx)
}

func (f *BuiltinCall) UpdateArgs(args []Value, type_ Type, ctx Context) *BuiltinCall {
	return NewBuiltinCall(f.name, args, type_, f.eval, ctx)
}

func (f *BuiltinCall) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, arg := range f.args {
		arg.CheckTypeNames(scope, ew)
	}
}
