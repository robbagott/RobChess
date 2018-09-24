package main

import (
	"math"
)

// Evaluate uses various heuristics to create a numeric evaluation of the position.
func Evaluate(p Position) float64 {
	// For now, let's play like a child. Maximize material.
	whitePieces := p.GetPieces(White)
	whiteSum := p.SumMaterial(whitePieces)

	blackPieces := p.GetPieces(Black)
	blackSum := p.SumMaterial(blackPieces)

	return whiteSum - blackSum
}

// Calculate is an implementation of negaMax. Perhaps someday it will implement negaScout.
func Calculate(p Position, side Side, depth int) float64 {
	// Evaluate the position if we're at the max depth.
	if depth == 0 {
		return Evaluate(p)
	}

	var oppSide Side
	if side == White {
		oppSide = Black
	} else {
		oppSide = White
	}

	// Find possible moves
	moves := p.GetMoves(side)
	bestSoFar := math.Inf(-1)

	// Calculate possible moves
	for i := range moves {
		newPos := *p.Copy()
		newPos.MakeMove(moves[i])
		bestSoFar = math.Max(bestSoFar, -Calculate(newPos, oppSide, depth-1))
	}

	// From possible moves, choose optimal move. Return the optimal move with its evaluation.
	return bestSoFar
}

// Think finds the best move according to the evaluation function.
func Think(p Position, side Side) Move {
	// Find possible moves
	moves := p.GetMoves(side)
	bestSoFar := math.Inf(-1)
	var bestMoveSoFar Move

	// Calculate possible moves
	for i := range moves {
		newPos := *p.Copy()
		newPos.MakeMove(moves[i])
		eval := Calculate(newPos, side, 3)
		if eval > bestSoFar {
			bestMoveSoFar = moves[i]
			bestSoFar = eval
		}
	}
	return bestMoveSoFar
}
