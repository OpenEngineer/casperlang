package main

import "strconv"

type BuiltinFunc struct {
	ValueData
	name        string
	argPatterns []Pattern
	linkReqs    []string
	links       map[string][]Func
	eval        func(self *BuiltinCall, ew ErrorWriter) Value
}

func NewBuiltinFunc(cfg BuiltinFuncConfig) *BuiltinFunc {
	name := cfg.Name

	argPatterns := cfg.ArgPatterns
	if argPatterns == nil {
		argPatterns = []Pattern{}

		for _, argType := range cfg.ArgTypes {
			argPatterns = append(argPatterns, NewTypePattern(NewBuiltinWord(argType)))
		}
	}

	eval := cfg.Eval

	return &BuiltinFunc{
		newValueData(NewBuiltinContext()),
		name,
		argPatterns,
		cfg.LinkReqs,
		make(map[string][]Func),
		eval,
	}
}

func (f *BuiltinFunc) TypeName() string {
	return "\\" + strconv.Itoa(f.NumArgs())
}

func (f *BuiltinFunc) Dump() string {
	return f.Name() + " <builtin>"
}

func (f *BuiltinFunc) Name() string {
	return f.name
}

func (f *BuiltinFunc) NumArgs() int {
	return len(f.argPatterns)
}

func (f *BuiltinFunc) header() *FuncHeader {
	return &FuncHeader{NewBuiltinWord(f.name), f.argPatterns}
}

func (f *BuiltinFunc) DumpHead() string {
	head := f.header()

	return head.Dump()
}

func (f *BuiltinFunc) ListHeaderTypes() []string {
	return []string{}
}

func (f *BuiltinFunc) Link(scope Scope, ew ErrorWriter) Value {
	// opportunity to get some constructors
	for _, k := range f.linkReqs {
		fns := scope.ListDispatchable(k, -1, ew)
		if len(fns) == 0 {
			ew.Add(f.Context().Error("\"" + f.Name() + "\" undefined"))
		}

		f.links[k] = fns
	}

	return f
}

func (f *BuiltinFunc) Dispatch(args []Value, ew ErrorWriter) *Dispatched {
	head := f.header()

	d := head.Destructure(args, ew)
	if d == nil {
		return nil
	}

	d.SetFunc(f)
	return d
}

// no: detach as regular, and get args from FuncScope
func (f *BuiltinFunc) EvalRhs(d *Dispatched) Value {
	return NewBuiltinCall(f.name, d.args, f.links, f.eval, d.ctx)
}
