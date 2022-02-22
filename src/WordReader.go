package main

type WordReader struct {
	TokenizerData
}

func NewWordReader(rs []Rune) *WordReader {
	return &WordReader{TokenizerData{rs}}
}

func (rr *WordReader) toWord() *Word {
	return NewWord(string(rr.Runes()), rr.Context())
}

func (rr *WordReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsWordChar():
		return NewWordReader(append(rr.buf, r)), nil
	default:
		return dispatchTokenizer(r, ew), []Token{rr.toWord()}
	}
}

func (rr *WordReader) Finalize(ew ErrorWriter) []Token {
	return []Token{rr.toWord()}
}
