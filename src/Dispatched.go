package main

type Dispatched struct {
	ctx      Context
	fn       Func
	distance []int
	args     []Value

	stack *Stack
}

func NewDispatched(args []Value, ctx Context) *Dispatched {
	cpy := make([]Value, len(args))
	for i, arg := range args {
		cpy[i] = arg
	}

	return &Dispatched{ctx, nil, []int{}, args, NewStack()}
}

func (d *Dispatched) UpdateArg(i int, des *Destructured) {
	if des == nil {
		panic("des can't be nil")
	}

	d.args[i] = des.arg

	if !d.Failed() && !des.Failed() {
		d.distance = append(d.distance, des.distance...)
		d.stack.Extend(des.stack)
	} else {
		d.distance = nil
	}
}

func (d *Dispatched) SetFunc(fn Func) {
	if d.fn != nil {
		panic("fn already set")
	}

	d.fn = fn
}

func (d *Dispatched) Failed() bool {
	return d.distance == nil
}

// second bool return value is true if there is some ambiguity
func (d *Dispatched) BetterThan(other *Dispatched) (bool, bool) {
	if other == nil {
		return true, false
	} else if len(d.distance) > len(other.distance) {
		return true, false
	} else if len(d.distance) == len(other.distance) {
		eachLE := true
		someLT := false

		for j, a := range d.distance {
			if a < other.distance[j] {
				someLT = true
			} else if a > other.distance[j] {
				eachLE = false
			}
		}

		if someLT && eachLE {
			return true, false
		} else if !someLT && !eachLE {
			return false, false
		} else {
			return false, true
		}
	} else {
		return false, false
	}
}

func (d *Dispatched) Eval() Value {
	if d.Failed() {
		panic("can't eval failed dispatch")
	}

	return d.fn.EvalRhs(d)
}

// nil return value means ambiguous
func PickBest(ds []*Dispatched) *Dispatched {
	if len(ds) == 0 {
		panic("no functions to pick from")
	}

	var (
		best *Dispatched = nil
	)

	for _, d := range ds {
		isBetter, isAmbig := d.BetterThan(best)

		if isAmbig {
			return nil
		} else if isBetter {
			best = d
		}
	}

	return best
}
