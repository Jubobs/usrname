package usrname

import (
	"fmt"
	"unicode"
)

type Violation interface{}

type TooShort struct {
	Min, Actual int
}

func (v *TooShort) String() string {
	const templ = "&TooShort{Min: %d, Actual: %d}"
	return fmt.Sprintf(templ, v.Min, v.Actual)
}

type TooLong struct {
	Max, Actual int
}

func (v *TooLong) String() string {
	const templ = "&TooLong{Max: %d, Actual: %d}"
	return fmt.Sprintf(templ, v.Max, v.Actual)
}

type IllegalSubstring struct {
	At      []int
	Pattern string
}

func (v *IllegalSubstring) String() string {
	const templ = "&IllegalSubstring{%q}"
	return fmt.Sprintf(templ, v.Pattern)
}

type IllegalPrefix struct {
	Pattern string
}

func (v *IllegalPrefix) String() string {
	const templ = "&IllegalPrefix{%q}"
	return fmt.Sprintf(templ, v.Pattern)
}

type IllegalSuffix struct {
	Pattern string
}

func (v *IllegalSuffix) String() string {
	const templ = "&IllegalSuffix{%q}"
	return fmt.Sprintf(templ, v.Pattern)
}

type IllegalChars struct {
	At        []int
	Whitelist *unicode.RangeTable
}

func (v *IllegalChars) String() string {
	const templ = "&IllegalChars{%v}"
	return fmt.Sprintf(templ, v.At)
}
