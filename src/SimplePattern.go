package main

import (
	"strings"
)

type SimplePattern struct {
	tName *Word
}

func NewSimplePattern(w *Word) *SimplePattern {
	return &SimplePattern{w}
}

func (p *SimplePattern) Dump() string {
	return p.tName.Dump()
}

func (p *SimplePattern) Context() Context {
	return p.tName.Context()
}

func (p *SimplePattern) CalcDistance(arg Value) []int {
	t := arg.Type()

	d := t.CalcNameDistance(p.tName)

	if d < 0 {
		return nil
	} else {
		return []int{d}
	}
}

func (p *SimplePattern) Destructure(arg Value, scope *FuncScope, ew ErrorWriter) *FuncScope {
	if p.CalcDistance(arg) == nil {
		ew.Add(arg.Context().Error("unable to match pattern \"" + p.Dump() + "\""))
	}

	return scope
}

func (p *SimplePattern) CheckTypeNames(scope Scope, ew ErrorWriter) {
	name := p.tName.Value()

	if name == "Int" || name == "Any" || name == "String" || name == "Float" || name == "IO" || strings.HasPrefix(name, "\\") || name == "[]" || name == "{}" {
		return
	}

	if len(scope.CollectFunctions(name)) == 0 {
		ew.Add(p.Context().Error("\"" + name + "\" undefined"))
	}
}

func (p *SimplePattern) ListTypes() []string {
	return []string{p.tName.Value()}
}
