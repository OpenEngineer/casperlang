package main

import "strings"

type StructType struct {
	ctx    Context
	parent *DictType
	keys   []*String
	vals   []Value
}

func NewStructType(keys []*String, vals []Value, ctx Context) *StructType {
	return &StructType{ctx, NewDictType(vals, ctx), keys, vals}
}

func (t *StructType) Dump() string {
	var b strings.Builder

	b.WriteString("{")

	for i, key := range t.keys {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(key.Dump())
		b.WriteString(":")
		b.WriteString(t.vals[i].Type().Dump())
	}

	b.WriteString("}")

	return b.String()
}

func (t *StructType) Parent() Type {
	return t.parent
}

func (t *StructType) CalcNameDistance(name *Word) int {
	d := t.parent.CalcNameDistance(name)
	if d < 0 {
		return -1
	} else {
		return d + 1
	}
}

func (t *StructType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	d := t.parent.CalcConstructorDistance(name, args)

	if d == nil {
		return nil
	} else {
		d[0] += 1
		return d
	}
}

func (t *StructType) CalcTupleDistance(args []Pattern) []int {
	return nil
}

func (t *StructType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	if len(t.keys) < len(keys) {
		return nil
	}

	d := []int{0}

	for i, k := range keys {
		j := searchStrings(t.keys, k)

		if j == -1 {
			return nil
		}

		valD := vals[i].CalcDistance(t.vals[j])
		if valD == nil {
			return nil
		}

		d = append(d, valD...)
	}

	return d
}

// not yet implemented
func (t *StructType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *StructType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *StructType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	for i, k := range keys {
		j := searchStrings(t.keys, k)

		if j >= 0 {
			scope = vals[i].Destructure(t.vals[j], scope, ew)
		}
	}

	return scope
}
