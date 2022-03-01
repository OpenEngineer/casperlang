package main

import "math"

var builtinMathFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()+AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "+",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()+AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(-AssertInt(self.args[0]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(-AssertFloat(self.args[0]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()-AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "-",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()-AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "*",
		ArgTypes: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()*AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "*",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()*AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "/",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()/AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "sqrt",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Sqrt(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "sin",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Sin(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "cos",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Cos(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "tan",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Tan(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "asin",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Asin(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "acos",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Acos(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "atan",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Atan(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "exp",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Exp(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "log",
		ArgTypes: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Log(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "pow",
		ArgTypes: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Pow(AssertFloat(self.args[0]).Value(), AssertFloat(self.args[1]).Value()), self.ctx)
		},
	},
}
