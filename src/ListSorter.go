package main

// can't sort on list directly because they can't be mutated
type ListSorter struct {
	ctx   Context
	comp  *AnonFunc
	items []Value
	ew    ErrorWriter
}

func NewListSorter(lst *List, comp *AnonFunc, ew ErrorWriter, ctx Context) *ListSorter {
	return &ListSorter{ctx, comp, lst.Items(), ew}
}

func (s *ListSorter) Len() int {
	return len(s.items)
}

func (s *ListSorter) Less(i, j int) bool {
	if !s.ew.Empty() {
		return true
	}

	args := []Value{s.items[i], s.items[j]}

	res := s.comp.EvalRhs(args, s.ew)

	if !s.ew.Empty() || res == nil {
		return true
	}

	lt, ok := GetBoolValue(res, s.ew)
	if ok && s.ew.Empty() {
		return lt
	} else {
		if !ok {
			s.ew.Add(res.Context().Error("expected Bool, got " + res.Dump()))
		}
		return true
	}
}

func (s *ListSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
