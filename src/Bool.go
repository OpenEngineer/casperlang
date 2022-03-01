package main

func GetBoolValue(v Value, ew ErrorWriter) (bool, bool) {
	isTrue, isBool := false, false

	EvalUntil(v, func(tn string) bool {
		if tn == "True" {
			isTrue = true
			isBool = true
			return true
		} else if tn == "False" {
			isTrue = false
			isBool = true
			return true
		} else {
			return false
		}
	}, ew)

	return isTrue, isBool
}
