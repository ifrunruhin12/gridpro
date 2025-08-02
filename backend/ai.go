package backend

import (
	"math/rand"
)

// GetAIMove selects the AI's move with a lookahead-1 strategy:
// 1. Win if possible
// 2. Block opponent's win
// 3. Avoid giving opponent an immediate win
// 4. Else, pick random valid move
func GetAIMove(board *Board) int {

	// 1. Try to win
	for col := 0; col < Cols; col++ {
		if !board.IsValidMove(col) {
			continue
		}
		sim := board.Clone()
		sim.Drop(col)
		if sim.CheckWin() == AI {
			return col
		}
	}

	// 2. Block human
	for col := 0; col < Cols; col++ {
		if !board.IsValidMove(col) {
			continue
		}
		sim := board.Clone()
		sim.CurrentTurn = Player
		sim.Drop(col)
		if sim.CheckWin() == Player {
			return col
		}
	}

	// 3. Safe move — avoid letting opponent win immediately
	safeMoves := []int{}
	for col := 0; col < Cols; col++ {
		if !board.IsValidMove(col) {
			continue
		}
		sim := board.Clone()
		sim.Drop(col)

		opponentCanWin := false
		for oppCol := 0; oppCol < Cols; oppCol++ {
			if !sim.IsValidMove(oppCol) {
				continue
			}
			oppSim := sim.Clone()
			oppSim.CurrentTurn = Player
			oppSim.Drop(oppCol)
			if oppSim.CheckWin() == Player {
				opponentCanWin = true
				break
			}
		}

		if !opponentCanWin {
			safeMoves = append(safeMoves, col)
		}
	}

	if len(safeMoves) > 0 {
		return safeMoves[rand.Intn(len(safeMoves))]
	}

	// 4. No safe moves? Pick random valid one
	validMoves := []int{}
	for col := 0; col < Cols; col++ {
		if board.IsValidMove(col) {
			validMoves = append(validMoves, col)
		}
	}

	if len(validMoves) > 0 {
		return validMoves[rand.Intn(len(validMoves))]
	}

	return -1
}

