package main

import (
	"errors"

	"github.com/openengineer/go-repl"
)

type ReplMode int

type Repl struct {
	p *Package
	f *File
	r *repl.Repl
}

func NewRepl(ew ErrorWriter) *Repl {
	p := LoadReplPackage(ew)
	if p == nil {
		return nil
	}

	f := p.modules[""].files[0]
	f.GetModules(p, []*Module{}, ew)
	if !ew.Empty() {
		return nil
	} else {
		p.RegisterFuncs(ew)
		return &Repl{p, f, nil}
	}
}

func (r *Repl) Prompt() string {
	return "> "
}

func (r *Repl) Tab(buffer string) string {
	return ""
}

func (r *Repl) Eval(line string) string {
	// first tokenize always
	ew := NewErrorWriter()
	s := NewSource("<stdin>", []byte(line))

	ts := Tokenize(s, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	ts = RemoveNLs(ts)

	if len(ts) == 0 {
		return ""
	}

	switch {
	case IsWord(ts[0], "import"):
		return r.evalImports(ts, ew)
	case ContainsSymbolBefore(ts[1:], "=", ";"):
		// regular function definition
		if IsOperatorSymbol(ts[0]) || (IsWord(ts[0]) && !IsSymbol(ts[1], "::")) {
			return r.defFunc(ts, ew)
		} else { // a destructure split at the first =
			return r.evalDestructure(ts, ew)
		}
	default:
		return r.evalExpr(ts, ew)
	}
}

func (r *Repl) linkAndEval(val Value, ew ErrorWriter) Value {
	val = val.Link(r.f, ew)
	if !ew.Empty() {
		return nil
	}

	val = EvalEager(val, ew)
	if !ew.Empty() {
		return nil
	} else if PANIC != "" {
		ew.Add(errors.New(PANIC))
		PANIC = ""
		return nil
	}

	return EvalPretty(val)
}

func (r *Repl) evalImports(ts []Token, ew ErrorWriter) string {
	strs := []*String{}

	for _, t := range ts[1:] {
		if str, ok := t.(*String); ok {
			strs = append(strs, str)
		} else {
			return t.Context().Error("invalid import statement").Error()
		}
	}

	return r.addImports(strs, ew)
}

func (r *Repl) addImports(paths []*String, ew ErrorWriter) string {
	for _, path := range paths {
		r.f.AddImport(path)
	}

	r.f.GetModules(r.p, []*Module{}, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	r.p.RegisterFuncs(ew)
	if !ew.Empty() {
		return ew.Dump()
	} else {
		return ""
	}
}

func (r *Repl) defFunc(ts []Token, ew ErrorWriter) string {
	fn := parseFunc(ts, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	headHash := fn.DumpPrettyHead()
	for _, check := range r.f.fns {
		if check.Name() == fn.Name() && check.NumArgs() == fn.NumArgs() && check.DumpPrettyHead() == headHash {
			ew.Add(errors.New("\"" + headHash + "\" already defined"))
			return ew.Dump()
		}
	}

	r.f.AddFunc(fn)
	r.p.RegisterFuncs(ew) // register everything in to have recursive funcs

	n := len(r.f.fns)
	// force linking now, not upon call
	if ew.Empty() {
		r.f.fns[n-1].fn = fn.Link(r.f, ew)
	}

	if !ew.Empty() {
		r.f.fns = r.f.fns[0 : n-1] // remove the last one
		r.p.RegisterFuncs(ew)
		return ew.Dump()
	}

	return ""
}

func (r *Repl) evalDestructure(ts []Token, ew ErrorWriter) string {
	pat, rhs := parseDestructure(ts, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	pat = pat.Link(NewFuncScope(r.f), ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	rhs = r.linkAndEval(rhs, ew)
	if !ew.Empty() {
		return ew.Dump()
	} else if rhs == nil {
		panic("unexpected")
	}

	d := pat.Destructure(rhs, ew)
	if !ew.Empty() {
		return ew.Dump()
	} else if d.Failed() {
		return "couldn't destructure"
	}

	for _, var_ := range d.stack.vars {
		name := var_.Name()

		for _, check := range r.f.fns {
			if check.Name() == name && check.NumArgs() == 0 {
				ew.Add(errors.New("\"" + name + "\" already defined"))
				return ew.Dump()
			}
		}
	}

	for i, var_ := range d.stack.vars {
		data := d.stack.data[i]
		name := var_.Name()

		fn := NewUserFunc(NewWord(name, var_.Context()), []Pattern{}, data, var_.Context())
		r.f.AddFunc(fn)
		r.p.RegisterFuncs(ew)
		if !ew.Empty() {
			return ew.Dump()
		}
	}

	return ""
}

func (r *Repl) evalExpr(ts []Token, ew ErrorWriter) string {
	// expect a regular expression
	val := ParseExpr(ts, ew)
	if !ew.Empty() {
		return ew.Dump()
	}

	out := ""
	val = r.linkAndEval(val, ew)
	if !ew.Empty() {
		return ew.Dump()
	} else {
		if io, ok := val.(*IO); ok {
			ioc := NewReplIOContext(r)
			val = io.Run(ioc)
			out += ioc.StdoutString()
		}

		if val != nil {
			if len(out) > 0 {
				out += "\n"
			}
			out += val.Dump()
		}

		return out
	}
}

func (r *Repl) Quit() {
	r.r.Quit()
}
