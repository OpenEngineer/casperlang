package main

import "strings"

type TuplePattern struct {
	TokenData
	items []Pattern
}

func NewTuplePattern(items []Pattern, ctx Context) *TuplePattern {
	return &TuplePattern{newTokenData(ctx), items}
}

func (t *TuplePattern) Name() string {
	return "[]"
}

func (t *TuplePattern) Dump() string {
	var b strings.Builder

	b.WriteString("[")
	for i, item := range t.items {
		b.WriteString(item.Dump())

		if i < len(t.items)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]")

	return b.String()
}

func (p *TuplePattern) ListTypes() []string {
	lst := []string{}

	for _, item := range p.items {
		lst = append(lst, item.ListTypes()...)
	}

	return lst
}

func (p *TuplePattern) ListNames() []*Word {
	lst := []*Word{}

	for _, item := range p.items {
		lst = append(lst, item.ListNames()...)
	}

	return lst
}

func (p *TuplePattern) ListVars() []*Variable {
	lst := []*Variable{}

	for _, item := range p.items {
		lst = append(lst, item.ListVars()...)
	}

	return lst
}

func (p *TuplePattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	items := []Pattern{}
	for _, item_ := range p.items {
		item := item_.Link(scope, ew)
		items = append(items, item)
	}

	return &TuplePattern{newTokenData(p.Context()), items}
}

func (p *TuplePattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, func(tn string) bool {
		return tn == "[]"
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil)
	}

	distance := []int{len(virt.Constructors())}

	if IsAll(concrete) {
		return NewDestructured(concrete, distance)
	}

	lst := AssertList(concrete)

	if lst.Len() != len(p.items) {
		return NewDestructured(concrete, nil)
	}

	lstItems := lst.Items()

	for i, pat := range p.items {
		d := pat.Destructure(lstItems[i], ew)

		if d.Failed() {
			return NewDestructured(
				NewList(lstItems, arg.Context()).SetConstructors(concrete.Constructors()),
				nil,
			)
		}

		distance = append(distance, d.distance...)
		lstItems[i] = d.arg
	}

	concrete = NewList(lstItems, arg.Context()).SetConstructors(concrete.Constructors())

	return NewDestructured(concrete, distance)
}
