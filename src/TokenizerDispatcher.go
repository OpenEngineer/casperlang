package main

type TokenizerDispatcher struct {
}

func dispatchTokenizer(r Rune, ew ErrorWriter) Tokenizer {
	switch {
	case r.IsNewline():
		return NewNLReader([]Rune{r})
	case r.IsWhitespace():
		return &TokenizerDispatcher{}
	case r.IsWordStart():
		return NewWordReader([]Rune{r})
	case r.IsSymbolStart():
		return NewSymbolReader([]Rune{r})
	case r.IsDigit():
		return NewNumberReader([]Rune{r})
	case r.IsCommentChar():
		return NewCommentReader([]Rune{r})
	case r.IsStringStart():
		return NewStringReader([]Rune{r})
	case r.IsEscape():
		return NewLambdaReader([]Rune{r})
	case r.IsDollar():
		return NewDollarReader([]Rune{r})
	default:
		ew.Add(r.ctx.Error("invalid syntax"))
		return &TokenizerDispatcher{}
	}
}

func (tr *TokenizerDispatcher) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	return dispatchTokenizer(r, ew), nil
}

func (tr *TokenizerDispatcher) Finalize(ew ErrorWriter) []Token {
	return nil
}
