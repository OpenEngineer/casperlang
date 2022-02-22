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
