package main

type AnyType struct {
}

func (t *AnyType) Dump() string {
	return "Any"
}

func (t *AnyType) Parent() Type {
	return nil
}

func (t *AnyType) CalcNameDistance(name *Word) int {
	if name.Value() == "Any" {
		return 0
	} else {
		return -1
	}
}

func (t *AnyType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	if len(args) == 0 && name.Value() == "Any" {
		return []int{0}
	} else {
		return nil
	}
}

func (t *AnyType) CalcTupleDistance(items []Pattern) []int {
	return nil
}

func (t *AnyType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

func (t *AnyType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *AnyType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *AnyType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
