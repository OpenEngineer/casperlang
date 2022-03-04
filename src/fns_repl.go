package main

import (
	"fmt"
)

var replHelpMessage = `help             display this message
import <string>  import a module
dir              list def names
quit             quit this program`

var builtinReplFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name:    "help",
		Args:    []string{},
		Targets: []string{"repl"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					fmt.Fprintf(ioc.Stdout(), "%s", replHelpMessage)
					return nil
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:    "dir",
		Args:    []string{},
		Targets: []string{"repl"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					// only unique names
					strs := AssertReplIOContext(ioc).ListNames()

					items := []Value{}

					for _, str := range strs {
						items = append(items, NewString(str, self.Context()))
					}

					return NewList(items, self.Context())
				},
				self.ctx,
			)
		},
	},
	BuiltinFuncConfig{
		Name:    "quit",
		Args:    []string{},
		Targets: []string{"repl"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewIO(
				func(ioc IOContext) Value {
					AssertReplIOContext(ioc).repl.Quit()
					return nil
				},
				self.ctx,
			)
		},
	},
}
