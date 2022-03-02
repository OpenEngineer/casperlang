package main

import (
	"strings"
)

type ConstructorPattern struct {
	TokenData
	name *Word
	args []Pattern
	fn   Func
}

func NewConstructorPattern(name *Word, args []Pattern, ctx Context) *ConstructorPattern {
	return &ConstructorPattern{newTokenData(ctx), name, args, nil}
}

func writeConstructorPattern(name *Word, args []Pattern) string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(name.Value())

	for _, arg := range args {
		b.WriteString(" ")
		b.WriteString(arg.Dump())
	}

	b.WriteString(")")

	return b.String()
}

func (t *ConstructorPattern) Name() string {
	return t.name.Value()
}

func (t *ConstructorPattern) NumArgs() int {
	return len(t.args)
}

func (t *ConstructorPattern) Dump() string {
	return writeConstructorPattern(t.name, t.args)
}

func (p *ConstructorPattern) ListTypes() []string {
	lst := []string{p.name.Value()}

	for _, arg := range p.args {
		lst = append(lst, arg.ListTypes()...)
	}

	return lst
}

func (p *ConstructorPattern) ListNames() []*Word {
	lst := []*Word{}

	for _, arg := range p.args {
		lst = append(lst, arg.ListNames()...)
	}

	return lst
}

func (p *ConstructorPattern) ListVars() []*Variable {
	lst := []*Variable{}

	for _, arg := range p.args {
		lst = append(lst, arg.ListVars()...)
	}

	return lst
}

func (p *ConstructorPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	name := p.Name()

	fns := scope.ListDispatchable(name, p.NumArgs(), ew)
	if len(fns) == 0 {
		ew.Add(p.name.Context().Error("\"" + name + "\" undefined"))
	} else if len(fns) > 1 {
		ew.Add(p.name.Context().Error("multiple definitions of \"" + name + "\""))
	}

	args := []Pattern{}
	for _, arg_ := range p.args {
		arg := arg_.Link(scope, ew)
		args = append(args, arg)
	}

	return &ConstructorPattern{newTokenData(p.Context()), p.name, args, fns[0]}
}

func (p *ConstructorPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, func(tn string) bool {
		return tn == p.Name()
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil)
	}

	distance := []int{len(virt.Constructors())}

	if IsAll(concrete) {
		return NewDestructured(concrete, distance)
	}

	call := AssertCall(concrete)

	if call.NumArgs() != len(p.args) {
		return NewDestructured(concrete, nil)
	}

	callArgs := call.Args()
	stack := NewStack()

	for i, pat := range p.args {
		d := pat.Destructure(callArgs[i], ew)
		if d.Failed() {

			if call.Name() == p.Name() {
				return NewDestructured(
					NewDisCall([]Func{p.fn}, callArgs, arg.Context()).SetConstructors(concrete.Constructors()),
					nil,
				)
			} else {
				return NewDestructured(concrete, nil)
			}
		}

		distance = append(distance, d.distance...)
		stack.Extend(d.stack)
		callArgs[i] = d.arg
	}

	if call.Name() == p.Name() {
		return NewDestructuredS(
			NewDisCall([]Func{p.fn}, callArgs, arg.Context()).SetConstructors(concrete.Constructors()),
			distance,
			stack,
		)
	} else {
		return NewDestructuredS(concrete, distance, stack)
	}
}
