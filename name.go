package main

// Name is an int representing a token name.
type name int

// These are the different valid names.
const (
	Line name = iota
	Dialogue
)

// These are strings used by serialize().
var nameString = []string{"LINE", "DIALOGUE"}
