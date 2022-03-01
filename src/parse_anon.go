package main

import "strconv"

func ParseAnon(gr *EscParens, ew ErrorWriter) *AnonFunc {
	ctx := gr.Context()

	args := []Pattern{}

	prev := 0

	for _, id := range gr.dollars {
		if id > prev+1 {
			for j := 0; j < id-(prev+1); j++ {
				args = append(args, NewNamedPattern(NewWord("_", ctx), NewAnyPattern(ctx), ctx))
			}

		}

		args = append(args, NewNamedPattern(NewWord("$"+strconv.Itoa(id), ctx), NewAnyPattern(ctx), ctx))

		prev = id
	}

	body := ParseExpr(gr.content, ew)

	return NewAnonFunc(args, body, ctx)
}
