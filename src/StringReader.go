package main

type StringReader struct {
	TokenizerData
}

type StringEscReader struct {
	TokenizerData
}

type StringFinReader struct {
	TokenizerData
	emptyAfterSub bool
}

type MaybeStringSubReader struct {
	TokenizerData
}

type StringSubReader struct {
	TokenizerData
	stack []Rune // can be '$', '{' or '"' or '\'
	// stack parts are added to end!
}

type StringSubPostReader struct {
	TokenizerData
}

func NewStringReader(rs []Rune) *StringReader {
	return &StringReader{TokenizerData{rs}}
}

func NewStringEscReader(rs []Rune) *StringEscReader {
	return &StringEscReader{TokenizerData{rs}}
}

func NewStringFinReader(rs []Rune, eas bool) *StringFinReader {
	return &StringFinReader{TokenizerData{rs}, eas}
}

func NewMaybeStringSubReader(rs []Rune) *MaybeStringSubReader {
	return &MaybeStringSubReader{TokenizerData{rs}}
}

func NewStringSubReader(rs []Rune, stack []Rune) *StringSubReader {
	return &StringSubReader{
		TokenizerData{rs},
		stack,
	}
}

func NewStringSubPostReader(rs []Rune) *StringSubPostReader {
	return &StringSubPostReader{TokenizerData{rs}}
}

func dispatchStringReader(rPrev []Rune, r Rune, ew ErrorWriter) Tokenizer {
	switch {
	case r.Is('\\'):
		return NewStringEscReader(append(rPrev, r))
	case r.Is('"'):
		return NewStringFinReader(append(rPrev, r), false)
	case r.Is('$'):
		return NewMaybeStringSubReader(append(rPrev, r))
	default:
		return NewStringReader(append(rPrev, r))
	}
}

func (rr *StringReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	return dispatchStringReader(rr.buf, r, ew), nil
}

func (rr *StringReader) Finalize(ew ErrorWriter) []Token {
	ctx := rr.buf[0].Context()
	ew.Add(ctx.Error("EOF while scanning string literal"))
	return nil
}

func (rr *StringEscReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	n := len(rr.buf)
	last := rr.last()
	buf := rr.buf[0 : n-1]
	newCtx := MergeContexts(last.Context(), r.Context())

	switch {
	case r.Is('"') || r.Is('$') || r.Is('\\'):
		return NewStringReader(append(buf, NewRune(r.r, newCtx))), nil
	case r.Is('n'):
		return NewStringReader(append(buf, NewRune('\n', newCtx))), nil
	case r.Is('t'):
		return NewStringReader(append(buf, NewRune('\t', newCtx))), nil
	default:
		ew.Add(newCtx.Error("unrecognized escape sequence"))
		return NewStringReader(buf), nil
	}
}

func (rr *StringEscReader) Finalize(ew ErrorWriter) []Token {
	ctx := rr.buf[0].Context()
	ew.Add(ctx.Error("EOF while scanning string literal"))
	return nil
}

func (rr *StringFinReader) toString() Token {
	if rr.emptyAfterSub {
		return nil
	} else {
		rs := rr.Runes()

		t := NewString(string(rs[1:len(rs)-1]), rr.Context())

		return t
	}
}

func (rr *StringFinReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	newTr := dispatchTokenizer(r, ew)

	t := rr.toString()

	return newTr, []Token{t}
}

func (rr *StringFinReader) Finalize(ew ErrorWriter) []Token {
	return []Token{rr.toString()}
}

func (rr *MaybeStringSubReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.Is('{'):
		rs := rr.Runes()

		ts := []Token{}

		if len(rs) > 2 {
			tStr := NewString(string(rs[1:len(rs)-1]), MergeContexts(rr.buf[0].Context(), rr.buf[len(rs)-2].Context()))
			plus := NewSymbol("+", rr.buf[len(rs)-1].Context())

			ts = append(ts, tStr, plus)
		}

		openPar := NewSymbol("(", rr.buf[len(rs)-1].Context())
		ts = append(ts, openPar)

		return NewStringSubReader([]Rune{r}, []Rune{}), ts
	default:
		return dispatchStringReader(rr.buf, r, ew), nil
	}
}

func (rr *MaybeStringSubReader) Finalize(ew ErrorWriter) []Token {
	ctx := rr.buf[0].Context()
	ew.Add(ctx.Error("EOF while scanning string literal"))
	return nil
}

func (rr *StringSubReader) stackHead() Rune {
	return rr.stack[len(rr.stack)-1]
}

func (rr *StringSubReader) stackHead2() Rune {
	return rr.stack[len(rr.stack)-2]
}

func (rr *StringSubReader) popStackHead() []Rune {
	n := len(rr.stack)

	return rr.stack[0 : n-1]
}

func (rr *StringSubReader) popStackHead2() []Rune {
	n := len(rr.stack)

	return rr.stack[0 : n-2]
}

func (rr *StringSubReader) pushStackHead(r Rune) []Rune {
	return append(rr.stack, r)
}

func (rr *StringSubReader) popPushStackHead(r Rune) []Rune {
	n := len(rr.stack)

	return append(rr.stack[0:n-1], r)
}

// this is a bit tricky because we want to allow nested string substitutions
func (rr *StringSubReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case len(rr.stack) == 0 && r.Is('}'):
		closePar := NewSymbol(")", r.Context())

		innerTokens := TokenizeRunes(rr.buf[1:], ew)

		if len(innerTokens) > 0 && IsNL(innerTokens[0]) {
			innerTokens = innerTokens[1:]
		}

		show := NewWord("show", r.Context())

		innerOpenPar := NewSymbol("(", r.Context())
		innerClosePar := NewSymbol(")", r.Context())

		ts := []Token{show, innerOpenPar}
		ts = append(ts, innerTokens...)
		ts = append(ts, innerClosePar, closePar)

		return NewStringSubPostReader([]Rune{NewRune('"', r.Context())}), ts
	case len(rr.stack) == 0:
		newStack := []Rune{}
		if r.Is('{') || r.Is('"') {
			newStack = []Rune{r}
		}

		return NewStringSubReader(append(rr.buf, r), newStack), nil
	case rr.stackHead().Is('{') && r.Is('}'):
		return NewStringSubReader(append(rr.buf, r), rr.popStackHead()), nil
	case rr.stackHead().Is('\\'):
		return NewStringSubReader(append(rr.buf, r), rr.popStackHead()), nil
	case rr.stackHead().Is('"') && r.Is('"'):
		return NewStringSubReader(append(rr.buf, r), rr.popStackHead()), nil
	case rr.stackHead().Is('"') && r.Is('\\'):
		return NewStringSubReader(append(rr.buf, r), rr.pushStackHead(r)), nil
	case rr.stackHead().Is('"') && r.Is('$'):
		return NewStringSubReader(append(rr.buf, r), rr.pushStackHead(r)), nil
	case rr.stackHead().Is('$') && r.Is('{'):
		return NewStringSubReader(append(rr.buf, r), rr.popPushStackHead(r)), nil
	case rr.stackHead().Is('$') && rr.stackHead2().Is('"') && r.Is('"'):
		return NewStringSubReader(append(rr.buf, r), rr.popStackHead2()), nil
	case rr.stackHead().Is('$'):
		return NewStringSubReader(append(rr.buf, r), rr.popStackHead()), nil
	default:
		return NewStringSubReader(append(rr.buf, r), rr.stack), nil
	}
}

func (rr *StringSubReader) Finalize(ew ErrorWriter) []Token {
	ctx := rr.buf[0].Context()
	ew.Add(ctx.Error("EOF while scanning string literal"))
	return nil
}

func (rr *StringSubPostReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	if r.Is('"') {
		return NewStringFinReader([]Rune{rr.buf[0], r}, true), nil
	} else {
		plus := NewSymbol("+", rr.buf[0].Context())

		return dispatchStringReader(rr.buf, r, ew), []Token{plus}
	}
}

func (rr *StringSubPostReader) Finalize(ew ErrorWriter) []Token {
	ctx := rr.buf[0].Context()
	ew.Add(ctx.Error("EOF while scanning string literal"))
	return nil
}
