package token

type Token uint64

func NewToken() Token {
	return Token(1)
}
