package main

import "fmt"

// there aren't really statements, but expressions separated by equals and semicolons
func parseStatements(ts []Token, ew ErrorWriter) Value {
	if DEBUG_PARSER {
		fmt.Printf("parseStatements IN: %s\n", DumpTokens(ts))
	}

	// lowest precedence operators are semicolons and equals
	var (
		parts []Token = []Token{}
	)

	prevSemicolon := -1
	groupCount := 0

Outer:
	for i, t := range ts {
		if IsGroupOpenSymbol(t) {
			groupCount += 1
		} else if IsGroupCloseSymbol(t) {
			groupCount -= 1
		}

		if groupCount != 0 {
			continue Outer
		}

		switch {
		case IsSymbol(t, "="):
			nextSemicolon := FindNextSymbol(ts, i+1, ";")
			if nextSemicolon == -1 {
				panic("unexpected")
			}

			fnBodyExpr := ParseExpr(ts[nextSemicolon+1:], ew)

			fnHeadPattern := ParsePattern(ts[prevSemicolon+1:i], ew)

			if fnBodyExpr == nil {
				ew.Add(t.Context().Error("unused assignment"))
			} else {
				checkArgNames([]Pattern{fnHeadPattern}, ew)
				fnExpr := NewSingleArgAnonFunc(fnHeadPattern, fnBodyExpr, t.Context())

				fnInput := parseExpr(ts[i+1:nextSemicolon], ew)

				if fnInput == nil {
					ew.Add(t.Context().Error("empty assignment rhs"))
				}

				parts = append(parts, NewNamedCall(AssertSymbol(t).ToWord(), []Value{fnInput, fnExpr}))
			}

			break Outer
		case IsSymbol(t, ";") || i == len(ts)-1:
			stop := i
			if !IsSymbol(t, ";") {
				stop = len(ts)
			}

			subExpr := parseExpr(ts[prevSemicolon+1:stop], ew)

			parts = append(parts, subExpr)

			prevSemicolon = i

			if IsSymbol(t, ";") {
				parts = append(parts, ts[prevSemicolon])
			}
		default:
			continue Outer
		}
	}

	if groupCount != 0 {
		panic("unmatched group should've been caught before")
	}

	var expr Value = nil

	for i := len(parts) - 1; i >= 0; i-- {
		part_ := parts[i]

		if IsValue(part_) {
			part := AssertValue(part_)

			if expr == nil {
				expr = part
			} else {
				expr = NewNamedCall(
					AssertSymbol(parts[i+1]).ToWord(),
					[]Value{part, expr},
				)
			}
		}
	}

	if DEBUG_PARSER {
		if expr != nil {
			fmt.Printf("parseStatements OUT: %s\n", expr.Dump())
		} else {
			fmt.Printf("parseStatements OUT: <empty>\n")
		}
	}

	return expr
}
