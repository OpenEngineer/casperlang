package main

import (
  "strconv"
)

type FilePos struct {
  char int
  col  int
  line int
}

func (fp FilePos) Advance(r rune) FilePos {

  if r == '\n' {
    return FilePos{fp.char+1, 0, fp.line+1}
  } else {
    return FilePos{fp.char+1, fp.col+1, fp.line}
  }
}

func (fp FilePos) ToString() string {
  return strconv.Itoa(fp.line+1) + ":" + strconv.Itoa(fp.col+1)
}

func (a FilePos) Before(b FilePos) bool {
  return a.char < b.char
}
