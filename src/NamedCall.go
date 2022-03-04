package main

type NamedCall struct {
	CallData
	name *Word
}

func NewNamedCall(name *Word, args []Value) *NamedCall {
	return &NamedCall{newCallData(args, name.Context()), name}
}

func (c *NamedCall) Name() string {
	return c.name.Value()
}

func (t *NamedCall) TypeName() string {
	if isConstructorName(t.Name()) {
		return t.Name()
	} else {
		return ""
	}
}

func (t *NamedCall) Dump() string {
	return t.CallData.dump(t.Name())
}

func (v *NamedCall) Link(scope Scope, ew ErrorWriter) Value {
	local := scope.GetLocal(v.Name())
	if local != nil {
		args := v.linkArgs(scope, ew)

		ref := NewVarRef(local, v.Context())
		if len(args) == 0 {
			return ref
		} else {
			return NewBlindCall(append([]Value{ref}, args...), v.name.Context())
		}
	} else {
		fns := scope.ListDispatchable(v.Name(), v.NumArgs(), ew)
		if len(fns) == 0 {

			fns = scope.ListDispatchable(v.Name(), -1, ew)

			if len(fns) == 0 {
				ew.Add(v.name.Context().Error("\"" + v.Name() + "\" undefined"))
				return v
			} else {
				ew.Add(v.name.Context().Error(badDispatchMessage(v.Name(), v.Args(), "unable to dispatch", fns)))
				return v
			}
		}

		args := v.linkArgs(scope, ew)

		return NewDisCall(fns, args, v.name.Context())
	}
}

func (t *NamedCall) SetConstructors(cs []Call) Value {
	panic("can't SetConstructors of NamedCall, should've been linked")
}

func (t *NamedCall) SubVars(stack *Stack) Value {
	panic("can't subVars NamedCall, should've been linked")
}

func (t *NamedCall) Eval(stack *Stack, ew ErrorWriter) Value {
	panic("can't eval NamedCall, should've been linked")
}
