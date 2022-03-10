package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"
)

type IOContext interface {
	Stdout() io.Writer
	ReadLine(echo bool) string
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

func (c *DefaultIOContext) ReadLine(echo bool) string {
	if echo {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		return scanner.Text()
	} else {
		txt, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return ""
		}

		return string(txt)
	}
}

func (c *ReplIOContext) Stdout() io.Writer {
	return c.stdout
}

func (c *ReplIOContext) ReadLine(echo bool) string {
	return c.repl.r.ReadLine(echo)
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
