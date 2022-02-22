package main

import (
	"fmt"
)

func ParseExpr(ts []Token, ew ErrorWriter) Value {
	if DEBUG_PARSER {
		fmt.Printf("ParseExpr IN: %s\n", DumpTokens(ts))
	}

	return parseStatements(ts, ew)
}

func ParseSingleTokenExpr(t_ Token, ew ErrorWriter) Value {
	switch t := t_.(type) {
	case *Word:
		return NewNamedCall(t, []Value{})
	case *Int:
		return t
	case *Float:
		return t
	case *String:
		return t
	case *List:
		return t
	case *Dict:
		return t
	case *AnonFunc:
		return t
	case *Call:
		return t
	default:
		ew.Add(t.Context().Error("invalid syntax"))
		return nil
	}
}

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

// no semicolons and no assignments
func parseExpr(ts []Token, ew ErrorWriter) Value {
	return parsePipes(parseGroups(ts, ew), ew)
}

func parseGroups(ts []Token, ew ErrorWriter) []Token {
	if DEBUG_PARSER {
		fmt.Printf("parseGroups IN: %s\n", DumpTokens(ts))
	}

	res := []Token{}

	for _, t := range ts {
		switch {
		case IsGroupSymbol(t):
			// should've been grouped during the tokenization step
			ew.Add(t.Context().Error("invalid syntax"))
		case IsSymbol(t, ","):
			// should've been grouped during the tokenization step
			ew.Add(t.Context().Error("invalid syntax"))
		case IsParens(t):
			p := AssertParens(t)

			innerExpr := ParseExpr(p.Content(), ew)
			if innerExpr == nil {
				ew.Add(p.Context().Error("empty parentheses"))
			} else {
				res = append(res, innerExpr)
			}
		case IsEscParens(t):
			innerExpr := ParseAnon(AssertEscParens(t), ew)
			if innerExpr != nil {
				res = append(res, innerExpr)
			}
		case IsBrackets(t):
			innerExpr := ParseList(AssertBrackets(t), ew)
			if innerExpr != nil {
				res = append(res, innerExpr)
			}
		case IsBraces(t):
			innerExpr := ParseDict(AssertBraces(t), ew)
			if innerExpr != nil {
				res = append(res, innerExpr)
			}
		default:
			res = append(res, t)
		}
	}

	if DEBUG_PARSER {
		fmt.Printf("parseGroups OUT: %s\n", DumpTokens(res))
	}

	return res
}

func parsePipes(ts []Token, ew ErrorWriter) Value {
	pivot := FindLastSymbol(ts, "|")

	if pivot == -1 {
		return ParseCalls(ts, ew)
	} else if pivot == len(ts)-1 {
		ew.Add(ts[pivot].Context().Error("empty pipe rhs"))
		return ParseCalls(ts[0:len(ts)-1], ew)
	} else if pivot == 0 {
		ew.Add(ts[pivot].Context().Error("empty pipe lhs"))
		return ParseCalls(ts[1:], ew)
	}

	lhs := parsePipes(ts[0:pivot], ew)

	return ParseCalls(append(ts[pivot+1:], lhs), ew)
}

func IsExprBoundary(a, b Token) bool {
	oa := IsOperatorSymbol(a)
	ob := IsOperatorSymbol(b)
	if !oa && !ob {
		return true
	} else {
		ina := IsInfixOperatorSymbol(a)
		inb := IsInfixOperatorSymbol(b)

		if oa {
			if !ob { // b isn't an operator
				return !ina && !IsPrefixOperatorSymbol(a)
			} else { // both are operators
				if ina && !inb && IsPrefixOperatorSymbol(b) {
					return false
				} else if inb && !ina && IsPostfixOperatorSymbol(a) {
					return false
				} else {
					return true
				}
			}
		} else { // a isn't an operator
			return !inb
		}
	}
}
