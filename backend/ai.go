package backend

import (
	"math/rand"
)


// GetGreedyAIMove picks the best column for the AI to play based on a greedy strategy:
// 1. Win if possible.
// 2. Block opponent's winning move.
// 3. Otherwise, pick a random valid column.
func GetGreedyAIMove(b *Board) int {
	// Try to win
	for col := 0; col < Cols; col++ {
		if !b.IsValidMove(col) {
			continue
		}
		sim := b.Clone()
		sim.CurrentTurn = AI
		sim.Drop(col)
		if sim.CheckWin() == AI {
			return col
		}
	}

	// Try to block Player
	for col := 0; col < Cols; col++ {
		if !b.IsValidMove(col) {
			continue
		}
		sim := b.Clone()
		sim.CurrentTurn = Player
		sim.Drop(col)
		if sim.CheckWin() == Player {
			return col
		}
	}

	// Fallback: pick a random valid move
	validCols := []int{}
	for col := 0; col < Cols; col++ {
		if b.IsValidMove(col) {
			validCols = append(validCols, col)
		}
	}
	if len(validCols) == 0 {
		return -1
	}
	return validCols[rand.Intn(len(validCols))]
}

