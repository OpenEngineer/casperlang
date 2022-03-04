package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Error struct {
	ctx Context
	end bool
	msg string
}

// Doing this in two steps incurs some overhead
func IsError(e error) bool {
	_, ok := e.(*Error)
	return ok
}

func AssertError(e_ error) *Error {
	e, ok := e_.(*Error)
	if ok {
		return e
	} else {
		panic("expected *Error")
	}
}

func (e *Error) Error() string {
	if e.ctx.src.Path() == "<stdin>" {
		return e.ctx.start.ToString() + ": " + e.msg
	} else if strings.HasPrefix(e.ctx.src.Path(), "<") {
		return writeError(e.ctx.src, e.ctx.end, e.msg, false)
	} else if e.end {
		return writeError(e.ctx.src, e.ctx.end, e.msg, true)
	} else {
		return writeError(e.ctx.src, e.ctx.start, e.msg, true)
	}
}

func abbreviatePath(p string) string {
	cwd, err := os.Getwd()
	if err == nil {
		alt, err := filepath.Rel(cwd, p)
		if err == nil && len(alt) < len(p) && !strings.HasPrefix(alt, "../../") {
			return alt
		} else {
			return p
		}
	} else {
		return filepath.Base(p)
	}
}

func writeError(src *Source, fp FilePos, msg string, writePos bool) string {
	var b strings.Builder

	b.WriteString(abbreviatePath(src.path))

	if writePos {
		b.WriteString(":")
		b.WriteString(fp.ToString())
	}

	if len(msg) > 0 {
		b.WriteString(": ")
		b.WriteString(msg)
	}

	return b.String()
}
