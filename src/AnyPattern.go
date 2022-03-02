package main

import (
	"math"
)

type AnyPattern struct {
	ctx Context
}

func NewAnyPattern(ctx Context) *AnyPattern {
	return &AnyPattern{ctx}
}

func (p *AnyPattern) Dump() string {
	return "Any"
}

func (p *AnyPattern) Context() Context {
	return p.ctx
}

func (p *AnyPattern) ListTypes() []string {
	return []string{}
}

func (p *AnyPattern) ListNames() []*Word {
	return []*Word{}
}

func (p *AnyPattern) ListVars() []*Variable {
	return []*Variable{}
}

func (p *AnyPattern) Link(scope *FuncScope, ew ErrorWriter) Pattern {
	return p
}

func (p *AnyPattern) Destructure(arg Value, stack *Stack, ew ErrorWriter) *Destructured {
	return NewDestructured(arg, []int{math.MaxInt32}, stack)
}
