package main

import "strings"

type Func interface {
	Token

	Name() string // empty for anonymous functions
	NumArgs() int
	DumpHead() string
	DumpPrettyHead() string
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

func DeferFunc(fn Func, args []Value, ctx Context) Value {
	return NewDisCall([]Func{fn}, args, ctx)
}
