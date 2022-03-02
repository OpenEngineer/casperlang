package main

type TypePattern struct {
	tName *Word
}

func NewTypePattern(w *Word) *TypePattern {
	if w.Value() == "Any" {
		panic("use AnyPattern instead")
	}
	return &TypePattern{w}
}

func (p *TypePattern) TypeName() string {
	return p.tName.Value()
}

func (p *TypePattern) Dump() string {
	return p.tName.Dump()
}

func (p *TypePattern) Context() Context {
	return p.tName.Context()
}

func (p *TypePattern) ListTypes() []string {
	return []string{p.TypeName()}
}

func (p *TypePattern) ListNames() []*Word {
	return []*Word{}
}

func (p *TypePattern) ListVars() []*Variable {
	return []*Variable{}
}

func (p *TypePattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	name := p.TypeName()

	if len(scope.ListDispatchable(name, -1, ew)) == 0 {
		ew.Add(p.Context().Error("\"" + name + "\" undefined"))
	}

	return p
}

func (p *TypePattern) Destructure(arg Value, stack *Stack, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, stack, func(tn string) bool {
		return tn == p.TypeName()
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil, nil)
	}

	distance := []int{len(virt.Constructors())}

	return NewDestructured(concrete, distance, stack)
}
