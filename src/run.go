package main

import (
	"fmt"
	"os"
)

// v might've been reduced by a previous call to Eval
func Run(v Value) {
	ew := NewErrorWriter()

	for {
		c, ok := v.(Call)
		if !ok || !ew.Empty() {
			break
		}

		v = c.Eval(ew)
	}

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
