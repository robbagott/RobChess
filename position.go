package main

import (
	"math"
	"strconv"
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

	boardPrint += "   ––––––––––––––––-----------------\n"
	for r := len(p.board) - 1; r >= 0; r-- {
		boardPrint += " " + strconv.Itoa(r+1) + " "
		for f := range p.board[r] {
			gamePiece := p.board[r][f].piece
			boardPrint += "| " + gamePiece.String() + " "
			if f == 7 {
				boardPrint += "|\n   ––––––––––––––––-----------------\n"
			}
		}
	}
	boardPrint += "     a   b   c   d   e   f   g   h\n"
	return boardPrint
}

// GetMoves returns the set of moves that are possible for the side indicated.
func (p *Position) GetMoves(side Side) []Move {
	moves := make([]Move, 0, 20)

	// Only pieces can make moves in chess, so we iterate through the board and check for pieces.
	// If a piece is found, we then check for legal moves for that piece.
	for r := range p.board {
		for f := range p.board[r] {
			moves = append(moves, p.GetMovesAt(f, r)...)
		}
	}
	return moves
}

// GetMovesAt returns the set of moves that are possible for the piece located at file f and rank r
func (p *Position) GetMovesAt(f, r int) []Move {
	moves := make([]Move, 0, 20)
	piece := p.board[r][f].piece
	switch piece.piece {
	case Pawn:
		moves = append(moves, p.getPawnMoves(f, r, piece.color)...)
	case Rook:
		// do something
	case Knight:
		// do something
	case Bishop:
		// do something
	case Queen:
		// do something
	case King:
		// do something
	}
	return moves
}

// TODO account for en passant
// Gets possible pawn moves starting at a specific square.
func (p *Position) getPawnMoves(f, r int, side Side) []Move {
	// A pawn can have a maximum of 4 moves (on promotion)
	moves := make([]Move, 0, 4)
	if side == White {
		// Possible forward moves
		if r == 1 {
			moves = append(moves, Move{f, r, f, r + 1, ""}, Move{f, r, f, r + 2, ""})
		} else if r > 1 && r < 6 {
			moves = append(moves, Move{f, r, f, r + 1, ""})
		}

		// Possible captures
		if f-1 >= 0 {
			if p.board[r+1][f-1].piece.piece != None && p.board[r+1][f-1].piece.color == Black {
				moves = append(moves, Move{f, r, f - 1, r + 1, ""})
			}
		} else if f+1 <= 7 {
			if p.board[r+1][f+1].piece.piece != None && p.board[r+1][f+1].piece.color == Black {
				moves = append(moves, Move{f, r, f + 1, r + 1, ""})
			}
		}
	} else if side == Black {
		// Possible forward moves
		if r == 1 {
			moves = append(moves, Move{f, r, f, r + 1, ""}, Move{f, r, f, r + 2, ""})
		} else if r > 1 && r < 6 {
			moves = append(moves, Move{f, r, f, r + 1, ""})
		}

		// Possible captures
		if f-1 >= 0 {
			if p.board[r-1][f-1].piece.piece != None && p.board[r-1][f-1].piece.color == Black {
				moves = append(moves, Move{f, r, f - 1, r - 1, ""})
			}
		} else if f+1 <= 7 {
			if p.board[r-1][f+1].piece.piece != None && p.board[r-1][f+1].piece.color == Black {
				moves = append(moves, Move{f, r, f + 1, r - 1, ""})
			}
		}
	}
	return moves
}

func (p *Position) getRookMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)

}

// GetPieces returns an array of pieces for the side indicated in a position.
func (p *Position) GetPieces(side Side) []GamePiece {
	pieces := make([]GamePiece, 0, 32)
	for r := range p.board {
		for f := range p.board[r] {
			piece := p.board[r][f].piece
			if piece.piece != None && piece.color == side {
				pieces = append(pieces, piece)
			}
		}
	}
	return pieces
}

// SumMaterial performs a rudimentary sum of the material using the classic chess piece values.
func (p *Position) SumMaterial(pieces []GamePiece) {
	var sum float64
	for i := range pieces {
		if pieces[i].piece != None && pieces[i].piece != King {
			sum += pieces[i].Value()
		}
	}
}

// Value returns the value of the piece. For now, the value is the classical chess piece value.
func (p GamePiece) Value() float64 {
	switch p.piece {
	case Pawn:
		return 1
	case Rook:
		return 5
	case Knight:
		return 3
	case Bishop:
		return 3
	case Queen:
		return 9
	case King:
		return math.Inf(1)
	default:
		return 1
	}
}
