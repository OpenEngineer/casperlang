package main

import "fmt"

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
