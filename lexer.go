package imdb

import (
	"io"

	"github.com/watuwat/imdb/lexer"
)

type Lexer struct {
	lex lexer.Lexer
}

func (l *Lexer) MainStateFn(emitter lexer.Emitter) lexer.StateFn {
	switch l.lex.Peek() {
	case 0:
		return nil
	case '\t':
		return l.TabStateFn
	case '\n':
		return l.ReturnStateFn
	default:
		return l.StringStateFn
	}
}

func (l *Lexer) TabStateFn(emitter lexer.Emitter) lexer.StateFn {
	l.lex.Next()
	l.lex.Ignore()

	emitter.Emit(&Token{
		typ:   Tab,
		value: "",
	})

	return l.MainStateFn
}

func (l *Lexer) ReturnStateFn(emitter lexer.Emitter) lexer.StateFn {
	l.lex.Next()
	l.lex.Ignore()

	emitter.Emit(&Token{
		typ:   Return,
		value: "",
	})

	return l.MainStateFn
}

func (l *Lexer) StringStateFn(emitter lexer.Emitter) lexer.StateFn {
	l.lex.AcceptRunUntil("\t\n")

	val := l.lex.CurrentString()
	l.lex.Ignore()

	emitter.Emit(&Token{
		typ:   Value,
		value: val,
	})

	return l.MainStateFn
}

func NewLexer(input io.Reader) *Lexer {
	return &Lexer{
		lex: lexer.NewStreamLexer(input, 1024),
	}
}
