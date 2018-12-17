package lexer

//Lexer base common interface
type Lexer interface {
	Ignore()
	Backup()
	Peek() rune
	PeekNth(int) rune
	Next() rune
	Current() rune
	CurrentString() string
	Accept(string) bool
	AcceptRun(string)
	AcceptRunUntil(string)
}

type Token interface {
	Type() int
	Value() string
}

type Emitter interface {
	Emit(Token)
}

// StateFn returns a lexer state function
type StateFn func(Emitter) StateFn

type emitter chan Token

func (e emitter) Emit(token Token) {
	e <- token
}

// Process process
func Process(stateFn StateFn) <-chan Token {
	tokens := make(chan Token, 1)

	emitter := emitter(tokens)

	go func() {
		defer close(tokens)
		for stateFn != nil {
			stateFn = stateFn(emitter)
		}
	}()

	return tokens
}
