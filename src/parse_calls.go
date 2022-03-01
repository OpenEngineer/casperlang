package main

import "fmt"

// at this point there should be no semicolons or assignments, only words, literals, groups and operators
func ParseCalls(ts []Token, ew ErrorWriter) Value {
	if DEBUG_PARSER {
		fmt.Printf("ParseCalls%s\n", DumpTokens(ts))
	}

	if len(ts) == 1 && IsValue(ts[0]) {
		return AssertValue(ts[0])
	}

	// divide the tokes into a list of expression

	parts := [][]Token{}

	prev := 0
	for i, t := range ts {
		if i == 0 {
			continue
		}

		if IsExprBoundary(ts[i-1], t) {
			parts = append(parts, ts[prev:i])

			prev = i
		}
	}

	if prev < len(ts) {
		parts = append(parts, ts[prev:])
	}

	// parse the operators of each part
	// alse detect that none of the parts, end in a nonpostfix operator, or start with a nonprefix operator
	exprs := []Value{}

	if DEBUG_PARSER {
		fmt.Printf("ParseCalls'[")
	}

	var name *Word = nil

	for i, part := range parts {
		if DEBUG_PARSER {
			fmt.Printf("%s,", DumpTokens(part))
		}

		if i == 0 && len(part) == 1 && IsWord(part[0]) {
			name = AssertWord(part[0])
		} else {
			if len(part) == 0 {
				panic("unexpected")
			} else if IsOperatorSymbol(part[0]) && len(part) == 1 {
				ew.Add(part[0].Context().Error("invalid syntax"))
			} else if IsOperatorSymbol(part[0]) && !IsPrefixOperatorSymbol(part[0]) {
				ew.Add(part[0].Context().Error("invalid syntax"))
				continue
			} else if IsOperatorSymbol(part[len(part)-1]) && !IsPostfixOperatorSymbol(part[len(part)-1]) {
				ew.Add(part[len(part)-1].Context().Error("invalid syntax"))
				continue
			} else {
				expr := ParseOperators(part, ew)

				if expr != nil {
					exprs = append(exprs, expr)
				}
			}
		}
	}

	if DEBUG_PARSER {
		fmt.Printf("]\n")
	}

	if name != nil {
		if name.Value() == "_" {
			ew.Add(name.Context().Error("invalid function call syntax"))
		}

		return NewNamedCall(name, exprs)
	} else {
		if len(exprs) == 0 {
			return nil
		} else if len(exprs) == 1 {
			return exprs[0]
		} else {
			if IsLiteral(exprs[0]) || IsContainer(exprs[0]) {
				ew.Add(exprs[0].Context().Error("invalid function call syntax"))
			}

			return NewBlindCall(exprs)
		}
	}
}
