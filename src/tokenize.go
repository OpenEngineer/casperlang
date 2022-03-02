package main

func Tokenize(s *Source, ew ErrorWriter) []Token {
	reader := SourceReader{s, FilePos{}}

	state := NewTokenizerState(FilePos{})

	var r Rune

	for {
		r, reader = reader.ReadNext()

		if r.IsEOF() {
			break
		} else {
			state.Ingest(r, ew)
		}
	}

	return state.Finalize(ew)
}

func TokenizeRunes(rs []Rune, ew ErrorWriter) []Token {
	state := NewTokenizerState(rs[0].Context().start)

	for _, r := range rs {
		state.Ingest(r, ew)
	}

	return state.Finalize(ew)
}

func TokenizeString(str string) []Token {
	rs := []Rune{}

	for _, c := range str {
		r := NewRune(c, NewBuiltinContext())
		rs = append(rs, r)
	}

	ew := NewErrorWriter()

	ts := TokenizeRunes(rs, ew)
	if !ew.Empty() {
		panic(ew.Dump())
	}

	if IsNL(ts[0]) {
		ts = ts[1:]
	}

	return ts
}
