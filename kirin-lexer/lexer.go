/*
The lexer program lexes the given KiriScript file and outputs the
result. It can also be loaded as a plugin for use in other
programs.

Command-line usage:

	lexer -i "path/to/input/file" -o "path/to/output/file"

can be used to lex the file at the given path. The lexer program also
supports reading from stdin; to do so, just omit the -i flag:

	lexer -o "path/to/output/file"

Omitting the -o flag will write the lexed content to stdout.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kirin-vn/lexer"
)

// These are our command-line flags.
var flags = struct {
	i string
	o string
}{}

// Initializes flags and returns a cmdFlags struct.
func initFlags() {
	flag.StringVar(&flags.i, "i", "", "Input file")
	flag.StringVar(&flags.o, "o", "", "Output file")
	flag.Parse()
}

// Serializes tokens in a readable format.
func serialize(tokenChannel <-chan lexer.Token, out *os.File) {
	var (
		outs  string
		token lexer.Token
		open  = true
		err   error
	)

	for {
		token, open = <-tokenChannel
		if !open {
			break
		}

		if int(token.Name) >= len(lexer.NameString) {
			panic(fmt.Errorf("%d is not a printable token name", token.Name))
		}
		outs = lexer.NameString[token.Name] + ": " + strings.Join(token.Args, ", ") + "\n"
		_, err = out.WriteString(outs)
		if err != nil {
			panic(err)
		}
	}
}

//  The main program. Not run when used as a plugin.
func main() {
	var (
		npt          *os.File
		out          *os.File
		err          error
		tokenChannel chan lexer.Token
	)

	// Initializes flags; see above:
	initFlags()

	// Getting our input File:
	if flags.i != "" {
		npt, err = os.Open(flags.i)
	} else {
		npt = os.Stdin
	}
	if err != nil {
		panic(err)
	}

	// Getting our output File:
	if flags.o != "" {
		out, err = os.Create(flags.o)
	} else {
		out = os.Stdout
	}
	if err != nil {
		panic(err)
	}

	// Tokenizing and serializing:
	tokenChannel = lexer.Tokenize(npt)
	serialize(tokenChannel, out)
}
