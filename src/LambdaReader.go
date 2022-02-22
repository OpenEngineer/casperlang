package main

import (
	"strconv"
)

type LambdaReader struct {
	TokenizerData
}

func NewLambdaReader(rs []Rune) *LambdaReader {
	return &LambdaReader{TokenizerData{rs}}
}

func (rr *LambdaReader) parseBuffer() (Token, error) {
	if rr.last().IsDigit() {
		i64, err := strconv.ParseInt(string(rr.Runes()[1:]), 10, 64)
		if err != nil {
			panic(err)
		}

		t := NewWord("\\"+strconv.Itoa(int(i64)), rr.Context())

		if i64 < 1 {
			return t, rr.Context().Error("invalid lambda pattern syntax")
		}

		return t, nil
	} else {
		return nil, rr.Context().Error("invalid syntax")
	}
}

func (rr *LambdaReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case rr.last().IsEscape() && r.Is('('):
		t := NewSymbol("\\(", MergeContexts(rr.Context(), r.Context()))
		return &TokenizerDispatcher{}, []Token{t}
	case r.IsDigit():
		return NewLambdaReader(append(rr.buf, r)), nil
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

func (rr *LambdaReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.parseBuffer()
	ew.Add(err)
	return []Token{t}
}
