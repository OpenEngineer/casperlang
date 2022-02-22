package main

import (
	"strconv"
)

type FloatSciReader struct {
	TokenizerData
}

func NewFloatSciReader(rs []Rune) *FloatSciReader {
	return &FloatSciReader{TokenizerData{rs}}
}

func (rr *FloatSciReader) parseBuffer() (Token, error) {
	if !rr.last().IsDigit() {
		return nil, rr.Context().Error("invalid float literal")
	} else {
		f64, err := strconv.ParseFloat(string(rr.Runes()), 64)
		if err != nil {
			return nil, rr.Context().Error("invalid float literal")
		} else {
			t := NewFloat(f64, rr.Context())
			return t, nil
		}
	}
}

func (rr *FloatSciReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsDigit():
		return NewFloatSciReader(append(rr.buf, r)), nil
	case (r.Is('+') || r.Is('-')) && rr.last().Is('e'):
		return NewFloatSciReader(append(rr.buf, r)), nil
	default:
		newTr := dispatchTokenizer(r, ew)

		t, err := rr.parseBuffer()
		if err != nil {
			ew.Add(err)
			return newTr, nil
		} else {
			return newTr, []Token{t}
		}
	}
}

func (rr *FloatSciReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
