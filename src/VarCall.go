package main

import "fmt"

type VarCall struct {
	CallData
	var_ *Variable
}

func NewVarCall(var_ *Variable, args []Value, ctx Context) *VarCall {
	return &VarCall{newCallData(args, ctx), var_}
}

func (c *VarCall) Name() string {
	return c.var_.Name()
}

func (c *VarCall) TypeName() string {
	return ""
}

func (t *VarCall) Dump() string {
	return t.CallData.dump(t.var_.Dump())
}

func (v *VarCall) SetConstructors(cs []Call) Value {
	return v.setConstructors(cs)
}

func (v *VarCall) setConstructors(cs []Call) Call {
	return &VarCall{v.CallData.setConstructors(cs), v.var_}
}

func (v *VarCall) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *VarCall) Eval(ew ErrorWriter) Value {
	data := v.var_.Data()
	if v.Name() == "$1" {
		fmt.Println("data in $1:", data.Dump())
	}

	if v.NumArgs() == 0 {
		return data
	}

	fn, ok := data.(Func)
	if !ok {
		ew.Add(v.Context().Error("not a function"))
		return nil
	} else {

		d := fn.Dispatch(v.args, ew)

		if !ew.Empty() {
			return nil
		} else if d == nil {
			ew.Add(v.Context().Error("unable to destructure"))
			return nil
		}

		res := d.Eval()

		return res
	}
}
