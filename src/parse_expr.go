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
	case *NamedCall:
		return t
	case *BlindCall:
		return t
	default:
		ew.Add(t.Context().Error("invalid syntax"))
		return nil
	}
}

// no semicolons and no assignments
func parseExpr(ts []Token, ew ErrorWriter) Value {
	return parsePipes(parseGroups(ts, ew), ew)
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
