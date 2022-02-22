package main

type TokenizerState struct {
	tokenizer Tokenizer
	ts        []Token
}

func NewTokenizerState(start FilePos) *TokenizerState {
	return &TokenizerState{
		&NLReader{TokenizerData{[]Rune{}}, start},
		[]Token{},
	}
}

func (s *TokenizerState) Ingest(r Rune, ew ErrorWriter) {
	var (
		ts []Token
	)

	s.tokenizer, ts = s.tokenizer.Ingest(r, ew)

	s.add(ts)
}

func (s *TokenizerState) add(ts []Token) {
	if ts != nil && len(ts) > 0 {
		for _, t := range ts {
			if t != nil {
				s.ts = append(s.ts, t)
			}
		}
	}
}

func (s *TokenizerState) Finalize(ew ErrorWriter) []Token {
	ts := s.tokenizer.Finalize(ew)

	s.add(ts)

	return GroupTokens(s.ts, ew)
}
