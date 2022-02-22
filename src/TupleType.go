package main

import (
	"strings"
)

type TupleType struct {
	ctx    Context
	parent *ListType
	items  []Value
}

func NewTupleType(items []Value, ctx Context) *TupleType {
	for _, item := range items {
		if item.Type() == nil {
			panic("item type inside tuple can't be nil")
		}
	}

	return &TupleType{ctx, NewListType(items, ctx), items}
}

func (t *TupleType) Dump() string {
	var b strings.Builder

	b.WriteString("[")

	for i, item := range t.items {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(item.Type().Dump())
	}

	b.WriteString("]")

	return b.String()
}

func (t *TupleType) Parent() Type {
	return t.parent
}

func (t *TupleType) CalcNameDistance(name *Word) int {
	d := t.parent.CalcNameDistance(name)
	if d < 0 {
		return -1
	} else {
		return d + 1
	}
}

func (t *TupleType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	d := t.parent.CalcConstructorDistance(name, args)

	if d == nil {
		return nil
	} else {
		d[0] += 1
		return d
	}
}

func (t *TupleType) CalcTupleDistance(args []Pattern) []int {
	if len(args) != len(t.items) {
		return nil
	}

	d := []int{0}

	for i, arg := range args {
		item := t.items[i]
		if item.Type() == nil {
			panic("item type in tuple type can't be nil")
		}

		argD := arg.CalcDistance(t.items[i])
		if argD == nil {
			return nil
		}

		d = append(d, argD...)
	}

	return d
}

func (t *TupleType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

func (t *TupleType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *TupleType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	for i, arg := range args {
		scope = arg.Destructure(t.items[i], scope, ew)
	}

	return scope
}

func (t *TupleType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
