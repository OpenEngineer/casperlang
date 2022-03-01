package main

import (
	"strings"
)

type CallData struct {
	ValueData
	args  []Value
	cache Value
}

func newCallData(args []Value, ctx Context) CallData {
	return CallData{newValueData(ctx), args, nil}
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
			panic("arg is nil")
		}
		b.WriteString(arg.Dump())
	}

	b.WriteString(")")

	return b.String()
}

func (c *CallData) linkArgs(scope Scope, ew ErrorWriter) []Value {
	args := []Value{}

	for _, arg_ := range c.args {
		arg := arg_.Link(scope, ew)
		args = append(args, arg)
	}

	return args
}

func (c *CallData) setConstructors(cs []Call) CallData {
	return CallData{ValueData{newTokenData(c.Context()), cs}, c.args, c.cache}
}

func (c *CallData) SetCache(v Value) {
	// cache can be overwritten with better values though
	c.cache = v
}

func (c *CallData) Args() []Value {
	// make a copy
	res := make([]Value, len(c.args))

	for i, arg := range c.args {
		res[i] = arg
	}

	return res
}

func SetCallDataCache(cd_ Value, v Value) {
	// XXX: how should we cache the results?
	//if cd, ok := cd_.(*CallData); ok {
	//cd.SetCache(v)
	//}
}
