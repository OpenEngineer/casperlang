package main

type Destructured struct {
	arg      Value
	distance []int

	stack *Stack
}

func NewDestructured(arg Value, distance []int) *Destructured {
	return &Destructured{arg, distance, NewStack()}
}

func NewDestructuredS(arg Value, distance []int, stack *Stack) *Destructured {
	return &Destructured{arg, distance, stack}
}

func (d *Destructured) AddVar(v *Variable, data Value) {
	d.stack.Set(v, data)
}

func (d *Destructured) Failed() bool {
	return d.distance == nil || d.arg == nil
}

func WorstDistance(ds []*Destructured) []int {
	worst := ds[0].distance

	for i := 1; i < len(ds); i++ {
		worst = worstDistance(ds[i].distance, worst)
	}

	return worst
}
