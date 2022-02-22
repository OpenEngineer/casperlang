package main

import (
	"fmt"
)

func GroupTokens(ts []Token, ew ErrorWriter) []Token {
	if DEBUG_PARSER {
		fmt.Printf("GroupTokens[%s]\n", DumpTokens(ts))
	}

	res := []Token{}

	for i := 0; i < len(ts); i++ {
		t := ts[i]
		switch {
		case IsGroupOpenSymbol(t):
			s := AssertSymbol(t)

			start := i
			stop := FindGroupMatch(ts, start, s, ew)

			if stop != -1 {
				grTs := RemoveNLs(ts[start : stop+1])

				var (
					innerT Token
				)

				switch {
				case s.Value() == "\\(":
					innerT = GroupEscParens(grTs, ew)
				case s.Value() == "(":
					innerT = GroupParens(grTs, ew)
				case s.Value() == "{":
					innerT = GroupBraces(grTs, ew)
				case s.Value() == "[":
					innerT = GroupBrackets(grTs, ew)
				default:
					panic("unexpected GroupOpenSymbol \"" + s.Value() + "\"")
				}

				res = append(res, innerT)
				i = stop // remember the i++ at the top!
			}
		case IsGroupCloseSymbol(t):
			s := AssertSymbol(t)
			ew.Add(s.Context().Error("unmatched \"" + s.Value() + "\""))
		default:
			res = append(res, t)
		}
	}

	return res
}
