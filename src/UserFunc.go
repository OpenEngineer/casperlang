package main

import (
	"strings"
)

type UserFunc struct {
	FuncData
	file *File // referene to file where it is located (acts as scope)
}

func NewUserFunc(name *Word, args []Pattern, body Value, ctx Context) *UserFunc {
	return &UserFunc{newFuncData(name, args, body, ctx), nil}
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

func (v *UserFunc) Update(type_ Type, ctx Context) Value {
	return &UserFunc{FuncData{newValueData(type_, ctx), v.head, v.body}, v.file}
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

func (f *UserFunc) Eval(scope Scope, ew ErrorWriter) Value {
	panic("can't eval *UserFunc, can only eval *AnonFunc")
}

func (f *UserFunc) CalcDistance(args []Value) []int {
	return f.head.CalcDistance(args)
}

// XXX: should argScope be attached to args?
func (f *UserFunc) Call(args []Value, argScope Scope, ctx Context, ew ErrorWriter) Value {
	if f.file == nil {
		panic("userfunc file not set")
	}

	var scope Scope = f.file
	res := f.call(scope, args, ctx, ew)
	if res == nil {
		return nil
	} //else if IsDeferredError(res) {
	//return res
	//}

	t := res.Type()
	if f.IsConstructor() {
		t = NewUserType(res.Type(), f.Name(), args, ctx)
	}

	return res.Update(t, ctx)
}

func (f *UserFunc) ListHeaderTypes() []string {
	return f.head.ListTypes()
}
