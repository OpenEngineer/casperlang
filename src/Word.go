package main

import "strings"

type Word struct {
	TokenData
	value string
}

func NewWord(value string, ctx Context) *Word {
	return &Word{newTokenData(ctx), value}
}

func NewBuiltinWord(value string) *Word {
	return &Word{newTokenData(NewBuiltinContext()), value}
}

func (t *Word) Value() string {
	return t.value
}

func (t *Word) Dump() string {
	return t.value
}

func (t *Word) ToString() *String {
	return NewString(t.Value(), t.Context())
}

func IsWord(t_ Token, v ...string) bool {
	t, ok := t_.(*Word)

	if len(v) == 0 {
		return ok
	} else if len(v) == 1 {
		return ok && t.value == v[0]
	} else {
		panic("expected 0 or 1 params to IsWord")
	}
}

func IsLowerCaseWord(t_ Token) bool {
	t, ok := t_.(*Word)

	if ok {
		return t.IsLowerCase()
	} else {
		return false
	}
}

func isLowerCase(s string) bool {
	s = strings.TrimLeft(s, "_")

	if len(s) == 0 {
		return true
	} else {
		fl := s[0]
		return fl == 95 || (fl >= 97 && fl <= 122)
	}
}

func (w *Word) IsLowerCase() bool {
	return isLowerCase(w.Value())
}

func isUpperCase(s string) bool {
	s = strings.TrimLeft(s, "_")
	if len(s) == 0 {
		return false
	} else {
		fl := s[0]

		return (fl >= 65 && fl <= 90)
	}
}

func isConstructorName(s string) bool {
	return isUpperCase(s)
}

func (w *Word) IsUpperCase() bool {
	return isUpperCase(w.Value())
}

func IsUpperCaseWord(t_ Token) bool {
	t, ok := t_.(*Word)

	if ok {
		return t.IsUpperCase()
	} else {
		return false
	}
}

func AssertWord(t_ Token) *Word {
	t, ok := t_.(*Word)
	if ok {
		return t
	} else {
		panic("expected *Word")
	}
}
