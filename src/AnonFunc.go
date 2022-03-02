package main

import (
	"strings"
)

type AnonFunc struct {
	FuncData
}

func NewAnonFunc(args []Pattern, body Value, ctx Context) *AnonFunc {
	return &AnonFunc{newFuncData(nil, args, body, ctx)}
}

func NewSingleArgAnonFunc(arg Pattern, body Value, ctx Context) *AnonFunc {
	return NewAnonFunc([]Pattern{arg}, body, ctx)
}

func NewNoArgAnonFunc(body Value, ctx Context) *AnonFunc {
	return NewAnonFunc([]Pattern{}, body, ctx)
}

func (f *AnonFunc) Dump() string {
	var b strings.Builder

	b.WriteString("\\(")

	if f.NumArgs() > 0 {
		b.WriteString(f.head.DumpArgs())
		b.WriteString("-> ")
	}
	b.WriteString(f.body.Dump())
	b.WriteString(")")

	return b.String()
}

func (f *AnonFunc) Link(scope Scope, ew ErrorWriter) Value {
	return &AnonFunc{f.linkArgs(scope, ew)}
}

func (f *AnonFunc) SetConstructors(cs []Call) Value {
	return &AnonFunc{f.setConstructors(cs)}
}

func (f *AnonFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	d := f.FuncData.dispatch(args, ew)
	d.SetFunc(f)
	return d
}
