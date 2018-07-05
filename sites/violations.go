package sites

import "unicode"

type Violation interface{}

type TooShort struct {
	Min, Actual int
}

type TooLong struct {
	Max, Actual int
}

type IllegalSubstring struct {
	Sub string
	At  int
}

type IllegalChars struct {
	At        []int
	Whitelist *unicode.RangeTable
}
