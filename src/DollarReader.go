package main

import (
	"strconv"
)

type DollarReader struct {
	TokenizerData
}

func NewDollarReader(rs []Rune) *DollarReader {
	return &DollarReader{TokenizerData{rs}}
}

func (rr *DollarReader) parseBuffer() (Token, error) {
	if rr.last().IsDigit() {
		i64, err := strconv.ParseInt(string(rr.Runes()[1:]), 10, 64)
		if err != nil {
			panic(err)
		}

		t := NewDollar(int(i64), rr.Context())

		return t, nil
	} else {
		return NewDollar(0, rr.Context()), nil
	}
}

func (rr *DollarReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsDigit():
		return NewDollarReader(append(rr.buf, r)), nil
	default:
		t, err := rr.parseBuffer()

		newTr := dispatchTokenizer(r, ew)

		if err != nil {
			ew.Add(err)
			return newTr, nil
		} else {
			return newTr, []Token{t}
		}
	}
}

func (rr *DollarReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
