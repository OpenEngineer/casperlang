package main

var DEBUG_PARSER = false

func ParseFile(path string, ew ErrorWriter) *File {
	s, err := ReadSource(path)
	if err != nil {
		ew.Add(err)
		return nil
	}

	return Parse(s, ew)
}

func Parse(s *Source, ew ErrorWriter) *File {
	ts := Tokenize(s, ew)

	if !ew.Empty() {
		panic("empty file")
		return &File{}
	}

	f := parseTokens(ts, ew)
	f.path = s.path

	return f
}

func parseTokens(ts []Token, ew ErrorWriter) *File {
	var fr FileReader = NewMaybeImportReader([]Token{})

	f := &File{}

	for _, t := range ts {
		fr, f = fr.Ingest(f, t, ew)
	}

	return fr.Finalize(f, ew)
}
