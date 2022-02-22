package main

type DictType struct {
	ctx  Context
	vals []Value
}

func NewDictType(vals []Value, ctx Context) *DictType {
	return &DictType{ctx, vals}
}

func (t *DictType) Dump() string {
	if len(t.vals) != 1 {
		return "({} <multi>)"
	} else {
		return "({} " + t.vals[0].Type().Dump() + ")"
	}
}

func (t *DictType) Parent() Type {
	return &AnyType{}
}

func AssertDictType(t_ Type) *DictType {
	for t_ != nil {
		if t, ok := t_.(*DictType); ok {
			return t
		} else {
			t_ = t_.Parent()
		}
	}

	panic("expected *ListType")
}

func MergeDictTypes(a_ Type, b_ Type, ctx Context) *DictType {
	a := AssertDictType(a_)
	b := AssertDictType(b_)

	vals := []Value{}

	for _, aVal := range a.vals {
		vals = append(vals, aVal)
	}

	for _, bVal := range b.vals {
		vals = append(vals, bVal)
	}

	return NewDictType(vals, ctx)
}

func (t *DictType) CalcNameDistance(name *Word) int {
	if name.Value() == "{}" {
		return 0
	} else if name.Value() == "Any" {
		return 1
	} else {
		return -1
	}
}

func (t *DictType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	if name.Value() == "{}" && len(args) == 1 {
		d := []int{0}

		worst := []int{}
		// each item must match
		for _, val := range t.vals {
			valD := args[0].CalcDistance(val)
			if valD == nil {
				return nil
			}

			worst = WorstDistance(worst, valD)
		}

		return append(d, worst...)
	} else {
		return nil
	}
}

func (t *DictType) CalcTupleDistance(args []Pattern) []int {
	return nil
}

func (t *DictType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

// not yet implemented
func (t *DictType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *DictType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *DictType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
