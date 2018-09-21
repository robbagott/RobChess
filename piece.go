package main

import "github.com/fatih/color"

// Piece represents a type of chess piece.
type Piece int

// Pawn, Rook, Knight, Bishop, Queen, and King are the values for a piece. None is provided for empty squares.
const (
	Pawn Piece = iota
	Rook
	Knight
	Bishop
	Queen
	King
	None
)

// Side represents a side in chess (white or black).
type Side int

// White and Black are the values for a piece color.
const (
	White Side = iota
	Black
)

// GamePiece represents a piece in a chess game. E.g. a black bishop.
type GamePiece struct {
	piece Piece
	color Side
}

func (p *GamePiece) String() string {
	var pieceStr string
	switch p.piece {
	case Pawn:
		pieceStr = "P"
	case Rook:
		pieceStr = "R"
	case Knight:
		pieceStr = "N"
	case Bishop:
		pieceStr = "B"
	case Queen:
		pieceStr = "Q"
	case King:
		pieceStr = "K"
	default:
		pieceStr = " "
	}

	if p.color == White {
		return color.GreenString(pieceStr)
	}
	return color.RedString(pieceStr)
}
