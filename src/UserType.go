package main

import "strings"

type UserType struct {
	ctx    Context
	parent Type
	name   string
	args   []Value
}

func NewUserType(parent Type, name string, args []Value, ctx Context) *UserType {
	if parent == nil {
		panic("parent can't be nil")
	}

	return &UserType{ctx, parent, name, args}
}

func AssertUserType(t_ Type, name string) *UserType {
	t, ok := t_.(*UserType)
	if ok && t.name == name {
		return t
	} else if t_ == nil {
		panic("expected *UserType named \"" + name + "\"")
	} else {
		return AssertUserType(t_, name)
	}
}

func (t *UserType) Parent() Type {
	return t.parent
}

func (t *UserType) Dump() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(t.name)

	for _, arg := range t.args {
		b.WriteString(" ")
		b.WriteString(arg.Type().Dump())
	}

	b.WriteString(")")

	return b.String()
}

func (t *UserType) CalcNameDistance(name *Word) int {
	if name.Value() == t.name {
		return 0
	} else {
		d := t.parent.CalcNameDistance(name)
		if d < 0 {
			return -1
		} else {
			return d + 1
		}
	}
}

func (t *UserType) CalcConstructorDistance(name *Word, args []Pattern) []int {
	var d []int = nil

	if t.name == name.Value() && len(args) == len(t.args) {
		d = []int{0}

		for i, arg := range args {
			argD := arg.CalcDistance(t.args[i])
			if argD == nil {
				d = nil
				break
			}

			d = append(d, argD...)
		}
	}

	if d == nil {
		d = t.parent.CalcConstructorDistance(name, args)

		if d == nil {
			return nil
		} else {
			d[0] += 1
			return d
		}
	} else {
		return d
	}
}

func (t *UserType) CalcTupleDistance(args []Pattern) []int {
	d := t.parent.CalcTupleDistance(args)

	if d == nil {
		return nil
	}

	d[0] += 1

	return d
}

func (t *UserType) CalcStructDistance(keys []*String, vals []Pattern) []int {
	d := t.parent.CalcStructDistance(keys, vals)

	if d == nil {
		return nil
	}

	d[0] += 1

	return d
}

// assume pattern match has already been checked so we can continue without errors
func (t *UserType) DestructureConstructor(name *Word, args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	for i, arg := range args {
		scope = arg.Destructure(t.args[i], scope, ew)
	}

	return scope
}

func (t *UserType) DestructureTuple(args []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}

func (t *UserType) DestructureStruct(keys []*String, vals []Pattern, scope *FuncScope, ew ErrorWriter) *FuncScope {
	return scope
}
