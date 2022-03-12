package main

import "reflect"

// first returned value is concrete value, second is virtual value
// if there was an error then both are nil
// if the condition failed, then virt will be nil, and concrete will the most recent eval
func EvalUntil(arg Value, cond func(string) bool, ew ErrorWriter) (Value, Value) {
	if arg == nil {
		return nil, nil
	}

	for _, c := range arg.Constructors() {
		if cond(c.TypeName()) {
			return arg, c
		}
	}

	if arg.TypeName() == "All" {
		return arg, arg
	}

	for {
		if cond(arg.TypeName()) {
			return arg, arg
		} else {
			call, ok := arg.(Call)
			if !ok {
				return nil, nil
			}

			arg = call.Eval(ew)
			if arg == nil || !ew.Empty() {
				return nil, nil
			}
		}
	}
}

func EvalEager(v Value, ew ErrorWriter) Value {
Outer:
	for v != nil {
		switch v_ := v.(type) {
		case Call:
			v = v_.Eval(ew)
		case *List:
			items := []Value{}
			for _, item := range v_.Items() {
				items = append(items, EvalEager(item, ew))
				if !ew.Empty() {
					return nil
				}
			}

			return NewList(items, v_.Context())
		case *Dict:
			vals := []Value{}
			for _, val := range v_.Values() {
				vals = append(vals, EvalEager(val, ew))
				if !ew.Empty() {
					return nil
				}
			}

			return NewDict(v_.Keys(), vals, v_.Context())
		default:
			break Outer
		}
	}

	return v
}

// eg. Just 4 shows as Just 4 instead of Any
func EvalPretty(v Value) Value {
	cs := v.Constructors()
	n := len(cs)

	switch v_ := v.(type) {
	case *Any:
		if n == 0 {
			return v_
		} else {
			v = cs[n-1]
			cs = cs[0 : n-1]
			return EvalPretty(v.SetConstructors(cs))
		}
	case Call:
		if v_.NumArgs() > 0 || n == 0 {
			return v_
		} else {
			v = cs[n-1]
			cs = cs[0 : n-1]
			return EvalPretty(v.SetConstructors(cs))
		}
	case *List:
		items := []Value{}
		for _, item := range v_.Items() {
			items = append(items, EvalPretty(item))
		}
		return NewList(items, v_.Context())
	case *Dict:
		vals := []Value{}
		for _, val := range v_.Values() {
			vals = append(vals, EvalPretty(val))
		}
		return NewDict(v_.Keys(), vals, v_.Context())
	default:
		return v
	}
}

func ConvertUntyped(x_ interface{}, ctx Context) Value {

	switch x := x_.(type) {
	case map[string]interface{}:
		keys := make([]*String, 0)
		vals := make([]Value, 0)

		for k, v := range x {
			keys = append(keys, NewString(k, ctx))
			vals = append(vals, ConvertUntyped(v, ctx))
		}

		return NewDict(keys, vals, ctx)
	case []interface{}:
		items := make([]Value, 0)

		for _, v := range x {
			items = append(items, ConvertUntyped(v, ctx))
		}

		return NewList(items, ctx)
	case string:
		return NewString(x, ctx)
	case float64:
		return NewFloat(x, ctx)
	case int:
		return NewInt(int64(x), ctx)
	case int64:
		return NewInt(x, ctx)
	case bool:
		if x {
			return NewInt(1, ctx)
		} else {
			return NewInt(0, ctx)
		}
	default:
		panic("unhandled type " + reflect.TypeOf(x_).String())
	}
}
