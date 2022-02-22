package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Source struct {
	path    string
	content []rune
}

func NewSource(path string, raw_ []byte) *Source {
	raw := []rune(string(raw_))

	n := len(raw)
	// remove the carriage returns and turn tabs into two spaces
	clean := []rune{}
	for i, r := range raw {
		switch r {
		case '\r':
			if (i > 0 && raw[i-1] != '\n') && (i < n-1 && raw[i+1] != '\n') {
				clean = append(clean, rune('\n'))
			}
		case '\t':
			clean = append(clean, rune(' '), rune(' '))
		default:
			clean = append(clean, r)
		}
	}

	return &Source{path, clean}
}

func ReadSource(path string) (*Source, error) {
	if !filepath.IsAbs(path) {
		panic("ReadSource(path) must be abs path")
	}

	//f, err := filepath.Abs(path)
	//if err != nil {
	//return nil, errors.New("invalid path \"" + path + "\"")
	//}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file \"" + path + "\" not found")
		} else {
			panic(err)
		}
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("unable to read \"" + path + "\"")
	}

	return NewSource(path, b), nil
}

func (s *Source) Get(fp FilePos) rune {
	if fp.char >= len(s.content) {
		return EOF
	} else {
		return s.content[fp.char]
	}
}

func (s *Source) Path() string {
	return s.path
}
