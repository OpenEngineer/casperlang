package main

import (
	"strconv"
	"strings"
)

type AnonFunc struct {
	FuncData
	scope Scope // scope is different from where it is called
}

func NewAnonFunc(args []Pattern, body Value, ctx Context) *AnonFunc {
	return &AnonFunc{newFuncData(nil, args, body, ctx), nil}
}

func NewSingleArgAnonFunc(arg Pattern, body Value, ctx Context) *AnonFunc {
	return NewAnonFunc([]Pattern{arg}, body, ctx)
}

func NewNoArgAnonFunc(body Value, ctx Context) *AnonFunc {
	return NewAnonFunc([]Pattern{}, body, ctx)
}

func IsAnonFunc(t Token) bool {
	_, ok := t.(*AnonFunc)
	return ok
}

func AssertAnonFunc(t_ Token) *AnonFunc {
	t, ok := t_.(*AnonFunc)

	if ok {
		return t
	} else {
		panic("expected *AnonFunc")
	}
}

func (v *AnonFunc) Update(type_ Type, ctx Context) Value {
	return &AnonFunc{FuncData{newValueData(type_, ctx), v.head, v.body}, v.scope}
}

func (f *AnonFunc) Dump() string {
	var b strings.Builder

	b.WriteString("\\(")

	if f.NumArgs() > 0 {
		b.WriteString(f.head.DumpArgs())
		b.WriteString("= ")
	}
	b.WriteString(f.body.Dump())
	b.WriteString(")")

	return b.String()
}

func ParseAnon(gr *EscParens, ew ErrorWriter) *AnonFunc {
	ctx := gr.Context()

	args := []Pattern{}

	prev := 0

	for _, id := range gr.dollars {
		if id > prev+1 {
			for j := 0; j < id-(prev+1); j++ {
				args = append(args, NewNamedPattern(NewWord("_", ctx), NewSimplePattern(NewWord("Any", ctx)), ctx))
			}

		}

		args = append(args, NewNamedPattern(NewWord("$"+strconv.Itoa(id), ctx), NewSimplePattern(NewWord("Any", ctx)), ctx))

		prev = id
	}

	body := ParseExpr(gr.content, ew)

	return NewAnonFunc(args, body, ctx)
}

func (v *AnonFunc) Eval(scope Scope, ew ErrorWriter) Value {
	if v.scope == nil {
		return &AnonFunc{v.FuncData, scope}
	} else {
		return v
	}
}

func (v *AnonFunc) Call(args []Value, argScope Scope, ctx Context, ew ErrorWriter) Value {
	res := v.FuncData.call(v.scope, args, ctx, ew)
	if res == nil {
		return nil
	}

	return res.Update(res.Type(), ctx)
}
