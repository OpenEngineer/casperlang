package main

import (
	"strconv"
)

type FloatReader struct {
	TokenizerData
}

func NewFloatReader(rs []Rune) *FloatReader {
	return &FloatReader{TokenizerData{rs}}
}

func (rr *FloatReader) parseBuffer() (*Float, error) {
	f64, err := strconv.ParseFloat(string(rr.Runes()), 64)
	if err != nil {
		return nil, rr.Context().Error("invalid float literal")
	} else {
		t := NewFloat(f64, rr.Context())
		return t, nil
	}
}

func (rr *FloatReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsDigit():
		return NewFloatReader(append(rr.buf, r)), nil
	case r.Is('e'):
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

func (rr *FloatReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
