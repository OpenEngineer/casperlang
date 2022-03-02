package main

type DictPattern struct {
	TokenData
	inner      Pattern
	innerNames []*Word
	innerVars  []*Variable
}

func NewDictPattern(inner Pattern, ctx Context) *DictPattern {
	innerNames := inner.ListNames()
	return &DictPattern{newTokenData(ctx), inner, innerNames, nil}
}

func (p *DictPattern) Dump() string {
	return "({} " + p.inner.Dump() + ")"
}

func (p *DictPattern) ListTypes() []string {
	return p.inner.ListTypes()
}

func (p *DictPattern) ListNames() []*Word {
	return p.inner.ListNames()
}

func (p *DictPattern) ListVars() []*Variable {
	return p.inner.ListVars()
}

func (p *DictPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	inner := p.inner.Link(scope, ew)

	return &DictPattern{newTokenData(p.Context()), inner, p.innerNames, inner.ListVars()}
}

func (p *DictPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	if arg == nil {
		panic("arg can't be nil")
	}

	concrete, virt := EvalUntil(arg, func(tn string) bool {
		return tn == "{}"
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil)
	}

	distance := []int{len(virt.Constructors())}

	if IsAll(concrete) {
		return NewDestructured(concrete, distance)
	}

	dict := AssertDict(arg)

	if dict.Len() == 0 {
		dAll := p.inner.Destructure(NewAll(p.Context()), ew)
		if dAll.Failed() {
			return NewDestructured(dict, nil)
		}

		d := NewDestructured(dict, append(distance, dAll.distance...))

		for _, innerVar := range p.innerVars {
			d.AddVar(innerVar, NewEmptyDict(arg.Context()))
		}

		return d
	} else {
		ds := []*Destructured{}
		items := dict.Values()
		keys := dict.Keys()

		for i, item := range items {
			d := p.inner.Destructure(item, ew)
			if d.Failed() {
				return NewDestructured(NewDict(keys, items, arg.Context()).SetConstructors(concrete.Constructors()), nil)
			}

			ds = append(ds, d)
			items[i] = d.arg
		}

		dFinal := NewDestructured(
			NewDict(keys, items, arg.Context()).SetConstructors(concrete.Constructors()),
			WorstDistance(ds),
		)

		for i, innerVar := range p.innerVars {
			innerItems := []Value{}

			for _, d := range ds {
				if d.vars[i] != innerVar {
					panic("unexpected")
				}

				innerItems = append(innerItems, d.data[i])
			}

			if len(keys) != len(innerItems) {
				panic("unexpected")
			}

			dFinal.AddVar(innerVar, NewDict(keys, innerItems, arg.Context()))
		}

		return dFinal
	}
}
