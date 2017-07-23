package lexer

// Name is an int representing a token name.
type Name int

// These are the different valid names.
const (
	Line Name = iota
	Dialogue
	NameMax = Dialogue
)

// These are strings used by serialize().
var nameString = []string{"LINE", "DIALOGUE"}

func (name Name) String() string {
	if name < 0 || name > NameMax {
		return "INVALID"
	}
	return nameString[name]
}
