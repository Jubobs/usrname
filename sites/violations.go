package sites

type Violation interface{}

type TooShort struct {
	Min, Actual int
}

type TooLong struct {
	Max, Actual int
}

type IllegalString struct {
	Lo, Hi int
}

type IllegalChars struct{} // TODO: refine later (indicate indices of illegal chars)
