package main

import (
	"reflect"
	"strings"
)

type String struct {
	ValueData
	value string
}

func NewString(v string, ctx Context) *String {
	return &String{newValueData(ctx), v}
}

func (t *String) Value() string {
	return t.value
}

func (v *String) TypeName() string {
	return "String"
}

func (t *String) Dump() string {
	var b strings.Builder

	rs := []rune(t.Value())

	b.WriteString("\"")
	for _, r := range rs {
		if r == '\n' {
			b.WriteString("\\n")
		} else if r == '\t' {
			b.WriteString("\\t")
		} else {
			b.WriteRune(r)
		}
	}
	b.WriteString("\"")

	return b.String()
}

func IsString(t Token) bool {
	_, ok := t.(*String)
	return ok
}

func AssertString(t_ interface{}) *String {
	t, ok := t_.(*String)
	if !ok {
		panic("expected *String, got " + reflect.TypeOf(t_).String())
	}

	return t
}

func extractStrings(vs_ interface{}) []string {
	res := []string{}
	switch vs := vs_.(type) {
	case []Value:
		for _, v := range vs {
			s := AssertString(v)
			res = append(res, s.Value())
		}
	case []Token:
		for _, v := range vs {
			s := AssertString(v)
			res = append(res, s.Value())
		}
	default:
		panic("expected []Value or []Token")
	}

	return res
}

func (t *String) ToWord() *Word {
	return NewWord(t.Value(), t.Context())
}

func (v *String) SetConstructors(cs []Call) Value {
	return &String{ValueData{newTokenData(v.Context()), cs}, v.value}
}

func (v *String) Link(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *String) Eval(stack *Stack, ew ErrorWriter) Value {
	return v
}

// -1 if not found
func searchStrings(keys []*String, key *String) int {
	for i, k := range keys {
		if k.Value() == key.Value() {
			return i
		}
	}

	return -1
}
