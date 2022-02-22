package main

type Pattern interface {
	Token

	CalcDistance(arg Value) []int

	// returns errors so it can be used in lhs of assignments
	Destructure(arg Value, fnScope *FuncScope, ew ErrorWriter) *FuncScope

	CheckTypeNames(scope Scope, ew ErrorWriter)

	ListTypes() []string
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

func WorstDistance(a []int, b []int) []int {
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
			return NewNamedPattern(t, NewSimplePattern(NewWord("Any", ctx)), ctx)
		} else {
			return NewSimplePattern(t)
		}
	case *Brackets:
		if t.Empty() {
			return NewSimplePattern(NewWord("[]", ctx))
		} else {
			return ParseTuplePattern(t, ew)
		}
	case *Braces:
		if t.Empty() {
			return NewSimplePattern(NewWord("{}", ctx))
		} else {
			return ParseStructPattern(t, ew)
		}
	case *Parens:
		if t.Empty() {
			ew.Add(ctx.Error("invalid pattern syntax"))
			return nil
		} else {
			return ParseConstructorPattern(t, ew)
		}
	case *NamedPattern:
		return t
	default:
		ew.Add(ctx.Error("invalid pattern syntax"))
		return nil
	}
}
