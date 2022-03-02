package main

import (
	"strings"
)

type Stack struct {
	vars []*Variable
	data []Value
}

func NewStack() *Stack {
	return &Stack{make([]*Variable, 0), make([]Value, 0)}
}

func (s *Stack) Get(v *Variable) Value {
	for i, v_ := range s.vars {
		if v_ == v {
			return s.data[i]
		}
	}

	return nil
}

func (s *Stack) Set(v *Variable, data Value) {
	for _, v_ := range s.vars {
		if v_ == v {
			panic("already set")
		}
	}

	if _, ok := data.(*VarRef); ok {
		panic("var data can't be another ref")
	}

	s.vars = append(s.vars, v)
	s.data = append(s.data, data)
}

func (s *Stack) Extend(other *Stack) {
	for i, v := range other.vars {
		s.Set(v, other.data[i])
	}
}

func (s *Stack) Dump() string {
	var b strings.Builder

	for i, v_ := range s.vars {
		b.WriteString("    ")
		b.WriteString(v_.Name())
		b.WriteString("=")
		b.WriteString(s.data[i].Dump())
		b.WriteString("\n")
	}

	return b.String()
}
