package main

import (
	"sort"
	"strings"
)

type FuncHeader struct {
	name *Word // anonymous functions will simple have name==nil
	args []Pattern
}

func NewFuncHeader(name *Word, args []Pattern) *FuncHeader {
	return &FuncHeader{name, args}
}

func (h *FuncHeader) IsAnon() bool {
	return h.name == nil
}

func (h *FuncHeader) Context() Context {
	if h.IsAnon() {
		return h.args[0].Context()
	} else {
		return h.name.Context()
	}
}

func (h *FuncHeader) IsConstructor() bool {
	if h.name != nil {
		if len(h.name.Value()) == 0 {
			panic("unexpected empty name")
		}

		return h.name.IsUpperCase()
	} else {
		return false
	}
}

func (h *FuncHeader) Name() string {
	if h.name == nil {
		return ""
	} else {
		return h.name.Value()
	}
}

func (h *FuncHeader) NumArgs() int {
	return len(h.args)
}

func (h *FuncHeader) Dump() string {
	var b strings.Builder

	if h.IsAnon() {
		b.WriteString("\\(")
	} else {
		b.WriteString(h.Name())
		b.WriteString(" ")
	}

	for i, arg := range h.args {
		if i > 0 {
			b.WriteString(" ")
		}

		b.WriteString(arg.Dump())
	}

	if h.IsAnon() {
		b.WriteString(")")
	}

	return b.String()
}

func (h *FuncHeader) DumpArgs() string {
	var b strings.Builder

	for _, arg := range h.args {
		b.WriteString(arg.Dump())
		b.WriteString(" ")
	}

	return b.String()
}

func (h *FuncHeader) Link(scope Scope, ew ErrorWriter) (*FuncHeader, *FuncScope) {
	fnScope := NewFuncScope(scope)

	args := []Pattern{}

	for _, arg_ := range h.args {
		arg := arg_.Link(fnScope, ew)
		args = append(args, arg)
	}

	return &FuncHeader{h.name, args}, fnScope
}

func (h *FuncHeader) Destructure(args []Value, ew ErrorWriter) *Dispatched {
	if len(args) != h.NumArgs() {
		panic("should've been caught higher up")
	}

	disp := NewDispatched(args, h.Context())

	for i, arg := range args {
		if arg == nil {
			panic("arg can't be nil")
		}

		des := h.args[i].Destructure(arg, ew)

		disp.UpdateArg(i, des)

		if disp.Failed() {
			break
		}
	}

	return disp
}

func sortUniqStrings(lst []string) []string {
	sort.Strings(lst)
	res := []string{}
	for i, x := range lst {
		if i == 0 {
			res = append(res, x)
		} else if x != lst[i-1] {
			res = append(res, x)
		}
	}

	return res
}

func (h *FuncHeader) ListTypes() []string {
	lst := []string{}

	for _, arg := range h.args {
		sub := arg.ListTypes()

		lst = append(lst, sub...)
	}

	return sortUniqStrings(lst)
}
