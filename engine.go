package main

import (
	"fmt"
	"math"
)

// Evaluate uses various heuristics to create a numeric evaluation of the position.
func Evaluate(p Position, side Side) float64 {
	var oppSide Side
	if side == White {
		oppSide = Black
	} else {
		oppSide = White
	}

	// For now, let's play like a child. Maximize material.
	sidePieces := p.GetPieces(side)
	sideSum := p.SumMaterial(sidePieces)

	oppPieces := p.GetPieces(oppSide)
	oppSum := p.SumMaterial(oppPieces)

	return sideSum - oppSum
}

// Calculate is an implementation of negaMax. Perhaps someday it will implement negaScout.
func Calculate(p Position, side Side, depth int) float64 {
	// Evaluate the position if we're at the max depth.
	if depth == 0 {
		return Evaluate(p, side)
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

	fmt.Println(moves)
	// Calculate possible moves
	for i := range moves {
		newPos := *p.Copy()
		newPos.MakeMove(moves[i])
		eval := Calculate(newPos, side, 2)
		if eval > bestSoFar {
			bestMoveSoFar = moves[i]
			bestSoFar = eval
		}
	}
	return bestMoveSoFar
}
