package common

type ValueSource interface {
	// Source is the name/key of the value
	Source() string
	// Orig is the original string value
	Orig() string
	// Values is the split list from original value
	Values() []string
}
