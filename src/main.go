package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/openengineer/go-repl"
)

const NAME = "casper"
const VERSION = "0.1.0"

var ARGS []string = nil

var SUBCMD_USAGE = map[string]string{
	"":         "<file> <args>",
	"tokenize": "tokenize <file>",
	"parse":    "parse    <file>",
	"profile":  "profile  <file> <args>",
	"version":  "version",
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

func parseArgs(args []string) ([]string, map[string]string) {
	options := make(map[string]string)

	return args, options
}

// XXX: arg parsing could be done with token parser itself
func main_() error {
	// TODO: extract options
	args, _ := parseArgs(os.Args[1:])

	if len(args) >= 1 && filepath.Ext(args[0]) == ".cas" {
		ARGS = args
		main_runFile(args[0])
		return nil
	} else if len(args) == 0 {
		return main_repl()
	} else if len(args) < 2 && args[0] != "version" {
		return genUsageError("")
	}

	cmd, subArgs := args[0], args[1:]

	switch cmd {
	case "tokenize":
		return main_tokenizeFile(subArgs)
	case "parse":
		return main_parseFile(subArgs)
	case "profile":
		return main_profiledRun(subArgs)
	case "version":
		fmt.Println(VERSION)
		return nil
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

var PROF_FILE *os.File = nil

func startProfiling(profFile string) {
	var err error
	PROF_FILE, err = os.Create(profFile)
	if err != nil {
		printErrorAndQuit(err.Error())
	}

	pprof.StartCPUProfile(PROF_FILE)

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		stopProfiling(profFile)

		os.Exit(1)
	}()
}

func stopProfiling(profFile string) {
	if PROF_FILE != nil {
		pprof.StopCPUProfile()

		// also write mem profile
		fMem, err := os.Create(profFile + ".mprof")
		if err != nil {
			printErrorAndQuit(err.Error())
		}

		pprof.WriteHeapProfile(fMem)
		fMem.Close()

		PROF_FILE = nil
	}
}

func main_profiledRun(args []string) error {
	startProfiling("profile.dat")

	ARGS = args[1:]
	main_runFile(args[0])

	stopProfiling("profile.dat")

	return nil
}

func main_runFile(path string) {
	ew := NewErrorWriter()

	RunPackage(path, ew)
	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	}
}

func main_repl() error {
	TARGET = "repl"

	ew := NewErrorWriter()
	h := NewRepl(ew)
	if !ew.Empty() {
		printErrorAndQuit(ew.Dump())
	}

	fmt.Printf("Welcome to Casper %s\n", VERSION)
	fmt.Println("Type \"help\" for more information")

	h.r = repl.NewRepl(h)

	return h.r.Loop()
}
