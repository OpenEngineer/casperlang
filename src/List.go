package main

import (
	"strings"
)

// XXX: do we also need a linked list type in case of concatenations?
type List struct {
	ValueData
	length int // explicit length in case append mess the items length
	items  []Value
}

func NewList(items []Value, ctx Context) *List {
	return &List{newValueData(ctx), len(items), items}
}

func NewEmptyList(ctx Context) *List {
	return NewList([]Value{}, ctx)
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

func (v *List) TypeName() string {
	return "[]"
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

func (v *List) SetConstructors(cs []Call) Value {
	return &List{ValueData{newTokenData(v.Context()), cs}, v.length, v.items}
}

func (v *List) Len() int {
	return v.length
}

func (v *List) Get(i int) Value {
	if i < 0 {
		i += v.length
	}

	if i < 0 || i >= v.length {
		return nil
	}

	return v.items[i]
}

func (v *List) Items() []Value {
	// a copy
	res := make([]Value, v.length)
	for i := 0; i < v.length; i++ {
		res[i] = v.items[i]
	}

	return res
}

func (v *List) Link(scope Scope, ew ErrorWriter) Value {
	items := []Value{}

	for _, item_ := range v.items {
		item := item_.Link(scope, ew)
		items = append(items, item)
	}

	return NewList(items, v.Context())
}

func (v *List) SubVars(stack *Stack) Value {
	items := []Value{}

	for _, item_ := range v.items {
		item := item_.SubVars(stack)
		items = append(items, item)
	}

	return NewList(items, v.Context())
}

func MergeLists(a_ Value, b_ Value, ctx Context) *List {
	a := AssertList(a_)
	b := AssertList(b_)

	items := append(a.items, b.items...)

	return NewList(items, ctx)
}
