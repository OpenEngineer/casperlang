package main

type Func interface {
	Value

	Name() string // empty for anonymous functions
	NumArgs() int

	// return a wrapped value of the call, which allows changing some underlying variables through destructure call
	// destructure doesn't necessarily need to be called (in which case there is no performance penalty)
	Call(args []Value, argsScope Scope, ctx Context, ew ErrorWriter) Value
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
