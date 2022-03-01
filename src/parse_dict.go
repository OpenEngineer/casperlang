package main

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

	return NewDict(keys, vals, gr.Context())
}
