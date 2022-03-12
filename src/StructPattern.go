package main

import (
	"strings"
)

type StructPattern struct {
	TokenData
	keys []*String
	vals []Pattern
}

func NewStructPattern(keys []*String, vals []Pattern, ctx Context) *StructPattern {
	return &StructPattern{newTokenData(ctx), keys, vals}
}

func (t *StructPattern) Name() string {
	return "{}"
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

func (t *StructPattern) DumpPretty() string {
	var b strings.Builder

	b.WriteString("{")

	for i, k := range t.keys {
		b.WriteString(k.Dump())
		b.WriteString(":")
		b.WriteString(t.vals[i].DumpPretty())

		if i < len(t.keys)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String()
}

func (p *StructPattern) ListTypes() []string {
	lst := []string{}

	for _, val := range p.vals {
		lst = append(lst, val.ListTypes()...)
	}

	return lst
}

func (p *StructPattern) ListNames() []*Word {
	lst := []*Word{}

	for _, val := range p.vals {
		lst = append(lst, val.ListNames()...)
	}

	return lst
}

func (p *StructPattern) ListVars() []*Variable {
	lst := []*Variable{}

	for _, val := range p.vals {
		lst = append(lst, val.ListVars()...)
	}

	return lst
}

func (p *StructPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	vals := []Pattern{}

	for _, val_ := range p.vals {
		val := val_.Link(scope, ew)
		vals = append(vals, val)
	}

	return &StructPattern{newTokenData(p.Context()), p.keys, vals}
}

func (p *StructPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, func(tn string) bool {
		return tn == "{}"
	}, ew)

	if virt == nil {
		return NewDestructured(concrete, nil)
	}

	distance := []int{len(virt.Constructors())}

	if IsAll(concrete) {
		return NewDestructured(concrete, distance)
	}

	dict := AssertDict(concrete)

	if dict.Len() < len(p.keys) {
		return NewDestructured(concrete, nil)
	}

	dictVals := dict.Values() // should be a copy, so we can mutate
	dictKeys := dict.Keys()
	stack := NewStack()

	for i, pat := range p.vals {
		key := p.keys[i]

		found := false
		for j, check := range dictKeys {
			if check.Value() == key.Value() {
				found = true

				d := pat.Destructure(dictVals[j], ew)
				if d.Failed() {
					return NewDestructured(
						NewDict(dictKeys, dictVals, arg.Context()).SetConstructors(concrete.Constructors()),
						nil,
					)
				}

				distance = append(distance, d.distance...)
				stack.Extend(d.stack)
				dictVals[j] = d.arg

				break
			}
		}

		if !found {
			return NewDestructured(
				NewDict(dictKeys, dictVals, arg.Context()).SetConstructors(concrete.Constructors()),
				nil,
			)
		}
	}

	concrete = NewDict(dictKeys, dictVals, arg.Context()).SetConstructors(concrete.Constructors())

	return NewDestructuredS(concrete, distance, stack)
}
