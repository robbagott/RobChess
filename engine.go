package main

import (
	"fmt"
	"math"
	"sort"
)

var numEvals = 0

// GameTree represents our series of calculations thus far. The tree also indexes its nodes for quick lookups.
type GameTree struct {
	parent   *GameTree
	children []*GameTree
	move     Move
	eval     float64
}

// GameContext Represents a chess game.
type GameContext struct {
	position Position
	moves    []Move
	gameTree *GameTree
}

// NewGame creates a new chess game.
func NewGame() *GameContext {
	gameTree := GameTree{nil, make([]*GameTree, 0), Move{}, 0}
	return &GameContext{*NewPosition(), make([]Move, 0), &gameTree}
}

// MakeMove makes a move in the game and records it.
func (g *GameContext) MakeMove(move Move) bool {
	if ok := g.position.MakeMove(move); ok {
		// Add move
		g.moves = append(g.moves, move)

		// Throw away unused game tree.
		for _, child := range g.gameTree.children {
			if child.move == move {
				g.gameTree = child
			}
		}
		return true
	}
	return false
}

// Evaluate uses various heuristics to create a numeric evaluation of the position.
func Evaluate(p *Position, side Side) float64 {
	numEvals++
	// For now, let's play like a child. Maximize material.
	sidePieces := p.GetPieces(side)
	sideSum := p.SumMaterial(sidePieces)

	oppPieces := p.GetPieces(side.OppSide())
	oppSum := p.SumMaterial(oppPieces)

	centralControl := p.CentralControl(side) * .1
	oppCentralControl := p.CentralControl(side.OppSide()) * .1

	return sideSum - oppSum + centralControl - oppCentralControl
}

// Calculate is an implementation of negaMax. Perhaps someday it will implement negaScout.
/* Alpha is like a higher order bestSoFar variable. For the maximizer, it is the minimum score we are assured in other branches that we have calculated in parent nodes.
Therefore, if the minimizer in the current branch assures a worse score for us with any of its replies, we can give up on the current branch altogether as the maximizer.
This logic is somewhat muddied by the negamax take on minimax. Alpha typically tracks the maximizer's assured score and beta typically tracks the
minimizer's assured score. In negaMax, we negate the minimizer's result in the call to Calculate() which allows us to share the calculate function
between the two. In order for alpha and beta to work, their values must match with whether the minimizer or the maximizer is evaluating. Now, when
we pass from the maximizer to the minimizer, we give the minimizer beta as its alpha and vice versa. */
func Calculate(p *Position, side Side, depth int, alpha, beta float64, node *GameTree) float64 {
	// Evaluate the position if we're at the max depth.
	if depth == 0 {
		return Evaluate(p, side)
	}

	// Check if there are moves on the node. If not, retrieve them and add them to the node.
	if len(node.children) == 0 {
		moves := p.GetMoves(side)
		for _, move := range moves {
			node.children = append(node.children, &GameTree{node, make([]*GameTree, 0), move, 0})
		}
	}

	// Calculate possible moves
	bestSoFar := math.Inf(-1)
	for _, child := range node.children {
		move := child.move
		// Keep track of move details so we can roll back.
		oldPiece := p.board[move.oRank][move.oFile].piece
		capturedPiece := p.board[move.nRank][move.nFile].piece

		p.MakeMove(move)
		child.eval = -Calculate(p, side.OppSide(), depth-1, -beta, -alpha, child)
		bestSoFar = math.Max(bestSoFar, child.eval)

		// Roll back move
		p.board[move.oRank][move.oFile].piece = oldPiece
		p.board[move.nRank][move.nFile].piece = capturedPiece

		alpha = math.Max(alpha, bestSoFar)
		if alpha >= beta {
			node.children = sortMoves(node.children)
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
	node.children = sortMoves(node.children)
	return bestSoFar
}

// Think finds the best move according to the evaluation function.
func Think(g GameContext, side Side) Move {
	// Get an initial move
	strongestMove := thinkDepth(g, side, 1)
	for i := 1; i < 7; i++ {
		fmt.Printf("Thinking to depth %d\n", i)
		strongestMove = thinkDepth(g, side, i)
	}
	fmt.Printf("Evaluated %d positions\n", numEvals)
	numEvals = 0
	return strongestMove
}

func thinkDepth(g GameContext, side Side, depth int) Move {
	p := g.position

	// Check if there are moves on the node. If not, retrieve them and add them to the node.
	if len(g.gameTree.children) == 0 {
		moves := p.GetMoves(side)
		for _, move := range moves {
			g.gameTree.children = append(g.gameTree.children, &GameTree{g.gameTree, make([]*GameTree, 0), move, 0})
		}
	}

	// Calculate possible moves
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	bestSoFar := math.Inf(-1)
	var bestMoveSoFar Move
	for _, child := range g.gameTree.children {
		move := child.move

		// Keep track of move details so we can roll back.
		oldPiece := p.board[move.oRank][move.oFile].piece
		capturedPiece := p.board[move.nRank][move.nFile].piece

		p.MakeMove(move)
		eval := -Calculate(&p, side.OppSide(), depth, -beta, -alpha, child)
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

func sortMoves(nodes []*GameTree) []*GameTree {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].eval > nodes[j].eval
	})
	return nodes
}
