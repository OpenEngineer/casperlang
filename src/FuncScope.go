package main

type FuncScope struct {
	parent Scope

	// use a list so we can attach pattern contexts, and don't have the overhead of map creation
	names []*Word
	vars  []*Variable
}

func NewFuncScope(parent Scope) *FuncScope {
	return &FuncScope{parent, []*Word{}, []*Variable{}}
}

func (s *FuncScope) find(name string) int {
	for i, check := range s.names {
		if check.Value() == name {
			return i
		}
	}

	return -1
}

func (s *FuncScope) Add(name *Word, var_ *Variable) {
	if s.find(name.Value()) != -1 {
		panic("arg already defined, should've been caught during parsing")
	}

	s.names = append(s.names, name)
	s.vars = append(s.vars, var_)
}

func (s *FuncScope) GetLocal(name string) *Variable {
	i := s.find(name)

	if i != -1 {
		return s.vars[i]
	} else {
		return s.parent.GetLocal(name)
	}
}

func (s *FuncScope) ListDispatchable(name string, nArgs int, ew ErrorWriter) []Func {
	return s.parent.ListDispatchable(name, nArgs, ew)
}
