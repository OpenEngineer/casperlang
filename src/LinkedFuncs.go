package main

type Linker struct {
	entry  Func          // linked entry point
	linked map[Func]Func // map from unlinked to linked
}

func NewLinker() *Linker {
	return &Linker{nil, make(map[Func]Func)}
}

func (s *Linker) LinkFunc(fn Func, scope Scope, ew ErrorWriter) Func {
	if l, ok := s.linked[fn]; ok {
		return l
	} else {
		// add something temporary to the list (for recursive links)
		tmp := NewWrappedFunc(nil)

		s.linked[fn] = tmp

		l_ := fn.Link(scope, ew) // recursion might happen here

		l, ok := l_.(Func)
		if !ok {
			panic("expected Func")
		}

		s.linked[fn] = l
		tmp.SetFunc(l)

		return l
	}
}
