package main

func NewMaybeType() *UserType {
	return NewUserType(&AnyType{}, "Maybe", []Value{}, NewBuiltinContext())
}

func NewJustValue(arg Value, ctx Context) Value {
	return NewValueData(NewUserType(NewMaybeType(), "Just", []Value{arg}, NewBuiltinContext()), ctx)
}

func NewNothingValue(ctx Context) Value {
	return NewValueData(NewUserType(NewMaybeType(), "Nothing", []Value{}, NewBuiltinContext()), ctx)
}

func GetMaybeValue(v Value) (Value, bool) {
	t_ := v.Type()

	for t_ != nil {
		if t, ok := t_.(*UserType); ok {
			if t.name == "Nothing" {
				return nil, true
			} else if t.name == "Just" {
				return t.args[0], true
			} else {
				t_ = t.Parent()
			}
		} else {
			break
		}
	}

	return nil, false
}
