package main

import (
	"strconv"
)

type OctalReader struct {
	TokenizerData
}

func NewOctalReader(rs []Rune) *OctalReader {
	return &OctalReader{TokenizerData{rs}}
}

func (rr *OctalReader) parseBuffer() (*Int, error) {
	if len(rr.buf) == 2 {
		return nil, rr.Context().Error("invalid octal literal")
	} else {
		i64, err := strconv.ParseInt(string(rr.Runes())[2:], 8, 64)
		if err != nil {
			return nil, rr.Context().Error("invalid octal literal")
		} else {
			return NewInt(i64, rr.Context()), nil
		}
	}
}

func (rr *OctalReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsOctalChar():
		return NewOctalReader(append(rr.buf, r)), nil
	case r.IsWordChar():
		newTr := dispatchTokenizer(r, ew)
		ew.Add(rr.Context().Error("invalid octal literal"))
		return newTr, nil
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

func (rr *OctalReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
