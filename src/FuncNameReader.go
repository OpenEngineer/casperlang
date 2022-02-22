package main

type FuncNameReader struct {
	FileReaderData
}

func NewFuncNameReader(ts []Token) *FuncNameReader {
	return &FuncNameReader{FileReaderData{ts}}
}

func (fr *FuncNameReader) isEmpty() bool {
	return IsNL(fr.last())
}

func (fr *FuncNameReader) Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) {
	switch {
	case IsNL(t) && !fr.isEmpty():
		// reset
		return NewFuncNameReader([]Token{t}), f
	case !fr.isEmpty():
		// ignore rest of line
		return fr, f
	case IsWord(t):
		w := AssertWord(t)
		return NewFuncHeaderReader(fr.indent(), w, []Token{}), f
	case IsOperatorSymbol(t):
		s := AssertSymbol(t)
		return NewFuncHeaderReader(fr.indent(), s.ToWord(), []Token{}), f
	case IsGroupSymbol(t):
		// ignore this whole line!
		if fr.isEmpty() {
			ew.Add(t.Context().Error("invalid function statement"))
			return NewFuncNameReader(append(fr.ts, t)), f
		} else {
			// error was already thrown before
			return fr, f
		}
	default:
		// bad name
		ew.Add(t.Context().Error("invalid function name"))
		return NewFuncHeaderReader(fr.indent(), NewWord("<invalid>", t.Context()), []Token{}), f
	}
}

func (fr *FuncNameReader) Finalize(f *File, ew ErrorWriter) *File {
	return f
}
