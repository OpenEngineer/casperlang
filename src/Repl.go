package main

import (
	"strings"

	"github.com/openengineer/go-terminal"
)

type ReplMode int

type Repl struct {
	p *Package
	f *File
	t *terminal.Terminal
}

func NewRepl(ew ErrorWriter) *Repl {
	p := LoadReplPackage(ew)
	if p == nil {
		return nil
	}

	f := p.modules[""].files[0]

	return &Repl{p, f, nil}
}

func (r *Repl) RegisterTerm(t *terminal.Terminal) {
	r.t = t
}

func (r *Repl) Eval(line string) (string, string) {
	// first tokenize always
	ew := NewErrorWriter()
	s := NewSource("<stdin>", []byte(line))

	ts := Tokenize(s, ew)
	if !ew.Empty() {
		return ew.Dump(), line
	}

	ts = RemoveNLs(ts)

	if len(ts) == 0 {
		return "", line
	}

	switch {
	case IsWord(ts[0], "import"):
		strs := []*String{}

		for _, t := range ts[1:] {
			if str, ok := t.(*String); ok {
				strs = append(strs, str)
			} else {
				return t.Context().Error("invalid import statement").Error(), line
			}
		}

		return r.evalImports(strs), line
	case strings.Contains(line, "="):
		return "assignments not yet implemented", line
	default:
		// expect a regular expression
		ew := NewErrorWriter()

		val := ParseExpr(ts, ew)
		if !ew.Empty() {
			return ew.Dump(), line
		}

		val = val.Link(r.f, ew)
		if !ew.Empty() {
			return ew.Dump(), line
		}

		val = EvalEager(val, ew)
		if !ew.Empty() {
			return ew.Dump(), line
		}

		if io, ok := val.(*IO); ok {
			ioc := NewReplIOContext()
			res := io.Run(ioc)

			out := ioc.StdoutString()
			if res != nil {
				out += "\n" + res.Dump()
			}
			return out, line
		} else {
			return val.Dump(), line
		}
	}
}

func (r *Repl) evalImports(paths []*String) string {
	return "import not yet implemented"
}
