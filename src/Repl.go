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
		ioc := NewReplIOContext()
		IO_CONTEXT = ioc

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
		} else if ioc.panicMsg != "" {
			return ioc.panicMsg, line
		}

		if io, ok := val.(*IO); ok {
			res := io.Run(ioc)

			out := ioc.StdoutString()
			if res != nil {
				if len(out) > 0 {
					out += "\n"
				}
				out += res.Dump()
			}
			return out, line
		} else {
			cs := val.Constructors()

			for len(cs) > 0 {
				n := len(cs)

				if val.TypeName() == "Any" || val.TypeName() == "Maybe" || val.TypeName() == "Path" {
					val = cs[n-1]
					cs = cs[0 : n-1]
				} else {
					break
				}
			}

			return val.Dump(), line
		}
	}
}

func (r *Repl) evalImports(paths []*String) string {
	for _, path := range paths {
		r.f.AddImport(path)
	}

	ew := NewErrorWriter()

	r.f.GetModules(r.p, []*Module{}, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	fillPackage(r.p, ew)

	if !ew.Empty() {
		return ew.Dump()
	} else {
		return ""
	}
}
