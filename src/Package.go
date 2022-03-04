package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var DEBUG_PKG_LOADING = false

type Package struct {
	dir       *String // dir of package (package.json could be in a default location)
	raw       *Dict
	consumers []*Package
	deps      map[string]*Package // loaded lazyly, circular import detection is done in during import phase
	modules   map[string]*Module  // also loaded lazyly during import phase, keys are paths wrt to roots
}

func isFile(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return !info.IsDir()
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return info.IsDir()
}

func userPackageConfig() (string, bool) {
	home := os.Getenv("HOME")

	path := filepath.Join(home, ".config", "casper", "package.json")

	if DEBUG_PKG_LOADING {
		fmt.Println("searching \"package.json\" in  \"" + filepath.Dir(path) + "\"")
	}

	if isFile(path) {
		return path, true
	} else {
		return "", false
	}
}

func searchPackageConfig(dir string) (string, bool) {
	path := filepath.Join(dir, "package.json")

	if DEBUG_PKG_LOADING {
		fmt.Println("searching \"package.json\" in  \"" + dir + "\"")
	}

	if isFile(path) {
		return path, true
	} else {
		if dir == "/" {
			return userPackageConfig()
		} else {
			return searchPackageConfig(filepath.Dir(dir))
		}
	}

	return path, true
}

func listPackageConsumers(downstream []*Package, this *String) string {
	var b strings.Builder

	b.WriteString("\n")
	for _, pkg := range downstream {
		b.WriteString("  ")
		b.WriteString(pkg.Dir())
		b.WriteString("\n")
	}
	b.WriteString("  ")
	b.WriteString(this.Value())

	return b.String()
}

func LoadConfig(dir *String, ew ErrorWriter) *Dict {
	var (
		path string
		ok   bool
	)

	if dir != nil {
		path, ok = searchPackageConfig(dir.Value())
	} else {
		path, ok = userPackageConfig()
	}

	if !ok {
		ew.Add(errors.New("no \"package.json\" found"))
		return nil
	}

	s, err := ReadSource(path)
	if err != nil {
		ew.Add(err)
		return nil
	}

	if DEBUG_PKG_LOADING {
		fmt.Println("config \"" + path + "\" loaded")
	}

	v := EvalJSON(s, ew)
	if !ew.Empty() {
		return nil
	}

	if !IsDict(v) {
		ew.Add(v.Context().Error("invalid data in package.json"))
		return nil
	}

	return AssertDict(v)
}

func LoadPackage(consumers []*Package, dir *String, ew ErrorWriter) *Package {
	if DEBUG_PKG_LOADING {
		fmt.Println("loading pkg \"" + dir.Value() + "\"")
	}

	for i, consumer := range consumers {
		if consumer.Dir() == dir.Value() {
			if i == len(consumers)-1 && consumer.Dir() == dir.Value() {
				ew.Add(dir.Context().Error("package imports self"))
				return nil
			} else {
				ew.Add(dir.Context().Error("circular package dependency:" + listPackageConsumers(consumers[i:], dir)))
				return nil
			}
		}
	}

	cfg := LoadConfig(dir, ew)
	if cfg == nil {
		return nil
	}

	if DEBUG_PKG_LOADING {
		fmt.Println("loaded package \"" + dir.Value() + "\" (unitialized)")
	}

	return &Package{
		dir,
		cfg,
		consumers,
		make(map[string]*Package),
		make(map[string]*Module),
	}
}

// this should probably be a constructor, so we keep access to consumers

func LoadEntryPackage(dir *String, ew ErrorWriter) *Package {
	p := LoadPackage([]*Package{}, dir, ew)
	if !ew.Empty() {
		return nil
	}

	m := LoadModule(p, []*Module{}, dir, ew)
	if !ew.Empty() {
		return nil
	}

	// last module to be added
	if m != nil {
		p.AddModule(m, ew)
	}

	p.RegisterFuncs(ew)
	if !ew.Empty() {
		return nil
	}

	return p
}

func (p *Package) RegisterFuncs(ew ErrorWriter) {
	// fresh linker
	linker := NewLinker()

	// TODO: search for relevant package.json file

	// link all files, modules and packages
	// the scopes will be based on this
	// what if File implements Scope?

	// 1. detect if entry point exists

	// 2. load module of entry point (i.e. parse all files in that dir)

	// 3. ask the entrypoint module for all out of package modules that must be imported and import those packages recursively (recursive function call of LoadPackage, pass in parent so that circular imports can be detected)

	// 4. ask the entrypoint module for all local modules that must be imported, and import those modules recursively (recursive with a module loader interface)

	p.MergeModuleFuncs()

	p.ImportFuncs(ew)

	p.SyncMethods(ew)

	gScope := NewGlobalScope(linker)

	p.BuildDBs(gScope)

	p.SetLinker(linker)
}

// contains exactly one module and exactly one file
func LoadReplPackage(ew ErrorWriter) *Package {
	cfg := LoadConfig(nil, ew)
	if cfg == nil {
		return nil
	}

	p := &Package{
		nil,
		cfg,
		[]*Package{},
		make(map[string]*Package),
		make(map[string]*Module),
	}

	m := LoadReplModule(p, ew)
	if !ew.Empty() {
		return nil
	}

	// last module to be added
	if m != nil {
		p.AddModule(m, ew)
	}

	p.RegisterFuncs(ew)
	if !ew.Empty() {
		return nil
	}

	return p
}

func (p *Package) Dir() string {
	if p.dir == nil {
		return "" // nowhere
	} else {
		return p.dir.Value()
	}
}

func (p *Package) GetLocalModule(path *String, consumers []*Module, ew ErrorWriter) *Module {
	if path.Value() == "/" {
		panic("invalid path")
	}
	//path = NewString(filepath.Clean(filepath.Join(p.Dir(), path.Value())), path.Context())

	m, ok := p.modules[path.Value()]
	if ok {
		return m
	} else {
		m := LoadModule(p, consumers, path, ew)
		if m != nil {
			p.modules[m.Dir()] = m
		}
		return m
	}
}

// TODO: download github
func (p *Package) depNameToPath(name *Word, ew ErrorWriter) *String {
	depsDict_, ok := p.raw.GetStrict("dependencies")
	if !ok {
		ew.Add(p.raw.Context().Error("dependencies not found in package.json"))
		return nil
	}

	depsDict, ok := depsDict_.(*Dict)
	if !ok {
		ew.Add(depsDict_.Context().Error("expected Dict"))
		return nil
	}

	depEntry_, ok := depsDict.GetStrict(name.Value())
	if !ok {
		ew.Add(name.Context().Error("package \"" + name.Value() + "\" not defined in \"" + p.raw.Context().Path() + "\""))
		return nil
	}

	depEntry, ok := depEntry_.(*Dict)
	if !ok {
		ew.Add(depEntry_.Context().Error("expected Dict"))
		return nil
	}

	// must have path declared
	depPath_, ok := depEntry.GetStrict("path")
	if !ok {
		return checkGitPackage(name, depEntry, ew)

		//ew.Add(depEntry.Context().Error("\"path\" undefined for \"" + name.Value() + "\""))
		//return nil
	}

	depPath, ok := depPath_.(*String)
	if !ok {
		ew.Add(depPath_.Context().Error("expected String"))
		return nil
	}

	depPathStr := depPath.Value()
	if isRelPath(depPathStr) {
		depPathStr = filepath.Clean(filepath.Join(p.raw.Context().Path(), depPathStr))
	}

	return NewString(depPathStr, name.Context())
}

func checkGitPackage(name *Word, raw *Dict, ew ErrorWriter) *String {
	url_, ok := raw.GetStrict("url")
	if !ok {
		ew.Add(raw.Context().Error("\"path\" nor \"url\" defined for \"" + name.Value() + "\""))
		return nil
	}

	url, ok := url_.(*String)
	if !ok {
		ew.Add(url_.Context().Error("expected String"))
		return nil
	}

	version_, ok := raw.GetStrict("version")
	if !ok {
		ew.Add(raw.Context().Error("\"version\" undefined for \"" + name.Value() + "\""))
		return nil
	}

	version, ok := version_.(*String)
	if !ok {
		ew.Add(version_.Context().Error("expected String"))
		return nil
	}

	home := os.Getenv("HOME")

	// TODO: XDG-style paths for Windows/Mac
	localPath := filepath.Join(home, ".cache", "casper", "pkg", url.Value())
	localPath = filepath.Join(localPath, version.Value())

	dst := NewString(localPath, raw.Context())

	if !isDir(dst.Value()) {
		if err := FetchGitRepo(url, version, "", dst); err != nil {
			ew.Add(err)
			return nil
		}
	}

	return NewString(localPath, raw.Context())
}

func (p *Package) GetPackage(name *Word, ew ErrorWriter) *Package {
	path := p.depNameToPath(name, ew)
	if path == nil || !ew.Empty() {
		return nil
	}

	if dep, ok := p.deps[path.Value()]; ok {
		return dep
	} else {
		consumers := []*Package{}
		consumers = append(consumers, p.consumers...)
		consumers = append(consumers, p)

		dep := LoadPackage(consumers, path, ew)
		if !ew.Empty() {
			return nil
		}

		p.deps[path.Value()] = dep
		return dep
	}
}

func (p *Package) GetExternalModule(path *String, consumers []*Module, ew ErrorWriter) *Module {
	parts := filepath.SplitList(path.Value())

	if len(parts) == 0 {
		ew.Add(errors.New("invalid empty import"))
		return nil
	}

	pName := parts[0]

	dep := p.GetPackage(NewWord(pName, path.Context()), ew)
	if !ew.Empty() {
		return nil
	}

	rest := NewString(filepath.Join(dep.Dir(), string([]byte{filepath.Separator})+filepath.Join(parts[1:]...)), path.Context())

	return dep.GetLocalModule(rest, consumers, ew)
}

func (p *Package) AddModule(m *Module, ew ErrorWriter) {
	if _, ok := p.modules[m.Dir()]; ok {
		ew.Add(errors.New("module " + m.Dir() + " already registered"))
	} else {
		p.modules[m.Dir()] = m
	}
}

func (p *Package) MergeModuleFuncs() {
	for _, dep := range p.deps {
		dep.MergeModuleFuncs()
	}

	for _, m := range p.modules {
		m.MergeFuncs()
	}
}

func (p *Package) ImportFuncs(ew ErrorWriter) {
	for _, dep := range p.deps {
		dep.ImportFuncs(ew)
	}

	for _, m := range p.modules {
		m.ImportFuncs(ew)
	}
}

func (p *Package) SyncMethods(ew ErrorWriter) {
	for _, m := range p.modules {
		m.SyncMethods(ew)
	}

	for _, dep := range p.deps {
		dep.SyncMethods(ew)
	}
}

func (p *Package) BuildDBs(gScope *GlobalScope) {
	for _, dep := range p.deps {
		dep.BuildDBs(gScope)
	}

	for _, m := range p.modules {
		m.BuildDBs(gScope)
	}
}

// check uniqueness and existence of the typenames
func (p *Package) SetLinker(linker *Linker) {
	for _, dep := range p.deps {
		dep.SetLinker(linker)
	}

	for _, m := range p.modules {
		m.SetLinker(linker)
	}
}

// a single run? so we avoid recursiveness problems?
func (p *Package) RunEntryPoint(path string, ew ErrorWriter) {
	// find out in which module it is

	m, ok := p.modules[filepath.Dir(path)]
	if !ok {
		ew.Add(errors.New("entry point \"" + path + "\" not found"))
		return
	}

	m.RunEntryPoint(path, ew)
}

func RunPackage(path string, ew ErrorWriter) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		ew.Add(err)
		return
	}

	dir := NewString(filepath.Dir(absPath), NewStdinContext())

	p := LoadEntryPackage(dir, ew)
	if !ew.Empty() {
		return
	}

	p.RunEntryPoint(absPath, ew)
}
