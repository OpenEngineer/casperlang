package main

import (
	"errors"

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
	ioc := NewReplIOContext()
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
	case ContainsSymbolBefore(ts[1:], "=", ";"):
		for len(ts) > 0 && IsSymbol(ts[len(ts)-1], ";") {
			ts = ts[0 : len(ts)-1]
		}

		if IsOperatorSymbol(ts[0]) || (IsWord(ts[0]) && !IsSymbol(ts[1], "::")) {
			fn := parseFunc(ts, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			}

			out_ := ""
			if fn.NumArgs() == 0 {
				val, out := r.evalReplValue(fn.body, ioc, ew)
				if !ew.Empty() {
					return ew.Dump(), line
				} else {
					if val == nil {
						return "rhs returns nil", line
					} else {
						fn.body = val
						out_ = out
					}
				}
			}

			r.f.AddFunc(fn)
			fillPackage(r.p, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			} else {
				return out_, line
			}
		} else if !ContainsSymbol(ts, ";") {
			pat, rhs := parseDestructure(ts, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			}

			pat = pat.Link(NewFuncScope(r.f), ew)
			if !ew.Empty() {
				return ew.Dump(), line
			}

			var out string
			rhs, out = r.evalReplValue(rhs, ioc, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			} else if rhs == nil {
				return "rhs is nil", line
			}

			d := pat.Destructure(rhs, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			} else if d.Failed() {
				return "couldn't destructure", line
			}

			for i, var_ := range d.stack.vars {
				data := d.stack.data[i]
				name := var_.Name()

				fn := NewUserFunc(NewWord(name, var_.Context()), []Pattern{}, data, var_.Context())
				r.f.AddFunc(fn)
				fillPackage(r.p, ew)
				if !ew.Empty() {
					return ew.Dump(), line
				}
			}

			return out, line
		} else {
			val := parseStatements(ts, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			}

			var out string
			val, out = r.evalReplValue(val, ioc, ew)
			if !ew.Empty() {
				return ew.Dump(), line
			}

			if val != nil {
				if len(out) > 0 {
					out += "\n"
				}
				out += val.Dump()
			}

			// should parse as an anonymous function
			return out, line
		}
	default:
		// expect a regular expression
		val := ParseExpr(ts, ew)
		if !ew.Empty() {
			return ew.Dump(), line
		}

		val, out := r.evalReplValue(val, ioc, ew)
		if !ew.Empty() {
			return ew.Dump(), line
		} else {
			if val != nil {
				if len(out) > 0 {
					out += "\n"
				}
				out += val.Dump()
			}

			return out, line
		}
	}
}

func (r *Repl) evalReplValue(val Value, ioc *ReplIOContext, ew ErrorWriter) (Value, string) {
	IO_CONTEXT = ioc

	val = val.Link(r.f, ew)
	if !ew.Empty() {
		return nil, ""
	}

	val = EvalEager(val, ew)
	if !ew.Empty() {
		return nil, ""
	} else if ioc.panicMsg != "" {
		ew.Add(errors.New(ioc.panicMsg))
		return nil, ""
	}

	if io, ok := val.(*IO); ok {
		res := io.Run(ioc)

		return res, ioc.StdoutString()
	} else {
		return prettyValue(val), ""
	}
}

func prettyValue(val Value) Value {
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

	return val
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
