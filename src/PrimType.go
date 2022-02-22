package main

import "strconv"

type PrimType struct {
	name string
}

func NewPrimType(name string) *PrimType {
	return &PrimType{name}
}

func NewStringType() *PrimType {
	return NewPrimType("String")
}

func NewIntType() *PrimType {
	return NewPrimType("Int")
}

func NewFloatType() *PrimType {
	return NewPrimType("Float")
}

func NewFuncType(nArgs int) *PrimType {
	return NewPrimType("\\" + strconv.Itoa(nArgs))
}

func (t *PrimType) Dump() string {
	return t.name
}

func (t *PrimType) Parent() Type {
	return &AnyType{}
}

func (t *PrimType) CalcNameDistance(name *Word) int {
	if name.Value() == t.name {
		return 0
	} else if name.Value() == "Any" {
		return 1
	} else {
		return -1
	}
}

func (t *PrimType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	return nil
}

func (t *PrimType) CalcTupleDistance(args []Pattern) []int {
	return nil
}

func (t *PrimType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

func (t *PrimType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *PrimType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *PrimType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
