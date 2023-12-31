package common

type Location struct {
	Original string
	Expanded string
}

type PathSetIn struct {
	What     Location
	WhereSet []Location // file it's potentially set in
}

type CandidateSource interface {
	WhereSet(somePath string) *PathSetIn
}
