package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// a module must load all files at once, can't be done lazily
type Module struct {
	dir      string
	p        *Package
	files    []*File
	all      []*ScopedFunc
	exported []*ScopedFunc
}

func LoadModule(p *Package, consumers []*Module, dir *String, ew ErrorWriter) *Module {
	if DEBUG_PKG_LOADING {
		fmt.Println("loading module \"" + dir.Value() + "\"")
	}

	for _, consumer := range consumers {
		if consumer.Dir() == dir.Value() {
			ew.Add(errors.New("circular import detected"))
			return nil
		}
	}

	// parse all files ending with .cas in the given directory
	infos, err := ioutil.ReadDir(dir.Value())
	if err != nil {
		ew.Add(err)
		return nil
	}

	files := []*File{}

	for _, info := range infos {
		if filepath.Ext(info.Name()) == ".cas" {
			path := filepath.Join(dir.Value(), info.Name())
			f := ParseFile(path, ew)
			if !ew.Empty() {
				return nil
			}

			files = append(files, f)
		}
	}

	if len(files) == 0 {
		ew.Add(errors.New("module " + dir.Value() + " doesn't contain any source files"))
		return nil
	}

	m := &Module{dir.Value(), p, files, nil, nil}

	consumers = append(consumers, m)

	for _, f := range m.files {
		f.GetModules(p, consumers, ew)
	}

	return m
}

func (m *Module) Dir() string {
	return m.dir
}

func isExportedName(name string) bool {
	return !strings.HasPrefix(name, "_")
}

func (m *Module) MergeFuncs() {
	allFns := []*ScopedFunc{}

	for _, f := range m.files {
		allFns = append(allFns, f.fns...)

	}

	exported := []*ScopedFunc{}
	for _, fn := range allFns {
		if isExportedName(fn.Name()) {
			exported = append(exported, fn)
		}
	}

	for _, f := range m.files {
		// make a copy of the list, so it can be mutated by files without affecting other files
		cpy := make([]*ScopedFunc, len(allFns))
		for i, fn := range allFns {
			cpy[i] = fn
		}

		f.fns = cpy
	}

	m.all = allFns
	m.exported = exported
}

func (m *Module) GetExportedFuncs() []*ScopedFunc {
	return m.exported
}

func (m *Module) ImportFuncs(ew ErrorWriter) {
	for _, f := range m.files {
		f.ImportFuncs(ew)
	}
}

func (m *Module) GetDependencies() []*Module {
	deps := []*Module{}

	for _, f := range m.files {
		for _, imodule := range f.imodules {
			firstTime := true
			for _, check := range deps {
				if check == imodule {
					firstTime = false
				}
			}

			if firstTime {
				deps = append(deps, imodule)
			}
		}
	}

	return deps
}

// also detect non-unique constructors
func (m *Module) ListLocalTypes(ew ErrorWriter) []string {
	lst := []string{}

Outer:
	for _, fn := range m.all {
		if isConstructorName(fn.Name()) {
			for _, check := range lst {
				if check == fn.Name() {
					ew.Add(fn.Context().Error("multiple definitions of \"" + fn.Name() + "\""))
					continue Outer
				}
			}

			lst = append(lst, fn.Name())
		}
	}

	return lst
}

// assume a and b are sorted
func haveCommonEntries(a []string, b []string) bool {
	i := 0
	j := 0

	for {
		if i >= len(a) || j >= len(b) {
			return false
		}

		if a[i] == b[j] {
			return true
		} else if a[i] > b[j] {
			j++
		} else {
			i++
		}
	}
}

func (m *Module) SyncMethods(ew ErrorWriter) {
	local := sortUniqStrings(m.ListLocalTypes(ew))

	// only the exported methods attached to any Constructors can be pushed upwards
	pushed := []*ScopedFunc{}

	for _, fn := range m.exported {
		fnTypes := fn.ListHeaderTypes()

		if haveCommonEntries(local, fnTypes) {
			pushed = append(pushed, fn)
		}
	}

	deps := m.GetDependencies()
	for _, dep := range deps {
		dep.PushMethods(pushed)
	}
}

func (m *Module) PushMethods(fns []*ScopedFunc) {
	for _, f := range m.files {
		f.PushMethods(fns)
	}

	deps := m.GetDependencies()
	for _, dep := range deps {
		dep.PushMethods(fns)
	}
}

func (m *Module) BuildDBs(gScope *GlobalScope) {
	for _, f := range m.files {
		f.BuildDB(gScope)
	}
}

func (m *Module) SetLinker(linker *Linker) {
	for _, f := range m.files {
		f.SetLinker(linker)
	}
}

func (m *Module) RunEntryPoint(path string, ew ErrorWriter) {
	// pick any file
	f := m.files[0]

	if DEBUG_PKG_LOADING {
		fmt.Println("fns available in \"" + path + "\":")
		fmt.Println(f.DumpFuncs())
	}

	name := "main"
	fns := f.ListDispatchable(name, 0, ew)
	if !ew.Empty() { // linking happens at the same time dispatching
		return
	}

	// now filter by path
	var fn Func
	for _, fn_ := range fns {
		if fn_.Context().Path() == path {
			if fn != nil {
				ew.Add(errors.New("multiple definitions of entry-point \"" + name + "\" in \"" + path + "\""))
				return
			} else {
				fn = fn_
			}
		}
	}

	if fn == nil {
		ew.Add(errors.New("\"" + name + "\" not found in \"" + path + "\""))
		return
	}

	f.linker.entry = fn

	if DEBUG_PKG_LOADING {
		fmt.Println("running entry point \"" + name + "\" in \"" + path + "\"")
	}

	retVal := RunFunc(fn, []Value{}, ew, fn.Context())
	if !ew.Empty() {
		return
	}

	Run(retVal)
}

func (m *Module) DumpFuncs() string {
	var b strings.Builder

	for _, fn := range m.all {
		b.WriteString(fn.Dump())
		b.WriteString("\n")
	}

	return b.String()
}
