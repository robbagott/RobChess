package main

import (
	"github.com/fatih/color"
)

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

// OppSide return the opposite side to the one given
func (s Side) OppSide() Side {
	if s == White {
		return Black
	}
	return White
}

func (p GamePiece) String() string {
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
		return 200
	default:
		return 1
	}
}
