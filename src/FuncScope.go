package main

import "strings"

type FuncScope struct {
	parent Scope

	// use a list so we can attach pattern contexts, and don't have the overhead of map creation
	argNames []*Word
	argVals  []Func
}

func (s *FuncScope) Parent() Scope {
	return s.parent
}

func (s *FuncScope) Dump() string {
	var b strings.Builder

	for i, name := range s.argNames {
		b.WriteString(name.Value())
		b.WriteString("=")
		b.WriteString(s.argVals[i].Dump())
		b.WriteString("\n")
	}

	return b.String()
}

func (s *FuncScope) find(name string) int {
	for i, argName := range s.argNames {
		if argName.Value() == name {
			return i
		}
	}

	return -1
}

func (s *FuncScope) add(name *Word, val Value) (*FuncScope, error) {
	if s.find(name.Value()) != -1 {
		return s, name.Context().Error("arg \"" + name.Value() + "\" already defined")
	} else {

		// wrap the val as a function if it is something simple
		fn, ok := val.(Func)
		if !ok {
			fn = NewNoArgAnonFunc(val, name.Context())
		}

		return &FuncScope{s.parent, append(s.argNames, name), append(s.argVals, fn)}, nil
	}
}

func (s *FuncScope) CollectFunctions(name string) []DispatchableFunc {
	return s.parent.CollectFunctions(name)
}

func (s *FuncScope) Dispatch(name *Word, args []Value, ew ErrorWriter) Func {
	// dispatchable functions are defined in parent scope or above, and are always shadowd by internal variables
	i := s.find(name.Value())

	if i == -1 {
		return s.parent.Dispatch(name, args, ew)
	} else {
		return s.argVals[i]
	}
}
