package main

type FuncData struct {
	ValueData
	head *FuncHeader
	body Value
}

func newFuncData(name *Word, args []Pattern, body Value, ctx Context) FuncData {
	return FuncData{newValueData(nil, ctx), NewFuncHeader(name, args), body}
}

func (v *FuncData) Type() Type {
	return NewFuncType(v.NumArgs())
}

func (f *FuncData) Name() string {
	return f.head.Name()
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

func (f *FuncData) call(scope Scope, args []Value, ctx Context, ew ErrorWriter) Value {
	//for _, arg := range args {
	//if IsDeferredError(arg) {
	//return arg
	//}
	//}

	subScope := f.head.DestructureArgs(scope, args, ctx, ew)
	if !ew.Empty() {
		return nil
	}

	v := f.body.Eval(subScope, ew)
	if !ew.Empty() {
		return nil
	}

	return v
}

func (f *FuncData) CheckTypeNames(scope Scope, ew ErrorWriter) {
	f.head.CheckTypeNames(scope, ew)
	f.body.CheckTypeNames(scope, ew)
}
