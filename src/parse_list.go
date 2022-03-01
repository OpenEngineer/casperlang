package main

import "fmt"

func ParseList(gr *Brackets, ew ErrorWriter) *List {
	if DEBUG_PARSER {
		fmt.Printf("ParseList IN: %s\n", gr.Dump())
	}

	items := []Value{}

	for _, field := range gr.values {
		item := ParseExpr(field, ew)
		items = append(items, item)
	}

	return NewList(items, gr.Context())
}
