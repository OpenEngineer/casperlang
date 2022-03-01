package main

import (
	"strings"
)

type UserFunc struct {
	FuncData
}

func NewUserFunc(name *Word, args []Pattern, body Value, ctx Context) *UserFunc {
	return &UserFunc{newFuncData(name, args, body, ctx)}
}

func IsUserFunc(t Token) bool {
	_, ok := t.(*UserFunc)
	return ok
}

func AssertUserFunc(t_ Token) *UserFunc {
	t, ok := t_.(*UserFunc)

	if ok {
		return t
	} else {
		panic("expected *UserFunc")
	}
}

func (f *UserFunc) Dump() string {
	var b strings.Builder

	b.WriteString(f.Name())
	b.WriteString(" ")
	b.WriteString(f.head.DumpArgs())
	b.WriteString("= ")
	b.WriteString(f.body.Dump())

	return b.String()
}

func (f *UserFunc) Link(scope Scope, ew ErrorWriter) Value {
	return &UserFunc{f.linkArgs(scope, ew)}
}

func (f *UserFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	d := f.FuncData.dispatch(args, ew)
	if !d.Failed() {
		d.SetFunc(f)
	}
	return d
}
