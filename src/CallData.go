package main

import (
	"strings"
)

type CallData struct {
	ValueData
	args []Value
}

func newCallData(args []Value, ctx Context) CallData {
	for _, arg := range args {
		if arg == nil {
			panic("arg is nil")
		}
	}

	return CallData{newValueData(ctx), args}
}

func (v *CallData) NumArgs() int {
	return len(v.args)
}

func (t *CallData) dump(first string) string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(first)

	for _, arg := range t.args {
		b.WriteString(" ")
		if arg == nil {
			b.WriteString("<nil-indicating-error>")
		} else {
			b.WriteString(arg.Dump())
		}
	}

	b.WriteString(")")

	return b.String()
}

func (c *CallData) linkArgs(scope Scope, ew ErrorWriter) []Value {
	args := []Value{}

	for _, arg_ := range c.args {
		arg := arg_.Link(scope, ew)
		if arg == nil {
			return []Value{}
		}
		args = append(args, arg)
	}

	return args
}

func (c *CallData) setConstructors(cs []Call) CallData {
	return CallData{ValueData{newTokenData(c.Context()), cs}, c.args}
}

func (c *CallData) subArgVars(stack *Stack) CallData {
	args := []Value{}

	for _, arg_ := range c.args {
		arg := arg_.SubVars(stack)
		args = append(args, arg)
	}

	return CallData{ValueData{newTokenData(c.Context()), c.constructors}, args}
}

func (c *CallData) Args() []Value {
	// make a copy
	res := make([]Value, len(c.args))

	for i, arg := range c.args {
		res[i] = arg
	}

	return res
}
