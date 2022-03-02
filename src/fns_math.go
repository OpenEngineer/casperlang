package main

import "math"

var builtinMathFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name: "+",
		Args: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()+AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "+",
		Args: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()+AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "-",
		Args: []string{"Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(-AssertInt(self.args[0]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "-",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(-AssertFloat(self.args[0]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "-",
		Args: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()-AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "-",
		Args: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()-AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "*",
		Args: []string{"Int", "Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(AssertInt(self.args[0]).Value()*AssertInt(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "*",
		Args: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()*AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "/",
		Args: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(AssertFloat(self.args[0]).Value()/AssertFloat(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "sqrt",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Sqrt(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "sin",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Sin(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "cos",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Cos(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "tan",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Tan(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "asin",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Asin(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "acos",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Acos(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "atan",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Atan(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "exp",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Exp(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "log",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Log(AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "pow",
		Args: []string{"Float", "Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewFloat(math.Pow(AssertFloat(self.args[0]).Value(), AssertFloat(self.args[1]).Value()), self.ctx)
		},
	},
}
