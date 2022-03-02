package main

import (
	"strconv"
	"strings"
)

// useful for builtinfuncs
func ParsePatternString(str string) Pattern {
	ts := TokenizeString(str)

	ew := NewErrorWriter()

	p := ParsePattern(ts, ew)

	if !ew.Empty() {
		panic(str + " " + ew.Dump())
	}

	return p
}

func ParsePattern(ts []Token, ew ErrorWriter) Pattern {
	ps := ParsePatterns(ts, ew)
	if !ew.Empty() {
		return nil
	}

	if len(ps) > 1 {
		ew.Add(ps[1].Context().Error("unexpected pattern"))
		return ps[0]
	}

	return ps[0]
}

func worstDistance(a []int, b []int) []int {
	if len(a) == 0 && len(b) > 0 {
		return b
	} else if len(a) > 0 && len(b) == 0 {
		return a
	} else if a == nil || b == nil {
		return nil
	} else if len(a) < len(b) {
		return a
	} else if len(a) == len(b) {
		c := make([]int, len(a))

		for i, x := range a {
			if x > b[i] {
				c[i] = x
			} else {
				c[i] = b[i]
			}
		}

		return c
	} else {
		return b
	}
}

// pattern tokens are split first
// first :: is parsed
func ParsePatterns(ts []Token, ew ErrorWriter) []Pattern {
	// look for the double colon symbols

	pivot := FindNextSymbol(ts, 0, "::")

	if pivot == -1 {
		res := []Pattern{}
		for _, t := range ts {
			p := parsePattern(t, ew)
			if p != nil {
				res = append(res, p)
			}
		}

		return res
	} else {
		res := []Token{}

		if len(ts) == 1 {
			ew.Add(ts[pivot].Context().Error("invalid syntax"))
			return []Pattern{}
		} else if pivot == 0 {
			ew.Add(ts[pivot].Context().Error("expected variable name on left"))
			return ParsePatterns(ts[pivot+1:], ew)
		} else if pivot == len(ts)-1 {
			ew.Add(ts[pivot].Context().Error("expected pattern on right"))
			return ParsePatterns(ts[:pivot], ew)
		}

		left := ts[pivot-1]
		right := ts[pivot+1]

		noErrors := true

		if !IsLowerCaseWord(left) {
			ew.Add(left.Context().Error("expected variable name"))
			noErrors = false
		}

		if IsLowerCaseWord(right) {
			ew.Add(right.Context().Error("invalid syntax"))
			noErrors = false
		}

		if noErrors {
			// do that actual parsing of the right
			rightP := parsePattern(right, ew)

			res = ts[0 : pivot-1]
			res = append(res, NewNamedPattern(AssertWord(left), rightP, ts[pivot].Context()))
			res = append(res, ts[pivot+2:]...)

			return ParsePatterns(res, ew)
		} else {
			res = ts[0 : pivot-1]
			res = append(res, ts[pivot+2:]...)

			return ParsePatterns(res, ew)
		}
	}
}

func parsePattern(t_ Token, ew ErrorWriter) Pattern {
	ctx := t_.Context()

	switch t := t_.(type) {
	case *Word:
		if t.IsLowerCase() {
			return NewNamedPattern(t, NewAnyPattern(ctx), ctx)
		} else if strings.HasPrefix(t.Value(), "\\") {
			na, err := strconv.Atoi(t.Value()[1:])
			if err != nil {
				panic(err)
			}

			return NewFuncPattern(na, t.Context())
		} else {
			switch t.Value() {
			case "Any":
				return NewAnyPattern(ctx)
			case "Int":
				return NewPrimPattern("Int", ctx)
			case "Float":
				return NewPrimPattern("Float", ctx)
			case "String":
				return NewPrimPattern("String", ctx)
			case "IO":
				return NewPrimPattern("IO", ctx)
			default:
				return NewTypePattern(t)
			}
		}
	case *Brackets:
		if t.Empty() {
			return NewPrimPattern("[]", ctx)
		} else {
			return ParseTuplePattern(t, ew)
		}
	case *Braces:
		if t.Empty() {
			return NewPrimPattern("{}", ctx)
		} else {
			return ParseStructPattern(t, ew)
		}
	case *Parens:
		if t.Empty() {
			ew.Add(ctx.Error("invalid pattern syntax2"))
			return nil
		} else {
			return ParseConstructorPattern(t, ew)
		}
	case *NamedPattern:
		return t
	default:
		ew.Add(ctx.Error("invalid pattern syntax1"))
		return nil
	}
}

func ParseStructPattern(gr *Braces, ew ErrorWriter) *StructPattern {
	keys := []*String{}
	vals := []Pattern{}

	for i, ts := range gr.vals {
		v := ParsePattern(ts, ew)
		if v != nil {
			vals = append(vals, v)
			keys = append(keys, gr.keys[i])
		}
	}

	return NewStructPattern(keys, vals, gr.Context())
}

func ParseTuplePattern(gr *Brackets, ew ErrorWriter) *TuplePattern {
	items := []Pattern{}

	for _, ts := range gr.values {
		p := ParsePattern(ts, ew)
		if p != nil {
			items = append(items, p)
		}
	}

	return NewTuplePattern(items, gr.Context())
}

func ParseConstructorPattern(p *Parens, ew ErrorWriter) Pattern {
	ctx := p.Context()

	ts := p.content

	t := ts[0]

	args := ParsePatterns(ts[1:], ew)

	var name *Word
	switch {
	case IsWord(t):
		name = AssertWord(t)
	case IsEmptyBraces(t):
		if len(args) == 1 {
			return NewDictPattern(args[0], ctx)
		} else {
			ew.Add(t.Context().Error("invalid dict pattern syntax"))
			return nil
		}
	case IsEmptyBrackets(t):
		if len(args) == 1 {
			return NewListPattern(args[0], ctx)
		} else {
			ew.Add(t.Context().Error("invalid list pattern syntax"))
			return nil
		}
	default:
		ew.Add(t.Context().Error("invalid constructor pattern syntax"))
		return nil
	}

	if name.Value() == "Int" ||
		name.Value() == "Float" ||
		name.Value() == "String" ||
		name.Value() == "Any" ||
		name.Value() == "IO" {
		ew.Add(name.Context().Error("can't apply constructor pattern to \"" + name.Value() + "\""))
		return nil
	}

	return NewConstructorPattern(name, args, ctx)
}

func checkArgNames(args []Pattern, ew ErrorWriter) {
	names := []*Word{}

	for _, arg := range args {
		names = append(names, arg.ListNames()...)
	}

	for i, name := range names {
		for j, check := range names {
			if j > i && name.Value() == check.Value() {
				ew.Add(check.Context().Error("duplicate argument \"" + check.Value() + "\""))
			}
		}
	}
}
