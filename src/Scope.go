package main

type Scope interface {
	GetLocal(name string) *Variable

	ListDispatchable(name string, nArgs int, ew ErrorWriter) []Func
}
