package main

import (
	"strconv"
)

type MaybeFloatReader struct {
	TokenizerData
}

func NewMaybeFloatReader(rs []Rune) *MaybeFloatReader {
	return &MaybeFloatReader{TokenizerData{rs}}
}

func (rr *MaybeFloatReader) parseBuffer() (Token, Rune) {
	rs := rr.Runes()
	i64, err_ := strconv.ParseInt(string(rs[0:len(rs)-1]), 10, 64)
	if err_ != nil {
		panic(err_)
	}

	iCtx := MergeContexts(rr.buf[0].Context(), rr.buf[len(rs)-2].Context())

	ti := NewInt(i64, iCtx)

	dot := rr.buf[len(rs)-1]

	return ti, dot
}

func (rr *MaybeFloatReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsDigit():
		return NewFloatReader(append(rr.buf, r)), nil
	default:
		ti, dot := rr.parseBuffer()

		newTr := dispatchTokenizer(dot, ew)

		newNewTr, t := newTr.Ingest(r, ew)

		ts := []Token{ti}
		ts = append(ts, t...)

		return newNewTr, ts
	}
}

func (rr *MaybeFloatReader) Finalize(ew ErrorWriter) []Token {
	ti, dot := rr.parseBuffer()

	return []Token{ti, NewSymbol(string([]rune{dot.r}), dot.Context())}
}
