package main

type Tokenizer interface {
	Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) // multiple tokens can be returned in case of some syntactic sugar
	Finalize(ew ErrorWriter) []Token
}

type TokenizerData struct {
	buf []Rune
}

func (rr *TokenizerData) Runes() []rune {
	res := make([]rune, len(rr.buf))
	for i, r := range rr.buf {
		res[i] = r.r
	}

	return res
}

func (rr *TokenizerData) Context() Context {
	ctx0 := rr.buf[0].ctx
	ctxe := rr.buf[len(rr.buf)-1].ctx

	return MergeContexts(ctx0, ctxe)
}

func (rr *TokenizerData) last() Rune {
	n := len(rr.buf)

	if n == 0 {
		panic("can't call last")
	}

	return rr.buf[n-1]
}
