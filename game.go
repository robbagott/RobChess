package main

import (
	"fmt"
)

// Start initiates a chess game with RobChess.
func Start() {
	fmt.Println("Welcome to RobChess! When entering moves, please use long algebraic chess notation.")
	fmt.Println(
		"Long algrebraic notation is the same as short algrebraic notation, except that instead\n" +
			"of entering a piece to move as the first symbol, the square that the piece resides on should\n" +
			"be used as an alternative. In addition, when a pawn is promoted, you must provide the piece the\n" +
			"pawn is being promoted to at the end (e.g. e7e8q).")
	posTest := NewPosition()
	fmt.Println(posTest)
	move := readMove()
	fmt.Println(move)
}

func readMove() Move {
	fmt.Print("Move: ")
	moveStr := scanMove()
	if move, ok := algebraicToMove(moveStr); ok {
		fmt.Println("What a nice move!")
		return move
	}
	fmt.Println("The move entered could not be understood. Please enter a move in algrebraic chess notation.")
	return readMove()
}

func scanMove() string {
	input := ""
	fmt.Scanln(&input)
	return input
}

func algebraicToMove(algMove string) (Move, bool) {
	// TODO implement
	return Move{}, false
}

// Move represents a move on the chess board. It encompasses a piece, the old square and the new square.
type Move struct {
	piece        GamePiece
	oRank, oFile int
	nRank, nFile int
}

// NewMove creates and initializes a new Move object.
func NewMove(piece GamePiece, oRank, oFile, nRank, nFile int) (Move, bool) {

	if oRank < 0 || oFile > 7 ||
		nRank < 0 || nFile > 7 {
		return Move{}, false
	}
	move := Move{
		piece,
		oRank,
		oFile,
		nRank,
		nFile}
	return move, true
}
