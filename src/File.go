package main

import (
	"path/filepath"
	"sort"
	"strings"
)

// syntax tree object, doesn't have a state
type File struct {
	path     string
	imports  []*String
	imodules []*Module // direct refernce
	fns      []*UserFunc
	db       map[string][]DispatchableFunc
}

func (f *File) Parent() Scope {
	return nil
}

func (f *File) CollectFunctions(name string) []DispatchableFunc {
	fns, ok := f.db[name]
	if !ok {
		return []DispatchableFunc{}
	} else {
		return fns
	}
}

// should be abs path
func (f *File) Path() string {
	return f.path
}

func (f *File) AddImport(imp *String) {
	f.imports = append(f.imports, imp)
}

func (f *File) AddFunc(fn *UserFunc) {
	f.fns = append(f.fns, fn)
	fn.file = f
}

func DumpFile(f *File) string {
	var b strings.Builder

	if len(f.imports) > 0 {
		b.WriteString("import ")

		for _, imp := range f.imports {
			b.WriteString("\"")
			b.WriteString(imp.Value())
			b.WriteString("\" ")
		}
		b.WriteString("\n")
	}

	for _, fn := range f.fns {
		b.WriteString(fn.Dump())
		b.WriteString("\n")
	}

	return b.String()
}

func isRelPath(path string) bool {
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		return true
	} else {
		return false
	}
}

func isAbsPath(path string) bool {
	return strings.HasPrefix(path, "/")
}

func (f *File) GetModules(p *Package, consumers []*Module, ew ErrorWriter) {
	for _, imp := range f.imports {
		var imodule *Module

		if isRelPath(imp.Value()) {
			impAbs := NewString(filepath.Clean(filepath.Join(filepath.Dir(f.path), imp.Value())), imp.Context())
			imodule = p.GetLocalModule(impAbs, consumers, ew)
		} else if isAbsPath(imp.Value()) {
			impAbs := NewString(filepath.Clean(filepath.Join(p.Dir(), imp.Value())), imp.Context())
			imodule = p.GetLocalModule(impAbs, consumers, ew)
		} else {
			imodule = p.GetExternalModule(imp, consumers, ew)
		}

		if imodule != nil {
			f.imodules = append(f.imodules, imodule)
		}
	}
}

func (f *File) CheckTypeNames(ew ErrorWriter) {
	for _, fn := range f.fns { // only the locally defined functions!
		if fn.file == f {
			fn.CheckTypeNames(f, ew)
		}
	}
}

func (f *File) ImportFuncs(ew ErrorWriter) {
	for _, imodule := range f.imodules {
		impFns := imodule.GetExportedFuncs()
		f.fns = append(f.fns, impFns...)
	}
}

func (f *File) PushMethods(fns []*UserFunc) {
Outer:
	for _, fn := range fns {
		for _, check := range f.fns {
			if check == fn {
				continue Outer
			}
		}

		f.fns = append(f.fns, fn)
	}
}

func (f *File) BuildDB(core map[string][]DispatchableFunc) {
	f.db = make(map[string][]DispatchableFunc)

	for key, fn := range core {
		f.db[key] = fn
	}

	registerUserFuncs(f.db, f.fns)
}

func (s *File) Dispatch(name *Word, args []Value, ew ErrorWriter) Func {
	if len(s.db) == 0 {
		panic("empty db? for " + s.Path())
	}

	// the connected scopes are bad this way
	fns := s.CollectFunctions(name.Value())

	if len(fns) == 0 {
		ew.Add(name.Context().Error("\"" + name.Value() + "\" undefined"))
		return nil
	}

	best, err := PickBest(fns, args, name.Context())
	if err != nil {
		ew.Add(err)
		return nil
	} else if best == nil {
		return nil
	} else {
		return best
	}
}

func (f *File) DumpFuncs() string {
	var b strings.Builder

	keys := []string{}
	for k, _ := range f.db {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, k := range keys {
		lst := f.db[k]
		for j, fn := range lst {
			b.WriteString("  ")
			b.WriteString(fn.Dump())

			if i < len(keys)-1 || j < len(lst)-1 {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}
