package main

type BuiltinFuncConfig struct {
	Name        string
	ArgTypes    []string
	ArgPatterns []Pattern
	LinkReqs    []string
	Eval        func(self *BuiltinCall, ew ErrorWriter) Value
}
