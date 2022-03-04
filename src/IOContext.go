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
}

func NewDefaultIOContext() *DefaultIOContext {
	return &DefaultIOContext{
		stdout: os.Stdout,
	}
}

func NewReplIOContext() *ReplIOContext {
	return &ReplIOContext{
		stdout: &bytes.Buffer{},
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
