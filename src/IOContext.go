package main

import (
	"bytes"
	"io"
	"os"
)

type IOContext interface {
	Stdout() io.Writer
}

type DefaultIOContext struct {
	stdout *os.File
}

type ReplIOContext struct {
	stdout *bytes.Buffer
	repl   *Repl
}

func NewDefaultIOContext() *DefaultIOContext {
	return &DefaultIOContext{
		stdout: os.Stdout,
	}
}

func NewReplIOContext(repl *Repl) *ReplIOContext {
	return &ReplIOContext{
		stdout: &bytes.Buffer{},
		repl:   repl,
	}
}

func (c *DefaultIOContext) Stdout() io.Writer {
	return c.stdout
}

func (c *ReplIOContext) Stdout() io.Writer {
	return c.stdout
}

func (c *ReplIOContext) StdoutString() string {
	return string(c.stdout.Bytes())
}

func (c *ReplIOContext) ListNames() []string {
	return c.repl.f.ListNames()
}

func AssertReplIOContext(ioc_ IOContext) *ReplIOContext {
	ioc, ok := ioc_.(*ReplIOContext)
	if ok {
		return ioc
	} else {
		panic("expected *ReplIOContext")
	}
}
