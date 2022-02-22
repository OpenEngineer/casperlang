package main

type NamedPattern struct {
	TokenData
	name    *Word
	pattern Pattern
}

func NewNamedPattern(name *Word, pattern Pattern, ctx Context) *NamedPattern {
	if pattern == nil {
		panic("pattern can't be nil")
	}

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

func (t *NamedPattern) Dump() string {
	return t.name.Value() + "::" + t.pattern.Dump()
}

func (p *NamedPattern) CalcDistance(arg Value) []int {
	return p.pattern.CalcDistance(arg)
}

func (p *NamedPattern) Destructure(arg Value, scope *FuncScope, ew ErrorWriter) *FuncScope {
	scope = p.pattern.Destructure(arg, scope, ew)

	scope, err := scope.add(p.name, arg)
	if err != nil {
		ew.Add(err)
	}

	return scope
}

func (p *NamedPattern) CheckTypeNames(scope Scope, ew ErrorWriter) {
	p.pattern.CheckTypeNames(scope, ew)
}

func (p *NamedPattern) ListTypes() []string {
	return p.pattern.ListTypes()
}
