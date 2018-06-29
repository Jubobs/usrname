package sites

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

type IllegalChars struct{} // TODO: refine later (indicate indices of illegal chars)
