package main

import (
  "strconv"
)

type NL struct {
  TokenData
  ind int
}

func NewNL(ind int, ctx Context) *NL {
  return &NL{newTokenData(ctx), ind}
}

func (t *NL) Dump() string {
  return "NL(" + strconv.Itoa(t.ind) + ")"
}

func (t *NL) Indent() int {
  return t.ind
}

func IsNL(t Token) bool {
  _, ok := t.(*NL)
  return ok
}

func AssertNL(t_ Token) *NL {
  t, ok := t_.(*NL)
  if !ok {
    panic("expected *NL")
  }

  return t
}

// newlines are irrelevant inside groups
func RemoveNLs(ts []Token) []Token {
  res := []Token{}
  for _, t := range ts {
    if !IsNL(t) {
      res = append(res, t)
    }
  }

  return res
}
