package main

type ListType struct {
	ctx   Context
	items []Value // nil or empty means type is unknown
}

func NewListType(items []Value, ctx Context) *ListType {
	return &ListType{ctx, items}
}

func (t *ListType) Dump() string {
	if len(t.items) != 1 {
		return "([] <multi>)"
	} else {
		return "([] " + t.items[0].Type().Dump() + ")"
	}
}

func (t *ListType) Parent() Type {
	return &AnyType{}
}

func AssertListType(t_ Type) *ListType {
	for t_ != nil {
		if t, ok := t_.(*ListType); ok {
			return t
		} else {
			t_ = t_.Parent()
		}
	}

	panic("expected *ListType")
}

func MergeListTypes(a_ Type, b_ Type, ctx Context) *ListType {
	a := AssertListType(a_)
	b := AssertListType(b_)

	items := []Value{}

	for _, item := range a.items {
		items = append(items, item)
	}

	for _, item := range b.items {
		items = append(items, item)
	}

	return NewListType(items, ctx)
}

func (t *ListType) CalcNameDistance(name *Word) int {
	if name.Value() == "[]" {
		return 0
	} else if name.Value() == "Any" {
		return 1
	} else {
		return -1
	}
}

func (t *ListType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	if name.Value() == "[]" && len(args) == 1 {
		d := []int{0}

		worst := []int{}
		// each item must match
		for _, item := range t.items {
			itemD := args[0].CalcDistance(item)
			if itemD == nil {
				return nil
			}

			worst = WorstDistance(worst, itemD)
		}

		return append(d, worst...)
	} else {
		return nil
	}
}

func (t *ListType) CalcTupleDistance(args []Pattern) []int {
	return nil
}

func (t *ListType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

// not yet implemented
func (t *ListType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *ListType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *ListType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
