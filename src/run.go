package main

import (
	"fmt"
	"os"
)

// v might've been reduced by a previous call to Eval
func Run(v Value, scope Scope) {
	ew := NewErrorWriter()

	v = v.Eval(scope, ew)

	if v != nil {
		if !IsIO(v) {
			ew.Add(v.Context().Error("expected IO, got " + v.Type().Dump()))
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
