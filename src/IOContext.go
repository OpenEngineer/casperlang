package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var IO_CONTEXT IOContext = nil

type IOContext interface {
	Stdout() io.Writer
	Panic(msg string)
}

type DefaultIOContext struct {
	stdout *os.File
}

type ReplIOContext struct {
	stdout   *bytes.Buffer
	panicMsg string
}

func NewDefaultIOContext() *DefaultIOContext {
	return &DefaultIOContext{
		stdout: os.Stdout,
	}
}

func NewReplIOContext() *ReplIOContext {
	return &ReplIOContext{
		stdout:   &bytes.Buffer{},
		panicMsg: "",
	}
}

func (c *DefaultIOContext) Stdout() io.Writer {
	return c.stdout
}

func (c *DefaultIOContext) Panic(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func (c *ReplIOContext) Stdout() io.Writer {
	return c.stdout
}

func (c *ReplIOContext) StdoutString() string {
	if c.panicMsg != "" {
		return c.panicMsg
	} else {
		return string(c.stdout.Bytes())
	}
}

func (c *ReplIOContext) Panic(msg string) {
	c.panicMsg = msg
}
