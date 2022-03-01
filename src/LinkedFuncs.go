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
		l_ := fn.Link(scope, ew)

		l, ok := l_.(Func)
		if !ok {
			panic("expeccted Func")
		}

		s.linked[fn] = l

		return l
	}
}
