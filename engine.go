package main

import (
	"math"
)

// Evaluate uses various heuristics to create a numeric evaluation of the position.
func Evaluate(p *Position, side Side) float64 {
	var oppSide = side.OppSide()

	// For now, let's play like a child. Maximize material.
	sidePieces := p.GetPieces(side)
	sideSum := p.SumMaterial(sidePieces)

	oppPieces := p.GetPieces(oppSide)
	oppSum := p.SumMaterial(oppPieces)

	return sideSum - oppSum
}

// Calculate is an implementation of negaMax. Perhaps someday it will implement negaScout.
/* Alpha is like a higher order bestSoFar variable. For the maximizer, It is the minimum score we are assured in other branches that we have calculated in parent nodes.
Therefore, if the minimizer in the current branch assures a worse score for us with any of its replies, we can give up on the current branch altogether as the maximizer.
This logic is somewhat muddied by the negamax take on minimax. Alpha typically tracks the maximizer's assured score and beta typically tracks the
minimizer's assured score. In negaMax, we negate the minimizer's result in the call to calculate which allows us to share the calculate function
between the two. In order for alpha and beta to work, their values must match with whether the minimizer or the maximizer is evaluating. Now, when
we pass from the maximizer to the minimizer, we give the minimizer beta as its alpha and vice versa. */
func Calculate(p *Position, side Side, depth int, alpha, beta float64) float64 {
	// Evaluate the position if we're at the max depth.
	if depth == 0 {
		return Evaluate(p, side)
	}

	var oppSide = side.OppSide()

	// Find possible moves
	moves := p.GetMoves(side)
	bestSoFar := math.Inf(-1)

	// Calculate possible moves
	for _, move := range moves {
		// Keep track of move details so we can roll back.
		oldPiece := p.board[move.oRank][move.oFile].piece
		capturedPiece := p.board[move.nRank][move.nFile].piece

		p.MakeMove(move)
		bestSoFar = math.Max(bestSoFar, -Calculate(p, oppSide, depth-1, -beta, -alpha))

		// Roll back move
		p.board[move.oRank][move.oFile].piece = oldPiece
		p.board[move.nRank][move.nFile].piece = capturedPiece

		alpha = math.Max(alpha, bestSoFar)
		if alpha >= beta {
			return bestSoFar
		}

		// if depth == 1 {
		// fmt.Printf("%v - %f\n", move, -Calculate(newPos, oppSide, depth-1))
		// }
		// if depth == 2 {
		// fmt.Printf("\t%v - %f\n", move, bestSoFar)
		// fmt.Println(newPos)
		// }
	}

	// From possible moves, choose optimal move. Return the optimal move with its evaluation.
	return bestSoFar
}

// Think finds the best move according to the evaluation function.
func Think(p Position, side Side) Move {
	var oppSide = side.OppSide()

	alpha := math.Inf(-1)
	beta := math.Inf(1)

	// Find possible moves
	moves := p.GetMoves(side)
	bestSoFar := math.Inf(-1)
	var bestMoveSoFar Move

	// Calculate possible moves
	for _, move := range moves {
		// Keep track of move details so we can roll back.
		oldPiece := p.board[move.oRank][move.oFile].piece
		capturedPiece := p.board[move.nRank][move.nFile].piece

		p.MakeMove(move)
		eval := -Calculate(&p, oppSide, 4, -beta, -alpha)
		if eval > bestSoFar {
			bestMoveSoFar = move
			bestSoFar = eval
		}

		alpha = math.Max(alpha, bestSoFar)

		// Roll back move.
		p.board[move.oRank][move.oFile].piece = oldPiece
		p.board[move.nRank][move.nFile].piece = capturedPiece
	}
	return bestMoveSoFar
}
