package main

type VarRef struct {
	ValueData
	var_ *Variable
}

func NewVarRef(var_ *Variable, ctx Context) *VarRef {
	return &VarRef{newValueData(ctx), var_}
}

func (v *VarRef) Dump() string {
	return "<reference::" + v.var_.Name() + ">"
}

func (v *VarRef) TypeName() string {
	return ""
}

func (v *VarRef) Name() string {
	return v.var_.Name()
}

func (v *VarRef) SetConstructors(cs []Call) Value {
	return &VarRef{ValueData{newTokenData(v.ctx), cs}, v.var_}
}

func (v *VarRef) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *VarRef) SubVars(stack *Stack) Value {
	res := stack.Get(v.var_)
	if res != nil {
		return res
	} else {
		return v
	}
}
