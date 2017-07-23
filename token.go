package lexer

//  Token identifies a token with a name and any number of arguments
type Token struct {
	Name Name
	Args []string
}
