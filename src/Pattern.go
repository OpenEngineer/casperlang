package main

type Pattern interface {
	Token

	// the caller isn't interested in the head variables, so these are not printed here
	DumpPretty() string

	ListTypes() []string

	ListNames() []*Word

	ListVars() []*Variable

	// mutates FuncScope
	Link(scope *FuncScope, ew ErrorWriter) Pattern

	Destructure(arg Value, ew ErrorWriter) *Destructured
}
