package main

type SourceReader struct {
  src *Source
  pos FilePos // pos of next rune to be read
}

func (s SourceReader) ReadNext() (Rune, SourceReader) {
  r := s.src.Get(s.pos)

  newPos := s.pos.Advance(r)

  return NewRune(r, Context{s.src, s.pos, newPos}), SourceReader{s.src, newPos}
}
