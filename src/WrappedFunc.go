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
	return fn.fn.Dump()
}

func (fn *WrappedFunc) Context() Context {
	return fn.fn.Context()
}

func (fn *WrappedFunc) Name() string {
	return fn.fn.Name()
}

func (fn *WrappedFunc) NumArgs() int {
	return fn.fn.NumArgs()
}

func (fn *WrappedFunc) DumpHead() string {
	return fn.fn.DumpHead()
}

func (fn *WrappedFunc) DumpPrettyHead() string {
	return fn.fn.DumpPrettyHead()
}

func (fn *WrappedFunc) ListHeaderTypes() []string {
	return fn.fn.ListHeaderTypes()
}

func (fn *WrappedFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	return fn.fn.Dispatch(args, ew)
}

func (fn *WrappedFunc) Link(scope Scope, ew ErrorWriter) Func {
	return fn
}

func (fn *WrappedFunc) EvalRhs(d *Dispatched) Value {
	return fn.fn.EvalRhs(d)
}
