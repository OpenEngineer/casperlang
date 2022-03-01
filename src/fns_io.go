package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// XXX: are the evals here redundant, because they are already done in Call?

var builtinIOFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "Path",
		ArgTypes: []string{"String"},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "File",
		ArgTypes: []string{"String"},
		LinkReqs: []string{"Path"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Path"][0], []Value{self.args[0]}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Dir",
		ArgTypes: []string{"String"},
		LinkReqs: []string{"Path"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Path"][0], []Value{self.args[0]}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "HttpReq",
		ArgTypes: []string{"String", "String", "String"},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "echo",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func() Value {
					str := AssertString(self.args[0]).Value()
					if len(str) > 0 {
						fmt.Printf("%s", str)
					}
					return nil
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     ";",
		ArgTypes: []string{"IO", "IO"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func() Value {
					a := AssertIO(self.args[0])

					aIO := a.Run()

					if aIO != nil {
						ew.Add(self.ctx.Error("unused return value of lhs"))
						return nil
					}

					return AssertIO(self.args[1]).Run()
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "=",
		ArgTypes: []string{"IO", "\\1"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			io := AssertIO(self.args[0])

			return NewIO(
				func() Value {
					runResult := io.Run()
					if runResult == nil {
						return nil
					}

					fn := AssertFunc(self.args[1])
					res := RunFunc(fn, []Value{runResult}, ew, self.ctx)
					if res == nil {
						return nil
					}

					res = ResolveIO(res, fn.Context(), ew)

					if res != nil {
						return AssertIO(res).Run()
					} else {
						return res
					}
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "readLine",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func() Value {
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()

					return NewString(scanner.Text(), self.ctx)
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:     "readArgs",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func() Value {
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
		ArgTypes: []string{"File"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			p := AssertString(a.Args()[0])

			fname := p.Value()

			return NewIO(
				func() Value {
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
		ArgTypes: []string{"File", "String"},
		LinkReqs: []string{"Error", "Ok"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			b := self.args[1]
			data := AssertString(b)
			fname := AssertString(a.Args()[0]).Value()

			return NewIO(
				func() Value {
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
		ArgTypes: []string{"File", "String"},
		LinkReqs: []string{"Error", "Ok"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			b := self.args[1]
			data := AssertString(b)
			fname := AssertString(a.Args()[0]).Value()

			return NewIO(
				func() Value {
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
		ArgTypes: []string{"HttpReq"},
		LinkReqs: []string{"Error"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			a := self.args[0].(Call)
			method := AssertString(a.Args()[0]).Value()
			url := AssertString(a.Args()[1]).Value()
			payload := AssertString(a.Args()[2]).Value()

			return NewIO(
				func() Value {
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
}
