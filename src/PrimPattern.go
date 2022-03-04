package main

type PrimPattern struct {
	name string
	ctx  Context
}

func NewPrimPattern(name string, ctx Context) *PrimPattern {
	return &PrimPattern{name, ctx}
}

func (p *PrimPattern) Dump() string {
	return p.name
}

func (p *PrimPattern) DumpPretty() string {
	return p.Dump()
}

func (p *PrimPattern) Context() Context {
	return p.ctx
}

func (p *PrimPattern) ListTypes() []string {
	return []string{}
}

func (p *PrimPattern) ListNames() []*Word {
	return []*Word{}
}

func (p *PrimPattern) ListVars() []*Variable {
	return []*Variable{}
}

func (p *PrimPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	return p
}

func (p *PrimPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, func(tn string) bool {
		return tn == p.name
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil)
	}

	distance := []int{len(virt.Constructors())}

	return NewDestructured(concrete, distance)
}
