package lexer

import (
	"bufio"
	"io"
	"strings"
)

/*
Tokenize reads a io.Reader line-by-line and creates a stream of tokens.
It returns a channel of Tokens. The channel closes when the io.Reader
has been thoroughly read.
*/
func Tokenize(reader io.Reader) (tokenChannel chan Token) {
	tokenChannel = make(chan Token)
	go tokenize(reader, tokenChannel)
	return
}

// This does the actual tokenizing.
func tokenize(reader io.Reader, tokenChannel chan<- Token) {
	var (
		scanner *bufio.Scanner
		scan    string
	)
	scanner = bufio.NewScanner(reader)
	for scanner.Scan() {
		scan = strings.TrimSpace(scanner.Text())
		if scan != "" {
			tokenChannel <- Token{Name: Line, Args: []string{scanner.Text()}}
		}
	}
	close(tokenChannel)
}
