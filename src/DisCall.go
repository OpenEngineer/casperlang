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

func badDispatchMessage(name string, args []Value, msg string, fns []Func) string {
	var b strings.Builder

	b.WriteString(msg)
	b.WriteString(" \"")
	b.WriteString(name)
	b.WriteString("\"\n  Want:\n    ")

	b.WriteString(name)
	for _, arg := range args {
		ew := NewErrorWriter()
		arg_ := EvalEager(arg, ew)
		if ew.Empty() {
			arg = EvalPretty(arg_)
		}
		b.WriteString(" ")
		b.WriteString(arg.Dump())
	}

	if fns != nil && len(fns) > 0 {
		b.WriteString("\n  Have:\n")

		for _, fn := range fns {
			b.WriteString("    ")
			b.WriteString(fn.DumpPrettyHead())
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *DisCall) badDispatchMessage(msg string, fns []Func) string {
	return badDispatchMessage(t.Name(), t.Args(), msg, fns)
}

func (t *DisCall) SubVars(stack *Stack) Value {
	return &DisCall{t.CallData.subArgVars(stack), t.fns}
}

func (t *DisCall) dispatch(ew ErrorWriter) *Dispatched {
	ds := []*Dispatched{}

	if !ew.Empty() {
		return nil
	}

	args := t.Args() // make a copy!
	for _, fn := range t.fns {
		if fn.NumArgs() == t.NumArgs() {
			d := fn.Dispatch(args, ew)
			if !ew.Empty() {
				return nil
			} else if !d.Failed() {
				for _, arg := range d.args {
					if arg == nil {
						panic("returned args can't be nil")
					}
				}

				args = d.args // update the args as much as possible

				ds = append(ds, d)
			}
		}
	}

	if len(ds) == 0 {
		ew.Add(t.Context().Error(t.badDispatchMessage("unable to dispatch", t.fns)))
		return nil
	}

	best := PickBest(ds)

	if !ew.Empty() {
		return nil
	} else if best == nil {
		ew.Add(t.Context().Error(t.badDispatchMessage("ambiguous dispatch", t.fns)))
		return nil
	}

	return best
}

func (t *DisCall) Eval(ew ErrorWriter) Value {
	d := t.dispatch(ew)
	if d == nil || d.Failed() {
		return nil
	} else {
		//fmt.Println(t.Name())
		d.ctx = t.Context()
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
