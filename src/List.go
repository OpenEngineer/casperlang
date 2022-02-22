package main

import (
	"fmt"
	"strings"
)

// XXX: do we also need a linked list type in case of concatenations?
type List struct {
	ValueData
	length int // explicit length in case append mess the items length
	items  []Value
}

func NewList(items []Value, ctx Context) *List {
	t := NewListType(items, ctx)
	return &List{newValueData(t, ctx), len(items), items}
}

// tuple is a literal list, can only be instantiated by parser
func NewTuple(items []Value, ctx Context) *List {
	return &List{newValueData(nil, ctx), len(items), items} // type is filled upon first eval
}

func IsList(t Token) bool {
	_, ok := t.(*List)
	return ok
}

func AssertList(t_ Token) *List {
	t, ok := t_.(*List)

	if ok {
		return t
	} else {
		panic("expected *List")
	}
}

func MergeLists(a_ Value, b_ Value, ctx Context) *List {
	a := AssertList(a_)
	b := AssertList(b_)

	items := append(a.items, b.items...)

	return NewList(items, ctx)
}

func (v *List) Update(type_ Type, ctx Context) Value {
	return &List{newValueData(type_, ctx), v.length, v.items}
}

func (t *List) Dump() string {
	var b strings.Builder

	b.WriteString("[")

	for i, item := range t.items {
		b.WriteString(item.Dump())

		if i < len(t.items)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String()
}

func ParseList(gr *Brackets, ew ErrorWriter) *List {
	if DEBUG_PARSER {
		fmt.Printf("ParseList IN: %s\n", gr.Dump())
	}

	items := []Value{}

	for _, field := range gr.values {
		item := ParseExpr(field, ew)
		items = append(items, item)
	}

	return NewTuple(items, gr.Context())
}

func (v *List) Eval(scope Scope, ew ErrorWriter) Value {
	items := []Value{}

	for _, item_ := range v.items {
		item := item_.Eval(scope, ew)
		if ew.Empty() {
			items = append(items, item)
		}
	}

	t := v.Type()
	if t == nil {
		t = NewTupleType(items, v.Context())
	}

	return &List{newValueData(t, v.Context()), len(items), items}
}

func (v *List) Get(i int, ctx Context) Value {
	if i < 0 {
		i += v.length
	}

	if i < 0 || i >= v.length {
		return NewNothingValue(ctx)
	}

	return NewJustValue(v.items[i], ctx)
}

func (v *List) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, item := range v.items {
		item.CheckTypeNames(scope, ew)
	}
}

func (v *List) Len() int {
	return v.length
}

func (v *List) Items() []Value {
	// a copy
	res := make([]Value, v.length)
	for i := 0; i < v.length; i++ {
		res[i] = v.items[i]
	}

	return res
}

// can't sort on list directly because it's supposed to immutable
type ListSorter struct {
	ctx                  Context
	insufficientTypeInfo bool
	comp                 Func
	items                []Value
	scope                Scope // used to dispatch '<'
	ew                   ErrorWriter
}

func NewListSorter(lst *List, scope Scope, comp Func, ew ErrorWriter, ctx Context) *ListSorter {
	return &ListSorter{ctx, false, comp, lst.Items(), scope, ew}
}

func (s *ListSorter) Len() int {
	return len(s.items)
}

func (s *ListSorter) Less(i, j int) bool {
	if !s.ew.Empty() || s.insufficientTypeInfo {
		return true
	}

	a := s.items[i]
	b := s.items[j]
	if a.Type() == nil || b.Type() == nil {
		s.insufficientTypeInfo = true
		return true
	}

	args := []Value{s.items[i], s.items[j]}

	res := s.comp.Call(args, s.scope, s.ctx, s.ew)
	if !s.ew.Empty() || res == nil {
		fmt.Println(s.ew.Dump())
		return true
	} else if IsDeferredError(res) {
		s.ew.Add(res.Context().Error("unable to dispatch comp (TODO: get rid of deferred error and do true lazy eval)"))
		return true
	}

	lt, ok := GetBoolValue(res)
	if ok {
		return lt
	} else {
		s.insufficientTypeInfo = true
		return true
	}
}

func (s *ListSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
