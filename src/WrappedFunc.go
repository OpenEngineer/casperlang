package main

type WrappedFunc struct {
	fn Func
}

func NewWrappedFunc(fn Func) *WrappedFunc {
	return &WrappedFunc{fn}
}

func (fn *WrappedFunc) SetFunc(fn_ Func) {
	fn.fn = fn_
}

func (fn *WrappedFunc) Dump() string {
	if fn.fn == nil {
		return "<wrappedfunc>"
	} else {
		return fn.fn.Dump()
	}
}

func (fn *WrappedFunc) Context() Context {
	if fn.fn == nil {
		return NewSpecialContext("<wrappedfunc>")
	} else {
		return fn.fn.Context()
	}
}

func (fn *WrappedFunc) TypeName() string {
	if fn.fn == nil {
		return "<wrappedfunc>"
	} else {
		return fn.fn.TypeName()
	}
}

func (fn *WrappedFunc) Link(scope Scope, ew ErrorWriter) Value {
	return fn
}

func (fn *WrappedFunc) Constructors() []Call {
	return fn.fn.Constructors()
}

func (fn *WrappedFunc) SetConstructors(cs []Call) Value {
	return fn.fn.SetConstructors(cs)
}

func (fn *WrappedFunc) Name() string {
	if fn.fn == nil {
		return "<wrappedfunc>"
	} else {
		return fn.fn.Name()
	}
}

func (fn *WrappedFunc) NumArgs() int {
	return fn.fn.NumArgs()
}

func (fn *WrappedFunc) DumpHead() string {
	return fn.fn.DumpHead()
}

func (fn *WrappedFunc) ListHeaderTypes() []string {
	return fn.fn.ListHeaderTypes()
}

func (fn *WrappedFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	return fn.fn.Dispatch(args, ew)
}

func (fn *WrappedFunc) EvalRhs(d *Dispatched) Value {
	return fn.fn.EvalRhs(d)
}
