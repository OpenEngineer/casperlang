package main

type Pattern interface {
	Token

	ListTypes() []string

	ListNames() []*Word

	ListVars() []*Variable

	// mutates FuncScope
	Link(scope *FuncScope, ew ErrorWriter) Pattern

	Destructure(arg Value, stack *Stack, ew ErrorWriter) *Destructured
}
