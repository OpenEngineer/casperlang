package main

import "strings"

type Func interface {
	Token

	Name() string // empty for anonymous functions
	NumArgs() int
	DumpHead() string
	ListHeaderTypes() []string

	// args should be detached at this point
	Dispatch(args []Value, ew ErrorWriter) *Dispatched // args here should already have all their values substituted
	Link(scope Scope, ew ErrorWriter) Func

	EvalRhs(d *Dispatched) Value
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
