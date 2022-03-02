package main

import "fmt"

type Variable struct {
	ValueData
	name string
	data Value
}

func NewVariable(name string, ctx Context) *Variable {
	return &Variable{newValueData(ctx), name, nil}
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

func (v *Variable) Data() Value {
	return v.data
}

func (v *Variable) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *Variable) SetData(val Value) {
	v.data = val
}
