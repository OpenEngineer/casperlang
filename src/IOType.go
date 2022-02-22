package main

type IOType struct {
	data Value
}

func NewIOType(data Value) *IOType {
	return &IOType{data}
}

//func NewIOValue(data Value, ctx Context) Value {
//return NewValueData(NewIOType(data), ctx)
//}

func (t *IOType) Dump() string {
	if t == nil || t.data == nil {
		return "IO"
	} else {
		return "(IO " + t.data.Dump() + ")"
	}
}

func (t *IOType) Parent() Type {
	return &AnyType{}
}

func (t *IOType) CalcNameDistance(name *Word) int {
	if name.Value() == "IO" {
		return 0
	} else if name.Value() == "Any" {
		return 1
	} else {
		return -1
	}
}

func (t *IOType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	return nil
}

func (t *IOType) CalcTupleDistance(args []Pattern) []int {
	return nil
}

func (t *IOType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	return nil
}

func (t *IOType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *IOType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *IOType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func IsIOType(t Type) bool {
	_, ok := t.(*IOType)
	if ok {
		return ok
	} else if t == nil {
		return false
	} else {
		return IsIOType(t.Parent())
	}
}

func IsVoidIOType(t_ Type) bool {
	t, ok := t_.(*IOType)
	if ok {
		return t.data == nil
	} else if t_ == nil {
		return false
	} else {
		return IsVoidIOType(t_.Parent())
	}
}

func AssertIOType(t_ Type) *IOType {
	t, ok := t_.(*IOType)
	if ok {
		return t
	} else if t_ == nil {
		panic("expected *IOType")
	} else {
		return AssertIOType(t_.Parent())
	}
}
