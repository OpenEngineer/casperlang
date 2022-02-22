package main

type BuiltinFuncConfig struct {
	Name        string
	ArgTypes    []string
	ArgPatterns []Pattern
	Eval        func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value
}
