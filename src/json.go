package main

// json is valid casper, so the parser can be used to read json into a value and give nice error messages in the proves

func EvalJSON(s *Source, ew ErrorWriter) Value {
	ts := Tokenize(s, ew)
	if !ew.Empty() {
		return nil
	}

	// strip initial NL
	if len(ts) > 0 && IsNL(ts[0]) {
		ts = ts[1:]
	}

	v := ParseExpr(ts, ew)
	if !ew.Empty() {
		return nil
	}

	scope := FillJSONScope()

	return v.Eval(scope, ew)
}
