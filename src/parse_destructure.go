package main

import "errors"

func parseDestructure(ts []Token, ew ErrorWriter) (Pattern, Value) {
	i := FindNextSymbol(ts, 0, "=")
	if i == -1 || i == len(ts)-1 {
		ew.Add(errors.New("invalid destructure"))
		return nil, nil
	}

	bef := ts[0:i]
	aft := ts[i+1:]

	pat := ParsePattern(bef, ew)
	rhs := ParseExpr(aft, ew)

	return pat, rhs
}
