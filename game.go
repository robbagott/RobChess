package main

import "fmt"

// Start initiates a chess game with RobChess.
func Start() {
	posTest := NewPosition()
	fmt.Println(posTest)

}

func readMove() string {
	return ""
}

func algebraicToMove() {

}

// Move represents a move on the chess board. It encompasses a piece, the old square and the new square.
type Move struct {
	piece     GamePiece
	newSquare Square
	oldSquare Square
}
