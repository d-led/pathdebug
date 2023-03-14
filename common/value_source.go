package common

type ValueSource interface {
	Source() string
	Orig() string
	Values() []string
}
