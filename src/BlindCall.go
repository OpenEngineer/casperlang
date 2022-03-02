package main

import (
	"fmt"
	"strings"
)

type BlindCall struct {
	CallData
	fn Value // evaluate this until TypeName begins with '\'
}

func NewBlindCall(args []Value, ctx Context) Call {
	fn := args[0]
	args = args[1:]

	return &BlindCall{newCallData(args, ctx), fn}
}

func (v *BlindCall) TypeName() string {
	return ""
}

func (v *BlindCall) Name() string {
	return ""
}

func (t *BlindCall) Dump() string {
	return t.CallData.dump(t.fn.Dump())
}

func (v *BlindCall) Link(scope Scope, ew ErrorWriter) Value {
	fn := v.fn.Link(scope, ew)
	args := v.linkArgs(scope, ew)

	return &BlindCall{newCallData(args, v.Context()), fn}
}

func (v *BlindCall) SetConstructors(cs []Call) Value {
	return &BlindCall{v.setConstructors(cs), v.fn}
}

func (t *BlindCall) Eval(stack *Stack, ew ErrorWriter) Value {
	fn_, _ := EvalUntil(t.fn, stack, func(tn string) bool {
		return strings.HasPrefix(tn, "\\")
	}, ew)

	if fn_ == nil {
		ew.Add(t.Context().Error("not a function"))
		return nil
	}

	fn := AssertAnonFunc(fn_)

	if fn.NumArgs() != t.NumArgs() {
		ew.Add(t.Context().Error(fmt.Sprintf("expected function with %d args, got %s", t.NumArgs(), fn.NumArgs())))
		return nil
	}

	v := fn.EvalRhs(t.args, stack, ew)
	v = v.SetConstructors(t.Constructors())

	return v
}
