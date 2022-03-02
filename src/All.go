package main

// used to find the pattern matching distance of empty lists
type All struct {
	ValueData
}

func NewAll(ctx Context) *All {
	return &All{newValueData(ctx)}
}

func IsAll(v Value) bool {
	_, ok := v.(*All)
	return ok
}

func (v *All) Dump() string {
	return "All"
}

func (v *All) TypeName() string {
	return "All"
}

func (v *All) SetConstructors(cs []Call) Value {
	return &All{ValueData{newTokenData(v.Context()), cs}}
}

func (v *All) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *All) SubVars(stack *Stack) Value {
	return v
}
