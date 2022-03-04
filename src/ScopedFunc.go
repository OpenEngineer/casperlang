package main

type ScopedFunc struct {
	fn    Func
	scope Scope // eg. File or GlobalScope
}

func NewScopedFunc(fn Func, scope Scope) *ScopedFunc {
	return &ScopedFunc{fn, scope}
}

func (f *ScopedFunc) Context() Context {
	return f.fn.Context()
}

func (f *ScopedFunc) Name() string {
	return f.fn.Name()
}

func (f *ScopedFunc) NumArgs() int {
	return f.fn.NumArgs()
}

func (f *ScopedFunc) Dump() string {
	return f.fn.Dump()
}

func (f *ScopedFunc) DumpHead() string {
	return f.fn.DumpHead()
}

func (f *ScopedFunc) DumpPrettyHead() string {
	return f.fn.DumpPrettyHead()
}

func (f *ScopedFunc) Link(_ Scope, ew ErrorWriter) Func {
	return f.fn.Link(f.scope, ew)
}

func (f *ScopedFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	return f.fn.Dispatch(args, ew)
}

func (f *ScopedFunc) EvalRhs(d *Dispatched) Value {
	panic("should've been turned into regular func")
}

func (f *ScopedFunc) ListHeaderTypes() []string {
	return f.fn.ListHeaderTypes()
}
