package main

import "strings"

type Func interface {
	Value

	Name() string // empty for anonymous functions
	NumArgs() int
	DumpHead() string
	ListHeaderTypes() []string

	// args should be detached at this point
	Dispatch(args []Value, ew ErrorWriter) *Dispatched

	EvalRhs(d *Dispatched) Value
}

func IsFunc(v Value) bool {
	_, ok := v.(Func)
	return ok
}

func AssertFunc(v_ Value) Func {
	v, ok := v_.(Func)
	if ok {
		return v
	} else {
		panic("expected Func")
	}
}

func listAmbiguousFuncs(fns []Func) string {
	var b strings.Builder

	b.WriteString("\n  Definitions:\n")

	for i, fn := range fns {
		b.WriteString("    ")
		b.WriteString(fn.DumpHead())
		b.WriteString(" (")
		b.WriteString(fn.Context().Error("").Error())
		b.WriteString(")")

		if i < len(fns)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func RunFunc(fn Func, args []Value, ew ErrorWriter, ctx Context) Value {
	d := fn.Dispatch(args, ew)

	if d == nil || d.Failed() {
		ew.Add(ctx.Error("failed to run func"))
		return nil
	}

	return d.Eval()
}

func DeferFunc(fn Func, args []Value, ctx Context) Value {
	return NewDisCall([]Func{fn}, args, ctx)
}
