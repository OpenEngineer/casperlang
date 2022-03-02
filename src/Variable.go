package main

type Variable struct {
	TokenData
	name string
}

func NewVariable(name string, ctx Context) *Variable {
	return &Variable{newTokenData(ctx), name}
}

func (v *Variable) Dump() string {
	return "<variable::" + v.Name() + ">"
}

func (v *Variable) Name() string {
	return v.name
}
