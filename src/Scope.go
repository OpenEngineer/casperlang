package main

type Scope interface {
	// XXX: is still function really necessary?
	Parent() Scope

	// return nil if undefined
	Dispatch(name *Word, args []Value, ew ErrorWriter) Func // also returns the scope in which the value was defined

	CollectFunctions(name string) []DispatchableFunc
}
