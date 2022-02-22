package main

const EOF = rune(0)

type Rune struct {
  r   rune
  ctx Context
}

func NewRune(r rune, ctx Context) Rune {
  return Rune{r, ctx}
}

func (r Rune) Context() Context {
  return r.ctx
}

func (r Rune) Is(c rune) bool {
  return r.r == c
}

func (r Rune) IsEOF() bool {
  return r.r == EOF
}

func (r Rune) IsNewline() bool {
  return r.r == '\n'
}

func (r Rune) IsWhitespace() bool {
  return r.IsNewline() || r.r == ' '
}

func (r Rune) IsWordStart() bool {
  return (r.r >= 65 && r.r <= 90) || (r.r == 95) || (r.r >= 97 && r.r <= 122)
}

func (r Rune) IsDigit() bool {
  return (r.r >= 48 && r.r <= 57)
}

func (r Rune) IsWordChar() bool {
  return r.IsDigit() || r.IsWordStart()
}

func (r Rune) IsGroupChar() bool {
  return r.r == '(' || r.r == ')' || r.r == '{' || r.r == '}' || r.r == '[' || r.r == ']'
}

func (r Rune) IsSymbolStart() bool {
  for _, validSymbol := range validSymbols {
    if rune(validSymbol[0]) == r.r {
      return true
    }
  }

  return false
}

func (r Rune) IsHexChar() bool {
  return (r.r >= 97 && r.r <= 102) || r.IsDigit()
}

func (r Rune) IsBinaryChar() bool {
  return r.r == '0' || r.r == '1'
}

func (r Rune) IsOctalChar() bool {
  return r.r >= 48 && r.r <= 55
}

func (r Rune) IsCommentChar() bool {
  return r.r == '#'
}

func (r Rune) IsStringStart() bool {
  return r.r == '"'
}

func (r Rune) IsEscape() bool {
  return r.r == '\\'
}

func (r Rune) IsDollar() bool {
  return r.r == '$'
}

func (r Rune) IsSimpleEscapable() bool {
  return r.r == '"' || r.r == '$' || r.IsEscape()
}
