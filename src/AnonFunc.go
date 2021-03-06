package main

import (
	"fmt"
	"reflect"
	"strings"
)

// implements both Func and Value interfaces
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

func AssertAnonFunc(t Value) *AnonFunc {
	f, ok := t.(*AnonFunc)
	if ok {
		return f
	} else {
		panic("expected *AnonFunc, got " + reflect.TypeOf(t).String())
	}
}

func (f *AnonFunc) Dump() string {
	var b strings.Builder

	b.WriteString("\\(")

	if f.NumArgs() > 0 {
		b.WriteString(f.head.DumpArgs())
		b.WriteString("-> ")
	}
	b.WriteString(unwrapParens(f.body.Dump()))
	b.WriteString(")")

	return b.String()
}

func (f *AnonFunc) SetConstructors(cs []Call) Value {
	return &AnonFunc{f.setConstructors(cs)}
}

func (f *AnonFunc) Link(scope Scope, ew ErrorWriter) Value {
	return &AnonFunc{f.linkArgs(scope, ew)}
}

func (f *AnonFunc) SubVars(stack *Stack) Value {
	return &AnonFunc{f.subRhsVars(stack)}
}

func (f *AnonFunc) EvalRhs(args []Value, ew ErrorWriter) Value {
	if f.NumArgs() != len(args) {
		panic("should've been checked by caller")
	}

	d := f.FuncData.dispatch(args, ew)

	if d == nil || d.Failed() {
		// find out exactly which argument failes
		if len(args) != f.NumArgs() {
			ew.Add(f.Context().Error("unable to destructure, number of args differ"))
		} else {
			for i, arg := range args {
				des := f.head.args[i].Destructure(arg, ew)
				if des == nil || des.Failed() {
					ew.Add(f.Context().Error(fmt.Sprintf("unable to destructure arg %d\n  Have:\n    %s\n  Expect:\n    %s", i+1, arg.Dump(), f.head.args[i].DumpPretty())))
					break
				}
			}
		}

		return nil
	}

	return f.FuncData.EvalRhs(d)
}
