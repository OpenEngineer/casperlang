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

var IO_ACTIVE = true

func NewPathType(str Value, ctx Context) Type {
	if str == nil {
		str = NewValueData(NewStringType(), ctx)
	}

	return NewUserType(&AnyType{}, "Path", []Value{str}, NewBuiltinContext())
}

func NewPathValue(str Value, ctx Context) Value {
	return NewValueData(NewPathType(str, ctx), ctx)
}

func NewFileType(str Value, ctx Context) Type {
	if str == nil {
		str = NewValueData(NewStringType(), ctx)
	}

	return NewUserType(NewPathType(str, ctx), "File", []Value{str}, NewBuiltinContext())
}

func NewFileValue(str Value, ctx Context) Value {
	return NewValueData(NewFileType(str, ctx), ctx)
}

func NewDirType(str Value, ctx Context) Type {
	if str == nil {
		str = NewValueData(NewStringType(), ctx)
	}

	return NewUserType(NewPathType(str, ctx), "Dir", []Value{str}, NewBuiltinContext())
}

func NewDirValue(str Value, ctx Context) Value {
	return NewValueData(NewDirType(str, ctx), ctx)
}

func NewHttpReqType(method Value, url Value, payload Value, ctx Context) Type {
	if method == nil {
		method = NewValueData(NewStringType(), ctx)
	}

	if url == nil {
		url = NewValueData(NewStringType(), ctx)
	}

	if payload == nil {
		payload = NewValueData(NewStringType(), ctx)
	}

	return NewUserType(&AnyType{}, "HttpReq", []Value{method, url, payload}, NewBuiltinContext())
}

func NewHttpReqValue(method Value, url Value, payload Value, ctx Context) Value {
	return NewValueData(NewHttpReqType(method, url, payload, ctx), ctx)
}

var builtinIOFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:     "Path",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}
			return NewPathValue(a, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "File",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}
			return NewFileValue(a, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Dir",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}
			return NewDirValue(a, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "HttpReq",
		ArgTypes: []string{"String", "String", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			return NewHttpReqValue(self.args[0], self.args[1], self.args[2], self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "echo",
		ArgTypes: []string{"String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			s := self.args[0]
			if IsDeferredError(s) {
				fmt.Println("deferred error should'be been caught earlier")
				return s
			}

			if IO_ACTIVE {
				return NewIO(
					nil,
					func() Value {
						str := AssertString(s.Eval(scope, ew)).Value()
						if len(str) > 0 {
							fmt.Printf("%s", str)
						}
						return nil
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs([]Value{s}, NewIOType(nil), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ";",
		ArgTypes: []string{"IO", "IO"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {

			if IO_ACTIVE {
				return NewIO(
					nil,
					func() Value {
						a := AssertIO(self.args[0].Eval(scope, ew))

						aIO := a.Run()

						if aIO != nil && !IsVoidIOType(aIO.Type()) {
							ew.Add(self.ctx.Error("unused return value of lhs"))
							return nil
						}

						return AssertIO(self.args[1].Eval(scope, ew)).Run()
					},
					self.ctx,
				)
			} else {
				lhs := self.args[0]
				rhs := self.args[1]

				return self.UpdateArgs([]Value{lhs, rhs}, rhs.Type(), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "=",
		ArgTypes: []string{"IO", "\\1"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			if IO_ACTIVE {
				if len(self.args) == 3 {
					a := self.args[0].Eval(scope, ew)
					if a == nil {
						return nil
					}

					io := AssertIO(a)
					v := AssertVariable(self.args[2])

					return NewIO(
						nil,
						func() Value {
							runResult := io.Run()
							if runResult == nil {
								return nil
							}

							v.SetValue(runResult)

							// should be the result now, fn should return a side-effect, otherwise it is pointless
							res := self.args[1].Eval(scope, ew)
							if res == nil {
								return nil
							}

							return AssertIO(res).Run()
						},
						self.ctx,
					)
				} else {
					a := self.args[0].Eval(scope, ew)
					if a == nil {
						return nil
					}

					io := AssertIO(a)

					return NewIO(
						nil,
						func() Value {
							runResult := io.Run()
							if runResult == nil {
								return nil
							}

							fn_ := self.args[1].Eval(scope, ew)
							if fn_ == nil {
								return nil
							}

							fn := AssertFunc(fn_)
							res := fn.Call([]Value{runResult}, scope, self.ctx, ew)
							if res == nil {
								return nil
							}

							if !IsIOType(res.Type()) {
								ew.Add(res.Context().Error("expected IO, got " + res.Type().Dump()))
								return nil
							}

							return AssertIO(res).Run()
						},
						self.ctx,
					)
				}
			} else {
				a := self.args[0].Eval(scope, ew)
				if a == nil {
					return nil
				}

				io := AssertIOType(a.Type())
				if io.data == nil {
					ew.Add(a.Context().Error("expected non-void IO, got void IO"))
					return nil
				}

				fn_ := self.args[1].Eval(scope, ew)
				if fn_ == nil {
					return nil
				}

				fn := AssertFunc(fn_)

				v := NewVariable(io.data.Type(), self.ctx)

				fnRes := fn.Call([]Value{v}, scope, self.ctx, ew)
				if fnRes == nil {
					return nil
				}

				return self.UpdateArgs([]Value{a, fnRes, v}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "readLine",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			if IO_ACTIVE {
				return NewIO(
					NewValueData(NewStringType(), self.ctx),
					func() Value {
						scanner := bufio.NewScanner(os.Stdin)
						scanner.Scan()

						return NewString(scanner.Text(), self.ctx)
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs([]Value{}, NewIOType(NewValueData(NewStringType(), self.ctx)), self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "readArgs",
		ArgTypes: []string{},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			if IO_ACTIVE {
				return NewIO(
					NewValueData(
						NewListType(
							[]Value{NewValueData(NewStringType(), self.ctx)},
							self.ctx,
						),
						self.ctx,
					),
					func() Value {
						items := []Value{}
						for _, arg := range ARGS {
							items = append(items, NewString(arg, self.ctx))
						}

						return NewList(items, self.ctx)
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs(
					[]Value{},
					NewIOType(
						NewValueData(
							NewListType(
								[]Value{NewValueData(NewStringType(), self.ctx)},
								self.ctx,
							),
							self.ctx,
						),
					),
					self.ctx,
				)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "read",
		ArgTypes: []string{"File"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			if a == nil {
				return nil
			}

			if IO_ACTIVE {
				f := AssertUserType(a.Type(), "File")
				p := AssertString(f.args[0])
				fname := p.Value()

				return NewIO(
					nil,
					func() Value {
						// check existence in a a
						if info, err := os.Stat(fname); os.IsNotExist(err) {
							return NewErrorValue(NewString("\""+fname+"\" not found", self.ctx), self.ctx)
						} else if err != nil {
							return NewErrorValue(NewString("\""+fname+"\" access error", self.ctx), self.ctx)
						} else if info.IsDir() {
							return NewErrorValue(NewString("\""+fname+"\" is a directory", self.ctx), self.ctx)
						} else {
							data, err := ioutil.ReadFile(fname)
							if err != nil {
								return NewErrorValue(NewString("\""+fname+"\" access error", self.ctx), self.ctx)
							} else {
								return NewString(string(data), self.ctx)
							}
						}
					},
					self.ctx,
				)

			} else {
				return self.UpdateArgs([]Value{a}, nil, self.ctx)
			}
		},
	},
	// TODO: custom permissions
	BuiltinFuncConfig{
		Name:     "write",
		ArgTypes: []string{"File", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if a == nil || b == nil {
				return nil
			}

			if IO_ACTIVE {
				f := AssertUserType(a.Type(), "File")
				data := AssertString(b)
				fname := AssertString(f.args[0]).Value()

				return NewIO(
					nil,
					func() Value {
						if info, err := os.Stat(fname); os.IsNotExist(err) {
							err := ioutil.WriteFile(fname, []byte(data.Value()), 0644)
							if err != nil {
								return NewErrorValue(NewString("couldn't write \""+fname+"\"", self.ctx), self.ctx)
							} else {
								return NewOkValue(self.ctx)
							}
						} else if err != nil {
							return NewErrorValue(NewString("can't write \""+fname+"\", access error", self.ctx), self.ctx)
						} else if info.IsDir() {
							return NewErrorValue(NewString("can't write \""+fname+"\", already exists as directory", self.ctx), self.ctx)
						} else {
							return NewErrorValue(NewString("can't write \""+fname+"\", already exists", self.ctx), self.ctx)
						}
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	// TODO: custom permissions
	BuiltinFuncConfig{
		Name:     "overwrite",
		ArgTypes: []string{"File", "String"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)
			b := self.args[1].Eval(scope, ew)

			if a == nil || b == nil {
				return nil
			}

			if IO_ACTIVE {
				f := AssertUserType(a.Type(), "File")
				data := AssertString(b)
				fname := AssertString(f.args[0]).Value()

				return NewIO(
					nil,
					func() Value {
						info, err := os.Stat(fname)
						if err != nil && !os.IsNotExist(err) {
							return NewErrorValue(NewString("can't write \""+fname+"\", access error", self.ctx), self.ctx)
						} else if info.IsDir() {
							return NewErrorValue(NewString("can't write \""+fname+"\", already exists as directory", self.ctx), self.ctx)
						} else {
							err := ioutil.WriteFile(fname, []byte(data.Value()), 0644)
							if err != nil {
								return NewErrorValue(NewString("couldn't write \""+fname+"\"", self.ctx), self.ctx)
							} else {
								return NewOkValue(self.ctx)
							}
						}
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs([]Value{a, b}, nil, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "send",
		ArgTypes: []string{"HttpReq"},
		Eval: func(self *BuiltinCall, scope Scope, ew ErrorWriter) Value {
			a := self.args[0].Eval(scope, ew)

			if a == nil {
				panic("expected to be caught by Call")
			}

			if IO_ACTIVE {
				r := AssertUserType(a.Type(), "HttpReq")
				method := AssertString(r.args[0]).Value()
				url := AssertString(r.args[1]).Value()
				payload := AssertString(r.args[2]).Value()

				return NewIO(
					nil,
					func() Value {
						if method != "GET" && method != "POST" && method != "PUT" && method != "HEAD" && method != "DELETE" && method != "TRACE" && method != "OPTIONS" && method != "CONNECT" {
							return NewErrorValue(NewString("unrecognized http method \""+method+"\"", r.args[0].Context()), r.args[0].Context())
						}

						var payloadBytes io.Reader = nil
						if payload != "" {
							payloadBytes = bytes.NewBuffer([]byte(payload))
						}

						req, err := http.NewRequest(method, url, payloadBytes)
						if err != nil {
							return NewErrorValue(NewString("invalid http request", self.ctx), self.ctx)
						}

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return NewErrorValue(NewString("invalid http request to \""+url+"\"", self.ctx), self.ctx)
						} else if resp.StatusCode != 200 {
							return NewErrorValue(NewString("http response error "+strconv.Itoa(resp.StatusCode), self.ctx), self.ctx)
						}

						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							return NewErrorValue(NewString("http response payload error", self.ctx), self.ctx)
						}

						return NewString(string(body), self.ctx)
					},
					self.ctx,
				)
			} else {
				return self.UpdateArgs([]Value{a}, nil, self.ctx)
			}
		},
	},
}
