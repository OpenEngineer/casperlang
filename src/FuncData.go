package main

import (
	"strconv"
)

type FuncData struct {
	ValueData
	head *FuncHeader
	body Value
}

func newFuncData(name *Word, args []Pattern, body Value, ctx Context) FuncData {
	return FuncData{newValueData(ctx), NewFuncHeader(name, args), body}
}

func (f *FuncData) Name() string {
	return f.head.Name()
}

func (f *FuncData) TypeName() string {
	return "\\" + strconv.Itoa(f.NumArgs())
}

func (f *FuncData) NumArgs() int {
	return f.head.NumArgs()
}

func (f *FuncData) IsConstructor() bool {
	return f.head.IsConstructor()
}

func (f *FuncData) DumpHead() string {
	return f.head.Dump()
}

func (f *FuncData) DumpPrettyHead() string {
	return f.head.DumpPretty()
}

func (f *FuncData) ListHeaderTypes() []string {
	return f.head.ListTypes()
}

func (f *FuncData) setConstructors(cs []Call) FuncData {
	return FuncData{ValueData{newTokenData(f.Context()), cs}, f.head, f.body}
}

func (f *FuncData) linkArgs(scope Scope, ew ErrorWriter) FuncData {
	head, fnScope := f.head.Link(scope, ew)
	body := f.body.Link(fnScope, ew)

	return FuncData{newValueData(f.Context()), head, body}
}

func (f *FuncData) subRhsVars(stack *Stack) FuncData {
	return FuncData{newValueData(f.Context()), f.head, f.body.SubVars(stack)}
}

// args should already have all their vars substituted
func (f *FuncData) dispatch(args []Value, ew ErrorWriter) *Dispatched {
	return f.head.Destructure(args, ew)
}

// how can modify this with a stack
func (f *FuncData) EvalRhs(d *Dispatched) Value {
	return f.body.SubVars(d.stack)
}
