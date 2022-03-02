package main

type Stack struct {
	parent *Stack
	vars   []*Variable
	data   []Value
}

func NewStack(parent *Stack) *Stack {
	return &Stack{parent, make([]*Variable, 0), make([]Value, 0)}
}

func (s *Stack) Get(v *Variable) Value {

	for i, v_ := range s.vars {
		if v_ == v {
			return s.data[i]
		}
	}

	if s.parent == nil {
		panic("var not found")
	} else {
		return s.parent.Get(v)
	}
}

func (s *Stack) Set(v *Variable, data Value) {
	for _, v_ := range s.vars {
		if v_ == v {
			panic("already set")
		}
	}

	s.vars = append(s.vars, v)
	s.data = append(s.data, data)
}

func (s *Stack) Extend(other *Stack) {
	for i, v := range other.vars {
		s.Set(v, other.data[i])
	}
}
