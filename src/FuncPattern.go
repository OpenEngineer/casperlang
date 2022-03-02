package main

import (
	"strconv"
)

type FuncPattern struct {
	nArgs int
	ctx   Context
}

func NewFuncPattern(nArgs int, ctx Context) *FuncPattern {
	return &FuncPattern{nArgs, ctx}
}

func (p *FuncPattern) Dump() string {
	return "\\" + strconv.Itoa(p.nArgs)
}

func (p *FuncPattern) Context() Context {
	return p.ctx
}

func (p *FuncPattern) ListTypes() []string {
	return []string{}
}

func (p *FuncPattern) ListNames() []*Word {
	return []*Word{}
}

func (p *FuncPattern) ListVars() []*Variable {
	return []*Variable{}
}

func (p *FuncPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	return p
}

func (p *FuncPattern) Destructure(arg Value, stack *Stack, ew ErrorWriter) *Destructured {
	concrete, virt := EvalUntil(arg, stack, func(tn string) bool {
		return tn == p.Dump()
	}, ew)

	if concrete == nil {
		return NewDestructured(arg, nil, nil)
	}

	distance := []int{len(virt.Constructors())}

	return NewDestructured(concrete, distance, stack)
}
