package main

import (
	"strings"
)

type DispatchableFunc interface {
	Func

	CalcDistance(args []Value) []int // the longer the better, and the lower the entries the better

	DumpHead() string
}

func listAmbiguousFuncs(fns []DispatchableFunc) string {
	var b strings.Builder

	b.WriteString("\n  Definitions:\n")

	for i, fn := range fns {
		b.WriteString("    ")
		b.WriteString(fn.DumpHead())
		b.WriteString(" (")
		b.WriteString(fn.Context().Error("").Error())
		b.WriteString(")")

		if i < len(fns)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func PickBest(fns []DispatchableFunc, args []Value, ctx Context) (DispatchableFunc, error) {
	if len(fns) == 0 {
		panic("no functions to pick from")
	}

	var (
		bestDist   []int              = nil
		bestEntryI                    = -1
		ambg       []DispatchableFunc = []DispatchableFunc{}
	)

	for i, fn := range fns {
		d := fn.CalcDistance(args)

		if d == nil {
			continue
		}

		if bestDist == nil {
			bestDist = d
			bestEntryI = i
		} else if len(d) > len(bestDist) {
			bestDist = d
			bestEntryI = i
		} else if len(d) == len(bestDist) {
			eachLE := true
			someLT := false

			for j, a := range d {
				if a < bestDist[j] {
					someLT = true
				} else if a > bestDist[j] {
					eachLE = false
				}
			}

			if someLT && eachLE {
				bestDist = d
				bestEntryI = i
			} else if !someLT && !eachLE {
				continue
			} else {
				if len(ambg) == 0 {
					ambg = append(ambg, fns[bestEntryI])
				}
				ambg = append(ambg, fns[i])
				continue
			}
		}
	}

	if len(ambg) > 0 {
		argInfo := ""
		if len(args) > 0 {
			argInfo = " " + DumpTypes(ToTypes(args))
		}

		return nil, ctx.Error("ambiguous dispatch of \"" + fns[0].Name() + argInfo + "\"" + listAmbiguousFuncs(ambg))
	}

	if bestDist == nil {
		return nil, nil
	} else {
		return fns[bestEntryI], nil
	}
}

func DispatchMessage(name string, args []Value, fns []DispatchableFunc) string {
	var b strings.Builder

	b.WriteString("unable to dispatch \"")
	b.WriteString(name)

	for _, arg := range args {
		b.WriteString(" ")
		b.WriteString(arg.Type().Dump())
	}
	b.WriteString("\"")

	if fns != nil && len(fns) > 0 {
		b.WriteString("\n  Have:\n")

		for _, fn := range fns {
			b.WriteString("    ")
			b.WriteString(fn.DumpHead())
			b.WriteString("\n")
		}
	}

	return b.String()
}
