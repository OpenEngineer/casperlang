package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// syntax tree object, doesn't have a state
type File struct {
	path     string
	imports  []*String
	imodules []*Module // direct refernce
	fns      []*ScopedFunc
	db       map[string][]*ScopedFunc
	linker   *Linker
}

func NewReplFile() *File {
	return &File{"", []*String{}, []*Module{}, []*ScopedFunc{}, make(map[string][]*ScopedFunc), nil}
}

func (f *File) GetLocal(name string) *Variable {
	return nil
}

func (f *File) ListDispatchable(name string, nArgs int, ew ErrorWriter) []Func {
	fns_, ok := f.db[name]
	if !ok {
		return []Func{}
	} else {
		fns := []Func{}

		for _, fn_ := range fns_ {
			if fn_.NumArgs() == nArgs || nArgs == -1 {
				// recursive calls of functions with same names and same number of arguments will give problems here.
				fn := f.linker.LinkFunc(fn_.fn, fn_.scope, ew)
				fns = append(fns, fn)
			}
		}

		return fns
	}
}

func (f *File) Dir() string {
	if f.path == "" {
		return f.Path()
	} else {
		return filepath.Dir(f.path)
	}
}

// should be abs path
func (f *File) Path() string {
	if f.path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		return pwd
	} else {
		return f.path
	}
}

func (f *File) AddImport(imp *String) {
	for _, check := range f.imports {
		if check.Value() == imp.Value() {
			return
		}
	}

	f.imports = append(f.imports, imp)
}

func (f *File) AddFunc(fn *UserFunc) {
	f.fns = append(f.fns, NewScopedFunc(fn, f))
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
	toRemove := []int{}
	for i, imp := range f.imports {
		var imodule *Module

		if imp.Value() == "." {
			imp = NewString("./", imp.Context())
		} else if imp.Value() == ".." {
			imp = NewString("../", imp.Context())
		}

		if isRelPath(imp.Value()) {
			impAbs := NewString(filepath.Clean(filepath.Join(f.Dir(), imp.Value())), imp.Context())
			imodule = p.GetLocalModule(impAbs, consumers, ew)
		} else if isAbsPath(imp.Value()) {
			impAbs := NewString(filepath.Clean(filepath.Join(p.Dir(), imp.Value())), imp.Context())
			imodule = p.GetLocalModule(impAbs, consumers, ew)
		} else {
			imodule = p.GetExternalModule(imp, consumers, ew)
		}

		if imodule != nil {
			found := false
			for _, check := range f.imodules {
				if check == imodule {
					found = true
				}
			}

			if !found {
				f.imodules = append(f.imodules, imodule)
			}
		} else {
			toRemove = append(toRemove, i)
		}
	}

	if len(toRemove) > 0 {
		filteredImports := []*String{}

	Outer:
		for i, imp := range f.imports {

			for _, check := range toRemove {
				if i == check {
					continue Outer
				}
			}

			filteredImports = append(filteredImports, imp)
		}

		f.imports = filteredImports
	}
}

func (f *File) ImportFuncs(ew ErrorWriter) {
	for _, imodule := range f.imodules {
		impFns := imodule.GetExportedFuncs()
		f.PushMethods(impFns)
	}
}

func (f *File) PushMethods(fns []*ScopedFunc) {
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

func (f *File) BuildDB(gScope *GlobalScope) {
	f.db = make(map[string][]*ScopedFunc)

	for _, fn := range f.fns {
		key := fn.Name()

		lst, ok := f.db[key]
		if ok {
			f.db[key] = append(lst, fn)
		} else {
			f.db[key] = []*ScopedFunc{fn}
		}
	}

	for key, fns_ := range gScope.db {
		fns := []*ScopedFunc{}

		for _, fn := range fns_ {
			fns = append(fns, NewScopedFunc(fn, gScope))
		}

		lst, ok := f.db[key]
		if ok {
			f.db[key] = append(lst, fns...)
		} else {
			f.db[key] = fns
		}
	}

}

func (f *File) SetLinker(linker *Linker) {
	f.linker = linker
}

func (f *File) ListNames() []string {
	keys := []string{}
	for k, _ := range f.db {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (f *File) DumpFuncs() string {
	keys := f.ListNames()

	var b strings.Builder
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
