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
	t := NewDictType(vals, ctx)

	return &Dict{newValueData(t, ctx), len(keys), keys, vals}
}

func NewStruct(keys []*String, vals []Value, ctx Context) *Dict {
	return &Dict{newValueData(nil, ctx), len(keys), keys, vals}
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

func (v *Dict) Update(type_ Type, ctx Context) Value {
	return &Dict{newValueData(type_, ctx), v.length, v.keys, v.vals}
}

func (t *Dict) Dump() string {
	var b strings.Builder

	b.WriteString("{")

	for i, val := range t.vals {
		key := t.keys[i]

		b.WriteString(key.Dump())
		b.WriteString(":")
		b.WriteString(val.Dump())

		if i < len(t.vals)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")

	return b.String()
}

func ParseDict(gr *Braces, ew ErrorWriter) *Dict {
	keys := []*String{}
	vals := []Value{}

	for i, field := range gr.vals {
		val := ParseExpr(field, ew)
		if val != nil {
			keys = append(keys, gr.keys[i])
			vals = append(vals, val)
		}

	}

	return NewStruct(keys, vals, gr.Context())
}

// XXX: type is lost?
func (v *Dict) Eval(scope Scope, ew ErrorWriter) Value {
	vals := []Value{}

	for _, val_ := range v.vals {
		val := val_.Eval(scope, ew)
		if val == nil {
			return nil
		}

		vals = append(vals, val)
	}

	t := v.Type()
	if t == nil {
		t = NewStructType(v.keys, vals, v.Context())
	}

	return &Dict{newValueData(t, v.Context()), v.length, v.keys, vals}
}

func (v *Dict) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, val := range v.vals {
		val.CheckTypeNames(scope, ew)
	}
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

func (v *Dict) Get(s string, ctx Context) Value {
	val, ok := v.GetStrict(s)
	if ok {
		return NewJustValue(val, ctx)
	} else {
		return NewNothingValue(ctx)
	}
}

func (v *Dict) Len() int {
	return v.length
}
