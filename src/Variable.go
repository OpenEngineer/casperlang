package main

// mutable Value (and the only mutable one!, and no other point should values be mutated)
type Variable struct {
	ValueData
	val Value
}

func NewVariable(type_ Type, ctx Context) *Variable {
	return &Variable{newValueData(type_, ctx), nil}
}

func IsVariable(v Value) bool {
	_, ok := v.(*Variable)
	return ok
}

func AssertVariable(v_ Value) *Variable {
	v, ok := v_.(*Variable)
	if ok {
		return v
	} else {
		panic("expected *Variable")
	}
}

// NOTE: this is the only mutating function in casperlang!
func (v *Variable) SetValue(val Value) {
	v.val = val
}

func (v *Variable) Dump() string {
	if v.Type() == nil {
		return "<var>"
	} else {
		return "<var::" + v.Type().Dump() + ">"
	}
}

func (v *Variable) Eval(scope Scope, ew ErrorWriter) Value {
	if v.val != nil {
		return v.val
	} else {
		return v
	}
}

// can't update variable by sending back new, because original context would get lost
func (v *Variable) Update(type_ Type, ctx Context) Value {
	v.type_ = type_
	v.ctx = ctx

	if v.val != nil {
		v.val = v.val.Update(type_, ctx)
	}

	return v
}
