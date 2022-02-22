package main

import "strings"

type Type interface {
	Dump() string

	Parent() Type

	// a negative number indicates no number
	CalcNameDistance(name *Word) int // name is the pattern

	// nil indicates no match
	CalcConstructorDistance(name *Word, args []Pattern) []int

	CalcTupleDistance(args []Pattern) []int

	CalcStructDistance(keys []*String, vals []Pattern) []int

	// ErrorWriter so duplicate names can be detected
	DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope

	DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope

	DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope
}

func ToTypes(vals []Value) []Type {
	ts := make([]Type, len(vals))

	for i, v := range vals {
		ts[i] = v.Type()
	}

	return ts
}

func DumpTypes(ts []Type) string {
	var b strings.Builder

	for i, t := range ts {
		if i > 0 {
			b.WriteString(" ")
		}

		b.WriteString(t.Dump())
	}

	return b.String()
}
