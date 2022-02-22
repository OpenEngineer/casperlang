package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const NAME = "casper"

var ARGS []string = nil

var SUBCMD_USAGE = map[string]string{
	"tokenize": "tokenize <file>",
	"parse":    "parse <file>",
	"":         "<file> <args>",
}

func printErrorAndQuit(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func printMessage(msg string) {
	fmt.Fprintf(os.Stdout, "%s\n", msg)
}

func main() {
	if err := main_(); err != nil {
		printErrorAndQuit(err.Error())
		os.Exit(1)
	}
}

func genUsageError(msg string) error {
	var b strings.Builder

	if msg != "" {
		b.WriteString(msg + "\n")
	}

	b.WriteString("Usage: ")
	b.WriteString(NAME)

	b.WriteString(" <command>\n\n")

	b.WriteString("Commands:\n")
	for _, v := range SUBCMD_USAGE {
		b.WriteString(v)
		b.WriteString("\n")
	}

	return errors.New(b.String())
}

func cmdUsageError(cmd string) error {
	var b strings.Builder

	b.WriteString("Usage: ")
	b.WriteString(NAME)

	subUsage, ok := SUBCMD_USAGE[cmd]
	if !ok {
		panic("not usage defined or subcmd " + cmd)
	}

	b.WriteString(subUsage)

	return errors.New(b.String())
}

// XXX: arg parsing could be done with token parser itself
func main_() error {
	args := os.Args[1:]

	if len(args) >= 1 && filepath.Ext(args[0]) == ".cas" {
		ARGS = args
		return main_runFile(args[0])
	}

	if len(args) < 2 {
		return genUsageError("")
	}

	cmd, subArgs := args[0], args[1:]

	switch cmd {
	case "tokenize":
		return main_tokenizeFile(subArgs)
	case "parse":
		return main_parseFile(subArgs)
	default:
		return genUsageError("Unrecognized command " + cmd)
	}
}

func main_tokenizeFile(args []string) error {
	if len(args) != 1 {
		return cmdUsageError("tokenize")
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	s, err := ReadSource(path)
	if err != nil {
		return err
	}

	ew := NewErrorWriter()
	ts := Tokenize(s, ew)

	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	} else {
		printMessage(DumpTokens(ts))
	}

	return nil
}

func main_parseFile(args []string) error {
	if len(args) != 1 {
		return cmdUsageError("parse")
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	s, err := ReadSource(path)
	if err != nil {
		return err
	}

	ew := NewErrorWriter()
	fi := Parse(s, ew)

	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	} else {
		printMessage(DumpFile(fi))
	}

	return nil
}

func main_evalJSON(args []string) error {
	if len(args) != 1 {
		return cmdUsageError("eval")
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	s, err := ReadSource(path)
	if err != nil {
		return err
	}

	ew := NewErrorWriter()
	v := EvalJSON(s, ew)
	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	} else {
		printMessage(v.Dump())
	}

	return nil
}

func main_runFile(path string) error {
	ew := NewErrorWriter()

	RunPackage(path, ew)
	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	}

	return nil
}
