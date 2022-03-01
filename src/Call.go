package main

type Call interface {
	Value

	Name() string
	NumArgs() int
	Args() []Value

	Eval(ew ErrorWriter) Value
}

func IsCall(t Token) bool {
	_, ok := t.(Call)
	return ok
}

func AssertCall(t_ Token) Call {
	t, ok := t_.(Call)
	if ok {
		return t
	} else {
		panic("expected Call")
	}
}
