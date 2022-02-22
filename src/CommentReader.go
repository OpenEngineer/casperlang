package main

// doesn't return a Token
type CommentReader struct {
	TokenizerData
}

func NewCommentReader(rs []Rune) *CommentReader {
	return &CommentReader{TokenizerData{rs}}
}

func (rr *CommentReader) Ingest(r Rune, ew ErrorWriter) (Tokenizer, []Token) {
	switch {
	case r.IsNewline():
		return NewNLReader([]Rune{r}), nil
	default:
		return NewCommentReader(append(rr.buf, r)), nil
	}
}

func (rr *CommentReader) Finalize(ew ErrorWriter) []Token {
	return nil
}
