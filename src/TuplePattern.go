package main

import "strings"

type TuplePattern struct {
	TokenData
	items []Pattern
}

func NewTuplePattern(items []Pattern, ctx Context) *TuplePattern {
	return &TuplePattern{newTokenData(ctx), items}
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

func ParseTuplePattern(gr *Brackets, ew ErrorWriter) *TuplePattern {
	items := []Pattern{}

	for _, ts := range gr.values {
		p := ParsePattern(ts, ew)
		if p != nil {
			items = append(items, p)
		}
	}

	return NewTuplePattern(items, gr.Context())
}

func (p *TuplePattern) CalcDistance(arg Value) []int {
	return arg.Type().CalcTupleDistance(p.items)
}

func (p *TuplePattern) Destructure(arg Value, scope *FuncScope, ew ErrorWriter) *FuncScope {
	t := arg.Type()

	d := t.CalcTupleDistance(p.items)
	if d == nil {
		ew.Add(arg.Context().Error("unable to match pattern \"" + p.Dump() + "\""))
		return scope
	}

	d0 := d[0]
	for d0 > 0 {
		t = t.Parent()
		d0 -= 1
	}

	return t.DestructureTuple(p.items, scope, ew)
}

func (p *TuplePattern) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, item := range p.items {
		item.CheckTypeNames(scope, ew)
	}
}

func (p *TuplePattern) ListTypes() []string {
	lst := []string{}

	for _, item := range p.items {
		lst = append(lst, item.ListTypes()...)
	}

	return lst
}
