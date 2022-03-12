package main

import (
	"sort"
	"strconv"
	"strings"
)

type EscParens struct {
	TokenData
	content []Token
	dollars []int
}

func NewEscParens(content []Token, ctx Context) *EscParens {
	return &EscParens{newTokenData(ctx), content, []int{}}
}

func IsEscParens(t Token) bool {
	_, ok := t.(*EscParens)
	return ok
}

func AssertEscParens(t_ Token) *EscParens {
	t, ok := t_.(*EscParens)
	if ok {
		return t
	} else {
		panic("expected *EscParens")
	}
}

func (t *EscParens) Dump() string {
	var b strings.Builder

	b.WriteString("\\(")

	for i, item := range t.content {
		b.WriteString(item.Dump())

		if i < len(t.content)-1 {
			b.WriteString(" ")
		}
	}

	b.WriteString(")")

	return b.String()
}

func GroupEscParens(ts []Token, ew ErrorWriter) Token {
	ctx := ts[0].Context()

	inner := ts[1 : len(ts)-1]

	//inner := GroupTokens(ts[1:len(ts)-1], ew)

	// now start auto naming the dollar signs
	impl := false
	expl := false
	ids := []int{}

	//
	for i := 0; i < len(inner); i++ {
		t := inner[i]
		if IsSymbol(t, "\\(") {
			// skip until the next group match
			i = FindGroupMatch(inner, i, AssertSymbol(t), ew)
		} else if IsDollar(t) {
			d := AssertDollar(t)

			if !impl && !expl {
				if d.id == 0 {
					impl = true
					id := len(ids) + 1
					inner[i] = NewWord("$"+strconv.Itoa(id), d.Context())
					ids = append(ids, id)
				} else {
					expl = true

					found := false
					for _, check := range ids {
						if check == d.id {
							found = true
							break
						}
					}

					if !found {
						ids = append(ids, d.id)
					}
					inner[i] = NewWord("$"+strconv.Itoa(d.id), d.Context())
				}
			} else if impl {
				if d.id == 0 {
					id := len(ids) + 1
					inner[i] = NewWord("$"+strconv.Itoa(id), d.Context())
					ids = append(ids, id)
				} else {
					ew.Add(d.Context().Error("mixed explicit and implicit lambda arg pos"))
				}
			} else if expl {
				if d.id != 0 {
					found := false
					for _, check := range ids {
						if check == d.id {
							found = true
							break
						}
					}

					if !found {
						ids = append(ids, d.id)
					}

					inner[i] = NewWord("$"+strconv.Itoa(d.id), d.Context())
				} else {
					ew.Add(d.Context().Error("mixed explicit and implicit lambda arg pos"))
				}
			} else {
				panic("unexpected")
			}
		}
	}

	if !impl && !expl {
		ew.Add(ctx.Error("lambda function without args"))
	}

	inner = GroupTokens(inner, ew)

	res := NewEscParens(inner, ctx)
	res.dollars = ids
	sort.Ints(res.dollars)

	return res
}
