package main

// asi is performed after full ingestion of body tokens
type FuncBodyReader struct {
	FileReaderData
	header     *FuncHeader
	initIndent int // to know when the function is finished, so that arguments can span multiple lines
}

func NewFuncBodyReader(header *FuncHeader, initIndent int) *FuncBodyReader {
	return &FuncBodyReader{
		FileReaderData{[]Token{}},
		header,
		initIndent,
	}
}

// groups are matched at this stage
func (fr *FuncBodyReader) Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) {
	switch {
	case IsGroupSymbol(t):
		panic("should've been handled during tokenization")
	case IsNL(t):
		nl := AssertNL(t)

		if nl.Indent() <= fr.initIndent {
			fn := fr.finalizeFunc(ew)

			f.AddFunc(fn)

			return NewFuncNameReader([]Token{t}), f
		}

		return &FuncBodyReader{
			FileReaderData{fr.ts},
			fr.header,
			fr.initIndent,
		}, f
	default:
		return &FuncBodyReader{
			FileReaderData{append(fr.ts, t)},
			fr.header,
			fr.initIndent,
		}, f
	}
}

func (fr *FuncBodyReader) finalizeTokens(ew ErrorWriter) []Token {
	if len(fr.ts) == 0 {
		ew.Add(fr.header.name.Context().Error("empty function body"))
		return nil
	}

	// also removes whitespace
	ts := RemoveNLs(fr.ts)

	return ts
}

func (fr *FuncBodyReader) finalizeFunc(ew ErrorWriter) *UserFunc {
	ts := fr.finalizeTokens(ew)
	if ts == nil || !ew.Empty() {
		return nil
	} else {
		fnBody := ParseExpr(ts, ew)
		checkArgNames(fr.header.args, ew)
		return NewUserFunc(fr.header.name, fr.header.args, fnBody, fr.header.name.Context())
	}
}

func (fr *FuncBodyReader) Finalize(f *File, ew ErrorWriter) *File {
	f.AddFunc(fr.finalizeFunc(ew))
	return f
}
