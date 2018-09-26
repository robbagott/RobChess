package main

import (
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

// Copy makes a copy of a position
func (p Position) Copy() *Position {
	newPos := Position{make([][]Square, len(p.board))}
	for i := range newPos.board {
		newPos.board[i] = make([]Square, len(p.board[i]))
		copy(newPos.board[i], p.board[i])
	}
	return &newPos
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

func (p Position) String() string {
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

// StringBlack prints the board from black perspective.
// TODO Currently breaks due to color formatting.
func (p Position) StringBlack() (result string) {
	for _, v := range p.String() {
		result = string(v) + result
	}
	return
}

// GetMoves returns the set of moves that are possible for the side indicated.
func (p *Position) GetMoves(side Side) []Move {
	moves := make([]Move, 0, 20)

	// Only pieces can make moves in chess, so we iterate through the board and check for pieces.
	// If a piece is found, we then check for legal moves for that piece.
	for r := range p.board {
		for f := range p.board[r] {
			if p.board[r][f].piece.color == side {
				moves = append(moves, p.GetMovesAt(f, r)...)
			}

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

// TODO account for en passant. Investigate if the engine knows how to move black pawns.
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

	// TODO remove and teach about promotion
	if r == 7 && side == White || r == 0 && side == Black {
		return moves
	}

	// Possible forward moves
	if r == 1 && side == White || r == 6 && side == Black {
		if p.board[r+rIncr][f].piece.piece == None {
			moves = append(moves, Move{f, r, f, r + rIncr, ""})
			if p.board[r+rIncr*2][f].piece.piece == None {
				moves = append(moves, Move{f, r, f, r + rIncr*2, ""})
			}
		}
	} else if r > 1 && side == White || r < 6 && side == Black {
		if p.board[r+rIncr][f].piece.piece == None {
			moves = append(moves, Move{f, r, f, r + rIncr, ""})
		}
	}

	// Possible captures
	if f-1 >= 0 {
		if p.board[r+rIncr][f-1].piece.piece != None && p.board[r+rIncr][f-1].piece.color != side {
			moves = append(moves, Move{f, r, f - 1, r + rIncr, ""})
		}
	}
	if f+1 <= 7 {
		if p.board[r+rIncr][f+1].piece.piece != None && p.board[r+rIncr][f+1].piece.color != side {
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
		if canMove, capture := canMoveToSquare(*p, i, r, side); canMove {
			moves = append(moves, Move{f, r, i, r, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}

	// Look left
	for i := f - 1; i >= 0; i-- {
		if canMove, capture := canMoveToSquare(*p, i, r, side); canMove {
			moves = append(moves, Move{f, r, i, r, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}

	// Look forward
	for i := r + 1; i <= 7; i++ {
		if canMove, capture := canMoveToSquare(*p, f, i, side); canMove {
			moves = append(moves, Move{f, r, f, i, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}

	// Look backward
	for i := r - 1; i >= 0; i-- {
		if canMove, capture := canMoveToSquare(*p, f, i, side); canMove {
			moves = append(moves, Move{f, r, f, i, ""})
			if capture {
				break
			}
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
	for i, j := r+1, f+1; i < 8 && j < 8; i, j = i+1, j+1 {
		if canMove, capture := canMoveToSquare(*p, j, i, side); canMove {
			moves = append(moves, Move{f, r, j, i, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}
	// Look diagonally forward-left
	for i, j := r+1, f-1; i < 8 && j >= 0; i, j = i+1, j-1 {
		if canMove, capture := canMoveToSquare(*p, j, i, side); canMove {
			moves = append(moves, Move{f, r, j, i, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}
	// Look diagonally backward-right
	for i, j := r-1, f+1; i >= 0 && j < 8; i, j = i-1, j+1 {
		if canMove, capture := canMoveToSquare(*p, j, i, side); canMove {
			moves = append(moves, Move{f, r, j, i, ""})
			if capture {
				break
			}
		} else {
			break
		}
	}
	// Look diagonally backward-left
	for i, j := r-1, f-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if canMove, capture := canMoveToSquare(*p, j, i, side); canMove {
			moves = append(moves, Move{f, r, j, i, ""})
			if capture {
				break
			}
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
	if canMove, _ := canMoveToSquare(*p, f+2, r+1, side); canMove {
		moves = append(moves, Move{f, r, f + 2, r + 1, ""})
	}
	if canMove, _ := canMoveToSquare(*p, f+2, r-1, side); canMove {
		moves = append(moves, Move{f, r, f + 2, r - 1, ""})
	}

	// Left L moves
	if canMove, _ := canMoveToSquare(*p, f-2, r+1, side); canMove {
		moves = append(moves, Move{f, r, f - 2, r + 1, ""})
	}
	if canMove, _ := canMoveToSquare(*p, f-2, r-1, side); canMove {
		moves = append(moves, Move{f, r, f - 2, r - 1, ""})
	}

	// Forward L moves
	if canMove, _ := canMoveToSquare(*p, f+1, r+2, side); canMove {
		moves = append(moves, Move{f, r, f + 1, r + 2, ""})
	}
	if canMove, _ := canMoveToSquare(*p, f-1, r+2, side); canMove {
		moves = append(moves, Move{f, r, f - 1, r + 2, ""})
	}

	// Backward L moves
	if canMove, _ := canMoveToSquare(*p, f+1, r-2, side); canMove {
		moves = append(moves, Move{f, r, f + 1, r - 2, ""})
	}
	if canMove, _ := canMoveToSquare(*p, f-1, r-2, side); canMove {
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
	if canMove, _ := canMoveToSquare(*p, f+1, r, side); canMove {
		moves = append(moves, Move{f, r, f + 1, r, ""})
	}
	// Back-right
	if canMove, _ := canMoveToSquare(*p, f+1, r-1, side); canMove {
		moves = append(moves, Move{f, r, f + 1, r - 1, ""})
	}
	// Back
	if canMove, _ := canMoveToSquare(*p, f, r-1, side); canMove {
		moves = append(moves, Move{f, r, f, r - 1, ""})
	}
	// Back-left
	if canMove, _ := canMoveToSquare(*p, f-1, r-1, side); canMove {
		moves = append(moves, Move{f, r, f - 1, r - 1, ""})
	}
	// Left
	if canMove, _ := canMoveToSquare(*p, f-1, r, side); canMove {
		moves = append(moves, Move{f, r, f - 1, r, ""})
	}
	// Forward-left
	if canMove, _ := canMoveToSquare(*p, f-1, r+1, side); canMove {
		moves = append(moves, Move{f, r, f - 1, r + 1, ""})
	}
	// Forward
	if canMove, _ := canMoveToSquare(*p, f, r+1, side); canMove {
		moves = append(moves, Move{f, r, f, r + 1, ""})
	}
	// Forward-right
	if canMove, _ := canMoveToSquare(*p, f+1, r+1, side); canMove {
		moves = append(moves, Move{f, r, f + 1, r + 1, ""})
	}

	return moves
}

// canMoveToSquare evaluates if a piece of a specific color can occupy the square specified.
func canMoveToSquare(p Position, f, r int, side Side) (canMove, capture bool) {
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return false, false
	}

	if p.board[r][f].piece.piece == None {
		return true, false
	} else if p.board[r][f].piece.color == side {
		return false, false
	} else {
		return true, true
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
func (p *Position) SumMaterial(pieces []GamePiece) float64 {
	var sum float64
	for _, piece := range pieces {
		if piece.piece != None {
			sum += piece.Value()
		}
	}
	return sum
}

// MakeMove modifies the given position to represent the position after the move is made.
func (p *Position) MakeMove(move Move) bool {
	of, or := move.oFile, move.oRank
	nf, nr := move.nFile, move.nRank
	if of > 7 || of < 0 || or > 7 || or < 0 ||
		nf > 7 || nf < 0 || nr > 7 || nr < 0 {
		return false
	}

	// Handle promotion

	// Make normal move
	piece := p.board[or][of].piece
	p.board[or][of].piece = GamePiece{None, White}
	p.board[nr][nf].piece = piece
	return true
}
