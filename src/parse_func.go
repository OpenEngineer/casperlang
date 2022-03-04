package main

func parseFunc(ts []Token, ew ErrorWriter) *UserFunc {
	var reader FileReader = NewFuncNameReader([]Token{NewNL(0, NewBuiltinContext())})

	dummy := NewReplFile()

	for _, t := range ts {
		if !IsNL(t) {
			reader, dummy = reader.Ingest(dummy, t, ew)
			if !ew.Empty() {
				return nil
			}
		}
	}

	dummy = reader.Finalize(dummy, ew)
	if !ew.Empty() {
		return nil
	}

	return dummy.fns[0].fn.(*UserFunc)
}
