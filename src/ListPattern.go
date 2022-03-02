package main

type ListPattern struct {
	TokenData
	inner      Pattern
	innerNames []*Word
	innerVars  []*Variable
}

func NewListPattern(inner Pattern, ctx Context) *ListPattern {
	innerNames := inner.ListNames()
	return &ListPattern{newTokenData(ctx), inner, innerNames, nil}
}

func (p *ListPattern) Dump() string {
	return "([] " + p.inner.Dump() + ")"
}

func (p *ListPattern) ListTypes() []string {
	return p.inner.ListTypes()
}

func (p *ListPattern) ListNames() []*Word {
	return p.inner.ListNames()
}

func (p *ListPattern) ListVars() []*Variable {
	return p.inner.ListVars()
}

func (p *ListPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	inner := p.inner.Link(scope, ew)

	return &ListPattern{newTokenData(p.Context()), inner, p.innerNames, inner.ListVars()}
}

func (p *ListPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
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

	if lst.Len() == 0 {
		// empty list should also set the inner scope
		// add another empty list for every name

		// distance is based on All type, which matches anything with distance 0
		dAll := p.inner.Destructure(NewAll(p.Context()), ew)
		if dAll.Failed() {
			return NewDestructured(lst, nil)
		}

		d := NewDestructured(lst, append(distance, dAll.distance...))

		for _, innerVar := range p.innerVars {
			d.AddVar(innerVar, NewEmptyList(arg.Context()))
		}

		return d
	} else {

		// apply the inner pattern to every list items, then merge
		// the worst non-fail distance is the final distance

		ds := []*Destructured{}
		items := lst.Items()

		for i, item := range items {
			d := p.inner.Destructure(item, ew)
			if d.Failed() {
				return NewDestructured(
					NewList(items, arg.Context()).SetConstructors(concrete.Constructors()),
					nil,
				)
			}

			ds = append(ds, d)
			items[i] = d.arg
		}

		dFinal := NewDestructured(
			NewList(items, arg.Context()).SetConstructors(concrete.Constructors()),
			WorstDistance(ds),
		)

		// create a bunch of lists
		for i, innerVar := range p.innerVars {
			innerItems := []Value{}

			for _, d := range ds {
				if d.stack.vars[i] != innerVar {
					panic("unexpected")
				}

				innerItems = append(innerItems, d.stack.data[i])
			}

			dFinal.AddVar(innerVar, NewList(innerItems, arg.Context()))
		}

		return dFinal
	}
}
