package main

import "fmt"

type GlobalScope struct {
	entries map[string][]DispatchableFunc
}

func (s *GlobalScope) Parent() Scope {
	return nil
}

func (s *GlobalScope) CollectFunctions(name string) []DispatchableFunc {
	fns, ok := s.entries[name]
	if ok {
		return fns
	} else {
		return []DispatchableFunc{}
	}
}

func (s *GlobalScope) Dispatch(name *Word, args []Value, ew ErrorWriter) Func {
	fns := s.CollectFunctions(name.Value())
	if len(fns) == 0 {
		fmt.Println("why are we here?")
		ew.Add(name.Context().Error("\"" + name.Value() + "\" undefined"))
		return nil
	}

	best, err := PickBest(fns, args, name.Context())

	if err != nil {
		ew.Add(err)
		return nil
	} else if best == nil {
		return nil
	} else {
		return best
	}
}

func registerBuiltinFuncs(db map[string][]DispatchableFunc, fns []BuiltinFuncConfig) {
	for _, fnCfg := range fns {
		name := fnCfg.Name
		lst, ok := db[name]
		if !ok {
			lst = []DispatchableFunc{}
		}

		lst = append(lst, NewBuiltinFunc(fnCfg))

		db[name] = lst
	}
}

func registerUserFuncs(db map[string][]DispatchableFunc, fns []*UserFunc) {
	for _, fn := range fns {
		name := fn.Name()

		lst, ok := db[name]
		if !ok {
			lst = []DispatchableFunc{}
		}

		lst = append(lst, fn)

		db[name] = lst
	}
}

func fillCoreDB() map[string][]DispatchableFunc {
	db := make(map[string][]DispatchableFunc)
	registerBuiltinFuncs(db, builtinCoreFuncs)
	registerBuiltinFuncs(db, builtinIOFuncs)
	registerBuiltinFuncs(db, builtinMathFuncs)
	return db
}

func FillGlobalScope() Scope {
	db := make(map[string][]DispatchableFunc)

	registerBuiltinFuncs(db, builtinCoreFuncs)
	registerBuiltinFuncs(db, builtinIOFuncs)
	registerBuiltinFuncs(db, builtinMathFuncs)

	return &GlobalScope{db}
}

func fillJsonDB() map[string][]DispatchableFunc {
	db := make(map[string][]DispatchableFunc)
	registerBuiltinFuncs(db, builtinJSONFuncs)
	return db
}

func FillJSONScope() Scope {
	db := make(map[string][]DispatchableFunc)

	registerBuiltinFuncs(db, builtinJSONFuncs)

	return &GlobalScope{db}
}
