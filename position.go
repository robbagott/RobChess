package main

import (
	"fmt"
	"strconv"
)

// Square represents a square in a chess position. Squares can have a piece placed on them.
type Square struct {
	file int
	rank int
}

// Position represents a chess position representation.
type Position struct {
	board               [][]GamePiece
	canCastleLongWhite  bool
	canCastleShortWhite bool
	canCastleLongBlack  bool
	canCastleShortBlack bool
}

// NewPosition creates and initializes a new Position with the starting arrangement of pieces.
func NewPosition() *Position {
	board := make([][]GamePiece, 8, 8)

	for i := range board {
		board[i] = make([]GamePiece, 8, 8)
	}

	pos := Position{board, true, true, true, true}

	pos.Reset()
	return &pos
}

// Copy makes a copy of a position
func (p Position) Copy() *Position {
	newPos := Position{make([][]GamePiece, len(p.board)), true, true, true, true}
	for i := range newPos.board {
		newPos.board[i] = make([]GamePiece, len(p.board[i]))
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
					p.board[r][f] = GamePiece{Rook, White}
				case f == 1 || f == 6:
					p.board[r][f] = GamePiece{Knight, White}
				case f == 2 || f == 5:
					p.board[r][f] = GamePiece{Bishop, White}
				case f == 3:
					p.board[r][f] = GamePiece{Queen, White}
				case f == 4:
					p.board[r][f] = GamePiece{King, White}
				}
			case r == 1:
				p.board[r][f] = GamePiece{Pawn, White}
			case r == 2 || r == 3 || r == 4 || r == 5:
				p.board[r][f] = GamePiece{None, White}
			case r == 6:
				p.board[r][f] = GamePiece{Pawn, Black}
			case r == 7:
				switch {
				case f == 0 || f == 7:
					p.board[r][f] = GamePiece{Rook, Black}
				case f == 1 || f == 6:
					p.board[r][f] = GamePiece{Knight, Black}
				case f == 2 || f == 5:
					p.board[r][f] = GamePiece{Bishop, Black}
				case f == 3:
					p.board[r][f] = GamePiece{Queen, Black}
				case f == 4:
					p.board[r][f] = GamePiece{King, Black}
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
			gamePiece := p.board[r][f]
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

	/* Only pieces can make moves in chess, so we iterate through the board and check for pieces.
	If a piece is found, we then check for legal moves for that piece. */
	for r := range p.board {
		for f := range p.board[r] {
			if p.board[r][f].color == side {
				moves = append(moves, p.GetMovesAt(f, r, side)...)
			}

		}
	}

	// Prune moves which lead to checks
	validMoves := make([]Move, 0, cap(moves))
	for _, move := range moves {
		if !causesCheck(p, move, side) {
			validMoves = append(validMoves, move)
		}
	}

	return validMoves
}

// GetMovesAt returns the set of moves that are possible for the piece located at file f and rank r
func (p *Position) GetMovesAt(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)
	piece := p.board[r][f]
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

func causesCheck(p *Position, move Move, side Side) bool {
	oldPiece := p.board[move.oRank][move.oFile]
	capturedPiece := p.board[move.nRank][move.nFile]
	toReturn := false

	p.MakeMove(move)
	kingSquare := getKingSquare(p, side)
	if inCheck(*p, kingSquare.file, kingSquare.rank, side) {
		toReturn = true
	}

	// Roll back move
	p.board[move.oRank][move.oFile] = oldPiece
	p.board[move.nRank][move.nFile] = capturedPiece
	return toReturn
}

func getKingSquare(p *Position, side Side) Square {
	for r := range p.board {
		for f := range p.board[r] {
			if p.board[r][f].piece == King && p.board[r][f].color == side {
				return Square{f, r}
			}
		}
	}
	panic(fmt.Sprintf("getKingSquare in %v returned no square for side %v.", p, side))
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
		if p.board[r+rIncr][f].piece == None {
			moves = append(moves, Move{f, r, f, r + rIncr, ""})
			if p.board[r+rIncr*2][f].piece == None {
				moves = append(moves, Move{f, r, f, r + rIncr*2, ""})
			}
		}
	} else if r > 1 && side == White || r < 6 && side == Black {
		if p.board[r+rIncr][f].piece == None {
			moves = append(moves, Move{f, r, f, r + rIncr, ""})
		}
	}

	// Possible captures
	if f-1 >= 0 {
		if p.board[r+rIncr][f-1].piece != None && p.board[r+rIncr][f-1].color != side {
			moves = append(moves, Move{f, r, f - 1, r + rIncr, ""})
		}
	}
	if f+1 <= 7 {
		if p.board[r+rIncr][f+1].piece != None && p.board[r+rIncr][f+1].color != side {
			moves = append(moves, Move{f, r, f + 1, r + rIncr, ""})
		}
	}

	return moves
}

// Get possible rook moves for a rook located at file r and rank f of color side.
func (p *Position) getRookMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)

	// Look right
	squares, piece := p.lookRight(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look left
	squares, piece = p.lookLeft(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look up
	squares, piece = p.lookUp(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look down
	squares, piece = p.lookDown(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}
	return moves
}

// Get possible bishop moves for a bishop located at file r and rank f of color side.
func (p *Position) getBishopMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 20)

	// Look up right
	squares, piece := p.lookUpRight(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look up right
	squares, piece = p.lookUpLeft(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look down-right
	squares, piece = p.lookDownRight(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}

	// Look down-left
	squares, piece = p.lookDownLeft(f, r)
	for i, s := range squares {
		if i == len(squares)-1 && piece.piece != None && piece.color == side {
			continue
		}
		moves = append(moves, Move{f, r, s.file, s.rank, ""})
	}
	return moves
}

// Get possible queen moves for a queen located at file r and rank f of color side.
func (p *Position) getKnightMoves(f, r int, side Side) []Move {
	moves := make([]Move, 0, 8)
	squares := p.lookL(f, r)
	for _, square := range squares {
		if canMove, _ := canMoveToSquare(*p, square.file, square.rank, side); canMove {
			moves = append(moves, Move{f, r, square.file, square.rank, ""})
		}
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

	// TODO Castle

	return moves
}

/* lookUp and other look functions look in a direction on the board from a starting square. When another piece is encountered, the function returns
with the squares traversed and the collision piece. */
func (p *Position) lookUp(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i := r + 1; i < 8; i++ {
		squares = append(squares, Square{f, i})
		if p.board[i][f].piece != None {
			piece = p.board[i][f]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookUpRight(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i, j := r+1, f+1; i < 8 && j < 8; i, j = i+1, j+1 {
		squares = append(squares, Square{j, i})
		if p.board[i][j].piece != None {
			piece = p.board[i][j]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookRight(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i := f + 1; i < 8; i++ {
		squares = append(squares, Square{i, r})
		if p.board[r][i].piece != None {
			piece = p.board[r][i]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookDownRight(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i, j := r-1, f+1; i >= 0 && j < 8; i, j = i-1, j+1 {
		squares = append(squares, Square{j, i})
		if p.board[i][j].piece != None {
			piece = p.board[i][j]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookDown(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i := r - 1; i >= 0; i-- {
		squares = append(squares, Square{f, i})
		if p.board[i][f].piece != None {
			piece = p.board[i][f]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookDownLeft(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i, j := r-1, f-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		squares = append(squares, Square{j, i})
		if p.board[i][j].piece != None {
			piece = p.board[i][j]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookLeft(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i := f - 1; i >= 0; i-- {
		squares = append(squares, Square{i, r})
		if p.board[r][i].piece != None {
			piece = p.board[r][i]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookUpLeft(f, r int) ([]Square, *GamePiece) {
	squares := make([]Square, 0)
	piece := GamePiece{None, White}
	for i, j := r+1, f-1; i < 8 && j >= 0; i, j = i+1, j-1 {
		squares = append(squares, Square{j, i})
		if p.board[i][j].piece != None {
			piece = p.board[i][j]
			break
		}
	}
	return squares, &piece
}

func (p *Position) lookL(f, r int) []Square {
	squares := make([]Square, 0)

	// Right L moves
	if f+2 < 8 && r+1 < 8 {
		squares = append(squares, Square{f + 2, r + 1})
	}
	if f+2 < 8 && r-1 >= 0 {
		squares = append(squares, Square{f + 2, r - 1})
	}

	// Left L moves
	if f-2 >= 0 && r+1 < 8 {
		squares = append(squares, Square{f - 2, r + 1})
	}
	if f-2 >= 0 && r-1 >= 0 {
		squares = append(squares, Square{f - 2, r - 1})
	}

	// Forward L moves
	if f+1 < 8 && r+2 < 8 {
		squares = append(squares, Square{f + 1, r + 2})
	}
	if f-1 >= 0 && r+2 < 8 {
		squares = append(squares, Square{f - 1, r + 2})
	}

	// Backward L moves
	if f+1 < 8 && r-2 >= 0 {
		squares = append(squares, Square{f + 1, r - 2})
	}
	if f-1 >= 0 && r-2 >= 0 {
		squares = append(squares, Square{f - 1, r - 2})
	}

	return squares
}

// canMoveToSquare evaluates if a piece of a specific color can occupy the square specified.
func canMoveToSquare(p Position, f, r int, side Side) (canMove, capture bool) {
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return false, false
	}

	if p.board[r][f].piece == None {
		return true, false
	} else if p.board[r][f].color == side {
		return false, false
	} else {
		return true, true
	}
}

// inCheck checks if the square at f, r is attacked by opponent's pieces.
func inCheck(p Position, f, r int, side Side) bool {
	if _, piece := p.lookUp(f, r); piece.color == side.OppSide() && (piece.piece == Rook || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookUpRight(f, r); piece.color == side.OppSide() && (piece.piece == Bishop || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookRight(f, r); piece.color == side.OppSide() && (piece.piece == Rook || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookDownRight(f, r); piece.color == side.OppSide() && (piece.piece == Bishop || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookDown(f, r); piece.color == side.OppSide() && (piece.piece == Rook || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookDownLeft(f, r); piece.color == side.OppSide() && (piece.piece == Bishop || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookLeft(f, r); piece.color == side.OppSide() && (piece.piece == Rook || piece.piece == Queen) {
		return true
	}
	if _, piece := p.lookUpLeft(f, r); piece.color == side.OppSide() && (piece.piece == Bishop || piece.piece == Queen) {
		return true
	}

	for _, square := range p.lookL(f, r) {
		if p.board[square.rank][square.file].color == side.OppSide() && p.board[square.rank][square.file].piece == Knight {
			return true
		}
	}
	return false
}

// GetPieces returns an array of pieces for the side indicated in a position.
func (p *Position) GetPieces(side Side) []GamePiece {
	pieces := make([]GamePiece, 0, 32)
	for r := range p.board {
		for f := range p.board[r] {
			piece := p.board[r][f]
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

// CentralControl returns a simple count of the number of squares in the center occupied by the specified side's pieces.
func (p *Position) CentralControl(side Side) float64 {
	count := 0.0
	for r := 2; r < 5; r++ {
		for f := 2; f < 5; f++ {
			if p.board[r][f].color == side {
				count++
			}
		}
	}
	return count
}

// MakeMove modifies the given position to represent the position after the move is made.
func (p *Position) MakeMove(move Move) bool {
	of, or := move.oFile, move.oRank
	nf, nr := move.nFile, move.nRank
	if of > 7 || of < 0 || or > 7 || or < 0 ||
		nf > 7 || nf < 0 || nr > 7 || nr < 0 {
		return false
	}

	// TODO handle promotion

	// Make normal move
	piece := p.board[or][of]
	p.board[or][of] = GamePiece{None, White}
	p.board[nr][nf] = piece
	return true
}
