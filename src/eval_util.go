package main

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
