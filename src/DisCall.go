package main

import (
	"strings"
)

type DisCall struct {
	CallData
	fns []Func
}

func NewDisCall(fns []Func, args []Value, ctx Context) *DisCall {
	return &DisCall{newCallData(args, ctx), fns}
}

func (t *DisCall) Name() string {
	return t.fns[0].Name()
}

func (t *DisCall) TypeName() string {
	if isConstructorName(t.Name()) {
		return t.Name()
	} else {
		return ""
	}
}

func (t *DisCall) Dump() string {
	return t.CallData.dump(t.Name())
}

func (v *DisCall) SetConstructors(cs []Call) Value {
	return v.setConstructors(cs)
}

func (v *DisCall) setConstructors(cs []Call) Call {
	return &DisCall{v.CallData.setConstructors(cs), v.fns}
}

func (t *DisCall) Link(scope Scope, ew ErrorWriter) Value {
	return t
}

func (t *DisCall) badDispatchMessage(ds []*Dispatched) string {
	var b strings.Builder

	b.WriteString("unable to dispatch \"")
	b.WriteString(t.Name())
	b.WriteString("\" ")
	b.WriteString(t.Dump())

	if ds != nil && len(ds) > 0 {
		b.WriteString("\n  Have:\n")

		for _, d := range ds {
			b.WriteString("    ")
			b.WriteString(d.fn.DumpHead())
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *DisCall) dispatch(ew ErrorWriter) *Dispatched {
	ds := []*Dispatched{}

	args := t.args
	for _, fn := range t.fns {
		d := fn.Dispatch(args, ew)
		if !ew.Empty() {
			return nil
		} else if d != nil {
			args = d.args // update the args as much as possible
			if !d.Failed() {
				ds = append(ds, d)
			}
		}
	}

	if len(ds) == 0 {
		ew.Add(t.Context().Error(t.badDispatchMessage(nil)))
		return nil
	}

	best := PickBest(ds)

	if best == nil {
		ew.Add(t.Context().Error(t.badDispatchMessage(ds)))
		return nil
	}

	return best
}

func (t *DisCall) Eval(ew ErrorWriter) Value {
	d := t.dispatch(ew)
	if d == nil {
		return nil
	} else {
		v := d.Eval()

		cs := t.Constructors()
		if isConstructorName(t.Name()) {
			cs = append(make([]Call, 0), cs...)
			cs = append(cs, t)
		}

		v = v.SetConstructors(cs)
		return v
	}
}
