package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Call struct {
	ValueData
	name *Word
	args []Value
}

func NewNamedCall(name *Word, args []Value) *Call {
	if name == nil {
		panic("name can't be nil, use NewBlindCall() instead")
	}

	return &Call{newValueData(nil, name.Context()), name, args}
}

func NewBlindCall(args []Value) *Call {
	return &Call{newValueData(nil, args[0].Context()), nil, args}
}

func IsCall(t Token) bool {
	_, ok := t.(*Call)
	return ok
}

func AssertCall(t_ Token) *Call {
	t, ok := t_.(*Call)
	if ok {
		return t
	} else {
		panic("expected *Call")
	}
}

func (v *Call) Update(type_ Type, ctx Context) Value {
	return v.UpdateArgs(v.args, type_, ctx)
}

func (v *Call) UpdateArgs(args []Value, type_ Type, ctx Context) *Call {
	return &Call{newValueData(type_, ctx), v.name, args}
}

func (t *Call) isBlind() bool {
	return t.name == nil
}

func (t *Call) Name() string {
	if t.isBlind() {
		return ""
	} else {
		return t.name.Value()
	}
}

func (t *Call) Dump() string {
	var b strings.Builder

	b.WriteString("(")
	if t.name != nil {
		b.WriteString(t.name.Value())

		for _, arg := range t.args {
			b.WriteString(" ")
			if arg != nil {
				b.WriteString(arg.Dump())
			} else {
				b.WriteString("<nil>")
			}
		}
	} else {
		for i, arg := range t.args {
			b.WriteString(arg.Dump())
			if i < len(t.args)-1 {
				b.WriteString(" ")
			}
		}
	}
	b.WriteString(")")

	return b.String()
}

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

func (v *Call) IsNamed() bool {
	return v.name != nil
}

func (v *Call) evalArgs(scope Scope, ew ErrorWriter) []Value {
	res := []Value{}

	for _, arg := range v.args {
		if arg == nil {
			panic("unexpected")
		}

		argVal := arg.Eval(scope, ew)
		if argVal == nil {
			return nil
		}

		res = append(res, argVal)
	}

	return res
}

func (v *Call) Eval(scope Scope, ew ErrorWriter) Value {
	// source of non-lazyness, how can we fix this? separate type eval?
	args := v.evalArgs(scope, ew)
	if args == nil || !ew.Empty() {
		return nil
	}

	var defErr Value
	for _, arg := range args {
		if arg.Type() == nil {
			// return self so we can try again in another eval round (with IO turned on)
			return v.UpdateArgs(args, v.Type(), v.Context())
		} else if IsDeferredError(arg) { // TODO: get rid of deferred errors by doing tru lazy evaluation
			defErr = arg
		}
	}

	if v.IsNamed() {
		fn := scope.Dispatch(v.name, args, ew)

		if !ew.Empty() {
			return nil
		} else if fn == nil {
			if defErr != nil {
				return defErr
			} else {
				return NewDeferredError(DispatchMessage(v.Name(), args, scope.CollectFunctions(v.Name())), v.Context())
			}
		} //else if defErr != nil {
		//return defErr
		//}

		return fn.Call(args, scope, v.name.Context(), ew)
	} else {
		fnVal := args[0]
		args = args[1:]

		if fnVal.Type() == nil {
			ew.Add(fnVal.Context().Error(fmt.Sprintf("expected function with %d args, got \\?", len(args))))
		} else if fnVal.Type().CalcNameDistance(NewBuiltinWord("\\"+strconv.Itoa(len(args)))) < 0 {
			ew.Add(fnVal.Context().Error(fmt.Sprintf("expected function with %d args, got %s", len(args), fnVal.Type().Dump())))
			return nil
		}

		if fn, ok := fnVal.(Func); ok {
			return fn.Call(args, scope, fnVal.Context(), ew)

		} else {
			return v.UpdateArgs(args, v.Type(), v.Context())
		}
	}
}

func (v *Call) CheckTypeNames(scope Scope, ew ErrorWriter) {
	for _, arg := range v.args {
		arg.CheckTypeNames(scope, ew)
	}
}
