package main

type FileReader interface {
	Ingest(f *File, t Token, ew ErrorWriter) (FileReader, *File) // updates file
	Finalize(f *File, ew ErrorWriter) *File
}

type FileReaderData struct {
	ts []Token
}

func (fr *FileReaderData) last() Token {
	return fr.ts[len(fr.ts)-1]
}

func (fr *FileReaderData) indent() int {
	if IsNL(fr.ts[0]) {
		return AssertNL(fr.ts[0]).Indent()
	} else {
		return 0
	}
}
