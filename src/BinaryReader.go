package main

import (
	"strconv"
)

type BinaryReader struct {
	TokenizerData
}

func NewBinaryReader(rs []Rune) *BinaryReader {
	return &BinaryReader{TokenizerData{rs}}
}

func (rr *BinaryReader) parseBuffer() (*Int, error) {
	if len(rr.buf) == 2 {
		return nil, rr.Context().Error("invalid binary literal")
	} else {
		i64, err := strconv.ParseInt(string(rr.Runes())[2:], 2, 64)
		if err != nil {
			return nil, rr.Context().Error("invalid binary literal")
		} else {
			return NewInt(i64, rr.Context()), nil
		}
	}
}

func (rr *BinaryReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsBinaryChar():
		return NewBinaryReader(append(rr.buf, r)), nil
	case r.IsWordChar():
		newTr := dispatchTokenizer(r, ew)
		ew.Add(rr.Context().Error("invalid binary literal"))
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

func (rr *BinaryReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
