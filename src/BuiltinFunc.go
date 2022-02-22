package main

type BuiltinFunc struct {
	ValueData
	name        string
	argPatterns []Pattern
	eval        func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value
}

func NewBuiltinFunc(cfg BuiltinFuncConfig) *BuiltinFunc {
	name := cfg.Name

	argPatterns := cfg.ArgPatterns
	if argPatterns == nil {
		argPatterns = []Pattern{}

		for _, argType := range cfg.ArgTypes {
			argPatterns = append(argPatterns, NewSimplePattern(NewBuiltinWord(argType)))
		}
	}

	eval := cfg.Eval

	return &BuiltinFunc{newValueData(nil, NewBuiltinContext()), name, argPatterns, eval}
}

func (v *BuiltinFunc) Update(type_ Type, ctx Context) Value {
	return &BuiltinFunc{newValueData(type_, ctx), v.name, v.argPatterns, v.eval}
}

func (f *BuiltinFunc) Name() string {
	return f.name
}

func (f *BuiltinFunc) CalcDistance(args []Value) []int {
	dummyHead := FuncHeader{NewBuiltinWord(f.name), f.argPatterns}

	return dummyHead.CalcDistance(args)
}

func (f *BuiltinFunc) DumpHead() string {
	dummyHead := FuncHeader{NewBuiltinWord(f.name), f.argPatterns}

	return dummyHead.Dump()
}

func (f *BuiltinFunc) NumArgs() int {
	return len(f.argPatterns)
}

func (f *BuiltinFunc) Dump() string {
	return f.Name() + " <builtin>"
}

func (f *BuiltinFunc) Type() Type {
	return NewFuncType(f.NumArgs())
}

func (f *BuiltinFunc) Eval(scope Scope, ew ErrorWriter) Value {
	panic("builtinfunc must be called before  eval")
}

func (f *BuiltinFunc) Call(args []Value, argScope Scope, ctx Context, ew ErrorWriter) Value {
	self := NewBuiltinCall(f.name, args, nil, f.eval, ctx)

	return f.eval(self, argScope, ew)
}
