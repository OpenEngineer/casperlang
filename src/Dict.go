package main

import (
	"strings"
)

type Dict struct {
	ValueData
	length int // explicit length, in case append messes up the slices
	keys   []*String
	vals   []Value
}

func NewDict(keys []*String, vals []Value, ctx Context) *Dict {
	return &Dict{newValueData(ctx), len(keys), keys, vals}
}

func NewEmptyDict(ctx Context) *Dict {
	return &Dict{newValueData(ctx), 0, []*String{}, []Value{}}
}

func IsDict(t Token) bool {
	_, ok := t.(*Dict)
	return ok
}

func AssertDict(t_ Token) *Dict {
	t, ok := t_.(*Dict)
	if ok {
		return t
	} else {
		panic("expected *Dict")
	}
}

func (t *Dict) TypeName() string {
	return "{}"
}

func (t *Dict) Dump() string {
	var b strings.Builder

	b.WriteString("{")

	for i, val := range t.vals {
		key := t.keys[i]

		b.WriteString(key.Dump())
		b.WriteString(":")
		b.WriteString(unwrapParens(val.Dump()))

		if i < len(t.vals)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String()
}

func (v *Dict) SetConstructors(cs []Call) Value {
	return &Dict{ValueData{newTokenData(v.Context()), v.constructors}, v.length, v.keys, v.vals}
}

func (v *Dict) Len() int {
	return v.length
}

func (v *Dict) GetStrict(s string) (Value, bool) {
	for i := 0; i < v.length; i++ {
		key := v.keys[i]
		if key.Value() == s {
			return v.vals[i], true
		}
	}

	return nil, false
}

func (v *Dict) Get(s string) Value {
	val, ok := v.GetStrict(s)
	if ok {
		return val
	} else {
		return nil
	}
}

func (v *Dict) Keys() []*String {
	// a copy
	res := make([]*String, v.length)
	for i := 0; i < v.length; i++ {
		res[i] = v.keys[i]
	}

	return res
}

func (v *Dict) Values() []Value {
	// a copy
	res := make([]Value, v.length)
	for i := 0; i < v.length; i++ {
		res[i] = v.vals[i]
	}

	return res
}

func (v *Dict) Link(scope Scope, ew ErrorWriter) Value {
	vals := []Value{}

	for _, val_ := range v.vals {
		val := val_.Link(scope, ew)
		vals = append(vals, val)
	}

	return NewDict(v.keys, vals, v.Context())
}

func (v *Dict) SubVars(stack *Stack) Value {
	vals := []Value{}

	for _, val_ := range v.vals {
		val := val_.SubVars(stack)
		vals = append(vals, val)
	}

	return NewDict(v.keys, vals, v.Context())
}

// keys in a are preferred to keys in b
func MergeDicts(a_ Value, b_ Value, ctx Context) *Dict {
	a := AssertDict(a_)
	b := AssertDict(b_)

	keys := []*String{}
	vals := []Value{}

	hasKey := func(s *String) bool {
		for _, k := range keys {
			if k.Value() == s.Value() {
				return true
			}
		}

		return false
	}

	for i, aKey := range a.keys {
		if !hasKey(aKey) {
			keys = append(keys, aKey)
			vals = append(vals, a.vals[i])
		}
	}

	for i, bKey := range b.keys {
		if !hasKey(bKey) {
			keys = append(keys, bKey)
			vals = append(vals, b.vals[i])
		}
	}

	return NewDict(keys, vals, ctx)
}
