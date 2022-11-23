package data

var (
	noPrefix    int = 10 // want "not prefixed"
	_withPrefix int = 10
	Upper       int = 15
)

var noPrefix6 int = 10 // want "not prefixed"

const (
	noPrefix3    int = 10 // want "not prefixed"
	_withPrefix4 int = 10
)

const noPrefix2 int = 10 // want "not prefixed"
