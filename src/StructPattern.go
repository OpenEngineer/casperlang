package main

import "strings"

type StructPattern struct {
	TokenData
	keys []*String
	vals []Pattern
}

func NewStructPattern(keys []*String, vals []Pattern, ctx Context) *StructPattern {
	return &StructPattern{newTokenData(ctx), keys, vals}
}

func (t *StructPattern) Dump() string {
	var b strings.Builder

	b.WriteString("{")

	for i, k := range t.keys {
		b.WriteString(k.Dump())
		b.WriteString(":")
		b.WriteString(t.vals[i].Dump())

		if i < len(t.keys)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String()
}

func ParseStructPattern(gr *Braces, ew ErrorWriter) *StructPattern {
	keys := []*String{}
	vals := []Pattern{}

	for i, ts := range gr.vals {
		v := ParsePattern(ts, ew)
		if v != nil {
			vals = append(vals, v)
			keys = append(keys, gr.keys[i])
		}
	}

	return NewStructPattern(keys, vals, gr.Context())
}

func (p *StructPattern) CalcDistance(arg Value) []int {
	return arg.Type().CalcStructDistance(p.keys, p.vals)
}

func (p *StructPattern) Destructure(arg Value, scope *FuncScope, ew ErrorWriter) *FuncScope {
	t := arg.Type()

	d := t.CalcStructDistance(p.keys, p.vals)
	if d == nil {
		ew.Add(arg.Context().Error("unable to match pattern \"" + p.Dump() + "\""))
		return scope
	}

	d0 := d[0]
	for d0 > 0 {
		t = t.Parent()
		d0 -= 1
	}

	return t.DestructureStruct(p.keys, p.vals, scope, ew)
}

func (p *StructPattern) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, val := range p.vals {
		val.CheckTypeNames(scope, ew)
	}
}

func (p *StructPattern) ListTypes() []string {
	lst := []string{}

	for _, val := range p.vals {
		lst = append(lst, val.ListTypes()...)
	}

	return lst
}
