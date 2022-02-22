package main

import "math"

var builtinMathFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				return NewInt(AssertInt(a).Value()+AssertInt(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				return NewFloat(AssertFloat(a).Value()+AssertFloat(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsInt(a) {
				return NewInt(-AssertInt(a).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(-AssertFloat(a).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				return NewInt(AssertInt(a).Value()-AssertInt(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				return NewFloat(AssertFloat(a).Value()-AssertFloat(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "*",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsInt(a) && IsInt(b) {
				return NewInt(AssertInt(a).Value()*AssertInt(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewIntType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "*",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				return NewFloat(AssertFloat(a).Value()*AssertFloat(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "/",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				return NewFloat(AssertFloat(a).Value()/AssertFloat(b).Value(), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "sqrt",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Sqrt(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "sin",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Sin(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "cos",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Cos(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "tan",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Tan(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "asin",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Asin(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "acos",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Acos(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "atan",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Atan(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "exp",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Exp(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "log",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if IsFloat(a) {
				return NewFloat(math.Log(AssertFloat(a).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a}, NewFloatType(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "pow",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if IsFloat(a) && IsFloat(b) {
				return NewFloat(math.Pow(AssertFloat(a).Value(), AssertFloat(b).Value()), self.ctx)
			} else {
				return self.UpdateArgs([]Value{a, b}, NewFloatType(), self.ctx)
			}
		},
	},
}
