package lexer

// Name is an enum representing a token name.
type Name int

// These are the different valid names.
const (
	Line Name = iota
	Dialogue
	nameMax = iota - 1
)

// These are strings used by String().
var nameString = []string{"LINE", "DIALOGUE"}

// String() gets the string representation for a Name.
func (name Name) String() string {
	if name < 0 || name > nameMax {
		return "INVALID"
	}
	return nameString[name]
}
