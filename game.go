package main

import (
	"fmt"
	"regexp"
	"strings"
)

var moveExp = *regexp.MustCompile("(?P<file1>[a-h])(?P<rank1>[1-8])(?P<file2>[a-h])(?P<rank2>[1-8])(?P<promotionPiece>[bnrq])?")

// StartUserSession initiates a chess game with RobChess.
func StartUserSession() {
	fmt.Println("Welcome to RobChess! When entering moves, please use long algebraic chess notation.")

	game := *NewGame()
	side := promptColor()
	var printFunc func() string
	if side == White {
		printFunc = game.position.String
	} else {
		// TODO Switch back to pos.StringBlack
		printFunc = game.position.String
	}

	fmt.Println(printFunc())
	gameLoop(White, side, game, printFunc)
}

func gameLoop(side Side, playerSide Side, g GameContext, printFunc func() string) {
	var oppSide = side.OppSide()

	if side == playerSide {
		move := readMove()
		if ok := g.MakeMove(move); !ok {
			fmt.Printf("Something went wrong processing move: %+v\n", move)
			return
		}
		fmt.Println(printFunc())
		gameLoop(oppSide, playerSide, g, printFunc)
	} else {
		fmt.Println("Engine is thinking...")
		engineMove := Think(g, side)
		fmt.Printf("Engine Move: %v\n", engineMove)
		g.MakeMove(engineMove)
		fmt.Println(printFunc())
		// fmt.Printf("I think your moves are %v", pos.GetMoves(oppSide))
		gameLoop(oppSide, playerSide, g, printFunc)
	}
}

func promptColor() Side {
	fmt.Println("Choose color ('w' or 'b' accepted)")
	var input string
	fmt.Scanln(&input)
	if input == "w" {
		fmt.Println("You chose white.")
		return White
	} else if input == "b" {
		fmt.Println("You chose black")
		return Black
	} else {
		fmt.Println("Could not understand input.")
		return promptColor()
	}
}
func readMove() Move {
	fmt.Print("Move: ")
	moveStr := scanMove()
	if move, ok := algebraicToMove(moveStr); ok {
		return move
	}
	fmt.Println("The move entered could not be understood. Please enter a move in long algrebraic chess notation.")
	return readMove()
}

func scanMove() string {
	moveStr := ""
	fmt.Scan(&moveStr)
	return moveStr
}

func algebraicToMove(algMove string) (Move, bool) {
	algMove = strings.ToLower(algMove)
	match := moveExp.FindStringSubmatch(algMove)
	// match will include the full string enclosing the submatches as the first element. Therefore, we look
	// for a minimum of length 5.
	if len(match) >= 5 {
		// Convert each match subgroup to a rune array. The first element in the rune array is the file/rank.
		// 'a' to '0' is a difference of 48 in html codes.
		file1Rune := []rune(match[1])[0]
		file1 := int(file1Rune - 48 - 48 - 1)
		rank1Rune := []rune(match[2])[0]
		rank1 := int(rank1Rune - 48 - 1)
		file2Rune := []rune(match[3])[0]
		file2 := int(file2Rune - 48 - 48 - 1)
		rank2Rune := []rune(match[4])[0]
		rank2 := int(rank2Rune - 48 - 1)

		// Optional promotion component
		var promoPiece string
		if len(match) == 6 {
			promoPiece = match[5]
		}

		move, ok := NewMove(file1, rank1, file2, rank2, promoPiece)
		return move, ok
	}
	return Move{}, false
}

// Move represents a move on the chess board. It encompasses a piece, the old square and the new square.
type Move struct {
	oFile, oRank int
	nFile, nRank int
	promoPiece   string
}

// NewMove creates and initializes a new Move object.
func NewMove(oFile, oRank, nFile, nRank int, promoPiece string) (Move, bool) {
	fmt.Println(oFile, oRank, nFile, nRank, promoPiece)
	if oRank < 0 || oFile > 7 ||
		nRank < 0 || nFile > 7 {
		return Move{}, false
	}
	move := Move{
		oFile,
		oRank,
		nFile,
		nRank,
		promoPiece}
	return move, true
}

func (m Move) String() string {
	return fmt.Sprintf("Move: %d%d %d%d", m.oFile, m.oRank, m.nFile, m.nRank)
}
