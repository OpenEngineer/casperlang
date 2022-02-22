package main

type FuncHeaderReader struct {
	FileReaderData
	indent int
	name   *Word
}

func NewFuncHeaderReader(indent int, name *Word, ts []Token) *FuncHeaderReader {
	return &FuncHeaderReader{FileReaderData{ts}, indent, name}
}

func (fr *FuncHeaderReader) Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) {
	switch {
	case IsSymbol(t, "="):
		args := ParsePatterns(fr.ts, ew)
		header := NewFuncHeader(fr.name, args)
		return NewFuncBodyReader(header, fr.indent), f
	case IsNL(t):
		newIndent := AssertNL(t).Indent()
		return NewFuncHeaderReader(newIndent, fr.name, fr.ts), f
	default:
		return NewFuncHeaderReader(fr.indent, fr.name, append(fr.ts, t)), f
	}
}

func (fr *FuncHeaderReader) Finalize(f *File, ew ErrorWriter) *File {
	ew.Add(fr.name.Context().EndError("expected function body"))
	return f
}
