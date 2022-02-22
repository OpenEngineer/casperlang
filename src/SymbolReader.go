package main

type SymbolReader struct {
	TokenizerData
}

func NewSymbolReader(rs []Rune) *SymbolReader {
	return &SymbolReader{TokenizerData{rs}}
}

func (rr *SymbolReader) toSymbol() (Token, error) {
	symbol := string(rr.Runes())
	if IsValidSymbol(symbol) {
		return NewSymbol(symbol, rr.Context()), nil
	} else {
		return nil, rr.Context().Error("invalid symbol")
	}
}

func (rr *SymbolReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	symbol := "" + string(rr.Runes()) + string([]rune{r.r})

	if MaybeValidSymbol(symbol) {
		return NewSymbolReader(append(rr.buf, r)), nil
	} else {
		t, err := rr.toSymbol()
		if err != nil {
			ew.Add(err)
		}

		newTr := dispatchTokenizer(r, ew)

		if err != nil {
			return newTr, []Token{}
		} else {
			return newTr, []Token{t}
		}
	}
}

func (rr *SymbolReader) Finalize(ew ErrorWriter) []Token {
	t, err := rr.toSymbol()

	if err != nil {
		ew.Add(err)
	}

	return []Token{t}
}
