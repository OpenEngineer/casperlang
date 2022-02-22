package main

import (
	"fmt"
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

func (h *FuncHeader) CalcDistance(args []Value) []int {
	if len(args) != h.NumArgs() {
		return nil
	}

	d := []int{}

	for i, arg := range args {
		p := h.args[i]
		subD := p.CalcDistance(arg)
		if subD == nil {
			return nil
		}

		d = append(d, subD...)
	}

	return d
}

func (h *FuncHeader) DestructureArgs(scope Scope, args []Value, ctx Context, ew ErrorWriter) *FuncScope {
	if len(args) != h.NumArgs() {
		ew.Add(ctx.Error(fmt.Sprintf("expected %d args, got %d", h.NumArgs(), len(args))))
		return nil
	}

	subScope := &FuncScope{scope, []*Word{}, []Func{}}

	for i, arg := range args {
		pat := h.args[i]

		subScope = pat.Destructure(arg, subScope, ew)
	}

	return subScope
}

func (h *FuncHeader) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, arg := range h.args {
		arg.CheckTypeNames(scope, ew)
	}
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
