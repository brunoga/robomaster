package token

import "sync/atomic"

type Generator struct {
	token Token
}

func NewGenerator() *Generator {
	return &Generator{
		token: NewToken(),
	}
}

func (g *Generator) Next() Token {
	for {
		currentValue := atomic.LoadUint64((*uint64)(&g.token))

		nextValue := currentValue + 1
		if nextValue == 0 {
			nextValue = 1
		}

		if atomic.CompareAndSwapUint64((*uint64)(&g.token), currentValue, nextValue) {
			return Token(currentValue)
		}

		// If we get here, the value was changed by another goroutine, so we
		// retry. Theoretically this could loop forever but, in practice, it
		// is unlikely to happen.
	}
}
