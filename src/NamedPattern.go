package main

type NamedPattern struct {
	TokenData
	name    *Word
	pattern Pattern
}

func NewNamedPattern(name *Word, pattern Pattern, ctx Context) *NamedPattern {
	return &NamedPattern{newTokenData(ctx), name, pattern}
}

func IsNamedPattern(t Token) bool {
	_, ok := t.(*NamedPattern)
	return ok
}

func AssertNamedPattern(t_ Token) *NamedPattern {
	t, ok := t_.(*NamedPattern)

	if ok {
		return t
	} else {
		panic("expected *NamedPattern")
	}
}

func (t *NamedPattern) Name() string {
	return t.name.Value()
}

func (t *NamedPattern) Dump() string {
	return t.Name() + "::" + t.pattern.Dump()
}

func (p *NamedPattern) ListTypes() []string {
	return p.pattern.ListTypes()
}

func (p *NamedPattern) ListNames() []*Word {
	return append([]*Word{p.name}, p.pattern.ListNames()...)
}

func (p *NamedPattern) ListVars() []*Variable {
	panic("should've been converted to a VarPattern")
}

func (p *NamedPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	var_ := NewVariable(p.Name(), p.name.Context())

	scope.Add(p.name, var_)

	pattern := p.pattern.Link(scope, ew)

	return NewVarPattern(var_, pattern, p.Context())
}

func (p *NamedPattern) Destructure(arg Value, ew ErrorWriter) *Destructured {
	panic("should've been converted into a VarPattern")
}
