package main

import (
	"sort"
	"strings"
)

type ErrorWriter interface {
	Add(err error)
	Empty() bool
	Dump() string
}

type ErrorWriterData struct {
	errs ErrorList
}

func NewErrorWriter() *ErrorWriterData {
	return &ErrorWriterData{ErrorList(make([]error, 0))}
}

func (ew *ErrorWriterData) Add(err error) {
	if err != nil {
		ew.errs = append(ew.errs, err)
	}
}

func (ew *ErrorWriterData) Empty() bool {
	return len(ew.errs) == 0
}

func (ew *ErrorWriterData) Dump() string {
	sort.Stable(ew.errs)

	var b strings.Builder

	for i, err := range ew.errs {
		if IsError(err) {
			if i > 0 {
				break
			}
		} else {
			b.WriteString("Error: ")
		}

		b.WriteString(err.Error())

		if i < len(ew.errs)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
