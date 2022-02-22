package main

import (
	"strings"
)

type ConstructorPattern struct {
	TokenData
	name *Word
	args []Pattern
}

func NewConstructorPattern(name *Word, args []Pattern, ctx Context) *ConstructorPattern {
	return &ConstructorPattern{newTokenData(ctx), name, args}
}

func writeConstructorPattern(name *Word, args []Pattern) string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(name.Value())

	for _, arg := range args {
		b.WriteString(" ")
		b.WriteString(arg.Dump())
	}

	b.WriteString(")")

	return b.String()
}

func (t *ConstructorPattern) Dump() string {
	return writeConstructorPattern(t.name, t.args)
}

func ParseConstructorPattern(p *Parens, ew ErrorWriter) *ConstructorPattern {
	ctx := p.Context()

	ts := p.content

	t := ts[0]

	var name *Word
	switch {
	case IsWord(t):
		name = AssertWord(t)
	case IsEmptyBraces(t):
		name = NewWord("{}", t.Context())
	case IsEmptyBrackets(t):
		name = NewWord("[]", t.Context())
	default:
		ew.Add(t.Context().Error("invalid constructor pattern syntax"))
		return nil
	}

	if name.Value() == "Int" ||
		name.Value() == "Float" ||
		name.Value() == "String" ||
		name.Value() == "Any" ||
		name.Value() == "IO" {
		ew.Add(name.Context().Error("can't apply constructor pattern to \"" + name.Value() + "\""))
		return nil
	}

	args := ParsePatterns(ts[1:], ew)

	return NewConstructorPattern(name, args, ctx)
}

func (p *ConstructorPattern) CalcDistance(arg Value) []int {
	return arg.Type().CalcConstructorDistance(p.name, p.args)
}

func (p *ConstructorPattern) Destructure(arg Value, scope *FuncScope, ew ErrorWriter) *FuncScope {
	t := arg.Type()

	d := t.CalcConstructorDistance(p.name, p.args)
	if d == nil {
		ew.Add(arg.Context().Error("unable to match pattern \"" + p.Dump() + "\""))
		return scope
	}

	d0 := d[0]
	for d0 > 0 {
		t = t.Parent()
		d0 -= 1
	}

	return t.DestructureConstructor(p.name, p.args, scope, ew)
}

func (p *ConstructorPattern) CheckTypeNames(scope Scope, ew ErrorWriter) {
	name := p.name.Value()

	if name == "[]" || name == "{}" {
		return
	}

	if len(scope.CollectFunctions(name)) == 0 {
		ew.Add(p.name.Context().Error("\"" + name + "\" undefined"))
	}

	for _, arg := range p.args {
		arg.CheckTypeNames(scope, ew)
	}
}

func (p *ConstructorPattern) ListTypes() []string {
	lst := []string{p.name.Value()}

	for _, arg := range p.args {
		lst = append(lst, arg.ListTypes()...)
	}

	return lst
}
