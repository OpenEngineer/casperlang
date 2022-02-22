package main

type OperatorSettings struct {
	Infix       bool
	Postfix     bool
	Prefix      bool
	Precedence  int
	LeftToRight bool // default is right to left
}

// |, ; and = are technically also operators, but are handled explicitely
var operators map[string]OperatorSettings = map[string]OperatorSettings{
	".":  OperatorSettings{Precedence: 18, Infix: true},
	"!":  OperatorSettings{Precedence: 17, Prefix: true},
	"*":  OperatorSettings{Precedence: 15, Infix: true, LeftToRight: true},
	"/":  OperatorSettings{Precedence: 15, Infix: true, LeftToRight: true},
	"-":  OperatorSettings{Precedence: 14, Infix: true, Prefix: true, LeftToRight: true},
	"+":  OperatorSettings{Precedence: 14, Infix: true, LeftToRight: true},
	">":  OperatorSettings{Precedence: 12, Infix: true},
	">=": OperatorSettings{Precedence: 12, Infix: true},
	"<":  OperatorSettings{Precedence: 12, Infix: true},
	"<=": OperatorSettings{Precedence: 12, Infix: true},
	"==": OperatorSettings{Precedence: 11, Infix: true},
	"!=": OperatorSettings{Precedence: 11, Infix: true},
	"&&": OperatorSettings{Precedence: 6, Infix: true},
	"||": OperatorSettings{Precedence: 5, Infix: true},
	";":  OperatorSettings{Precedence: 0, Infix: false}, // treated differently!
	"=":  OperatorSettings{Precedence: 0, Infix: false}, // treated differently!
}

func FindLowestInfixOp(ts []Token) int {
	lowestI := -1
	lowestP := 0

	for i, t := range ts {
		if IsInfixOperatorSymbol(t) {
			s := AssertSymbol(t)
			p := operators[s.Value()].Precedence
			l2r := operators[s.Value()].LeftToRight
			if lowestI == -1 {
				lowestI = i
				lowestP = p
			} else if p < lowestP {
				lowestI = i
				lowestP = p
			} else if p == lowestP && l2r && i > lowestI {
				lowestI = i
			}
		}
	}

	return lowestI
}

// [op] a op b op c op d [op]
func ParseOperators(ts []Token, ew ErrorWriter) Value {
	// find the lowest precedence operator
	pivotI := FindLowestInfixOp(ts)

	if pivotI != -1 && len(ts) > 2 {
		lhTs := ts[0:pivotI]
		rhTs := ts[pivotI+1:]

		lhs := ParseOperators(lhTs, ew)
		rhs := ParseOperators(rhTs, ew)

		return NewNamedCall(AssertSymbol(ts[pivotI]).ToWord(), []Value{lhs, rhs})
	} else if len(ts) == 1 {
		t := ts[0]
		return ParseSingleTokenExpr(t, ew)
	} else if len(ts) == 2 {
		if IsPrefixOperatorSymbol(ts[0]) {
			arg := ParseOperators(ts[1:], ew)
			return NewNamedCall(AssertSymbol(ts[0]).ToWord(), []Value{arg})
		} else if IsPostfixOperatorSymbol(ts[1]) {
			arg := ParseOperators(ts[0:1], ew)
			return NewNamedCall(AssertSymbol(ts[1]).ToWord(), []Value{arg})
		} else {
			ew.Add(ts[0].Context().Merge(ts[1].Context()).Error("invalid syntax"))
			return nil
		}
	} else {
		panic("unexpected")
	}
}
