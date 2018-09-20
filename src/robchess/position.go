package robchess

import (
	"math/rand"
)

// Position represents a chess position representation.
type Position struct {
	something int
}

// Evaluate takes a position and returns an evaluation of the position.
func (p *Position) Evaluate() float32 {
	return float32(rand.Intn(100))
}
