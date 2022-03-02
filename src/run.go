package main

import (
	"fmt"
	"os"
)

func EvalNonLazy(v Value, ew ErrorWriter) Value {
	for v != nil {
		call, ok := v.(Call)
		if ok {
			v = call.Eval(ew)
		} else {
			break
		}
	}

	return v
}

// v might've been reduced by a previous call to Eval
func Run(v Value) {
	ew := NewErrorWriter()

	v = EvalNonLazy(v, ew)

	if v != nil && ew.Empty() {
		if !IsIO(v) {
			ew.Add(v.Context().Error("expected IO, got " + v.Dump()))
		} else {
			io := AssertIO(v)

			v = io.Run()

			if v != nil {
				ew.Add(v.Context().Error("unused IO result"))
			}
		}
	}

	if !ew.Empty() {
		fmt.Fprintf(os.Stdout, "%s\n", ew.Dump())
	}
}
