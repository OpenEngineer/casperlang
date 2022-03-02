package main

import "fmt"

type Variable struct {
	ValueData
	name string
}

func NewVariable(name string, ctx Context) *Variable {
	return &Variable{newValueData(ctx), name}
}

func (v *Variable) Dump() string {
	return "<variable::" + v.name + ">"
}

func (v *Variable) TypeName() string {
	if v.data == nil {
		fmt.Println("variable with unset data returning type name")
		return ""
	} else {
		return v.data.TypeName()
	}
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *Variable) Eval(stack *Stack, ew ErrorWriter) Value {
	panic("should be wrapped by VarCall")
}
