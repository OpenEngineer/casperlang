package main

type GlobalScope struct {
	db     map[string][]Func
	linker *Linker
}

func (s *GlobalScope) GetLocal(name string) *Variable {
	return nil
}

func (s *GlobalScope) ListDispatchable(name string, numArgs int, ew ErrorWriter) []Func {
	fns_, ok := s.db[name]
	if ok {
		fns := []Func{}

		for _, fn_ := range fns_ {
			if fn_.NumArgs() == numArgs || numArgs == -1 {
				fn := s.linker.LinkFunc(fn_, s, ew)
				fns = append(fns, fn)
			}
		}

		return fns
	} else {
		return []Func{}
	}
}

func registerBuiltinFuncs(db map[string][]Func, fns []BuiltinFuncConfig) {
	for _, fnCfg := range fns {
		if fnCfg.Allowed() {
			name := fnCfg.Name
			lst, ok := db[name]
			if !ok {
				lst = []Func{}
			}

			lst = append(lst, NewBuiltinFunc(fnCfg))

			db[name] = lst
		}
	}
}

func registerUserFuncs(db map[string][]Func, fns []*UserFunc) {
	for _, fn := range fns {
		name := fn.Name()

		lst, ok := db[name]
		if !ok {
			lst = []Func{}
		}

		lst = append(lst, fn)

		db[name] = lst
	}
}

func fillCoreDB() map[string][]Func {
	db := make(map[string][]Func)

	registerBuiltinFuncs(db, builtinCoreFuncs)
	registerBuiltinFuncs(db, builtinIOFuncs)
	registerBuiltinFuncs(db, builtinMathFuncs)
	registerBuiltinFuncs(db, builtinReplFuncs)

	return db
}

func NewGlobalScope(linker *Linker) *GlobalScope {
	return &GlobalScope{fillCoreDB(), linker}
}
