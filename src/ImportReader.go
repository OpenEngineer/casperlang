package main

type ImportReader struct {
	FileReaderData
}

func NewImportReader(ts []Token) *ImportReader {
	return &ImportReader{FileReaderData{ts}}
}

func (fr *ImportReader) isEmpty() bool {
	return IsWord(fr.last(), "import")
}

func (fr *ImportReader) isValid() bool {
	return fr.isEmpty() || IsString(fr.last())
}

func (fr *ImportReader) Context() Context {
	return fr.ts[1].Context()
}

func (fr *ImportReader) Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) {
	switch {
	case IsNL(t):
		nl := AssertNL(t)
		if nl.Indent() > fr.indent() {
			return fr, f
		} else {
			if fr.isEmpty() {
				ew.Add(fr.Context().Error("empty import statement"))
				return fr, f
			} else {
				newFr := NewMaybeImportReader([]Token{t})
				return newFr, f
			}
		}
	case fr.isValid() && IsString(t):
		f.AddImport(AssertString(t))
		return NewImportReader(append(fr.ts, t)), f
	default:
		if fr.isValid() {
			ew.Add(t.Context().Error("invalid import statement"))
			return NewImportReader(append(fr.ts, t)), f
		} else {
			// error was thrown before
			return fr, f
		}
	}
}

func (fr *ImportReader) Finalize(f *File, ew ErrorWriter) *File {
	if fr.isEmpty() {
		ew.Add(fr.Context().Error("empty import statement"))
		return f
	} else {
		return f
	}
}
