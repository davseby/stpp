package data

var (
	noPrefix    int = 10 // want "not prefixed with \"_\""
	_withPrefix int = 10
	Upper       int = 15
)

var noPrefix6 int = 10 // want "not prefixed with \"_\""

const (
	noPrefix3    int = 10 // want "not prefixed with \"_\""
	_withPrefix4 int = 10
)

const noPrefix2 int = 10 // want "not prefixed with \"_\""

func name() int {
	var a = 5
	b := 5

	return a + b
}
