package main

import (
	"strconv"
)

type Dollar struct {
	TokenData
	id int
}

func NewDollar(id int, ctx Context) *Dollar {
	return &Dollar{newTokenData(ctx), id}
}

func (t *Dollar) Dump() string {
	if t.id <= 0 {
		return "Dollar$"
	} else {
		return "Dollar$" + strconv.Itoa(t.id)
	}
}

func IsDollar(t Token) bool {
	_, ok := t.(*Dollar)
	return ok
}

func AssertDollar(t_ Token) *Dollar {
	t, ok := t_.(*Dollar)

	if ok {
		return t
	} else {
		panic("expected *Dollar")
	}
}
