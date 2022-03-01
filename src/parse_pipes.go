package main

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
