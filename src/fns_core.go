package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func NewErrorType(str Value, ctx Context) Type {
	if str == nil {
		str = NewValueData(NewStringType(), ctx)
	}

	return NewUserType(&AnyType{}, "Error", []Value{str}, NewBuiltinContext())
}

func NewErrorValue(str Value, ctx Context) Value {
	return NewValueData(NewErrorType(str, ctx), ctx)
}

func NewOkType() Type {
	return NewUserType(&AnyType{}, "Ok", []Value{}, NewBuiltinContext())
}

func NewOkValue(ctx Context) Value {
	return NewValueData(NewOkType(), ctx)
}

// basic builtin manipulation functions, should be avaiable on all target platforms

var builtinCoreFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "Any",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewValueData(&AnyType{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Bool",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewValueData(NewBoolType(), NewBuiltinContext())
		},
	},
	BuiltinFuncConfig{
		Name:     "True",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewTrueValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "False",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewFalseValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Maybe",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewValueData(NewMaybeType(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Just",
		ArgTypes: []string{"Any"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}
			return NewJustValue(a, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Nothing",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewNothingValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Error",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}
			return NewErrorValue(a, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Ok",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewOkValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "if",
		ArgTypes: []string{"Bool", "Any", "Any"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			b := self.args[0].Eval(scope, ew)

			if bv, ok := GetBoolValue(b); ok {
				if bv {
					return self.args[1].Eval(scope, ew)
				} else {
					return self.args[2].Eval(scope, ew)
				}
			} else {
				return self.UpdateArgs([]Value{b, self.args[1], self.args[2]}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "show",
		ArgTypes: []string{"Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsInt(a) {
				return NewString(strconv.FormatInt(AssertInt(a).Value(), 10), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewStringType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "showf",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewString(fmt.Sprintf("%f", AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewStringType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsString(a) && IsString(b) {
				return NewString(AssertString(a).Value()+AssertString(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewStringType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"[]", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsList(a) && IsList(b) {
				return MergeLists(a, b, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, MergeListTypes(a.Type(), b.Type(), self.ctx), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"{}", "{}"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsDict(a) && IsDict(b) {
				return MergeDicts(a, b, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, MergeDictTypes(a.Type(), b.Type(), self.ctx), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "toInt",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsString(a) {
				i, err := strconv.ParseInt(AssertString(a).Value(), 10, 64)
				if err != nil {
					return NewNothingValue(self.ctx)
				}

				return NewJustValue(NewInt(i, self.ctx), self.ctx)
			} else {
				// TODO: union type of Nothing and Just value
				return self.UpdateArgs([]Value{a}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "toInt",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				i := int64(math.Round(AssertFloat(a).Value()))

				return NewInt(i, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "toFloat",
		ArgTypes: []string{"Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsInt(a) {
				f := float64(AssertInt(a).Value())

				return NewFloat(f, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "toFloat",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsString(a) {
				f, err := strconv.ParseFloat(AssertString(a).Value(), 64)
				if err != nil {
					return NewNothingValue(self.ctx)
				}

				return NewJustValue(NewFloat(f, self.ctx), self.ctx)
			} else {
				// TODO: make an Or type of Just and Nothing
				return self.UpdateArgs([]Value{a}, nil, self.ctx)
			}
		},
	},
	// returns Maybe
	BuiltinFuncConfig{
		Name:     "get",
		ArgTypes: []string{"[]", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsList(a) && IsInt(b) {
				lst := AssertList(a)
				i := AssertInt(b)

				return lst.Get(int(i.Value()), self.Context())
			} else {
				// TODO: make an Or type of Just or Nothing
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "get",
		ArgTypes: []string{"{}", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsDict(a) && IsString(b) {
				dict := AssertDict(a)
				k := AssertString(b)

				return dict.Get(k.Value(), self.Context())
			} else {
				// TODO: make an Or type of Just or Nothing
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "get",
		ArgTypes: []string{"String", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsInt(b) {
				str := AssertString(a)
				i := AssertInt(b).Value()

				sub := str.Value()[i : i+1]

				return NewString(sub, self.ctx)
			} else {
				// TODO: make an Or type of Just or Nothing
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "len",
		ArgTypes: []string{"[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if IsList(a) {
				return NewInt(int64(AssertList(a).Len()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "len",
		ArgTypes: []string{"{}"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if IsList(a) {
				return NewInt(int64(AssertDict(a).Len()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "len",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			if IsString(a) {
				return NewInt(int64(len(AssertString(a).Value())), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "toInts",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			if IsString(a) {
				items := []Value{}
				str := AssertString(a).Value()

				for _, r := range []rune(str) {
					items = append(items, NewInt(int64(r), self.ctx))
				}

				return NewList(items, self.ctx)
			} else {
				return self.UpdateArgs(
					[]Value{a},
					NewListType(
						[]Value{NewValueData(NewIntType(), self.ctx)},
						self.ctx,
					),
					self.ctx,
				)
			}
		},
	},
	BuiltinFuncConfig{
		Name: "toString",
		ArgPatterns: []Pattern{
			NewConstructorPattern(
				NewBuiltinWord("[]"),
				[]Pattern{
					NewSimplePattern(NewBuiltinWord("Int")),
				},
				NewBuiltinContext(),
			),
		},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			if IsList(a) {
				lst := AssertList(a)

				rs := []rune{}

				items := lst.Items()

				for _, i_ := range items {
					i := AssertInt(i_).Value()
					rs = append(rs, rune(i))
				}

				return NewString(string(rs), self.ctx)
			} else {
				return self.UpdateArgs(
					[]Value{a},
					NewStringType(),
					self.ctx,
				)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "map",
		ArgTypes: []string{"\\1", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}

			b := self.args[1].Eval(scope, ew)
			if b == nil {
				return b
			}

			if IsFunc(a) && IsList(b) {
				fn := AssertFunc(a)
				lst := AssertList(b)

				oldItems := lst.Items()

				newItems := []Value{}

				for _, oldItem := range oldItems {
					newItem := fn.Call([]Value{oldItem}, scope, self.ctx, ew)
					if newItem == nil {
						return nil
					} else if IsDeferredError(newItem) {
						return newItem
					}

					newItems = append(newItems, newItem)
				}

				return NewList(newItems, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "map",
		ArgTypes: []string{"\\2", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}

			b := self.args[1].Eval(scope, ew)
			if b == nil {
				return b
			}

			if IsFunc(a) && IsList(b) {
				fn := AssertFunc(a)
				lst := AssertList(b)

				oldItems := lst.Items()

				newItems := []Value{}

				for i, oldItem := range oldItems {
					newItem := fn.Call([]Value{NewInt(int64(i), self.ctx), oldItem}, scope, self.ctx, ew)
					if newItem == nil {
						return nil
					} else if IsDeferredError(newItem) {
						return newItem
					}

					newItems = append(newItems, newItem)
				}

				return NewList(newItems, self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	// returns Just or Nothing
	BuiltinFuncConfig{
		Name:     "fold",
		ArgTypes: []string{"\\2", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsFunc(a) && IsList(b) {
				fn := AssertFunc(a)
				lst := AssertList(b)

				var acc Value = nil

				oldItems := lst.Items()

				for _, oldItem := range oldItems {
					if acc == nil {
						acc = oldItem
					} else {
						acc = fn.Call([]Value{acc, oldItem}, scope, self.ctx, ew)
						if acc == nil {
							return nil
						} else if IsDeferredError(acc) {
							return acc
						}
					}
				}

				if acc == nil {
					return NewNothingValue(self.ctx)
				} else {
					return NewJustValue(acc, self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "fold",
		ArgTypes: []string{"\\2", "Any", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)
			c := self.args[2].Eval(scope, ew)

			if a == nil || b == nil || c == nil {
				return nil
			}

			if IsFunc(a) && IsList(c) {
				fn := AssertFunc(a)
				acc := b
				lst := AssertList(c)

				oldItems := lst.Items()

				for _, oldItem := range oldItems {
					acc = fn.Call([]Value{acc, oldItem}, scope, self.ctx, ew)
					if acc == nil {
						return nil
					} else if IsDeferredError(acc) {
						return acc
					}
				}

				return acc
			} else {
				return self.UpdateArgs([]Value{a, b, c}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() == AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() != AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() < AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() > AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() <= AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				if AssertInt(a).Value() >= AssertInt(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() == AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() != AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() < AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() > AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() <= AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsFloat(a) && IsFloat(b) {
				if AssertFloat(a).Value() >= AssertFloat(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() == AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() != AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() < AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() > AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() <= AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		ArgTypes: []string{"String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0]
			b := self.args[1]

			if IsString(a) && IsString(b) {
				if AssertString(a).Value() >= AssertString(b).Value() {
					return NewTrueValue(self.ctx)
				} else {
					return NewFalseValue(self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{a, b}, NewBoolType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "panic",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			s := self.args[0].Eval(scope, ew)
			if s == nil {
				panic("expected to be caught by Call")
			}

			if IO_ACTIVE {
				// TODO: print stack trace somehow
				fmt.Fprintf(os.Stderr, "%s\n", self.ctx.Error(AssertString(s).Value()).Error())
				os.Exit(1)
				return nil
			} else {
				return self.UpdateArgs([]Value{s}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "sort",
		ArgTypes: []string{"\\2", "[]"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			comp_ := self.args[0]
			lst_ := self.args[1]

			if IO_ACTIVE {
				AssertList(lst_)
			}

			var t Type = nil
			if lst_.Type() != nil {
				// what to do about the tuple type?
				t = AssertListType(lst_.Type())
			}
			comp := AssertFunc(comp_)

			if IsList(lst_) {
				s := NewListSorter(AssertList(lst_), scope, comp, ew, self.ctx)
				sort.Stable(s)

				if !s.ew.Empty() {
					return nil
				} else {
					return NewList(s.items, self.ctx)
				}
			} else {
				return self.UpdateArgs([]Value{comp_, lst_}, t, self.ctx)
			}
		},
	},
}
