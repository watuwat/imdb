package imdb

import "fmt"

type Type int

const (
	_ Type = iota
	Value
	Tab
	Return
)

// Token is a struct to store basic information about each lexical
// information
type Token struct {
	typ   Type
	value string
}

// Value returns the value of token
func (t Token) Value() string {
	return t.value
}

// Type returns the token type in an integer value
func (t Token) Type() int {
	return int(t.typ)
}

func (t Token) String() string {
	typ := ""
	val := ""

	switch t.typ {
	case Return:
		typ = "RETURN"
		val = ""
	case Value:
		typ = "VALUE"
		val = t.value
	case Tab:
		typ = "TAB"
		val = ""

	}

	return fmt.Sprintf("Type: %s, Value: %s", typ, val)
}

// NewToken returns new tokens based on basic info
func NewToken(typ Type, value string) *Token {
	return &Token{
		typ:   typ,
		value: value,
	}
}
