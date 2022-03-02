package main

type BuiltinFuncConfig struct {
	Name     string
	Args     []string
	LinkReqs []string
	Eval     func(self *BuiltinCall, ew ErrorWriter) Value
}
