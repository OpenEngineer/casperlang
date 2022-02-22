package main

import (
	"strconv"
)

type NumberReader struct {
	TokenizerData
}

func NewNumberReader(rs []Rune) *NumberReader {
	return &NumberReader{TokenizerData{rs}}
}

func (rr *NumberReader) isSingleZero() bool {
	return len(rr.buf) == 1 && rr.buf[0].Is('0')
}

func (rr *NumberReader) parseBuffer() *Int {
	i64, err := strconv.ParseInt(string(rr.Runes()), 10, 64)
	if err != nil {
		panic(err)
	}

	t := NewInt(i64, rr.Context())

	return t
}

func (rr *NumberReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	t := rr.parseBuffer()

	switch {
	case r.Is('x') && rr.isSingleZero():
		return NewHexReader(append(rr.buf, r)), nil
	case r.Is('o') && rr.isSingleZero():
		return NewOctalReader(append(rr.buf, r)), nil
	case r.Is('b') && rr.isSingleZero():
		return NewBinaryReader(append(rr.buf, r)), nil
	case r.Is('e'):
		return NewFloatSciReader(append(rr.buf, r)), nil
	case r.IsWordStart():
		ew.Add(rr.Context().Error("invalid number syntax"))
		return NewWordReader([]Rune{r}), []Token{t}
	case r.IsDigit():
		return NewNumberReader(append(rr.buf, r)), nil
	case r.Is('.'):
		return NewMaybeFloatReader(append(rr.buf, r)), nil
	default:
		newTr := dispatchTokenizer(r, ew)
		return newTr, []Token{t}
	}
}

func (rr *NumberReader) Finalize(ew ErrorWriter) []Token {
	return []Token{rr.parseBuffer()}
}
