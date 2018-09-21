package main

import (
	"math/rand"
)

// Square represents a square in a chess position. Squares can have a piece placed on them.
type Square struct {
	piece GamePiece
}

// Position represents a chess position representation.
type Position struct {
	board [][]Square
}

// NewPosition creates and initializes a new Position with the starting arrangement of pieces.
func NewPosition() *Position {
	board := make([][]Square, 8, 8)

	for i := range board {
		board[i] = make([]Square, 8, 8)
	}

	pos := Position{board}

	pos.Reset()
	return &pos
}

// Evaluate takes a position and returns an evaluation of the position as a floating point number.
func (p *Position) Evaluate() float32 {
	return float32(rand.Intn(100))
}

// Reset resets the chess position to the starting chess arrangement.
func (p *Position) Reset() {
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			switch {
			case r == 0:
				switch {
				case f == 0 || f == 7:
					p.board[r][f].piece = GamePiece{Rook, White}
				case f == 1 || f == 6:
					p.board[r][f].piece = GamePiece{Knight, White}
				case f == 2 || f == 5:
					p.board[r][f].piece = GamePiece{Bishop, White}
				case f == 3:
					p.board[r][f].piece = GamePiece{Queen, White}
				case f == 4:
					p.board[r][f].piece = GamePiece{King, White}
				}
			case r == 1:
				p.board[r][f].piece = GamePiece{Pawn, White}
			case r == 2 || r == 3 || r == 4 || r == 5:
				p.board[r][f].piece = GamePiece{None, White}
			case r == 6:
				p.board[r][f].piece = GamePiece{Pawn, Black}
			case r == 7:
				switch {
				case f == 0 || f == 7:
					p.board[r][f].piece = GamePiece{Rook, Black}
				case f == 1 || f == 6:
					p.board[r][f].piece = GamePiece{Knight, Black}
				case f == 2 || f == 5:
					p.board[r][f].piece = GamePiece{Bishop, Black}
				case f == 3:
					p.board[r][f].piece = GamePiece{Queen, Black}
				case f == 4:
					p.board[r][f].piece = GamePiece{King, Black}
				}
			}
		}
	}
}

func (p *Position) String() string {
	boardPrint := ""

	boardPrint += "––––––––––––––––-----------------\n"
	for r := len(p.board) - 1; r >= 0; r-- {
		for f := range p.board[r] {
			gamePiece := p.board[r][f].piece
			boardPrint += "| " + gamePiece.String() + " "
			if f == 7 {
				boardPrint += "|\n––––––––––––––––-----------------\n"
			}
		}
	}
	return boardPrint
}
