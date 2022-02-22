package main

type MaybeImportReader struct {
	FileReaderData
}

func NewMaybeImportReader(ts []Token) *MaybeImportReader {
	return &MaybeImportReader{FileReaderData{ts}}
}

func (fr *MaybeImportReader) Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) {
	switch {
	case IsNL(t):
		return NewMaybeImportReader([]Token{t}), f
	case IsLiteral(fr.last()):
		// error already thrown
		return fr, f
	case IsWord(t, "import"):
		return NewImportReader(append(fr.ts, t)), f
	case IsLiteral(t):
		ew.Add(t.Context().Error("invalid statement"))
		return NewMaybeImportReader(append(fr.ts, t)), f
	default:
		newFr := NewFuncNameReader(fr.ts)
		return newFr.Ingest(f, t, ew)
	}
}

func (fr *MaybeImportReader) Finalize(f *File, ew ErrorWriter) *File {
	if len(fr.ts) > 1 {
		// error already thrown
		return f
	} else if len(fr.ts) == 1 && !IsNL(fr.ts[0]) {
		ew.Add(fr.ts[0].Context().Error("invalid statement"))
		return f
	} else {
		return f
	}
}
