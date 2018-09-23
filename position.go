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
		moves = append(moves, p.getRookMoves(f, r, piece.color)...)
	case Knight:
		moves = append(moves, p.getKnightMoves(f, r, piece.color)...)
	case Bishop:
		moves = append(moves, p.getBishopMoves(f, r, piece.color)...)
	case Queen:
		moves = append(moves, p.getQueenMoves(f, r, piece.color)...)
	case King:
		moves = append(moves, p.getKingMoves(f, r, piece.color)...)
	}
	return moves
}

// TODO account for en passant
// Gets possible pawn moves starting at a specific square.
func (p *Position) getPawnMoves(f, r int, side Side) []Move {
	// A pawn can have a maximum of 4 moves (on promotion)
	moves := make([]Move, 0, 4)

	// Define rank increment direction.
	var rIncr int
	if side == White {
		rIncr = 1
	} else {
		rIncr = -1
	}

	// Possible forward moves
	if r == 1 {
		moves = append(moves, Move{f, r, f, r + rIncr, ""}, Move{f, r, f, r + rIncr*2, ""})
	} else if r > 1 && r < 6 {
		moves = append(moves, Move{f, r, f, r + rIncr, ""})
	}

	// Possible captures
	if f-1 >= 0 {
		if p.board[r+1][f-1].piece.piece != None && p.board[r+1][f-1].piece.color != side {
			moves = append(moves, Move{f, r, f - 1, r + rIncr, ""})
		}
	} else if f+1 <= 7 {
		if p.board[r+1][f+1].piece.piece != None && p.board[r+1][f+1].piece.color != side {
			moves = append(moves, Move{f, r, f + 1, r + rIncr, ""})
		}
	}
	return moves
}

// Get possible rook moves for a rook located at file r and rank f of color side.
func (p *Position) getRookMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)

	// Look right
	for i := f + 1; i <= 7; i++ {
		if canMoveToSquare(*p, i, r, side) {
			moves = append(moves, Move{f, r, i, r, ""})
		} else {
			break
		}
	}

	// Look left
	for i := f - 1; i >= 0; i-- {
		if canMoveToSquare(*p, i, r, side) {
			moves = append(moves, Move{f, r, i, r, ""})
		} else {
			break
		}
	}

	// Look forward
	for i := r + 1; i <= 7; i++ {
		if canMoveToSquare(*p, f, i, side) {
			moves = append(moves, Move{f, r, f, i, ""})
		} else {
			break
		}
	}

	// Look backward
	for i := r - 1; i >= 0; i-- {
		if canMoveToSquare(*p, f, i, side) {
			moves = append(moves, Move{f, r, f, i, ""})
		} else {
			break
		}
	}
	return moves
}

// Get possible bishop moves for a bishop located at file r and rank f of color side.
func (p *Position) getBishopMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)
	// Look diagonally forward-right
	for i, j := r+1, f+1; i < 8 && j < 8; i, j = r+1, f+1 {
		if canMoveToSquare(*p, j, i, side) {
			moves = append(moves, Move{f, r, j, i, ""})
		} else {
			break
		}
	}
	// Look diagonally forward-left
	for i, j := r+1, f-1; i < 8 && j >= 0; i, j = r+1, f-1 {
		if canMoveToSquare(*p, j, i, side) {
			moves = append(moves, Move{f, r, j, i, ""})
		} else {
			break
		}
	}
	// Look diagonally backward-right
	for i, j := r-1, f+1; i >= 0 && j < 8; i, j = r-1, f+1 {
		if canMoveToSquare(*p, j, i, side) {
			moves = append(moves, Move{f, r, j, i, ""})
		} else {
			break
		}
	}
	// Look diagonally backward-left
	for i, j := r-1, f-1; i >= 0 && j >= 0; i, j = r-1, f-1 {
		if canMoveToSquare(*p, j, i, side) {
			moves = append(moves, Move{f, r, j, i, ""})
		} else {
			break
		}
	}
	return moves
}

// Get possible queen moves for a queen located at file r and rank f of color side.
func (p *Position) getKnightMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 8)

	// Right L moves
	if canMoveToSquare(*p, f+2, r+1, side) {
		moves = append(moves, Move{f, r, f + 2, r + 1, ""})
	}
	if canMoveToSquare(*p, f+2, r-1, side) {
		moves = append(moves, Move{f, r, f + 2, r - 1, ""})
	}

	// Left L moves
	if canMoveToSquare(*p, f-2, r+1, side) {
		moves = append(moves, Move{f, r, f - 2, r + 1, ""})
	}
	if canMoveToSquare(*p, f-2, r-1, side) {
		moves = append(moves, Move{f, r, f - 2, r - 1, ""})
	}

	// Forward L moves
	if canMoveToSquare(*p, f+1, r+2, side) {
		moves = append(moves, Move{f, r, f + 1, r + 2, ""})
	}
	if canMoveToSquare(*p, f-1, r+2, side) {
		moves = append(moves, Move{f, r, f - 1, r + 2, ""})
	}

	// Backward L moves
	if canMoveToSquare(*p, f+1, r-2, side) {
		moves = append(moves, Move{f, r, f + 1, r - 2, ""})
	}
	if canMoveToSquare(*p, f-1, r-2, side) {
		moves = append(moves, Move{f, r, f - 1, r - 2, ""})
	}
	return moves
}

// Get possible queen moves for a queen located at file r and rank f of color side.
func (p *Position) getQueenMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 30)
	moves = append(moves, p.getRookMoves(f, r, side)...)
	moves = append(moves, p.getBishopMoves(f, r, side)...)
	return moves
}

// Get possible king moves for a king located at file r and rank f of color side.
func (p *Position) getKingMoves(f, r int, side Side) []Move {
	// Look at adjacent squares
	// TODO Evaluate checks
	moves := make([]Move, 0, 8)
	// Right
	if canMoveToSquare(*p, f+1, r, side) {
		moves = append(moves, Move{f, r, f + 1, r, ""})
	}
	// Back-right
	if canMoveToSquare(*p, f+1, r-1, side) {
		moves = append(moves, Move{f, r, f + 1, r - 1, ""})
	}
	// Back
	if canMoveToSquare(*p, f, r-1, side) {
		moves = append(moves, Move{f, r, f, r - 1, ""})
	}
	// Back-left
	if canMoveToSquare(*p, f-1, r-1, side) {
		moves = append(moves, Move{f, r, f - 1, r - 1, ""})
	}
	// Left
	if canMoveToSquare(*p, f-1, r, side) {
		moves = append(moves, Move{f, r, f - 1, r, ""})
	}
	// Forward-left
	if canMoveToSquare(*p, f-1, r+1, side) {
		moves = append(moves, Move{f, r, f - 1, r + 1, ""})
	}
	// Forward
	if canMoveToSquare(*p, f, r+1, side) {
		moves = append(moves, Move{f, r, f, r + 1, ""})
	}
	// Forward-right
	if canMoveToSquare(*p, f+1, r+1, side) {
		moves = append(moves, Move{f, r, f + 1, r + 1, ""})
	}

	return moves
}

// canMoveToSquare evaluates if a piece of a specific color can occupy the square specified.
func canMoveToSquare(p Position, f, r int, side Side) bool {
	if f < 0 || f > 7 || r < 0 || f > 7 {
		return false
	}
	if p.board[r][f].piece.piece == None {
		return true
	} else if p.board[r][f].piece.color == side {
		return false
	} else {
		return true
	}
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
