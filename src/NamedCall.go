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

		if len(args) == 0 {
			return local
		} else {
			return NewBlindCall(append([]Value{local}, args...), v.name.Context())
		}
	} else {
		fns := scope.ListDispatchable(v.Name(), v.NumArgs(), ew)
		if len(fns) == 0 {
			ew.Add(v.name.Context().Error("\"" + v.Name() + "\" undefined"))
			return v
		}

		args := v.linkArgs(scope, ew)

		return NewDisCall(fns, args, v.name.Context())
	}
}

func (t *NamedCall) Eval(ew ErrorWriter) Value {
	panic("can't eval NamedCall, should've been linked")
}
