package imdb

import (
	"io"

	"github.com/watuwat/imdb/lexer"
)

func Parser(input io.ReadCloser, builder func() Assigner) <-chan Assigner {
	lex := NewLexer(input)
	tokens := lexer.Process(lex.MainStateFn)
	values := make(chan Assigner, 1)

	keys := make([]string, 0)
	firstLine := true
	idx := 0

	var assigner Assigner

	go func() {
		defer close(values)
		defer input.Close()

		for token := range tokens {
			switch Type(token.Type()) {
			case Value:
				if firstLine {
					keys = append(keys, token.Value())
				} else {
					key := keys[idx]
					value := token.Value()

					assigner.Assign(key, value)

					idx++
				}
			case Tab:
				// ignore
			case Return:
				idx = 0
				if !firstLine {
					values <- assigner
				}

				assigner = builder()
				firstLine = false
			}
		}
	}()

	return values
}
