package main

func NewBoolType() *UserType {
	return NewUserType(&AnyType{}, "Bool", []Value{}, NewBuiltinContext())
}

func NewTrueValue(ctx Context) Value {
	return NewValueData(NewUserType(NewBoolType(), "True", []Value{}, NewBuiltinContext()), ctx)
}

func NewFalseValue(ctx Context) Value {
	return NewValueData(NewUserType(NewBoolType(), "False", []Value{}, NewBuiltinContext()), ctx)
}

func GetBoolValue(v Value) (bool, bool) {
	t_ := v.Type()

	for t_ != nil {
		if t, ok := t_.(*UserType); ok {
			if t.name == "True" {
				return true, true
			} else if t.name == "False" {
				return false, true
			} else {
				t_ = t.Parent()
			}
		} else {
			break
		}

	}

	return false, false
}
