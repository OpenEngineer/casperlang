package main

var builtinJSONFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "true",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewTrueValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "false",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewFalseValue(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "null",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewNothingValue(self.ctx)
		},
	},
}
