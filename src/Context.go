package main

import (
	"path/filepath"
	"strings"
)

type Context struct {
	src   *Source
	start FilePos
	end   FilePos // exclusive
}

func NewBuiltinContext() Context {
	return NewSpecialContext("<builtin>")
}

func NewStdinContext() Context {
	return NewSpecialContext("<stdin>")
}

func NewSpecialContext(name string) Context {
	return Context{
		&Source{name, []rune{}},
		FilePos{},
		FilePos{},
	}
}

func (ctx Context) Path() string {
	return ctx.src.Path()
}

func (ctx Context) Dir() string {
	return filepath.Dir(ctx.Path())
}

func (ctx Context) Error(msg string) error {
	return &Error{ctx, false, msg}
}

func (ctx Context) EndError(msg string) error {
	return &Error{ctx, true, msg}
}

func (ctx Context) WriteLocation() string {
	return ctx.Path() + ":" + ctx.start.ToString()
}

func MergeContexts(a Context, b Context) Context {
	start := a.start
	end := b.end

	return Context{a.src, start, end}
}

func (a Context) Merge(b Context) Context {
	return MergeContexts(a, b)
}

func (a Context) Before(b Context) bool {
	pathComp := strings.Compare(a.src.path, b.src.path)

	if pathComp < 0 {
		return true
	} else if pathComp > 0 {
		return false
	} else {
		return a.start.Before(b.start)
	}
}
