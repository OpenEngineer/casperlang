package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// XXX: are the evals here redundant, because they are already done in Call?

var builtinIOFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "Path",
		Args:     []string{"String"},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "File",
		Args:     []string{"String"},
		LinkReqs: []string{"Path"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Path"][0], []Value{self.args[0]}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Dir",
		Args:     []string{"String"},
		LinkReqs: []string{"Path"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Path"][0], []Value{self.args[0]}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "HttpReq",
		Args:     []string{"String", "String", "String"},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:    "exit",
		Args:    []string{"Int"},
		Targets: []string{"default"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					i := AssertInt(self.args[0]).Value()
					os.Exit(int(i))
					return nil
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name: "echo",
		Args: []string{"String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					str := AssertString(self.args[0]).Value()
					if len(str) > 0 {
						fmt.Fprintf(ioc.Stdout(), "%s", str)
					}
					return nil
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name: ";",
		Args: []string{"IO", "IO"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					a := AssertIO(self.args[0])

					aIO := a.Run(ioc)

					if aIO != nil {
						ew.Add(self.ctx.Error("unused return value of lhs"))
						return nil
					}

					return AssertIO(self.args[1]).Run(ioc)
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name: "=",
		Args: []string{"IO", "\\1"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			io := AssertIO(self.args[0])

			return NewIO(
				func(ioc IOContext) Value {
					runResult := io.Run(ioc)
					if runResult == nil {
						return nil
					}

					fn := AssertAnonFunc(self.args[1])
					res := fn.EvalRhs([]Value{runResult}, ew)
					if res == nil {
						return nil
					}

					res = ResolveIO(res, fn.Context(), ew)

					if res != nil {
						return AssertIO(res).Run(ioc)
					} else {
						return res
					}
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name: "readLine",
		Args: []string{},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()

					return NewString(scanner.Text(), self.ctx)
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name: "readArgs",
		Args: []string{},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					items := []Value{}
					for _, arg := range ARGS {
						items = append(items, NewString(arg, self.ctx))
					}

					return NewList(items, self.ctx)
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "read",
		Args:     []string{"File"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			p := AssertString(a.Args()[0])

			// TODO: should be relative to current path?
			fname := p.Value()

			return NewIO(
				func(ioc IOContext) Value {
					// check existence in a a
					if info, err := os.Stat(fname); os.IsNotExist(err) {
						return DeferFunc(self.links["Error"][0], []Value{NewString("\""+fname+"\" not found", self.ctx)}, self.ctx)
					} else if err != nil {
						return DeferFunc(self.links["Error"][0], []Value{NewString("\""+fname+"\" access error", self.ctx)}, self.ctx)
					} else if info.IsDir() {
						return DeferFunc(self.links["Error"][0], []Value{NewString("\""+fname+"\" is a directory", self.ctx)}, self.ctx)
					} else {
						data, err := ioutil.ReadFile(fname)
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString("\""+fname+"\" access error", self.ctx)}, self.ctx)
						} else {
							return NewString(string(data), self.ctx)
						}
					}
				},
				self.ctx,
			)
		},
	},
	// TODO: custom permissions
	BuiltinFuncConfig{
		Name:     "write",
		Args:     []string{"File", "String"},
		LinkReqs: []string{"Error", "Ok"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			b := self.args[1]
			data := AssertString(b)
			fname := AssertString(a.Args()[0]).Value()

			return NewIO(
				func(ioc IOContext) Value {
					if info, err := os.Stat(fname); os.IsNotExist(err) {
						err := ioutil.WriteFile(fname, []byte(data.Value()), 0644)
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString("couldn't write \""+fname+"\"", self.ctx)}, self.ctx)
						} else {
							return DeferFunc(self.links["Ok"][0], []Value{}, self.ctx)
						}
					} else if err != nil {
						return DeferFunc(self.links["Error"][0], []Value{NewString("can't write \""+fname+"\", access error", self.ctx)}, self.ctx)
					} else if info.IsDir() {
						return DeferFunc(self.links["Error"][0], []Value{NewString("can't write \""+fname+"\", already exists as directory", self.ctx)}, self.ctx)
					} else {
						return DeferFunc(self.links["Error"][0], []Value{NewString("can't write \""+fname+"\", already exists", self.ctx)}, self.ctx)
					}
				},
				self.ctx,
			)
		},
	},
	// TODO: custom permissions
	BuiltinFuncConfig{
		Name:     "overwrite",
		Args:     []string{"File", "String"},
		LinkReqs: []string{"Error", "Ok"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			b := self.args[1]
			data := AssertString(b)
			fname := AssertString(a.Args()[0]).Value()

			return NewIO(
				func(ioc IOContext) Value {
					info, err := os.Stat(fname)
					if err != nil && !os.IsNotExist(err) {
						return DeferFunc(self.links["Error"][0], []Value{NewString("can't write \""+fname+"\", access error", self.ctx)}, self.ctx)
					} else if info.IsDir() {
						return DeferFunc(self.links["Error"][0], []Value{NewString("can't write \""+fname+"\", already exists as directory", self.ctx)}, self.ctx)
					} else {
						err := ioutil.WriteFile(fname, []byte(data.Value()), 0644)
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString("couldn't write \""+fname+"\"", self.ctx)}, self.ctx)
						} else {
							return DeferFunc(self.links["Ok"][0], []Value{}, self.ctx)
						}
					}
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "send",
		Args:     []string{"HttpReq"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			method := AssertString(a.Args()[0]).Value()
			url := AssertString(a.Args()[1]).Value()
			payload := AssertString(a.Args()[2]).Value()

			return NewIO(
				func(ioc IOContext) Value {
					if method != "GET" && method != "POST" && method != "PUT" && method != "HEAD" && method != "DELETE" && method != "TRACE" && method != "OPTIONS" && method != "CONNECT" {
						return DeferFunc(self.links["Error"][0], []Value{NewString("unrecognized http method \""+method+"\"", a.Args()[0].Context())}, a.Args()[0].Context())
					}

					var payloadBytes io.Reader = nil
					if payload != "" {
						payloadBytes = bytes.NewBuffer([]byte(payload))
					}

					req, err := http.NewRequest(method, url, payloadBytes)
					if err != nil {
						return DeferFunc(self.links["Error"][0], []Value{NewString("invalid http request", self.ctx)}, self.ctx)
					}

					req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return DeferFunc(self.links["Error"][0], []Value{NewString("invalid http request to \""+url+"\"", self.ctx)}, self.ctx)
					} else if resp.StatusCode != 200 {
						return DeferFunc(self.links["Error"][0], []Value{NewString("http response error "+strconv.Itoa(resp.StatusCode), self.ctx)}, self.ctx)
					}

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return DeferFunc(self.links["Error"][0], []Value{NewString("http response payload error", self.ctx)}, self.ctx)
					}

					return NewString(string(body), self.ctx)
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "run",
		Args:     []string{"String"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					cmdRaw := AssertString(self.args[0]).Value()

					fields := strings.Fields(cmdRaw)
					if len(fields) < 1 {
						return DeferFunc(self.links["Error"][0], []Value{NewString("empty cmd", self.ctx)}, self.ctx)
					} else {
						cmdName := fields[0]
						args := fields[1:]

						if isRelPath(cmdName) {
							cmdName = filepath.Join(filepath.Dir(self.ctx.Path()), cmdName)
						}

						cmd := exec.Command(cmdName, args...)

						out, err := cmd.CombinedOutput()
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString(err.Error(), self.ctx)}, self.ctx)
						} else {
							return NewString(string(out), self.ctx)
						}
					}
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "run",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					cmdRaw := AssertString(self.args[0]).Value()
					stdinRaw := AssertString(self.args[1]).Value()

					fields := strings.Fields(cmdRaw)
					if len(fields) < 1 {
						return DeferFunc(self.links["Error"][0], []Value{NewString("empty cmd", self.ctx)}, self.ctx)
					} else {
						cmdName := fields[0]
						args := fields[1:]

						if isRelPath(cmdName) {
							cmdName = filepath.Join(filepath.Dir(self.ctx.Path()), cmdName)
						}

						cmd := exec.Command(cmdName, args...)
						stdin, err := cmd.StdinPipe()
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString(err.Error(), self.ctx)}, self.ctx)
						}

						go func() {
							defer stdin.Close()
							io.WriteString(stdin, stdinRaw)
						}()

						out, err := cmd.CombinedOutput()
						if err != nil {
							return DeferFunc(self.links["Error"][0], []Value{NewString(err.Error(), self.ctx)}, self.ctx)
						} else {
							return NewString(string(out), self.ctx)
						}
					}
				},
				self.ctx,
			)
		},
	},
}
