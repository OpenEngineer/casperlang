package main

type Any struct {
	ValueData
}

func NewAny(ctx Context) *Any {
	return &Any{newValueData(ctx)}
}

func (t *Any) Dump() string {
	return "Any"
}

func (v *Any) TypeName() string {
	return "Any"
}

func (v *Any) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *Any) SetConstructors(cs []Call) Value {
	return &Any{ValueData{newTokenData(v.Context()), cs}}
}
