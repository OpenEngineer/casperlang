package main

type NLReader struct {
	TokenizerData
	startPos FilePos
}

func NewNLReader(rs []Rune) *NLReader {
	if len(rs) == 0 {
		panic("can't extract start")
	}

	return &NLReader{TokenizerData{rs}, rs[0].Context().end}
}

func (rr *NLReader) isFirstLine() bool {
	return len(rr.buf) == 0 || !rr.buf[0].IsNewline()
}

func (rr *NLReader) indent() int {
	if rr.isFirstLine() {
		return len(rr.buf)
	} else {
		return len(rr.buf) - 1
	}
}

func (rr *NLReader) Context(src *Source) Context {
	if len(rr.buf) > 0 {
		return rr.TokenizerData.Context()
	} else {
		return Context{src, rr.startPos, rr.startPos.Advance(' ')}
	}
}

func (rr *NLReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsNewline():
		return &NLReader{TokenizerData{[]Rune{}}, r.Context().end}, nil
	case r.IsWhitespace():
		return NewNLReader(append(rr.buf, r)), nil
	case r.IsCommentChar():
		return NewCommentReader([]Rune{r}), nil
	default:
		newTr := dispatchTokenizer(r, ew)

		t := NewNL(rr.indent(), rr.Context(r.Context().src))

		return newTr, []Token{t}
	}
}

func (rr *NLReader) Finalize(ew ErrorWriter) []Token {
	return nil
}
