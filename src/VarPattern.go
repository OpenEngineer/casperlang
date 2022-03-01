package main

type VarPattern struct {
	TokenData
	var_    *Variable
	pattern Pattern
}

func NewVarPattern(var_ *Variable, pattern Pattern, ctx Context) *VarPattern {
	return &VarPattern{newTokenData(ctx), var_, pattern}
}

func (t *VarPattern) Dump() string {
	return t.var_.Dump() + "::" + t.pattern.Dump()
}

func (p *VarPattern) ListTypes() []string {
	return p.pattern.ListTypes()
}

func (p *VarPattern) ListNames() []*Word {
	panic("should still be NamedPattern at this stage")
}

func (p *VarPattern) ListVars() []*Variable {
	return append([]*Variable{p.var_}, p.pattern.ListVars()...)
}

func (p *VarPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	return p
}

func (p *VarPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	d := p.pattern.Destructure(arg, ew)

	if !d.Failed() {
		d.AddVar(p.var_, d.arg)
	}

	return d
}
