package main

import (
	"strconv"
)

type HexReader struct {
	TokenizerData
}

func NewHexReader(rs []Rune) *HexReader {
	return &HexReader{TokenizerData{rs}}
}

func (rr *HexReader) parseBuffer() (*Int, error) {
	if len(rr.buf) == 2 {
		return nil, rr.Context().Error("invalid hexadecimal literal")
	} else {
		i64, err := strconv.ParseInt(string(rr.Runes())[2:], 16, 64)
		if err != nil {
			return nil, rr.Context().Error("invalid hexadecimal literal")
		} else {
			return NewInt(i64, rr.Context()), nil
		}
	}
}

func (rr *HexReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsHexChar():
		return NewHexReader(append(rr.buf, r)), nil
	case r.IsWordChar():
		newTr := dispatchTokenizer(r, ew)
		ew.Add(rr.Context().Error("invalid hexadecimal literal"))
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

func (rr *HexReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
