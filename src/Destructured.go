package main

type Destructured struct {
	arg      Value
	distance []int
	vars     []*Variable
	data     []Value
}

func NewDestructured(arg Value, distance []int) *Destructured {
	return &Destructured{arg, distance, []*Variable{}, []Value{}}
}

func (d *Destructured) AddVar(v *Variable, data Value) {
	d.vars = append(d.vars, v)
	d.data = append(d.data, data)
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
