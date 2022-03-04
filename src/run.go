package main

import (
	"fmt"
	"os"
)

func EvalEager(v Value, ew ErrorWriter) Value {
Outer:
	for v != nil {
		switch v_ := v.(type) {
		case Call:
			v = v_.Eval(ew)
		case *List:
			items := []Value{}
			for _, item := range v_.Items() {
				items = append(items, EvalEager(item, ew))
				if !ew.Empty() {
					return nil
				}
			}

			return NewList(items, v_.Context())
		case *Dict:
			vals := []Value{}
			for _, val := range v_.Values() {
				vals = append(vals, EvalEager(val, ew))
				if !ew.Empty() {
					return nil
				}
			}

			return NewDict(v_.Keys(), vals, v_.Context())
		default:
			break Outer
		}
	}

	return v
}

// v might've been reduced by a previous call to Eval
func Run(v Value) {
	ew := NewErrorWriter()

	IO_CONTEXT = NewDefaultIOContext()

	v = EvalEager(v, ew)

	if v != nil && ew.Empty() {
		if !IsIO(v) {
			ew.Add(v.Context().Error("expected IO, got " + v.Dump()))
		} else {
			io := AssertIO(v)

			v = io.Run(IO_CONTEXT)

			if v != nil {
				ew.Add(v.Context().Error("unused IO result"))
			}
		}
	}

	if !ew.Empty() {
		fmt.Fprintf(os.Stdout, "%s\n", ew.Dump())
	}
}
