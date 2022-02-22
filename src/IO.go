package main

import (
	"fmt"
	"os"
	"reflect"
)

type IO struct {
	TokenData
	isDefErr bool
	type_    *IOType // XXX: is the data type in *IOType still necessary?, can we just use PrimType("IO") instead?
	Run      func() Value
}

func NewIO(innerType Value, Run func() Value, ctx Context) *IO {
	return &IO{newTokenData(ctx), false, NewIOType(innerType), Run}
}

// the inner one should throw the first error
func NewDeferredError(msg string, ctx Context) *IO {
	return &IO{
		newTokenData(ctx),
		true,
		nil,
		func() Value {
			fmt.Fprintf(os.Stdout, "%s\n", msg)
			os.Exit(1)
			return nil
		},
	}
}

func (v *IO) Dump() string {
	if v.isDefErr {
		return "<deferred-error>"
	} else {
		return v.type_.Dump()
	}
}

func (v *IO) Type() Type {
	return v.type_
}

func (v *IO) CheckTypeNames(scope Scope, ew ErrorWriter) {
}

func (v *IO) Eval(scope Scope, ew ErrorWriter) Value {
	return v
}

func (v *IO) Update(type_ Type, ctx Context) Value {
	return &IO{newTokenData(ctx), v.isDefErr, AssertIOType(type_), v.Run}
}

func IsIO(v_ Value) bool {
	_, ok := v_.(*IO)
	return ok
}

func AssertIO(v_ Value) *IO {
	v, ok := v_.(*IO)
	if ok {
		return v
	} else {
		panic("expected *IO, got " + reflect.TypeOf(v_).String())
	}
}

func IsDeferredError(v_ Value) bool {
	v, ok := v_.(*IO)
	if ok {
		return v.isDefErr
	} else {
		return false
	}
}
