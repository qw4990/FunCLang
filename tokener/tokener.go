package tokener

import "io"

// Rule x
type Rule interface {
	RegExp() string
	Type() int
}

// Token x
type Token interface {
	Type() int
	Token() string
}

//Tokener x
type Tokener interface {
	Tokenize(input io.Reader) error // tokenizes data from input, all tokens uncustomed in cache would lose
	HasNext() bool                  // query if there are more tokens
	Next() (Token, error)           // Next pops and returns the next Token
	Peekth(kth int) (Token, error)  // Peekth returns the kth next Token, without pop

	/*
		 Seek sets the offset for the next call of Next(),
		 interpreted according to whence:
		 	0 means relative to the origin of the file,
			1 means relative to the current offset,
			2 means relative to the end;
	*/
	Seek(offset int, whence int) error
	Offset() int
}
