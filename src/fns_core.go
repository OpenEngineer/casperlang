package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// possible targets:
//  * default
//  * repl
var (
	TARGET = "default"
	PANIC  = ""
)

// basic builtin manipulation functions, should be avaiable on all target platforms
var builtinCoreFuncs []BuiltinFuncConfig = []BuiltinFuncConfig{
	BuiltinFuncConfig{
		Name: "Any",
		Args: []string{},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewAny(self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Bool",
		Args:     []string{},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "True",
		Args:     []string{},
		LinkReqs: []string{"Bool"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Bool"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "False",
		Args:     []string{},
		LinkReqs: []string{"Bool"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Bool"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Maybe",
		Args:     []string{},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Just",
		Args:     []string{"Any"},
		LinkReqs: []string{"Maybe"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Maybe"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Nothing",
		Args:     []string{},
		LinkReqs: []string{"Maybe"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Maybe"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Error",
		Args:     []string{"String"},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "Ok",
		Args:     []string{},
		LinkReqs: []string{"Any"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return DeferFunc(self.links["Any"][0], []Value{}, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "show",
		Args: []string{"Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewString(strconv.FormatInt(AssertInt(self.args[0]).Value(), 10), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "showf",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewString(fmt.Sprintf("%f", AssertFloat(self.args[0]).Value()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "+",
		Args: []string{"String", "String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewString(AssertString(self.args[0]).Value()+AssertString(self.args[1]).Value(), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "+",
		Args: []string{"[]", "[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return MergeLists(self.args[0], self.args[1], self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "+",
		Args: []string{"{}", "{}"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return MergeDicts(self.args[0], self.args[1], self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "toInt",
		Args:     []string{"String"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			i, err := strconv.ParseInt(AssertString(self.args[0]).Value(), 10, 64)
			if err != nil {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["Just"][0], []Value{NewInt(i, self.ctx)}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name: "toInt",
		Args: []string{"Float"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			i := int64(math.Round(AssertFloat(self.args[0]).Value()))
			return NewInt(i, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "toFloat",
		Args: []string{"Int"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			f := float64(AssertInt(self.args[0]).Value())
			return NewFloat(f, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name:     "toFloat",
		Args:     []string{"String"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			f, err := strconv.ParseFloat(AssertString(self.args[0]).Value(), 64)
			if err != nil {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["Just"][0], []Value{NewFloat(f, self.ctx)}, self.ctx)
			}
		},
	},
	// returns Maybe
	BuiltinFuncConfig{
		Name:     "get",
		Args:     []string{"[]", "Int"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			lst := AssertList(self.args[0])
			i := AssertInt(self.args[1])
			item := lst.Get(int(i.Value()))
			if item == nil {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["Just"][0], []Value{item}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "get",
		Args:     []string{"{}", "String"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			dict := AssertDict(self.args[0])
			k := AssertString(self.args[1])
			item := dict.Get(k.Value())
			if item == nil {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["Just"][0], []Value{item}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "get",
		Args:     []string{"String", "Int"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			str := AssertString(self.args[0])
			n := int64(len(str.Value()))
			i := AssertInt(self.args[1]).Value()
			if i < 0 {
				i += n
			}
			if i < 0 || i > n-1 {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				sub := str.Value()[i : i+1]
				return DeferFunc(self.links["Just"][0], []Value{NewString(sub, self.ctx)}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name: "len",
		Args: []string{"[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(int64(AssertList(self.args[0]).Len()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "len",
		Args: []string{"{}"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(int64(AssertDict(self.args[0]).Len()), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "len",
		Args: []string{"String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			return NewInt(int64(len(AssertString(self.args[0]).Value())), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "toInts",
		Args: []string{"String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			items := []Value{}
			str := AssertString(self.args[0]).Value()

			for _, r := range []rune(str) {
				items = append(items, NewInt(int64(r), self.ctx))
			}

			return NewList(items, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "toString",
		Args: []string{"([] Int)"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			lst := AssertList(self.args[0])

			rs := []rune{}

			items := lst.Items()

			for _, i_ := range items {
				i := AssertInt(i_).Value()
				rs = append(rs, rune(i))
			}

			return NewString(string(rs), self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "map",
		Args: []string{"\\1", "[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			fn := AssertAnonFunc(self.args[0])
			lst := AssertList(self.args[1])

			oldItems := lst.Items()

			newItems := []Value{}

			for _, oldItem := range oldItems {
				// a deferred call is fine here
				newItem := fn.EvalRhs([]Value{oldItem}, ew)
				if newItem == nil {
					return nil
				}

				newItems = append(newItems, newItem)
			}

			return NewList(newItems, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "map",
		Args: []string{"\\2", "[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			fn := AssertAnonFunc(self.args[0])
			lst := AssertList(self.args[1])

			oldItems := lst.Items()

			newItems := []Value{}

			for i, oldItem := range oldItems {
				newItem := fn.EvalRhs([]Value{NewInt(int64(i), self.ctx), oldItem}, ew)
				if newItem == nil {
					return nil
				}

				newItems = append(newItems, newItem)
			}

			return NewList(newItems, self.ctx)
		},
	},
	BuiltinFuncConfig{
		Name: "map",
		Args: []string{"\\2", "{}"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			fn := AssertAnonFunc(self.args[0])
			dct := AssertDict(self.args[1])

			keys := dct.Keys()
			vals := dct.Values()

			newItems := []Value{}

			for i, val := range vals {
				newItem := fn.EvalRhs([]Value{keys[i], val}, ew)
				if newItem == nil {
					return nil
				}

				newItems = append(newItems, newItem)
			}

			return NewList(newItems, self.ctx)
		},
	},
	// returns Just or Nothing
	BuiltinFuncConfig{
		Name:     "fold",
		Args:     []string{"\\2", "[]"},
		LinkReqs: []string{"Just", "Nothing"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			fn := AssertAnonFunc(self.args[0])
			lst := AssertList(self.args[1])

			var acc Value = nil

			oldItems := lst.Items()

			// XXX: in the lazy case this doesn't work because a the variables are overwritten
			for _, oldItem := range oldItems {
				if acc == nil {
					acc = oldItem
				} else {
					acc = EvalEager(fn.EvalRhs([]Value{acc, oldItem}, ew), ew)
					if acc == nil {
						return nil
					}
				}
			}

			if acc == nil {
				return DeferFunc(self.links["Nothing"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["Just"][0], []Value{acc}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name: "fold",
		Args: []string{"\\2", "Any", "[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			fn := AssertAnonFunc(self.args[0])
			acc := self.args[1]
			lst := AssertList(self.args[2])

			oldItems := lst.Items()

			for _, oldItem := range oldItems {
				acc = EvalEager(fn.EvalRhs([]Value{acc, oldItem}, ew), ew)
				if acc == nil {
					return nil
				}
			}

			return acc
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() == AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() != AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() < AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() > AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() <= AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		Args:     []string{"Int", "Int"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertInt(self.args[0]).Value() >= AssertInt(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() == AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() != AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() < AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() > AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() <= AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		Args:     []string{"Float", "Float"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertFloat(self.args[0]).Value() >= AssertFloat(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "==",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() == AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "!=",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() != AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() < AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() > AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     "<=",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() <= AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name:     ">=",
		Args:     []string{"String", "String"},
		LinkReqs: []string{"False", "True"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			if AssertString(self.args[0]).Value() >= AssertString(self.args[1]).Value() {
				return DeferFunc(self.links["True"][0], []Value{}, self.ctx)
			} else {
				return DeferFunc(self.links["False"][0], []Value{}, self.ctx)
			}
		},
	},
	BuiltinFuncConfig{
		Name: "panic",
		Args: []string{"String"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			msg := self.ctx.Error(AssertString(self.args[0]).Value()).Error()
			if TARGET == "repl" {
				PANIC = msg
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", msg)
				os.Exit(1)
			}
			return nil
		},
	},
	BuiltinFuncConfig{
		Name: "sort",
		Args: []string{"\\2", "[]"},
		Eval: func(self *BuiltinCall, ew ErrorWriter) Value {
			comp := AssertAnonFunc(self.args[0])
			lst := AssertList(self.args[1])

			s := NewListSorter(lst, comp, ew, self.ctx)
			sort.Stable(s)

			if !s.ew.Empty() {
				return nil
			} else {
				return NewList(s.items, self.ctx)
			}
		},
	},
}
