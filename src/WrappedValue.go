package main

type WrappedValue struct {
	inner Value
	stack *Stack
}

func NewWrappedValue(inner Value, stack *Stack) *WrappedValue {
	return &WrappedValue{inner, stack}
}

func (t *WrappedValue) Context() Context {
	return t.inner.Context()
}

func (t *WrappedValue) Dump() string {
	return t.inner.Dump()
}

func (t *WrappedValue) TypeName() string {
	return t.inner.TypeName()
}

func (t *WrappedValue) Constructors() []Call {
	return t.inner.Constructors()
}

func (t *WrappedValue) SetConstructors(cs []Call) Value {
	return NewWrappedValue(t.inner.SetConstructors(cs), t.stack)
}

func (t *WrappedValue) Link(scope Scope, ew ErrorWriter) Value {
	return NewWrappedValue(t.inner.Link(scope, ew), t.stack)
}

func (t *WrappedValue) Eval(_ *Stack, ew ErrorWriter) Value {
	return t.inner.Eval(t.stack, ew)
}
