package main

type ErrorList []error

func (l ErrorList) Len() int {
	return len(l)
}

func (l ErrorList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ErrorList) Less(i, j int) bool {
	iIsError := IsError(l[i])
	jIsError := IsError(l[j])

	if iIsError && jIsError {
		ei := AssertError(l[i])
		ej := AssertError(l[j])

		return ei.ctx.Before(ej.ctx)
	} else if iIsError {
		return false
	} else {
		// regular errors come first
		return true
	}
}
